package golem

import (
	"fmt" 
	"testing" 
)

func GenerateFileDefault() {
	fp := "./datos/default" 
	sz := 10050
	nc := 20 
	rfr := Pair{float32(3.4), float32(10.5)} 
	bv := make([]int, 0) 
	GenerateRandomCSVDataToFile(fp, sz, nc, rfr, bv)	
}

func GenerateFileDefault2() {
	fp := "./datos/default2" 
	p := Pair{a: float32(3.5), b: float32(12.5)} 
	bv := []int{0,1,3}
	sv := []int{1,4,6,7} 
	GenerateRandomCSVDataToFileExtended(fp, 100, 10, p, bv, sv, nil, true) 
}

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