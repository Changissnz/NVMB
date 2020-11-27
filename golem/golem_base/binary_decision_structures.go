package golem

import (
	"fmt"
)

type Split struct {
	splitType string // discrete, continuous
	splitDC []interface{}
}

func OneSplit(splitIntervals []interface{}, colType string) *Split {
	return &Split{splitType: colType, splitDC: splitIntervals} 
}

func (s *Split) Display() {
	switch {
	case s.splitType == "continuous": 
		p := InterfaceSliceToPairSlice(s.splitDC)
		DisplayPairSlice(p) 
	default: 
		fmt.Println(s.splitDC)		

	}
}

func InterfaceSliceToFloatSlice(p []interface{}) []float64 {
	output := make([]float64,0) 

	for _,p_ := range p {
		output = append(output, p_.(float64)) 
	}
	return output
}

type SplitData struct {
	node *BDNode
	split *Split
	leftScore float64 
	rightScore float64
	leftMismatch float64
	rightMismatch float64
	inheritance string
	column int
	active bool
}

/*
*/
func SplitPairDataToSplitData(p *Pair, n *BDNode, inheritance string, ml,mr float64) *SplitData {
	var s *SplitData 

	if p != nil {
		s = &SplitData{node: n, split: p.a.(*Pair).a.(*Split), leftScore: p.b.(*Pair).a.(float64), rightScore: p.b.(*Pair).b.(float64), inheritance: inheritance, column: p.a.(*Pair).b.(int), leftMismatch: ml, rightMismatch: mr, active: true}
	}

	return s
}


/*
about struct BDNode 

- decision-making 
	left is no
	right is yes
- each Node stores the indices of the sampling data it considers
WARNING: may lead to memory problems
*/ 
type BDNode struct {
	key int // column that node operates on
	value int

	left *BDNode 
	right *BDNode 
	parent *BDNode
	inheritance string
	split *Split
	theFunction func(float64) bool

	leftIndices []int 
	rightIndices []int
	availableColumnIndices []int

	/// use one of the below variables to get error of split
	errorAt float64
	dataSize float64 /// int 
	dataSizeChange float64 /// int

	leftErr float64
	rightErr float64
	leftMismatch float64 
	rightMismatch float64

	depth int
	finished bool /// TODO: delete this 
	leftFin bool
	rightFin bool 
}

/*
*/
func OneBDNode(k int, sd *SplitData, r float64, depth int) *BDNode {///, scoreLeft float64, scoreRight float64) *BDNode {
	var tf func(x float64) bool

	switch {

		case sd.split.splitType == "discrete":
			tf = func (x float64) bool {
				for _,q := range sd.split.splitDC {
					if x == q.(float64) {
						return true 
					}
				}
				return false 
			}
		
		//// TODO: group of x needs to be known 
		case sd.split.splitType == "continuous":
			tf = func (x float64) bool {

				for _,q := range sd.split.splitDC {
					if x >= q.(*Pair).a.(float64) && x <= q.(*Pair).b.(float64) {
						return true 
					}
				}
				return false 
			}

		default: 
			panic("invalid split type")

	}

	n := &BDNode{key: k, value: sd.column, split: sd.split, theFunction: tf, dataSize: r, dataSizeChange: -1, depth: depth, finished: false, leftFin: false, rightFin: false}  
	n.leftErr, n.rightErr = sd.leftScore, sd.rightScore 
	n.leftMismatch, n.rightMismatch = sd.leftMismatch, sd.rightMismatch
	return n
}

/// TODO: unused
func (n *BDNode) IsRelationNull(pos string) bool { 

	switch {

		case pos == "l": 
			if (n.left == nil) {
				return true
			}
			return false 

		case pos == "r": 
			if (n.right == nil) {
				return true
			}
			return false 

		case pos == "p": 
			if (n.parent == nil) {
				return true
			}
			return false 

		default: 
			panic(fmt.Sprintf("invalid position %s", pos))
	}
}

/*
outputs qualifying indices given column info 
*/
func (n *BDNode) ProcessData(column []float64) []int {
	qualIndices := make([]int,0) 
	for i, c := range column {
		if n.theFunction(c) {
			qualIndices = append(qualIndices, i)
		}
	}
	return qualIndices
}