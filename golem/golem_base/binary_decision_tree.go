package golem 

import (
	"gonum.org/v1/gonum/mat"
	"fmt"
)

type BDTree struct {
	identifier int
	nodeKeyAI int 
	ruut *BDNode

	costError float64

	outcomeClassValue float64 
	outcomeClassWeight float64
	openNodes []*BDNode
	values []int
	size int
}

/*
*/ 
func OneBDTree(identifier int, outcomeClassValue float64) *BDTree {
	return &BDTree{identifier: identifier, outcomeClassValue: outcomeClassValue, openNodes: []*BDNode{}, outcomeClassWeight: 1.0}  
}

///////////////////////////////////// start: data processing 

/// TODO: design choice, could also store the indices at each node initially for future use
/*
*/
///// TODO: delete this, wrong!
func (bdt *BDTree) InOrderTraversalForData(n *BDNode, m *mat.Dense, indices []int, captureOpen bool) { // 
	
	if n == nil {
		return
	}

	// collect all the values that match
		/// TODO: 2 iterations for collection inefficient 
	column, stat := MatrixColToFloat64Slice(m, n.value)
	column = Float64SliceCollectByIndices(column, indices, false)

	if !stat {
		panic("invalid key for node, data column does not exist!")
	}

	//////// version 0: this version will run calculations only on leaves
	if n.left == nil && n.right == nil {
		qualIndices := bdt.ProcessDataAtNode(n, column, indices)
		n.rightIndices = qualIndices
		n.leftIndices = BayesianIdentityOfIndexSubset(&Pair{a: 0, b: len(column)}, qualIndices, true)
	} else {
		bdt.InOrderTraversalForData(n.left, m, n.leftIndices, captureOpen)
		bdt.InOrderTraversalForData(n.right, m, n.rightIndices, captureOpen)
	}

	if ((n.right == nil && len(n.rightIndices) != 0) || (n.left == nil && len(n.leftIndices) != 0)) && captureOpen {
		bdt.openNodes = append(bdt.openNodes, n) 
	} 

	bdt.values = append(bdt.values, n.value)
}

/*
*/ 
func (bdt *BDTree) InOrderTraversalForData_(m *mat.Dense, captureOpen bool) {

	// make indices for 
	r,_ := m.Dims() 
	is := IntRange(0, r)
	bdt.openNodes = make([]*BDNode,0)
	bdt.values = make([]int,0)

	if bdt.ruut == nil {
		var n *BDNode
		bdt.openNodes = append(bdt.openNodes, n)
	}

	bdt.InOrderTraversalForData(bdt.ruut, m, is, captureOpen)
}
 
/*
*/
func (bdt *BDTree) ProcessDataAtNode(n *BDNode, column []float64, sortedQualIndices []int) []int {
	return bdt.ProcessDataWithFunction(n.theFunction, column, sortedQualIndices)  
}

/*
*/
func (bdt *BDTree) ProcessDataWithFunction(f func(float64)bool, column []float64, sortedQualIndices []int) []int {
	qualIndices := make([]int,0) 

	for i, c := range column {

		if len(sortedQualIndices) == 0 {
			break
		}

		if i != sortedQualIndices[0] {
			continue
		}

		if f(c) {
			qualIndices = append(qualIndices, i) 
			sortedQualIndices = sortedQualIndices[1:]
		}
	}

	return qualIndices 
}

///////////////////////////////////// end: data processing 

/*
*/
func (bdt *BDTree) OperatePreOrderTraversal(n *BDNode, f func(n *BDNode)) {
	if n == nil {
		return
	}

	f(n) 

	bdt.OperatePreOrderTraversal(n.left, f) 
	bdt.OperatePreOrderTraversal(n.right, f) 
}

/*
*/
func (bdt *BDTree) OperatePreOrderTraversal_(f func(n *BDNode)) {
	// clear all data first 
	bdt.openNodes = make([]*BDNode, 0) 
	bdt.values = make([]int, 0)
	bdt.costError = 0
	bdt.OperatePreOrderTraversal(bdt.ruut, f) 
}

/*
*/
func (b *BDTree) DisplayNodeData(n *BDNode) {
	fmt.Println("key: ", n.key, " value: ", n.value) 
	fmt.Println("left size: ", len(n.leftIndices)) 
	fmt.Println("left err: ", n.leftErr)
	fmt.Println("left mismatch: ", n.leftMismatch)
	fmt.Println("right size: ", len(n.rightIndices)) 
	fmt.Println("right err: ", n.rightErr)
	fmt.Println("right mismatch: ", n.rightMismatch)
	fmt.Println("size: ", n.dataSize)
	fmt.Println("size change: ", n.dataSizeChange) 
	fmt.Println("\tsplit")
	n.split.Display()
	fmt.Println("-----------------------------------")
}

/*
*/
func (bdt *BDTree) SizeAt(n *BDNode) int {
	f := func(n *BDNode) {bdt.size++}
	bdt.size = 0
	bdt.OperatePreOrderTraversal(n, f)
	return bdt.size 
}

/*
*/
func (bdt *BDTree) Size() int {
	return bdt.SizeAt(bdt.ruut)
}

/*
*/
func (bdt *BDTree) MaxDepthOfOpenNodes() int {
	depth := 0
	for _, n := range bdt.openNodes {
		if n.depth > depth {
			depth = n.depth
		}
	}
	return depth
}

/*
*/
func (bdt *BDTree) OpenNodes() {
	bdt.openNodes = make([]*BDNode,0)
	if bdt.ruut == nil {
		var n *BDNode
		bdt.openNodes = append(bdt.openNodes, n) 
		return
	}

	bdt.OperatePreOrderTraversal(bdt.ruut, bdt.FetchOpenNode)

}

/*
*/
func (bdt *BDTree) FetchOpenNode(n *BDNode) {

	if n.leftFin && n.rightFin {
		return
	}

	if (n.left == nil && len(n.leftIndices) > 1) || (n.right == nil && len(n.rightIndices) > 1) {
		bdt.openNodes = append(bdt.openNodes, n)
	}

	if n.left != nil {
		n.leftFin = true 
	} 

	if len(n.leftIndices) == 0 {
		n.leftFin = true
	}

	if n.right != nil {
		n.rightFin = true
	}

	if len(n.rightIndices) == 0 {
		n.rightFin = true
	}
	
} 

/*
*/
func (bdt *BDTree) FetchValues(n *BDNode) {
	bdt.values = append(bdt.values, n.value)
}

/// TODO: unnecessary
/*
*/
func (bdt *BDTree) FetchError(n *BDNode) {
	if n == nil {
		return
	}

	if n.right == nil {
		bdt.costError += n.rightErr
	}

	if n.left == nil {
		bdt.costError += n.leftErr
	}
}

/// TODO refactor w/ ProcessDataAtNode
/*
*/
/////
func (b *BDTree) Predict(m *mat.Dense, predictOutput string) ([]float64, []int) {
	r,_ := m.Dims() 
	results := make([]float64,0) 
	ri := make([]int,0) 
	for i := 0; i < r; i++ {
		v,_ := MatrixRowToFloat64Slice(m, i)
		output := b.PredictOne(b.ruut, v) 

		switch {

		case predictOutput == "full": 
			results = append(results, output)
		case predictOutput == "pos-index": 
			if output != -1 {
				ri = append(ri, i)
			}
		case predictOutput == "neg-index": 
			if output == -1 {
				ri = append(ri, i)
			}
		default: 
			panic("invalid predict output")
		}

	}

	return results, ri 
}

/// look to refactoring this with `ProcessDataWithFunction` above. 
func (b *BDTree) PredictOne(n *BDNode, x []float64) float64 {

	if n == nil {
		panic("cannot perform prediction with null node")
	}

	right := float64(-1) /// CAUTION: default NOT-value is -1
	if n.theFunction(x[n.value]) {
		right = b.outcomeClassValue
	}

	switch {
		case right != b.outcomeClassValue:  
			if n.left == nil {
				return right
			}
			return b.PredictOne(n.left, x) 
		default: 
			if n.right == nil {
				return right
			}
			return b.PredictOne(n.right, x)
	}
}

/*
*/
func (b *BDTree) MeasureDataSizeChange() {
	b.OperatePreOrderTraversal(b.ruut, b.DataSizeChange) 
}

/*
*/ 
func (b *BDTree) DataSizeChange(n *BDNode) {
	if n == nil {
		return
	}

	if n != b.ruut {
		n.dataSizeChange = n.parent.dataSize - n.dataSize
	}
}

///////////////////////////////////////////// prediction metric

/*
*/
type PredictionMetric struct {
	totalError float64
	errorPerClass map[float64]float64
	classCounts map[float64]float64
}

/*
*/
func (pm *PredictionMetric) Display() {
	fmt.Println("* total error: ", pm.totalError)
	fmt.Println("* error per class") 
	fmt.Println(pm.errorPerClass)
	fmt.Println() 
	fmt.Println("* class counts") 
	fmt.Println(pm.classCounts)
}

/*
prediction metric is used for binary classification purposes
*/
func RunPredictionMetric(actual []float64, predicted []float64, outcomeClassValue float64) *PredictionMetric {

	if outcomeClassValue == -1 {
		panic("value -1 already occupied!")
	}

	if len(predicted) != len(actual) {
		panic("predicted and actual must be correct size")
	}

	totalError := float64(0)
	errorPerClass := make(map[float64]float64)
	classCounts := make(map[float64]float64, 0)

	for i,x := range actual {
		switch {
		case x == outcomeClassValue: 
			classCounts[x]++
			if predicted[i] != x {
				errorPerClass[x]++
				totalError++ 
			}
		
		default: 
			classCounts[float64(-1)]++
			if predicted[i] != float64(-1) {
				errorPerClass[float64(-1)]++
				totalError++
			}
		}
	}

	totalError /= float64(len(actual))
	
	for k,v := range classCounts {
		errorPerClass[k] = ZeroDivStandard(errorPerClass[k],v) 
	}
	
	return &PredictionMetric{totalError: totalError, errorPerClass: errorPerClass, classCounts: classCounts}
}

/*
multi-class classification
*/
func RunPredictionMetricMC(actual []float64, predicted []float64) *PredictionMetric {
	if len(predicted) != len(actual) {
		panic("predicted and actual must be correct size")
	}

	totalError := float64(0)
	errorPerClass := make(map[float64]float64)
	classCounts := make(map[float64]float64, 0)

	for i,x := range actual {
		classCounts[x]++
		if predicted[i] != x {
			errorPerClass[x]++
			totalError++ 
		}
	}

	totalError /= float64(len(actual)) 
	for k,v := range classCounts {
		errorPerClass[k] = ZeroDivStandard(errorPerClass[k],v) 
	}
	return &PredictionMetric{totalError: totalError, errorPerClass: errorPerClass, classCounts: classCounts}
}