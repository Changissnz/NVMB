package golem

import (
	"fmt"
	"testing"
	"gonum.org/v1/gonum/mat"
	//"time"
	//"reflect"
	//"strconv"
	//"encoding/csv"
	//"io"
)

func Test_MatrixRowRange(t *testing.T) {

	var f = []float64{0,2,3,4,5,15,16,17,18,19,40,41,51,67,80} 
	m := mat.NewDense(3,5,f) 

	var f2 = []float64{0,2,3,4,5,15,16,17,18,19}
	answer := mat.NewDense(2,5,f2) 
	m2 := MatrixRowRange(m, 0, 2)	
	if !mat.Equal(m2,answer) {
		panic("incorrect matrix row range")
	}
}

func Test_MatrixStack(t *testing.T) {

	var f = []float64{0,2,3,4,5,15,16,17,18,19,40,41,51,67,80} 
	m := mat.NewDense(3,5,f) 

	var f2 = []float64{20,32,43,64,75,115,216,917,418,919}
	m2 := mat.NewDense(2,5,f2) 

	var f3 = []float64{0,2,3,4,5,15,16,17,18,19,40,41,51,67,80,20,32,43,64,75,115,216,917,418,919} 
	answer := mat.NewDense(5,5,f3) 

	m3 := StackTwoMatrices(m, m2)

	r,_  := m3.Dims() 

	for i := 0; i < r; i++ {
		r1,_ := MatrixRowToFloat64Slice(answer,i)  
		r2,_ := MatrixRowToFloat64Slice(m3,i) 

		if !EqualFloat64Slices(r1,r2,2) {
			panic(fmt.Sprintf("incorrect stack at index %d", i)) 
		}
	}
}