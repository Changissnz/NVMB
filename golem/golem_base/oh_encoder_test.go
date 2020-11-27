package golem

// refer to file `oh_encoder` and method `func (c *ConcurrentFileBlockReader) ConvertPartitionBlockToRegularMatrix() {` 

import (
	"fmt"
	"testing"
)

func InputForEncoder() (*Block, map[int]map[string]float64) {

	// elephant, eleph, ant, phant, lephant, ele
	// shoe, sh, sho, she, 
	// we, were, wew, wow, waz
	data := [][]string{{"elephant", "shoe", "we"},
	{"eleph", "sh", "were"},
	{"elephant", "shoe", "we"},
	{"ant", "sho", "wew"},
	{"lephant", "she", "wow"},
	{"ele", "shoe", "wew"},
	{"eleph", "sh", "waz"},
	{"phant", "sho", "we"},
	{"lephant", "she", "were"},
	{"ele", "shoe", "wew"},
	{"eleph", "sh", "wow"},
	{"ant", "sho", "waz"},
	{"lephant", "she", "were"},
	{"phant", "shoe", "wow"}} 
	
	pe := make(map[int]map[string]float64, 0) 
	pe[0] = map[string]float64{"elephant" : float64(0), "eleph": float64(1), "ant": float64(2)}  
	pe[1] = map[string]float64{"shoe": float64(0), "sh": float64(1)} 
	pe[2] = map[string]float64{"we": float64(0), "were": float64(1)}  
	
	b := OneBlock()
	for _, x := range data {
		b.AddOne(x) 
	}
	return b, pe
}

/// TODO make assertions
func Test_OHEncodeBlock_Dumb_UsesPreviousEncoding(t *testing.T) {
	b, m := InputForEncoder()

	m2 := OHEncodeBlock_Dumb_UsesPreviousEncoding(b,m)

	// display encoding map m afterwards
	fmt.Println("MAP AFTERWARDS") 
	fmt.Println(m)
	fmt.Println("---------------------------------")

	// iterate through matrix and display encoding
	fmt.Println("MATRIX AFTERWARDS")	
	r,_ := m2.Dims() 
	for i := 0; i < r; i++ {
		row,_ := MatrixRowToFloat64Slice(m2,i)
		fmt.Println("row: ", row)
	} 
}