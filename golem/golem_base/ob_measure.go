/*
file contains observation measures 
*/ 
package golem

import (
	"fmt"
	//"math" 
	//"gonum.org/v1/gonum/mat"
)

//////////// revisions to delta functions 
func Delta_Ends(f []float64) float64 {
	x := f[len(f) -1] - f[0] 
	return ZeroDiv(x, f[0],float64(1), float64(-1))
}

/*
*/
func Delta_Mean2End(f []float64) float64 {
	mean := MeanIterable(f) 
	x := f[len(f) - 1] - mean
	return ZeroDiv(x, mean,float64(1), float64(-1))
}

func Delta_Mean2Start(f []float64) float64 {
	mean := MeanIterable(f) 
	x := mean - f[0]
	return ZeroDiv(x, mean, float64(1), float64(-1))
}

/*
0-div -> 100% delta
*/ 
func Delta_Linear(f []float64) float64 {
	j := len(f) - 1

	deltaSum := float64(0); 

	for i := 0; i < j; i++ {
		deltaSum += ZeroDiv(f[i + 1] - f[i], f[i], float64(-1), float64(1))
	}

	return deltaSum / float64(j)
}


// TODO: rename this.
/*
*/
func ZeroDiv(delta float64, reference float64, undefPos float64, undefNeg float64) float64 {

	if reference == float64(0) {

		if delta < reference {
			return undefNeg
		}

		return undefPos
	}
	return delta / reference
}

// TODO: error-check
func DeltaOpOnSequence(f []float64, opType string) float64 {

	if len(f) < 2 {
		panic("slice must be at least size 2")
	}

	if opType == "ends" {
		return Delta_Ends(f)
	}

	switch opType{
		case "ends": 
			return Delta_Ends(f) 
		case "mean-to-end":
			return Delta_Mean2End(f) 
		case "mean-to-start": 
			return Delta_Mean2Start(f) 
		case "linear": 
			return Delta_Linear(f)
		default: 
			panic(fmt.Sprintf("invalid op type %s", opType))
	}
}