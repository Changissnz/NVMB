/*
This file is to process a file

NOTE: csv must have header column
*/

package golem

import (
	"strconv"
	"encoding/csv"
	"errors"
	"os"
	"github.com/stoicperlman/fls"
	"fmt"
	"io"
	"sync"
	"sort"
	"gonum.org/v1/gonum/mat"
	"bufio"
	"bytes"
)

var WRITESIZE int64 = 10000 // numRows * numCols
var DATA_LENGTH_READ_CAP int64 = 2000000
var DEFAULT_TIME_COLUMN_LABEL string = "time"

/*
calculates sampling size given the `size`  
*/
func GetSampleSizeDefault(size int64) int {
	if size < 1000 {
		return int(size)
	}

	if float64(size)*0.1 < 10000 {
		return 10000
	} else {
		return int(float64(size) * 0.1)
	}
}

/*
this is a block reader that reads data in increments (blocks), and
outputs analyses of the blocks by timestamp-based analysis
*/
type ConcurrentFileBlockReader struct {
	fp string
	ofp string

	marker int
	lineMarker int64 

	blockData         []*Block
	blockDataSize     int64
	blockMapCollector *MapCollector

	fileLength int64
	partitionLength int64

	analyses []string
	columns  []string
	excludedColumnIndices []int ///@ 

	// key: column index
	// value: string|int|int32|int64|float32|float64
	columnTypes map[int]string
	stringColumnEncodings map[int]map[string]float64
	oi *BasicSet
	flsFileObj *fls.File
}

/*
makes a ConcurrentFileBlockReader 
*/
func OneConcurrentFileBlockReader(fp string) *ConcurrentFileBlockReader {

	c := ConcurrentFileBlockReader{fp: fp, ofp: "nonatos", marker: 0,
		blockData:   make([]*Block, 0),
		analyses:    make([]string, 0),
		//columns:     make([]string, 0),
		columnTypes: make(map[int]string, 0),
	}

	if fp != "" {
		c.ReadHeader()
	}

	return &c
}

/*
*/
func (c *ConcurrentFileBlockReader) ReadHeader() {
	
	err := errors.New("foo")
	c.flsFileObj, err = fls.OpenFile(c.fp, os.O_RDONLY, 0400)
	if err != nil {
		panic("invalid output file path")
	}

	reader := csv.NewReader(c.flsFileObj)
	c.columns, err = reader.Read()

	if err != nil {
		panic("cannot read file column header")
	}

	c.SetReaderToBeginning()
}

/*
transfers chunk of file data to block
*/
func (c *ConcurrentFileBlockReader) ReadBlockAtLine(bl int) (*Block, int64, bool) {

	_, err := c.flsFileObj.SeekLine(int64(bl), io.SeekStart)
	if err != nil {
		panic(fmt.Sprintf("cannot read line %d", bl))
	}

	return c.ReadBlockAtSpot(CPARTSIZE, true) 
}

func (c *ConcurrentFileBlockReader) SetReaderToBeginning() { 
	c.flsFileObj.SeekLine(int64(1), io.SeekStart)
	c.lineMarker = 1
}

/*
scans through file and determines its length
*/ 
func (c *ConcurrentFileBlockReader) PrelimRead() {
	c.fileLength = 0
	c.partitionLength = 0
	
    buf := make([]byte, bufio.MaxScanTokenSize)
    for {
        bufferSize, err := c.flsFileObj.Read(buf)
        if err != nil && err != io.EOF {
			panic("invalid file, could not read line")
        }

        var buffPosition int
        for {
            i := bytes.IndexByte(buf[buffPosition:], '\n')
            if i == -1 || bufferSize == buffPosition {
                break
            }
            buffPosition += i + 1
			c.fileLength++
		}
		
		/// case: file exceeds DATA_LENGTH_READ_CAP, stop here.
		if c.fileLength >= DATA_LENGTH_READ_CAP {
			break
		}

        if err == io.EOF {
            break
        }
    }

	c.partitionLength = int64(int(c.fileLength) / CPARTSIZE)
	c.SetReaderToBeginning()
}

/*
reads a block of length `bs` at current file location

return: 
- *Block: the block object starting at current file location 
- size of block 
- finished reading
*/
func (c *ConcurrentFileBlockReader) ReadBlockAtSpot(bs int, collectData bool) (*Block, int64, bool) {
	
	c.flsFileObj.SeekLine(c.lineMarker, io.SeekStart)	
	reader := csv.NewReader(c.flsFileObj)
	b := OneBlock()

	var blockSize int = 0
	for i := 0; i < bs; i++ {

		switch {
		/// TODO: refactor below
		case collectData == true: 
			out, err := reader.Read()
			if err != nil {
				c.lineMarker += int64(i + 1) 
				return b, int64(i), true
			}
			///fmt.Println("LEN OUT [0]: ", len(out))
			if len(c.excludedColumnIndices) != 0 {
				///fmt.Println("Y")
				out = StringSliceCollectByIndices(out, c.excludedColumnIndices, true)
			}
			///fmt.Println("LEN OUT ", len(out))
			b.AddOne(out)
			blockSize += len(c.columns)
		default: 
			_, err := reader.Read()
			if err != nil {
				c.lineMarker += int64(i + 1) 
				return b, int64(i), true
			}
			blockSize += len(c.columns)

		}
	}

	c.lineMarker += int64(bs)
	return b, int64(blockSize), false
}



/*
*/
func (c *ConcurrentFileBlockReader) ReadBlockAtSpotFullTimestamp(bs int) (*Block, int64, bool) {
	block, bs_, stat := c.ReadBlockAtSpot(bs, true) 

	// case: finished 
	if stat {
		return block,bs_,stat 
	}

	// case: empty block 
	if bs_ == 0 {
		return block,bs_,stat
	}

	// get last timestamp
	index := StringIndexInSlice(c.columns, DEFAULT_TIME_COLUMN_LABEL)
	r,_ := block.Dims() 
	ts,_ := strconv.Atoi(block.GetAtOne([]int{r - 1, index}))
	c.flsFileObj.SeekLine(c.lineMarker, io.SeekStart)
	reader := csv.NewReader(c.flsFileObj)
	i := 0
	for {

		// case: end of file 
		out, err := reader.Read()
		if err != nil {
			c.lineMarker += int64(i + 1) 
			return block, int64(i), true
		}

		// case: timestamp does not match
			// parse timestamp
		ts2,_ := strconv.Atoi(out[index]) 

		if ts != ts2 {
			break
		}

		block.AddOne(out) 
		bs_ += int64(len(c.columns))
		i++ 
	}

	c.lineMarker += int64(i)  
	return block,bs_,false
}

/*
reads a partition's worth of data from `flsFileObj` starting at its
current pointer

arguments: 
- readType: full|exact

return: 
- size of partition 
*/
func (c *ConcurrentFileBlockReader) ReadPartition(readType string) int64 {
	c.blockData = make([]*Block, 0)

	var x int64 = 0
	var b *Block 
	var sz int64
	var stat bool

	for {
		///fmt.Println("read one block")
		switch {
		case readType == "full":
			b, sz, stat = c.ReadBlockAtSpotFullTimestamp(CPARTSIZE)
		case readType == "exact": 
			b,sz,stat = c.ReadBlockAtSpot(CPARTSIZE, true) 
		default: 
			panic("invalid read type")
		}

		x += sz
		if sz > 0 {
			c.blockData = append(c.blockData, b)
		}
		if stat {
			break
		}

		if x >= WRITESIZE {
			break
		}
	}
	c.blockDataSize = x
	return x
}

//////// Start: Column-type deduction /////////////////////

/*
deduces column types based on sampling data
*/
func (c *ConcurrentFileBlockReader) DeduceColumnTypes() {
	if len(c.blockData) == 0 {
		c.ReadPartition("exact")
	}

	c.blockMapCollector = OneMapCollector()
	c.ChooseRandomRowDataIndicesDefault()

	// convert c.oi to slice[][]int
	coi := c.oi.ToSliceIndexFormat()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for _, r_ := range coi {
			c.Read2DeduceOneRowDataIndex(r_, true)
		}
		wg.Done()
	}()

	wg.Wait()

	c.blockMapCollector.Predict("majority")
	c.columnTypes = c.blockMapCollector.predictedTypes
}

func (c *ConcurrentFileBlockReader) AllColumnTypesKnown() bool {

	for _, v := range c.columnTypes {
		if v == "?" {
			return false 
		} 
	}
	return true 
}

func (c *ConcurrentFileBlockReader) GetColumnTypes() *BasicSet {
	bs := OneBasicSet() 
	for _, v := range c.columnTypes {
		bs.AddOne(v) 
	}
	return bs
}

/*
outputs a map with key denoting data type, and value denoting the slice of indices belonging to that data type
*/ 
func (c *ConcurrentFileBlockReader) ColumnTypeMapToTypeNotation() map[string][]int{

	TYPE_KEYS := []string{"string", "int", "float", "vector", "undef"} 
	q := make(map[string][]int,0)

	for k, v := range c.columnTypes {
		if v == TYPE_KEYS[0] {
			q[TYPE_KEYS[0]] = append(q[TYPE_KEYS[0]], k)
		} else if v == TYPE_KEYS[1] {
			q[TYPE_KEYS[1]] = append(q[TYPE_KEYS[1]], k)
		} else if v == TYPE_KEYS[2] {
			q[TYPE_KEYS[2]] = append(q[TYPE_KEYS[2]], k)
		} else {
			q[TYPE_KEYS[3]] = append(q[TYPE_KEYS[3]], k)
		}
	}

	for _, k := range TYPE_KEYS { 
		v := q[k] 
		sort.Ints(v) 
		q[k] = v 
	}

	return q 
}

/*
helper method for below. 
Index is 3-d index of the form (partition, row, column)
*/
func (c *ConcurrentFileBlockReader) Read2DeduceOneRowDataIndex_(index []int) string {
	s := c.GetBlockDataAt3dIndex(index)
	return DeduceBasicStringType(s)
}

/// WARNING: not thoroughly checked yet. 
/*
Performs deduction on all elements of index's row. 

arguments: 
- index - 2-d or 3-d
*/
func (c *ConcurrentFileBlockReader) Read2DeduceOneRowDataIndex(index []int, fullRowRead bool) {

	m := make(map[int]string, 0)

	if fullRowRead == true {

		for i := 0; i < len(c.columns); i++ {
			q := []int{index[0], index[1], i}
			m[i] = c.Read2DeduceOneRowDataIndex_(q)
		}
	} else {
		m[index[1]] = c.Read2DeduceOneRowDataIndex_(index)
	}

	c.blockMapCollector.AddOne(m, true, true)
}

/*
chooses a random 3-d index, represented as a size-3 int slice
*/
func (c *ConcurrentFileBlockReader) ChooseRandomIndex() []int {

	if len(c.blockData) == 0 {
		panic("cannot operate on empty block data")
	}

	f := RandomIntInRange(0, len(c.blockData))
	f2 := RandomIntInRange(0, c.blockData[f].Length())
	f3 := RandomIntInRange(0, c.blockData[f].Width())

	return []int{f, f2, f3}
}

/*
outputs single string value from 3-index
*/
func (c *ConcurrentFileBlockReader) GetBlockDataAt3dIndex(x []int) string {

	if len(x) != 3 {
		panic("invalid index size ,!")
	}

	return c.blockData[x[0]].GetAtOne(x[1:])
}

/*
chooses n random 3-indices in partition

return:
- number of iterations used for output

NOTE: efficiency needs work
*/
func (c *ConcurrentFileBlockReader) ChooseRandomRowDataIndices(n int) int {
	c.oi = OneBasicSet()
	var x []int
	var lastChangedTerminate int
	if len(c.blockData)*100 >= 10000000000 {
		lastChangedTerminate = int(10000000000)
	} else {
		// partition
		lastChangedTerminate = len(c.blockData) * 100
	}

	var lastChanged int = 0
	var prevLen int
	c_ := 0
	for {

		// got req
		if c.oi.Len() == n {
			break
		}

		// stall, done
		if lastChanged >= lastChangedTerminate {
			break
		}

		if prevLen != -1 {
			// compare
			if c.oi.Len() == prevLen {
				lastChanged++
			} else {
				lastChanged = 0
			}
			prevLen = c.oi.Len()

		} else {
			prevLen = c.oi.Len()

		}

		x = c.ChooseRandomIndex()
		c.oi.AddOne(DefaultIntSliceToString(x, DEFAULT_DELIMITER))
		c_ += 1
	}

	return c_
}

/*
chooses n random row data indices by default sampling size calculation
*/
func (c *ConcurrentFileBlockReader) ChooseRandomRowDataIndicesDefault() int {
	q := int(GetSampleSizeDefault(c.blockDataSize))
	return c.ChooseRandomRowDataIndices(q)
}
 
func (c *ConcurrentFileBlockReader) ManualColumnTypeSet(ct map[string]string) {
	for i, col := range c.columns {
		t, ok := ct[col] 
		if (!ok) {
			panic(fmt.Sprintf("missing key in typemap: %s", col))
		}
		c.columnTypes[i] = t
	}
}

func (c *ConcurrentFileBlockReader) IndicesToColumnLabels(indices []int) []string {
	output := make([]string,0)
	for _, v := range indices {
		output = append(output, c.columns[v]) 
	}
	return output
}

func (c *ConcurrentFileBlockReader) ConvertPartitionBlockToMatrix(blockIndex int) *CFBRDataMatrix {

	// null case
	if len(c.blockData) <= blockIndex {
		fmt.Println("YA")
		return nil
	}

	// TODO save below as struct var
	q := c.ColumnTypeMapToTypeNotation()
	
	mfloat := c.blockData[blockIndex].FetchColumns(q["float"])
	fl := c.IndicesToColumnLabels(q["float"])

	mInt := c.blockData[blockIndex].FetchColumns(q["int"])
	il := c.IndicesToColumnLabels(q["int"])

	mstring := c.blockData[blockIndex].FetchColumns(q["string"])
	sl := c.IndicesToColumnLabels(q["string"])

	mvec := c.blockData[blockIndex].FetchColumns(q["vector"]) 
	vl := c.IndicesToColumnLabels(q["vector"])

	/// TODO: this declaration is not complete; modify accordingly. 
	return &CFBRDataMatrix{floatData: mfloat, floatDataColumnKeys: q["float"],
		floatDataColumnLabels: fl, intData: mInt, intDataColumnKeys: q["int"], 
		intDataColumnLabels: il, stringData: mstring, stringDataColumnKeys: q["string"],
		stringDataColumnLabels: sl, vectorData: mvec, vectorDataColumnKeys: q["vector"],
		vectorDataColumnLabels: vl}
}

func (c *ConcurrentFileBlockReader) EncodeStringData(columnIndices []int, blockIndex int, ohc string) *mat.Dense {

	mstring := c.blockData[blockIndex].FetchColumns(columnIndices)

	if mstring == nil {
		return nil 
	}

	if c.stringColumnEncodings == nil {
		c.stringColumnEncodings = make(map[int]map[string]float64, 0) 
	}

	var x *mat.Dense 
	switch {

	case (ohc == "dumb"): 
		x = OHEncodeBlock_Dumb_UsesPreviousEncoding(mstring, c.stringColumnEncodings) 

	default: 
		panic("! code the rest of the encoders !") 
	}

	return x 
}

/// 
/*
Converts partition block data to a float matrix.
Matrix column order is identical to that of input data

CAUTION: 
	will only work if vector and undef data are null

arguments: 
- ohc : dumb| frequency | class
*/ 
func (c *ConcurrentFileBlockReader) ConvertPartitionBlockToRegularMatrix(blockIndex int, ohc string) *mat.Dense {

	q := c.ColumnTypeMapToTypeNotation()

	// throw error for invalid data 
	if (len(q["undef"]) > 0) {
		panic("undefined columns, cannot convert!") 
	}

	if (len(q["vector"]) > 0) {
		panic("vector columns, cannot convert!") 
	}  
	 
	// encode string data
	x := c.EncodeStringData(q["string"], blockIndex, ohc) 	

	// make into regular matrix
	nr,nc := c.blockData[blockIndex].Dims() 
	stringEncodeDataIndex := 0
	output := mat.NewDense(nr,nc, nil)  

	for i := 0; i < nc; i++ {
		t := c.columnTypes[i]
		var data []float64
		switch {
		case t == "int" || t == "float": 
			data = DefaultStringSliceToFloat64Slice(c.blockData[blockIndex].GetAtCol(i)) 
		case t == "string": 
			data, _ = MatrixColToFloat64Slice(x, stringEncodeDataIndex)
			stringEncodeDataIndex++
		default: 
			panic("[CFBR] invalid data type")
		}

		output.SetCol(i, data)
	}
	
	return output
}

//// TODO: this needs to be checked!
func (c *ConcurrentFileBlockReader) SetColumnsForRemoval(outputIndices []int, removeIndices []int) []int {
	/// check arguments 
	for _,r := range removeIndices {
		if r >= len(c.columns) || r < 0 {
			panic("invalid index to remove")
		}
	}
	c.excludedColumnIndices = removeIndices
		
	// get new columns 
	l := len(c.columns)
	c.columns = StringSliceCollectByIndices(c.columns, c.excludedColumnIndices, true) 

	// shift all columns
	shifted:= make(map[int]int,0) // old index -> new index
	for i := 0; i < l; i++ {
		shifted[i] = i
	}

	for _,r := range removeIndices {
		for i := r +1; i <l; i++ {
			shifted[i]--
		}
	}

	// calibrate output indices 
	newOutput := make([]int,0)
	for _, oi := range outputIndices {
		newOutput = append(newOutput, shifted[oi])
	}
	
	// calibrate col types map
	newColTypes := make(map[int]string,0)
	for k,v := range shifted {
		if IntSInSlice(removeIndices, k) {
			continue
		}

		newColTypes[v] = c.columnTypes[k]
	}

	c.columnTypes = newColTypes
	return newOutput
}

/*
performs shutdown operation on file 
*/ 
func (c *ConcurrentFileBlockReader) Shutdown() {
	c.flsFileObj.Close()
}