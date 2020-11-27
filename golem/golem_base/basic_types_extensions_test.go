package golem

import (
	"fmt"
	"testing" 
)

func Test_DefaultStringToFloat64Slice(t *testing.T) {

	ts1 := " 1.000000 0.000000 0.000000 5.000000 0.000000 0.000000  "
	floatSlice, stat := DefaultStringToFloat64Slice(ts1, ' ')

	if (!stat) {
		panic("incorrect slice conversion")
	}

	floatAns := []float64{1,0,0,5,0,0}
	var res bool = EqualFloat64Slices(floatSlice, floatAns, 1) 

	if (res != true) {
		panic("String to float64 slice conversion failure") 
	}
} 

func Test_SubstringInStringAll(t *testing.T) {
	test1 := "CHECKEN IN THE BOXBOX CHECKEN BOX" 
	test1s := "BOX" 
	x := SubstringInStringAll(test1,test1s)
	l := len(test1s)

	for _, x_ := range x {
		if test1[x_: x_ + l] != "BOX" {
			panic("incorrect substring")
		} 
	}

	if len(x) != 3 {
		panic("incorrect number of substrings")
	}
}

func Test_SortStringSliceAscending(t *testing.T) {

	target := []string{"alabama", "ARM", "arm", "um", "beasten", "9-oh"} 
	answer := []string{"9-oh", "ARM", "alabama", "arm", "beasten", "um"}

	result := SortStringSliceAscending(target) 

	for i, v := range answer {
		if v != result[i] {
			panic("error: wrong sort") 
		}
	}	
}

/*
*/
func Test_BayesianIdentityOfFloatSliceSubset(t *testing.T) {
	floatRange := &Pair{a: 3.0, b: 105.0} 
	usedIndexRanges := []*Pair{&Pair{a: 3.0, b: 5.0},
				&Pair{a: 8.0, b: 10.0}, &Pair{a: 12.0, b: 12.5},
				&Pair{a: 18.0, b: 22.1}, &Pair{a: 29.0, b: 45.0},
				&Pair{a: 103.0, b: 105.0}}

	x := BayesianIdentityOfFloatSliceSubset(floatRange, usedIndexRanges)
	expected := []*Pair{&Pair{a: 5.0, b: 8.0},
	&Pair{a: 10.0, b: 12.0}, &Pair{a: 12.5, b: 18.0},
	&Pair{a: 22.1, b: 29.0}, &Pair{a: 45.0, b: 103.0}}

	for i,x_ := range x {
		if x_.a.(float64) != expected[i].a.(float64) {
			panic("incorrect bayesian identity for float slice subset!")
		}
	}
}

func Test_SampleSubrangesForRange(t *testing.T) {
	p := &Pair{a: 0.0, b: 100.0}
	output := SampleSubrangesForRange(p, nil, 30, []*Pair{}, 3, true)
	fmt.Println("OUTPUT") 
	DisplayPairSlice(output)
}

func Test_InterfaceSliceToIntSlice(t *testing.T) {
	px := []interface{}{1.0,2.0,3.0} 
	py := InterfaceSliceToIntSlice(px) 
	fmt.Println("INTERFACE SLICE TO INT SLICE ", py)
}