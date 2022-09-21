package golem

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	//"log"
)

/// TODO: data generation could use some work.
/*
*/ 

var CPARTSIZE int = 1000 // concurrency reader partition size
const DEFAULT_MAX_CLASSES int = 100

/*
*/
type Pair struct {
	a, b interface{}
}

func DisplayPair(p *Pair) {
	fmt.Println(p.a, " ", p.b) 
}

func DisplayPairSlice(p []*Pair) {
	for _,p_ := range p {
		DisplayPair(p_) 
		fmt.Println("----------------------")
	}
}


func containsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

/*
opens a file at path
*/
func Openness(fp string) *os.File {

	fi, err := os.OpenFile(fp, os.O_APPEND|os.O_RDWR, 0666)
	//fi, err := os.OpenFile(fp, os.O_TRUNC, 0666)
	if err != nil {
		panic("invalid output file path")
	}
	return fi
}

/*
// description
creates new file at path if it does not exist

// return
new or existing file object
*/
func Newness(fp string) *os.File {

	fi, err := os.OpenFile(fp, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fi2,_ := os.Create(fp)
		return fi2
	}
	return fi
}

/*
*/
func GenerateDefaultRowLabels(numRows int) []string {

	var x []string = make([]string, 0)

	for i := 0; i < numRows; i++ {
		x = append(x, "var_"+strconv.Itoa(i))
	}
	return x
}

/*
generates n random string classes starting from
*/
func GenerateRandomStringClassesDefault(n int, startIndex int) []string {

	x := make([]string, 0)
	endo := n + startIndex
	for i := startIndex; i < endo; i++ {
		q := fmt.Sprintf("DefClass_%d", i)
		x = append(x, q)
	}
	return x
}

/*
makes map of keys (`indices`) to slice of string types starting at startIndex
*/
func GenerateRandomStringClassesDefaultToMap(indices []int, startIndex int) map[int][]string {
	x := make(map[int][]string, 0)

	for _, y := range indices {

		// get random number of key types
		q := RandomInt(1, int32(DEFAULT_MAX_CLASSES+1))
		s := GenerateRandomStringClassesDefault(int(q), startIndex)
		startIndex += len(s)
		x[y] = s
	}

	return x
}

///////////////////////////// START: generation methods for simple (numerical) data

/*
*/
func GenerateRandomRowData(numColumns int, randomFloatRange Pair, booleanVariables []int) []string {

	output := make([]string, 0)
	var b float32
	var s string

	for i := 0; i < numColumns; i++ {
		if containsInt(booleanVariables, i) {
			b = RandomBoolean()
			s = strconv.Itoa(int(b))
			output = append(output, s)
			continue
		}

		b = RandomFloat(randomFloatRange.a.(float32), randomFloatRange.b.(float32))
		s = DefaultFloatToString(b)
		output = append(output, s)
	}

	return output
}

/*
generates `l` rows of random numerical data
*/
func GenerateRandomRowDatas(l int, numColumns int, randomFloatRange Pair, booleanVariables []int) [][]string {

	output := make([][]string, 0)
	for i := 0; i < l; i++ {
		output = append(output, GenerateRandomRowData(numColumns, randomFloatRange, booleanVariables))
	}

	return output
}

/*
this generator method will generate a file with only numerical values
that denote either bools or floats
*/
func GenerateRandomCSVDataToFile(fp string, size int, numColumns int, randomFloatRange Pair, booleanVariables []int) {

	// check if file exists
	fi := Newness(fp)

	// write the rows first
	rowLabels := GenerateDefaultRowLabels(numColumns)
	writer := csv.NewWriter(fi)
	writer.Write(rowLabels)

	// partition the number of columns
	var p Pair = GetPartitionInfo(size)

	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	// write the first x full partitions
	for i := 0; i < p.a.(int); i++ {
		wg.Add(1)
		go func() {

			// generate data
			for j := 0; j < CPARTSIZE; j++ {
				gen := GenerateRandomRowData(numColumns, randomFloatRange, booleanVariables)
				mutex.Lock()
				writer.Write(gen)
				writer.Flush()
				mutex.Unlock()
			}

			wg.Done()
		}()
	}

	wg.Wait()

	last := GenerateRandomRowDatas(p.b.(int), numColumns, randomFloatRange, booleanVariables)
	writer.WriteAll(last)
	writer.Flush()
}

///////////////////////////// END: generation methods for simple (numerical) data

///////////////////////////// START: generation methods for extended (numerical/string) data

/*
generates one row of random numerical/string data

---arguments
numColumns := size of row
randomFloatRange := Pair of two float32's
booleanValues := indices that are booleans
stringVariables := indices that are strings
stringVarValueMap :=
invalidColumns := indices that are invalid values
*/
func GenerateRandomRowDataExtended(numColumns int, randomFloatRange Pair, booleanVariables []int,
	stringVariables []int, stringVarValueMap map[int][]string) []string {

	datos := make([]string, 0)
	for i := 0; i < numColumns; i++ {
		q := GenerateOneRandomCSVDataExtended(i, randomFloatRange, booleanVariables,
			stringVariables, stringVarValueMap, make([]int, 0))
		datos = append(datos, q)
	}

	return datos
}

/*
helper method
generates data, outputs its string representation
*/
func GenerateOneRandomCSVDataExtended(index int, randomFloatRange Pair, booleanVariables []int,
	stringVariables []int, stringVarValueMap map[int][]string, invalidColumns []int) string {

	var estos string
	if containsInt(invalidColumns, index) {
		return "NULLOS"
	}

	// string var
	if containsInt(stringVariables, index) {
		// random string
		estos, _ = RandomElementInStringSlice(stringVarValueMap[index])
	} else if containsInt(booleanVariables, index) {
		estos = strconv.Itoa(int(RandomBoolean()))
	} else {
		estos = DefaultFloatToString(RandomFloat(randomFloatRange.a.(float32), randomFloatRange.b.(float32)))
	}

	return estos
}

/*
generates key type data for each key
*/
// TODO: need to add probability distribution to random selection
func GenerateRandomColumnClasses(keys []int, maxClasses int) map[int][]string {

	maposMarkos := make(map[int][]string, 0)
	indexios := 0
	for _, k := range keys {
		q := RandomInt(0, int32(maxClasses))
		g := GenerateRandomStringClassesDefault(int(q), indexios)
		maposMarkos[k] = g
		indexios += len(g)
	}

	return maposMarkos
}

/*
one random output value by designation
*/
func GenerateOneRandomCSVDataExtendedFlawInbut(index int, randomFloatRange Pair, booleanVariables []int,
	stringVariables []int, stringVarValueMap map[int][]string, flawPerCol map[int]float32) string {

	// get flawed columns
	if _, ok := flawPerCol[index]; ok {
		q := RandomFloat(0.0, 1.0)
		if q >= flawPerCol[index] {
			return "NULLOS"
		}
	}

	return GenerateOneRandomCSVDataExtended(index, randomFloatRange, booleanVariables, stringVariables, stringVarValueMap, make([]int, 0))
}

/*
one row by designation
*/
func GenerateRandomRowDataExtendedFlawInbut(numColumns int, randomFloatRange Pair, booleanVariables []int,
	stringVariables []int, stringVarValueMap map[int][]string, flawPerCol map[int]float32) []string {

	datos := make([]string, 0)
	var q string
	for i := 0; i < numColumns; i++ {
		q = GenerateOneRandomCSVDataExtendedFlawInbut(i, randomFloatRange, booleanVariables,
			stringVariables, stringVarValueMap, flawPerCol)
		datos = append(datos, q)
	}
	return datos
}

/*
generates l rows by above method
*/
// will use concurrency
func GenerateRandomRowDatasExtendedFlawInbut(l int, numColumns int, randomFloatRange Pair, booleanVariables []int, stringVariables []int, stringVarValueMap map[int][]string, flawPerCol map[int]float32) [][]string {

	datos := make([][]string, 0)
	for i := 0; i < l; i++ {
		r := GenerateRandomRowDataExtendedFlawInbut(numColumns, randomFloatRange, booleanVariables, stringVariables, stringVarValueMap, flawPerCol)
		datos = append(datos, r)
	}

	return datos
}

/*
This is the extended version of above method. Capable of generating
csv with string data columns.

Will generate data based on unclean measure, a map that determines the number of invalid values at each column.
Invalid values will be either mismatched types or null.

if numStrings != -1|null, then

flawsos := if nil, no error
*/
func GenerateRandomCSVDataToFileExtended(fp string, size int, numColumns int, randomFloatRange Pair, booleanVariables []int, stringVariables []int, flawso map[int]float32, headerOn bool) {

	// generate header
	fi := Newness(fp)

	// write the headers first
	rowLabels := GenerateDefaultRowLabels(numColumns)
	writer := csv.NewWriter(fi)
	writer.Write(rowLabels)

	// partition the number of columns
	var p Pair = GetPartitionInfo(size)

	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	// generate variable string variables
	stringus := GenerateRandomStringClassesDefaultToMap(stringVariables, 0)

	// write the first x full partitions
	for i := 0; i < p.a.(int); i++ {
		wg.Add(1)
		go func() {

			// generate data
			for j := 0; j < CPARTSIZE; j++ {
				gen := GenerateRandomRowDataExtendedFlawInbut(numColumns, randomFloatRange,
					booleanVariables, stringVariables, stringus, flawso)

				mutex.Lock()
				writer.Write(gen)
				writer.Flush()
				mutex.Unlock()
			}

			wg.Done()
			// write to file
		}()
	}

	wg.Wait()

	// write last partition
	last := GenerateRandomRowDatasExtendedFlawInbut(p.b.(int), numColumns, randomFloatRange, booleanVariables, stringVariables, stringus, flawso)
	writer.WriteAll(last)
	writer.Flush()
}

///////////////////////////// END: generation methods for extended (numerical/string) data; unclean

//////////////////////////// START: sample dataset generation //////////////////////

/*
generates a flawed file by specifications.

null/error probability rate for each column is at least 25%
*/
func MakeFlawedCSVFile(fp string, sizos int) ([]int32, []int32, map[int]float32) {
	os.Create(fp)

	nc := 20
	p := Pair{a: float32(3.5), b: float32(12.5)}

	// types of columns
	/// 25% boolean
	/// 25% string
	/// 50% float
	bv := RandomUniqueIntSlice(int32(0), int32(20), 5, make([]int32, 0))
	sv := RandomUniqueIntSlice(int32(0), int32(20), 5, bv)
	bv_ := Int32SliceToIntSlice(bv)
	sv_ := Int32SliceToIntSlice(sv)
	GenerateRandomCSVDataToFileExtended(fp, sizos, nc, p, bv_, sv_, nil, true)
	return bv, sv, nil
}

//////////////////////////// END: sample dataset generation //////////////////////

/*
calculates partitioning info based on size
returns a Pair{# of full partitions, size of last partition}
*/
func GetPartitionInfo(size int) Pair {

	numX := size / CPARTSIZE
	rem := size - (numX * CPARTSIZE)

	return Pair{numX, rem}
}

///////////////////////////////// START: write to csv file ////////////////////
/*
writes csv string to file in append mode
*/
func CSVDataToFile(fp string, cs [][]string) bool {

	// check if file exists
	fi, err := os.Open(fp)

	if err != nil {
		fmt.Printf("invalid file path @ %s\n", fp)
		return false
	}

	writer := csv.NewWriter(fi)
	writer.WriteAll(cs)
	return true
}

func CSVDataToFileObject(fi *os.File, cs [][]string) {
	writer := csv.NewWriter(fi)
	writer.WriteAll(cs)
}

///////////////////////////////// END: write to csv file ////////////////////


func GenerateFileDefaultX() {
	fp := "./golem_base/datos/defaultx" 
	sz := 1000
	nc := 20 
	rfr := Pair{float32(3.4), float32(10.5)} 
	bv := make([]int, 0) 
	GenerateRandomCSVDataToFile(fp, sz, nc, rfr, bv)	
}

func GenerateFileDefault() {
	/*
	path, err := os.Getwd()
	if err != nil {
    	log.Println(err)
	}
	
	fmt.Println("PATHOS")
	fmt.Println(path) 

	fp := path + "/golem_base/datos/default" 
	fmt.Println("PATHOS2")
	fmt.Println(fp) 
	*/ 
	fp := "./golem_base/datos/default" 

	sz := 10050
	nc := 20 
	rfr := Pair{float32(3.4), float32(10.5)} 
	bv := make([]int, 0) 
	GenerateRandomCSVDataToFile(fp, sz, nc, rfr, bv)	
}

func GenerateFileDefault2() {
	fp := "./golem_base/datos/default2" 
	p := Pair{a: float32(3.5), b: float32(12.5)} 
	bv := []int{0,1,3}
	sv := []int{1,4,6,7} 
	GenerateRandomCSVDataToFileExtended(fp, 100, 10, p, bv, sv, nil, true) 
}

func GenerateFileDefault3() {
	fp := "./golem_base/datos/default3" 
	p := Pair{a: float32(3.5), b: float32(12.5)} 
	bv := []int{}
	sv := []int{4,6,8,9,11}

	/*
	var answers map[int]string = map[int]string{0: "float",
	1: "float",
	2: "float",
	3: "float",
	4: "string",
	5: "float",
	6: "string",
	7: "float",
	8: "string",
	9: "string",
	10: "float",
	11: "string",
	12: "float",
	13: "float",
	14: "float",
	15: "float",
	16: "float",
	17: "float",
	18: "float",
	19: "float"}
	*/

	var flawsos map[int]float32 = map[int]float32{0:0.2,
		1: 0.2,
		2: 0.2,
		3: 0.2,
		4: 0.2,
		5: 0.2,
		6: 0.2,
		7: 0.2,
		8: 0.2,
		9: 0.2,
		10: 0.2,
		11: 0.2,
		12: 0.2,
		13: 0.2,
		14: 0.2,
		15: 0.2,
		16: 0.2,
		17: 0.2,
		18: 0.2,
		19: 0.2}
	
	GenerateRandomCSVDataToFileExtended(fp,100,10,p,bv,sv, flawsos,true)
}