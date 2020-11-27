/*
this is where supplementary code for gonum.Matrix data structure goes. 
*/ 

package golem 

import (
	"sync"
	"fmt"
	"gonum.org/v1/gonum/mat"
)

func TransposeMatrix(m *mat.Dense) *mat.Dense {

	l, c := m.Dims() 
	m2 := mat.NewDense(c, l, nil) 
	
	for i := 0; i < c; i++ {
		newRow,_ := MatrixColToFloat64Slice(m, i) 
		m2.SetRow(i, newRow)  
	}

	return m2
}

/*
*/
func MaxMatrixAtColumnOrRow(m *mat.Dense, rc string, index int) float64 {
	
	r,c := m.Dims()
	if r == 0 || c == 0 {
		panic("number of rows cannot be 0") 
	}
	
	var t int
	if rc == "row" {
		t = r
	} else if rc == "column" {
		t = c
	} else {
		panic("arg. must be row or col.")
	}

	x := m.At(0, index) 
	for i := 0; i < t; i++ {
		var x2 float64
		switch {
			case rc == "row": 
				x2 = m.At(index, i)
			default: 
				x2 = m.At(i, index)  	
		}

		if x2 > x {
			x = x2
		}
	}

	return x
}

/*
*/
func MaxMatrixEachColumnOrRow(m *mat.Dense, rc string) []float64 {

	_, c := m.Dims()
	x := make([]float64, c)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < c; i++ {
			maxx := MaxMatrixAtColumnOrRow(m, rc, i) 
			x[i] = maxx
		}
		wg.Done() 
	}() 
	wg.Wait()

	return x
}

/*
*/
func MeanMatrixAtColumnOrRow(m *mat.Dense, rc string, index int) float64 {
	s := float64(0)
	r,c := m.Dims()
	if r == 0 || c == 0 {
		panic("number of rows cannot be 0") 
	}

	var t int

	if rc == "row" {
		t = r
	} else if rc == "column" {
		t = c
	} else {
		panic("arg. must be row or col.")
	}

	for i := 0; i < t; i++ {

		switch {

		case rc == "row": 
			s += m.At(index, i) 
	
		default: 
			s += m.At(i, index)  	
		}
	}

	return s / float64(t) 
}

/*
*/
func MeanMatrixEachColumnOrRow(m *mat.Dense, rc string) []float64 {

	_, c := m.Dims()
	x := make([]float64, c) 

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < c; i++ {
			mean := MeanMatrixAtColumnOrRow(m, rc, i) 
			x[i] = mean
		}
		wg.Done() 
	}() 
	wg.Wait()

	return x
}

///////// END re-factor below

/*
*/ 
func MatrixRowToFloat64Slice(m *mat.Dense, row int) ([]float64, bool) {
	l, c := m.Dims()
	
	if row >= l {
		return nil, false 
	}
	
	r := make([]float64,0)
	for i := 0; i < c; i++ {
		r = append(r, m.At(row, i)) 
	}
	return r, true 
}

/*
*/
func MatrixColToFloat64Slice(m *mat.Dense, col int) ([]float64, bool) {
	l, c := m.Dims()
	
	if col >= c {
		return nil, false 
	}
	
	r := make([]float64,0)
	for i := 0; i < l; i++ {
		r = append(r, m.At(i, col)) 
	}
	return r, true 
}

/*
*/
func MatrixRowRange(m *mat.Dense, rowStart int, rowEnd int) *mat.Dense {
	
	if rowStart >= rowEnd {
		panic("start must be less than end")
	}
	
	if rowStart < 0 {
		panic("start cannot be negative")
	}

	output := make([]float64,0)
	_, c := m.Dims()
	for i := rowStart; i < rowEnd; i++ {
		for j := 0; j < c; j++ {
			output = append(output, m.At(i, j))
		}
	}
	return mat.NewDense(rowEnd - rowStart, c, output) 
}

/*
CAUTION: index values may be unordered 
*/ 
func MatrixCollectByIndices(m *mat.Dense, indices []int) *mat.Dense {

	if len(indices) == 0 {
		return nil
	}

	output := make([]float64,0)
	for _, i := range indices {
		r, f := MatrixRowToFloat64Slice(m, i)
		
		if !f {
			panic(fmt.Sprintf("invalid index %d", i)) 
		}
		
		output = append(output, r...)
	}

	_, c := m.Dims()
	x := mat.NewDense(len(indices), c, output)  
	return x
}

/*
*/ 
func MatrixTo1DFloat(m *mat.Dense) []float64 {

	r,_ := m.Dims() 
	output := make([]float64,0) 

	for i := 0; i < r; i++ {
		q, stat := MatrixRowToFloat64Slice(m, i) //m *mat.Dense, row int) ([]float64, bool) {
		
		if stat == false {
			break 
		} 

		output = append(output, q...) 
	}

	return output
} 

/*
stacks two matrices on top of each other 
*/ 
func StackTwoMatrices(m1, m2 *mat.Dense) *mat.Dense {
	r, c := m1.Dims() 
	r2, c2 := m2.Dims() 
	if (c != c2) {
		panic("unequal num. of cols.") 
	}
	output := mat.NewDense(r + r2, c, nil) 
	output.Stack(m1,m2) 
	return output 
}

/*
stacks slice of matrices on top of each other in sequential order 
*/ 
func StackMatrices(m []*mat.Dense) *mat.Dense {
	l := len(m) - 1 

	if (l == -1) {
		return nil 
	}

	if (l == 0) {
		return m[0] 
	}

	var output *mat.Dense
	for i := 0; i < l; i++ {
		output = StackTwoMatrices(m[i], m[i +1]) 
	}
	return output 
} 

// TODO: consider concurrency 
/*
given a slice of matrices b, each element of equal size, 
outputs a matrix in which each row is element @ index oi 
of elements in b 
*/ 
func BooleanMatrixSliceToMatrix(b []*mat.Dense, oi int) *mat.Dense {
	if len(b) == 0 {
		return nil 
	}

	_, c := b[0].Dims() 
	output := mat.NewDense(len(b),c,nil)   

	for i, b_ := range b {
		m, stat := MatrixRowToFloat64Slice(b_, oi)
		if stat {
			output.SetRow(i, m) 
		}
	}
	
	return output 
}



/*
given s1 and s2 of equal length, outputs row-wise mult. op. of equal length 
*/ 
func MultiplyTwoSlices(s1 []float64, s2 []float64) []float64 {
	
	if (len(s1) != len(s2) || len(s1) == 0) {
		panic("invalid slice input for mult.")
	}
	
	q := make([]float64,0) 
	for i, x := range s1 {
		q = append(q, x * s2[i])
	}
	return q 
}

/*
*/
func DisplayMatrix(m *mat.Dense) {

	r,_ := m.Dims() 

	for i := 0; i < r; i++ {
		row,_ := MatrixRowToFloat64Slice(m,i)
		fmt.Println(row)  
	}
} 

/*
*/
func DisplayStringMatrix(s [][]string) {

	for _,s_ := range s {
		fmt.Println(s_) 
	}

}

/*
*/
func MatricesAreEqual(m1 *mat.Dense, m2 *mat.Dense) bool {
	r1,c1 := m1.Dims()
	r2,c2 := m2.Dims()

	if r1 != r2 {
		return false 
	}

	if c1 != c2 {
		return false
	}

	for i := 0; i < r1; i++ {
		row1,_ := MatrixRowToFloat64Slice(m1, i)
		row2,_ := MatrixRowToFloat64Slice(m2, i) 

		if !EqualFloat64Slices(row1, row2, 1) {
			return false 
		}
	}

	return true 
}