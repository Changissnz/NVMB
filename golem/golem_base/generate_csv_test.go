package golem

import (
	"fmt" 
	"testing" 
)


/*
*/ 
func Test_GenerateRandomRowDataExtendedFlawInbut(t *testing.T) {

	nc := 20 
	rfg := Pair{a: float32(5.0), b: float32(30.0)}
	bv := []int{0,1,4} 
	sv := []int{6,8,9}
	svvm := GenerateRandomColumnClasses(sv, 20)
	fpc := map[int]float32{1: float32(2)}

	out := GenerateRandomRowDataExtendedFlawInbut(nc, rfg, bv,
		sv, svvm, fpc)

	fmt.Println("**TestGenerateRandomRowDataExtendedFlawInbut**")
	fmt.Println()
	fmt.Println(out)
	fmt.Println()
} 