package golem

import (
	"gonum.org/v1/gonum/mat"
	"fmt"
	"math"
	"sync"
)

var RANGE_SEARCH_SPACE_MAX_THRESHOLD int = 1000

////// PRESENT
/*
- work on tree pruning
- review CAUTION at AddNodeToTree
- PredictMC needs work!
*/

var DEFAULT_SPLIT_TYPE string = "range"

type DTree struct {
	samplr *Sampler
	excludeCols []string 
	outputCols []string
	bdts []*BDTree
	data *mat.Dense
	classOrder []float64
	classIndices map[int][]int
	bpm map[int]map[float64]map[string]float64

	weightedOutput bool 
}

func OneDTree(s *Sampler, excludeCols []string, outputCols []string) *DTree { 
	dt := &DTree{samplr: s, bdts: make([]*BDTree,0), excludeCols: excludeCols, outputCols: outputCols}  
	return dt 
}

/*
preprocesses Sampler
*/ 
func (dt *DTree) Preprocess(keyOrder string) {
	// samplr processes
	dt.samplr.RunSampling2(dt.excludeCols, dt.outputCols)
	/// TODO: check fetch sampling no refresh
	dt.FetchSamplingData(false)

	// determine outcome classes and partition frequencies
	/// TODO: delete this, run bayesian frequency instead
	if dt.samplr.imbal2 {
		dt.samplr.DetermineOutcomeClassesOfData("imbal")
	} else {
		dt.samplr.DetermineOutcomeClassesOfData("sample")
	}
	
	// assign to trees
	dt.classOrder = dt.GatherKeys(keyOrder)
	dt.AssignOutcomeClassesToBDT()
}

////// TODO: needs to follow ordering set by MC. 
/*
declares each BDTree and assigns a class to them
*/
func (dt *DTree) AssignOutcomeClassesToBDT() {

	if len(dt.samplr.outputColumnIndices) != 1 {
		panic(fmt.Sprintf("cannot perform decision tree on multi-column outcome classes %s", "sadfadfa"))
	}

	dt.bdts = make([]*BDTree,len(dt.classOrder) - 1)
	for i := 0; i < len(dt.classOrder) - 1; i++ {
		dt.bdts[i] = OneBDTree(i, dt.classOrder[i])

		if dt.weightedOutput {
			r,_ := dt.data.Dims() 
			dt.bdts[i].outcomeClassWeight = float64(r) / dt.samplr.outcomeDataFrequencies[DefaultFloat64ToString(dt.classOrder[i])] 
		}
	}
}

/*
*/ 
func (dt *DTree) IndexOfTreeOfClass(class float64) int { 
	for i,b := range dt.bdts {
		if b.outcomeClassValue == class {
			return i
		}
	}	
	return -1
}

//// TODO: new here. 
func SplitCostStandard(r float64, yx float64, x float64, notYX float64, notX float64, f func(float64) float64) (float64,float64) {

	var q float64
	if x == 0 {
		q = 0 
	} else {
		q = yx /x 
	}

	one := x / r * f(q)

	if notX == 0 {
		q = 0 
	} else {
		q = notYX / notX 
	}
	two := notX / r * f(q) 
	return two,one
}

/*
*/
func StandardCost(x float64) float64 {

	if x < 0.0 || x > 1.0 {
		panic("invalid prob. measure") 
	}

	if x > 1.0 - x {
		return 1.0 - x
	} 

	return x 
} 

/*
*/
func EntropyCost(x float64) float64 {
	if x == 0.0 || x == 1.0 {
		return 1.0 
	}

	return -x * math.Log(x) - (1.0 - x) * math.Log(1.0 - x) 
}

/*
*/
func (dt *DTree) Make(maxDepth int, splitFunc func(float64)float64) { 
	dt.classIndices = make(map[int][]int,0) 
	data := dt.data
	var indices []int
	for i,k := range dt.classOrder {
		k_ := dt.IndexOfTreeOfClass(k)

		if i == len(dt.classOrder) - 1 {
			break 
		}

		dt.MakeTree(data, k_, maxDepth, splitFunc)
		// update data here
		indices = make([]int,0) 
		r,_ := data.Dims() 

		for i := 0; i < r; i++ {
			row,_ := MatrixRowToFloat64Slice(data, i) 
			if row[dt.samplr.outputColumnIndices[0]] != dt.bdts[k_].outcomeClassValue {
				indices = append(indices, i) 
			}
		}
		data = MatrixCollectByIndices(data, indices)
	}
}

/*
*/
func (dt *DTree) PredictMC() *PredictionMetric {
	/// TODO: check
	data := dt.samplr.MatrixToGroupEncoding(dt.data, "spare") 
	r,_ := data.Dims()
	actual,_ := MatrixColToFloat64Slice(data, dt.samplr.outputColumnIndices[0])
	results := make([]float64, r)
	
	wg := sync.WaitGroup{}
	for i := 0; i < r; i++ {

		wg.Add(1) 
		go func(i int) {
			row,_ := MatrixRowToFloat64Slice(data, i)
			results[i] = dt.PredictOneMC(row)
		}(i)
		wg.Done()
	}
	wg.Wait() 

	return RunPredictionMetricMC(actual, results)
}

/// TODO: method's output scheme is based on ordering, 
///       look into others
/*
*samples for each tree*
elimination: NOT-samples eliminated after each tree
full: all samples run through each 

*output* 
return the first
return a random positive
*/
func (dt *DTree) PredictOneMC(x []float64) float64 { 
	co := dt.classOrder[:len(dt.classOrder) - 1]

	//// TODO: alternative prediction here 
	///possible := make([]float64,0)
	for _,k := range co {
		k_ := dt.IndexOfTreeOfClass(k)
		output := dt.bdts[k_].PredictOne(dt.bdts[k_].ruut, x)
		if output != -1 {
			///possible = append(possible, output)
			return output
		}
	}

	return dt.classOrder[len(dt.classOrder) - 1]
}

func (dt *DTree) GatherKeys(makeMode string) []float64 {
	keyOrder := make([]float64,0)
	oc := dt.samplr.CollectOutcomeClasses()

	switch {
		case makeMode == "in-order": 
			for _,c := range oc {
				f,_ := DefaultStringToFloat64(c)
				keyOrder = append(keyOrder, f) 
			}

		case makeMode == "increasing" || makeMode == "decreasing": 
			// determine the size of each class
			x := SortKeysByValueMapStringToFloat64(dt.samplr.outcomeDataFrequencies, makeMode)
			for _,x_ := range x {
				f,_ := DefaultStringToFloat64(x_.a.(string))
				keyOrder = append(keyOrder, f) 
			}
			
		case makeMode == "random":
			panic("random mode not yet implemented!")
		default: 
			panic("invalid make-mode")
	}

	return keyOrder
}

// TODO: make depth value for node
func (dt *DTree) MakeTree(data *mat.Dense, treeIndex int, maxDepth int, splitFunc func(float64)float64) { 
	///// used for termination by 0-delta size
	/*
	var quit bool = true 
	var lastSize int = 0 /// 
	var roundsUnchanged int = 0 ///
	*/

	/// clear tree 
	dt.bdts[treeIndex].ruut = nil

	for {

		///quit = true 		
		dt.bdts[treeIndex].OpenNodes()

		// terminating condition 
		if len(dt.bdts[treeIndex].openNodes) == 0 {
			break
		}
		///fmt.Println("\t\tITERATING")
		// add to open nodes open nodes
		for _,x := range dt.bdts[treeIndex].openNodes {
			
			if x != nil {
				///// used for termination problems
				/*
				fmt.Println("depth of node ", x.key, "  ", x.depth)
				fmt.Println("SIZE ", len(x.leftIndices), " ", len(x.rightIndices))
				fmt.Println("LR ", x.leftFin, x.rightFin)
				*/
				if x.depth >= maxDepth {
					x.leftFin,x.rightFin = true,true
					continue
				}
			}
			dt.SplitNode(x, splitFunc, treeIndex)
			///quit = false 
		}
		
		///// TODO: delete below, terminating conditions
		/*
		if dt.bdts[treeIndex].nodeKeyAI == lastSize {
			roundsUnchanged++ 
			if roundsUnchanged > 15 {
				break 
			}
		} else {
			roundsUnchanged = 0 
			lastSize = dt.bdts[treeIndex].nodeKeyAI
		}
		
		
		/*
		////
		if dt.bdts[treeIndex].Size() >= 1 {
			break
		}
		////
		*/
	}
}

/*
*/ 
func (dt *DTree) AddNodeToTree(sd *SplitData, treeIndex int) bool {

	if sd == nil {
		return false
	}

	if dt.AddNodeToTree_(sd, treeIndex) {
		dt.bdts[treeIndex].nodeKeyAI++ 
		return true
	}
	return false 
} 

/*
*/
func (dt *DTree) AddNodeToTree_(sd *SplitData, treeIndex int) bool { 
	
	if sd.node != nil {
		if !dt.SplitIsGood(sd) {
			return false
		}
	}

	dt.MakeNewNode(sd, treeIndex)
	return true
}

/*
*/
func (dt *DTree) MakeNewNode(sd *SplitData, treeIndex int) { 

	var r, depth int
	var targetIndices, availableIndices []int
	
	switch {
	case sd.node == nil: 
		r,_ = dt.data.Dims() 
		targetIndices = IntRange(0,r)
		availableIndices = IntSliceDifference(dt.samplr.InputColumnIndices(), []int{sd.column})
		depth = 0 
	default: 
		if sd.inheritance == "l" {
			r = len(sd.node.leftIndices)
			targetIndices = sd.node.leftIndices
		} else {
			r = len(sd.node.rightIndices)
			targetIndices = sd.node.rightIndices
		}

		availableIndices = IntSliceDifference(sd.node.availableColumnIndices, []int{sd.column})
		depth = sd.node.depth + 1
	}

	n := OneBDNode(dt.bdts[treeIndex].nodeKeyAI, sd, float64(r), depth) 
	n.availableColumnIndices = availableIndices
	targetColumn,_ := MatrixColToFloat64Slice(dt.data, sd.column) 

	switch {
		case dt.samplr.colDC[sd.column] == "continuous": 
			ps := InterfaceSliceToPairSlice(sd.split.splitDC)
			n.leftIndices, n.rightIndices = dt.IndicesOfColumnInQualRange(targetColumn, targetIndices, sd.column, ps)
		default: 
			ps := InterfaceSliceToFloatSlice(sd.split.splitDC)
			n.leftIndices, n.rightIndices = dt.IndicesOfColumnInQualRangeDiscrete(targetColumn, targetIndices, sd.column, ps)
	}

	switch {
	case sd.node == nil: 
		dt.bdts[treeIndex].ruut = n 

	default: 
		if sd.inheritance == "l" {
			sd.node.left = n 
		} else {
			sd.node.right = n 
		}
		n.parent = sd.node
	}
}

/*
*/ 
func (dt *DTree) SplitIsGood(sd *SplitData) bool { 
	if sd.node == nil {
		panic("comparison on nil nodilanco!")
	}

	if sd.inheritance == "l" {
		if sd.node.leftMismatch < sd.leftMismatch + sd.rightMismatch {
			sd.node.leftFin = true 
			return false
		}
	} else {
		if sd.node.rightMismatch < sd.leftMismatch + sd.rightMismatch {
			sd.node.rightFin = true 
			return false
		}
	}

	return true
}

/*
*/
func (dt *DTree) SplitNode(n *BDNode, splitFunc func(float64)float64, treeIndex int) {
	var s1, s2 *SplitData
	if n == nil {
		s2 = dt.SplitNodeChild(n, "r", splitFunc, treeIndex, 0)
		dt.AddNodeToTree(s2, treeIndex)
		return 
	}

	if n.left == nil && !n.leftFin {
		s1 = dt.SplitNodeChild(n, "l", splitFunc, treeIndex, 0)
		dt.AddNodeToTree(s1, treeIndex)
	}

	if n.right == nil && !n.rightFin {
		s2 = dt.SplitNodeChild(n, "r", splitFunc, treeIndex, 0)
		dt.AddNodeToTree(s2, treeIndex)
	}
}

/*
*/
func (dt *DTree) SplitNodeChild(n *BDNode, inheritance string, splitFunc func(float64)float64, treeIndex int, verbose int) *SplitData {

	if verbose < 0 || verbose > 2 {
		panic("invalid verbose mode")
	}

	// fetch remaining column indices to consider
	var availableIndices []int
	if n != nil {
		availableIndices = n.availableColumnIndices
	} else {
		availableIndices = dt.samplr.InputColumnIndices()
	}
	data := dt.CollectSamplingData(n, inheritance)

	///// TODO: mod here!
	var bestP *SplitData

	if data == nil || len(availableIndices) == 0 {
		switch {
		case inheritance == "l": 
			n.leftFin = true 
		default: 
			n.rightFin = true 
		}
		return bestP
	}

	r,_ := data.Dims()

	/// TODO: modification here 
	bpm := dt.samplr.BayesProbOverInputAllColumns(data)
	for _,c := range availableIndices {
		p, ml, mr := dt.RandomGreedyRangeFinder(bpm[c], splitFunc, DefaultFloat64ToString(dt.bdts[treeIndex].outcomeClassValue), treeIndex, c, float64(r)) 
		if p == nil {
			continue
		}

		sd := SplitPairDataToSplitData(p, n, inheritance, ml, mr)
		if verbose == 2 {
			if n != nil {
				fmt.Println("\t node ", n.key, " split")
			} else {
				fmt.Println("\t root node split")
			}
			fmt.Println("best split for column ", c)
			fmt.Println("split line: ", " left ", sd.leftScore, " right ", sd.rightScore) /// sd.split.floatSplit,
			fmt.Println("---------------------------------")
		}
		
		switch {
		case bestP == nil:
			bestP = sd
		default:
			///$ TODO : check comparison method
			if sd.leftScore + sd.rightScore < bestP.leftScore + bestP.rightScore {
				bestP = sd 
			} 
		}
	}

	if verbose > 0 {
		if bestP != nil {
			fmt.Println()
			fmt.Println("\tSPLIT DATA:")
			bestP.split.Display()
			fmt.Println("\tOTHER:")
			fmt.Println("left score ", bestP.leftScore, " \trightscore ", bestP.rightScore) 
			fmt.Println("left mismatch ", bestP.leftMismatch, " \trightscore ", bestP.rightMismatch) 
			fmt.Println("inheritance ", bestP.inheritance, " \tcolumn ", bestP.column)
			fmt.Println("$$$$$")
		} else {
			fmt.Println("best split is nil")
		}
	}
	
	if bestP == nil {
		switch {
			case inheritance == "l": 
				n.leftFin = true 
			default: 
				n.rightFin = true 
			}
	}

	return bestP 
}

/*
*/
func (dt *DTree) CollectSamplingData(x *BDNode, inheritance string) *mat.Dense {

	var data *mat.Dense
	if x != nil {
		switch {
			case inheritance == "l":
				data = MatrixCollectByIndices(dt.data, x.leftIndices)
			default: 
				data = MatrixCollectByIndices(dt.data, x.rightIndices)
			}
	} else {
		data = dt.data
	}

	return data 
} 

/*
*/
func (dt *DTree) FetchSamplingData(refresh bool) {
	switch {
	case dt.samplr.imbal2:
		dt.samplr.MakeDataForImbalanced(false, refresh)
		dt.data = dt.samplr.imbal2SamplingData 
	default:
		if refresh {
			dt.samplr.DrawSamplesConcurrent(false, "uniform", dt.samplr.sampleSize)
		}
		dt.data = dt.samplr.samplingData
	}
}