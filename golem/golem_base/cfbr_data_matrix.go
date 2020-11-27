/*
this matrix is to be used alongside concurrent file block reader. 
*/ 
package golem

import (
	"gonum.org/v1/gonum/mat"
	"fmt"
	"sync"
	"math"
)

type CFBRDataMatrix struct {
	floatData *Block
	floatDataColumnKeys []int
	floatDataColumnLabels []string 
	floatDim Pair

	intData           *Block
	intDataColumnKeys []int
	intDataColumnLabels []string
	intDim Pair

	stringData *Block
	stringDataColumnKeys []int 
	stringDataColumnLabels []string
	stringDim Pair

	vectorData *Block
	vectorDataColumnKeys []int 
	vectorDataColumnLabels []string 
	vectorDim Pair

	convertedData map[string]*mat.Dense

	nullKey string

	/// TODO: make sure these are used
	undefinedKeys []int
	undefinedLabels []string 
}

/*
return: 
- (int::index, string::dataType)
*/ 
func (c *CFBRDataMatrix) ColumnLabelToIndexData(columnLabel string) (int, string) {

	floatIndex := StringIndexInSlice(c.floatDataColumnLabels, columnLabel) 
	if (floatIndex != -1) {
		return floatIndex, "float"
	}

	intIndex := StringIndexInSlice(c.intDataColumnLabels, columnLabel) 
	if (intIndex != -1) {
		return intIndex, "int"
	}
	
	stringIndex := StringIndexInSlice(c.stringDataColumnLabels, columnLabel) 
	if (stringIndex != -1) {
		return stringIndex, "string"
	}
	
	vecIndex := StringIndexInSlice(c.vectorDataColumnLabels, columnLabel) 
	if (vecIndex != -1) {
		return vecIndex, "vector"
	}

	panic(fmt.Sprintf("column label %s could not be found", columnLabel))
} 

/*
*/
func (c *CFBRDataMatrix) ValidIndex(index int, rowOrColumn string, dataType string) bool {
	
	if index < 0 {
		return false
	}

	if rowOrColumn != "r" && rowOrColumn != "c" {
		panic("invalid argument!")
	}

	if dataType != "int" && dataType != "float" {
		panic("invalid argument!")
	}
	
	var p Pair 
	if dataType == "int" {
		p = c.intDim 
	} else {
		p = c.floatDim
	}

	if rowOrColumn == "r" {
		if index >= p.a.(int) {
			return false
		}
	} else {
		if index >= p.b.(int) {
			return false
		}
	}

	return true
}

/// CAUTION: method not fully tested 
func (c *CFBRDataMatrix) FetchValueNumerical(r int, col int, dataType string) (float64, bool) {

	if !c.ValidIndex(r, "r", dataType) {
		return float64(-1), false 	
	}

	if !c.ValidIndex(col, "c", dataType) {
		return float64(-1), false 
	}

	return c.convertedData[dataType].At(r,col), true 
}

// TODO: finish | delete 
func (c *CFBRDataMatrix) FetchValueString() {
}

/*
*/
func (c *CFBRDataMatrix) FetchRowNumerical(r int, dataType string) ([]float64, bool) {
	if dataType != "int" && dataType != "float" {
		panic("[cfbr data matrix] invalid data type")
	}
	return MatrixRowToFloat64Slice(c.convertedData[dataType], r)
}

/*
*/
func (c *CFBRDataMatrix) FetchRowRangeNumerical(startRange int, endRange int, dataType string) (*mat.Dense, bool) {

	if (endRange <= startRange) {
		panic("start has to be less than end")
	}

	output := make([]float64,0) 
	for i := startRange; i < endRange; i++ {
		toAdd, stat := c.FetchRowNumerical(i, dataType) 
		if (!stat) {
			return nil, stat
		}
		output = append(output, toAdd...) 
	}

	var colSize int 
	if dataType == "int" {
		colSize = c.intDim.b.(int) 
	} else {
		colSize = c.floatDim.b.(int)
	}

	return mat.NewDense(endRange - startRange, colSize, output), true
}

/*
*/
func (c *CFBRDataMatrix) FetchRowString(r int, dataType string) ([]string,bool) {
	if dataType != "string" && dataType != "vector" {
		panic("invalid data type")
	}

	if dataType == "string" {
		return c.stringData.FetchElementsAtRow(r,nil)
	}

	return c.vectorData.FetchElementsAtRow(r,nil)
}

/*
*/
func (c *CFBRDataMatrix) FetchRowRangeString(startRange int, endRange int, dataType string) ([][]string, bool) {

	if dataType != "string" && dataType != "vector" {
		panic("invalid data type")
	}

	if dataType == "string" {
		return c.stringData.FetchElementsAtRowRange(startRange, endRange, nil)
	}

	return c.vectorData.FetchElementsAtRowRange(startRange, endRange, nil)
}

/*
*/
func (c *CFBRDataMatrix) BlockToMatrix(blockType string) *mat.Dense {
	// set block of target
	if blockType != "int" && blockType != "float" {
		panic(fmt.Sprintf("invalid block type %s", blockType))
	}

	var b *Block
	if blockType == "int" {
		b = c.intData
	} else {
		b = c.floatData
	}

	// construct matrix
	l,w := b.Dims()
	output := mat.NewDense(l, w, nil) 
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for i := 0; i < l; i++ {
			for j := 0; j < w; j++ {
				index := []int{i,j}
				v := b.GetAtOne(index) 
				var f float64
				if v == c.nullKey {
					f = math.NaN()
				} else {
					f = float64(DefaultStringToFloat(v))
				}
				output.Set(i,j,f) 
			}
		}
		wg.Done() 
	}()

	wg.Wait() 

	return output
}

/*
*/
func (c *CFBRDataMatrix) BlocksToMatrices() {
	c.convertedData = make(map[string]*mat.Dense,0)
	c.convertedData["int"] = c.BlockToMatrix("int")
	c.convertedData["float"] = c.BlockToMatrix("float")
	c.intData.Clear() 
	c.floatData.Clear()
}

/*
converts float and int blocks to matrices
*/
func (c *CFBRDataMatrix) Preprocess() {
	c.BlocksToMatrices()
	c.ClearFloatAndInt()
	c.UpdateDims() 
}

func (c *CFBRDataMatrix) Dims() (int,int) {
	r := c.intDim.a.(int) 

	c_ := c.intDim.b.(int) 
	c_ += c.floatDim.b.(int)
	c_ += c.stringDim.b.(int) 
	c_ += c.vectorDim.b.(int) 

	return r, c_
}

/*
sets dimensions for each partition of data
*/
func (c *CFBRDataMatrix) UpdateDims() {
	ax, bx := c.convertedData["int"].Dims()   
	c.intDim = Pair{a: ax, b: bx} 

	ax,bx = c.convertedData["float"].Dims() 
	c.floatDim = Pair{a: ax, b: bx}

	ax, bx = c.stringData.Dims()
	c.stringDim = Pair{a: ax, b: bx}

	ax, bx = c.vectorData.Dims()
	c.vectorDim = Pair{a: ax, b: bx}
}

func (c *CFBRDataMatrix) ClearFloatAndInt() {
	c.floatData.Clear()
	c.intData.Clear() 
}

/*
CAUTION: use only after calling .Preprocess() 
*/ 
func (c *CFBRDataMatrix) IndexRange(start int, end int) *CFBRDataMatrix {

	// collect string and vector data 
	bs, stat := c.stringData.FetchElementsAtRowRange(start,end,nil) 

	if (!stat) {
		panic("could not collect string data") 
	}
	bd := OneBlock() 
	bd.datos = bs

	vs, stat := c.vectorData.FetchElementsAtRowRange(start,end,nil) 
	if (!stat) {
		panic("could not collect vector data") 
	}
	vd := OneBlock() 
	vd.datos = vs

	// collect int and float data 
	mi := MatrixRowRange(c.convertedData["int"], start, end)
	fi := MatrixRowRange(c.convertedData["float"], start, end)

	/// copy data over to new instance 
	c2 := &CFBRDataMatrix{} 
	c2.convertedData = make(map[string]*mat.Dense,0)
	c2.convertedData["float"] = fi 
	c2.floatDataColumnKeys = c.floatDataColumnKeys
	c2.floatDataColumnLabels = c.floatDataColumnLabels

	c2.convertedData["int"] = mi
	c2.intDataColumnKeys = c.intDataColumnKeys
	c2.intDataColumnLabels = c.intDataColumnLabels

	c2.stringData = bd 
	c2.stringDataColumnKeys = c.stringDataColumnKeys
	c2.stringDataColumnLabels = c.stringDataColumnLabels

	c2.vectorData = vd 
	c2.vectorDataColumnKeys = c.vectorDataColumnKeys
	c2.vectorDataColumnLabels = c.vectorDataColumnLabels

	c2.nullKey = c.nullKey 

	c2.undefinedKeys = c.undefinedKeys
	c2.undefinedLabels = c.undefinedLabels

	c2.UpdateDims()
	return c2
}

/// TODO CAUTION: has not been tested
func (c *CFBRDataMatrix) CollectByIndices(indices []int) *CFBRDataMatrix {
	c2 := &CFBRDataMatrix{} 

	// string data 
	sd := c.stringData.GetAtByIndices(indices)
	c2.stringData = sd 
	c2.stringDataColumnKeys = c.stringDataColumnKeys
	c2.stringDataColumnLabels = c.stringDataColumnLabels
	

	// vector data 
	vd := c.vectorData.GetAtByIndices(indices)
	c2.vectorData = vd 
	c2.vectorDataColumnKeys = c.vectorDataColumnKeys
	c2.vectorDataColumnLabels = c.vectorDataColumnLabels

	// float and int data 
	c2.convertedData = make(map[string]*mat.Dense,0)

	c2.convertedData["float"] = MatrixCollectByIndices(c.convertedData["float"], indices) 
	c2.floatDataColumnKeys = c.floatDataColumnKeys
	c2.floatDataColumnLabels = c.floatDataColumnLabels

	c2.convertedData["int"] = MatrixCollectByIndices(c.convertedData["int"], indices) 
	c2.intDataColumnKeys = c.intDataColumnKeys
	c2.intDataColumnLabels = c.intDataColumnLabels

	return c2
} 



/*
stacks instance on top of c2. 

CAUTION: assumes both instance and c2 have been preprocessed. 
CAUTION: assumes column keys, column labels, and dims are the same. 
*/ 
func (c *CFBRDataMatrix) StackMatrix(c2 *CFBRDataMatrix) {
	c.stringData.StackBlock(c2.stringData) 
	c.vectorData.StackBlock(c2.vectorData) 
	c.convertedData["int"] = StackTwoMatrices(c.convertedData["int"], c2.convertedData["int"])
	c.convertedData["float"] = StackTwoMatrices(c.convertedData["float"], c2.convertedData["float"])
	c.UpdateDims() 
}

/*
*/
func (c *CFBRDataMatrix) Length() int {
	return c.intDim.a.(int)
}

/*
*/
func (c *CFBRDataMatrix) Display() {
	fmt.Println("*\tINT DATA")
	DisplayMatrix(c.convertedData["int"]) 
	fmt.Println()

	fmt.Println("*\tFLOAT DATA")
	DisplayMatrix(c.convertedData["float"]) 
	fmt.Println()

	fmt.Println("*\tSTRING DATA")
	DisplayStringMatrix(c.stringData.datos) 
	fmt.Println()

	fmt.Println("*\tVECTOR DATA")
	DisplayStringMatrix(c.vectorData.datos) 
	fmt.Println() 
}