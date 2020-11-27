package golem

/*
this file contains accuracy as well as time tests for
class ConcurrentFileBlockReader
*/

import (
	"fmt"
	"testing"
	"time"
	"reflect"
)

const TEST_FP string = "./datos/defaultx"

/*
Tests that declaration of `ConcurrentFileBlockReader` is correct. 
*/ 
func Test_ConcurrentFileBlockReader_Declaration(t *testing.T) {

	answer := []string{"var_0", "var_1", "var_2", "var_3",
		"var_4", "var_5", "var_6", "var_7", "var_8",
		"var_9", "var_10", "var_11", "var_12", "var_13",
		"var_14", "var_15", "var_16", "var_17", "var_18",
		"var_19"}

	x := OneConcurrentFileBlockReader(TEST_FP)

	// assert class has correct column names 
	if !EqualStringSlices(x.columns, answer) {
		panic("header column read is incorrect")
	}

	x.Shutdown()
}

/*
Tests the method <ConcurrentFileBlockReader.ReadBlockAtSpot>
*/ 
func Test_ConcurrentFileBlockReader_ReadBlockAtSpot(t *testing.T) {
	///// case 1 
	GenerateFileDefault() 
	fp := "./datos/default"
	x := OneConcurrentFileBlockReader(fp)

	b, sz, stat := x.ReadBlockAtSpot(3,true)

	if stat {
		panic("file not done")
	}

	if len(b.datos) != 3 {
		panic("incorrect length")
	}

	if sz != 60 {
		panic("incorrect size")
	}

	x.Shutdown() 

	///// case 2
	x2 := OneConcurrentFileBlockReader(TESTFP_1)	
	c := 0 
	for {
		_, _, stat := x2.ReadBlockAtSpot(3,true)
		c++ 
		if stat {
			break 
		}
	}

	if c != 117 {
		panic(fmt.Sprintf("invalid number of blocks read %d, want 117", c))
	}
}

/*
time tests reading partition 
*/ 
func Test_ConcurrentFileBlockReader_TimeCapture(t *testing.T) {
	fmt.Println("** TIME-TEST: partition read **")
	x := OneConcurrentFileBlockReader(TEST_FP)

	fmt.Println("###############################################")
	fmt.Println("COLS ARE:\t", x.columns)
	start := time.Now()
	//... operation that takes 20 milliseconds ...
	x.ReadPartition("exact")
	ti := time.Now()
	elapsed := ti.Sub(start)
	fmt.Println(fmt.Sprintf("elapsed time non-concurrent %v", elapsed))

	x.Shutdown()
}

/*
Tests for correct partition dimensions of read by DEFAULT values. 
TODO : more tests on different WRITESIZE, CPARTSIZE
*/
func Test_ConcurrentFileBlockReader_ReadPartition(t *testing.T) {
	x := OneConcurrentFileBlockReader(TEST_FP)

	x.ReadPartition("exact")

	if len(x.blockData) != 1 {
		panic("incorrect # of blocks")
	}

	b := x.blockData[0]

	if b.Length() != 1000 {
		panic("incorrect length")
	}

	if b.Width() != 20 {
		panic("incorrect width")
	}

	if b.Size() != 20000 {
		panic("incorrect size")
	}

	x.Shutdown() 
}

/*
 */
func Test_ChooseRandomRowDataIndices(t *testing.T) {

	x := OneConcurrentFileBlockReader(TEST_FP)

	// read the first partition
	x.ReadPartition("exact")
	x.ChooseRandomRowDataIndices(20)

	if x.oi.Len() != 20 {
		panic("Error getting random row indices")
	}
}

func Test_ConcurrentFileBlockReader_Read2DeduceOneRowDataIndex(t *testing.T) { //index []int) string {

	i1 := []int{0, 0, 0}
	i2 := []int{0, 1, 0}

	/// CASE 1: default TEST_FP
	x := OneConcurrentFileBlockReader(TEST_FP)
	x.ReadPartition("exact")
	x.blockMapCollector = OneMapCollector()
	x.Read2DeduceOneRowDataIndex(i1, true)
	x.Read2DeduceOneRowDataIndex(i2, true)

	for k, v := range x.blockMapCollector.datos {
		_, ok := v["float"]

		if !ok {
			panic(fmt.Sprintf("invaid type deduction for %d of test case 1", k))
		}
	}
}

/*
*/
func TestDeduceBasicStringType(t *testing.T) {
	q := DeduceBasicStringType("-42.13")
	q2 := DeduceBasicStringType("141334")
	q3 := DeduceBasicStringType("afl;jsd@#!%")

	if q != "float" {
		panic("error type [0]")
	}

	if q2 != "int" {
		panic("error type [1]")
	}

	fmt.Println("ERROR TYPE [2] ", q3)
	if q3 != "string" {
		panic("error type [2]")
	}
}

func Test_ConcurrentFileBlockReader__MapCollect__PredictKeyOfTypeByRequiredThreshold(t *testing.T) {

	fmt.Println("** ConcurrentFileBlockReader__MapCollect__PredictKeyOfTypeByRequiredThreshold **")

	//// test file case 0 
	x := OneConcurrentFileBlockReader(TEST_FP)
	x.DeduceColumnTypes()
	for _, v := range x.columnTypes {
		if v != "float" {
			panic(fmt.Sprintf("invalid type %s", v))
		}
	}

	//// test file case 1
	fmt.Println("CASE 1")
	fp3 := "./datos/default3"
	x = OneConcurrentFileBlockReader(fp3)
	x.DeduceColumnTypes() 

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
	
	// check default3 by minimum threshold 
	for k,_ := range x.blockMapCollector.datos {
		q := x.blockMapCollector.PredictKeyOfTypeByMode("threshold", k, 0.001)
		if q != answers[k] {
			panic(fmt.Sprintf("[1.0] wrong prediction %s for index %d, want %s", q, k, answers[k]))
		}
	}

	// check default3 by DEFAULT_DOMINANCE_RATIO
	for k,_ := range x.blockMapCollector.datos {
		q := x.blockMapCollector.PredictKeyOfTypeByMode("threshold", k, DEFAULT_DOMINANCE_RATIO)
		if q != "?" {
			panic(fmt.Sprintf("[1.1] wrong prediction %s for index %d, want %s", q, k, answers[k]))
		}
	}
}

// TODO: incomplete test
func Test_ConcurrentFileBlockReader_ConvertPartitionBlockToMatrix(t *testing.T) {

	fmt.Println("** ConvertPartitionBlockToMatrix **")
	fp3 := "./datos/default"
	x := OneConcurrentFileBlockReader(fp3) 
	x.ReadPartition("exact") 
	x.DeduceColumnTypes()
	
	q := x.ColumnTypeMapToTypeNotation()
	x.blockData[0].FetchColumns(q["float"]) //?
	x.ConvertPartitionBlockToMatrix(0)
}


////////////////////////////////////////////////////////////////////////////////////////////////////////
func Test_ConcurrentFileBlockReader__DeduceColumnTypes__DataIsCreditCard(t *testing.T) {

	// case 1: 
	fp3 := "./datos/default3"
	x1 := OneConcurrentFileBlockReader(fp3) 

	fmt.Println("trial")
	x1.DeduceColumnTypes()
	fmt.Println("column types for: default3")
	fmt.Println(x1.columnTypes)
	fmt.Println(reflect.TypeOf(x1))

	// case 2: credit card dataset  
	fp := "./datos/other_data/creditcard.csv" 
	x := OneConcurrentFileBlockReader(fp)

	x.DeduceColumnTypes() 
	
	answers := map[int]string{0:"int", 1:"float", 2:"float", 3:"float", 4:"float",
	 5:"float", 6:"float", 7:"float", 8:"float", 9:"float",
	 10:"float", 11:"float", 12:"float", 13:"float", 14:"float",
	 15:"float", 16:"float", 17:"float", 18:"float", 19:"float",
	 20:"float", 21:"float", 22:"float", 23:"float", 24:"float",
	 25:"float", 26:"float", 27:"float", 28:"float", 29:"float", 30:"int"}

	 for k, v := range answers {
		if v != x.columnTypes[k] {
			panic(fmt.Sprintf("incorrect variable value for %d, want %s got %s", k, v, x.columnTypes[k]))
		}		
	}
} 

/*
*/ 
func Test_ConcurrentFileBlockReader_ReadPartition__DeletedRow__DataIsCreditCard(t *testing.T) {

	fp := "./datos/other_data/creditcard.csv" 
	x1 := OneConcurrentFileBlockReader(fp) 
	x1.DeduceColumnTypes()

	// remove columns here
	x1.SetColumnsForRemoval([]int{}, []int{0, 3, 5}) 
	answer := map[int]string{
		0 :  "float",
		1 :  "float",
		2 :  "float",
		3 :  "float",
		4 :  "float",
		5 :  "float",
		6 :  "float",
		7 :  "float",
		8 :  "float",
		9 :  "float",
		10 :  "float",
		11 :  "float",
		12 :  "float",
		13 :  "float",
		14 :  "float",
		15 :  "float",
		16 :  "float",
		17 :  "float",
		18 :  "float",
		19 :  "float",
		20 :  "float",
		21 :  "float",
		22 :  "float",
		23 :  "float",
		24 :  "float",
		25 :  "float",
		26 :  "float",
		27 :  "int"}

	for k,v := range x1.columnTypes {
		
		if v != answer[k] {
			panic(fmt.Sprintf("incorrect value for key %d", k))
		}
	}

	/// test reading one partition
	x1.ReadPartition("exact") 
	_,c := x1.blockData[0].Dims()

	if c != len(x1.columnTypes) {
		panic("partition number of columns incorrect!") 
	}

}

/*
check to see that null values such as ' ', NaN, null, etc. are registered appropriated in matrix
representation of data. 

This method will be the starting point to converting nvmb initial output data to readable
matrix data for training. 
*/ 
func Test_ConcurrentFileBlockReader__DataIsNode0(t *testing.T) {
	fp := "./datos/node_data/node_1" 
	x1 := OneConcurrentFileBlockReader(fp) 
	x1.DeduceColumnTypes()

	fmt.Println("predicted column types") 
	fmt.Println(x1.columnTypes)
}