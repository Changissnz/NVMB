/*
file contains observation measures 
*/ 
package golem

import (
	"fmt"
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

/*
divides delta by reference, undef(Pos|Neg) are floats to output
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

/*
divides delta by reference, outputs 0 for zero-division
*/ 
func ZeroDivStandard(delta float64, reference float64) float64 {
	if reference == float64(0) {
		return reference 
	}

	return delta / reference
}

// TODO: error-check
func DeltaOpOnSequence(f []float64, opType string) float64 {

	if len(f) < 2 {
		panic("slice must be at least size 2")
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

/// TODO: code "ends"
/*
*/ 
func DeltaOpOnDiscreteSequence(s []string) string {
	output := make(map[string]float64,0) 

	for _,s_ := range s {
		output[s_]++ 
	}

	_,max := MinAndMaxKeyByValueOfMapStringToFloat64(output)
	return max
}

/*
given a sequence of at least size 2,

return: 
- |set_difference(s[0], s[-1])| : sub
- |set_difference(s[-1], s[0])| : add
*/ 
func DeltaOpOnSetSequence(s []string) (int,int) {
	if len(s) < 2 {
		panic("cannot perform delta op. on set seq. less than length 2")
	}

	s0,_ := DefaultStringToStringSlice(s[0], DEFAULT_DELIMITER) 
	s1,_ := DefaultStringToStringSlice(s[len(s) - 1], DEFAULT_DELIMITER) 

	return len(StringSliceDifference(s0,s1)), len(StringSliceDifference(s1,s0)) 
}