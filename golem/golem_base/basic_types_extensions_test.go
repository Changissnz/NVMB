package golem

import (
	//"fmt"
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