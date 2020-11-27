package golem

/*
this file contains accuracy as well as time tests for
class ConcurrentFileBlockReader
*/

import (
	"fmt"
	"testing"
)

var TESTFP_1 string = "./datos/node_data/node_0"
var TESTFP_5 string = "./datos/node_data/node_5"

/*
*/ 
func SetUpMatrix1(fp string) *CFBRDataMatrix {
	x := OneConcurrentFileBlockReader(fp)
	x.ManualColumnTypeSet(NVMB_NODEDATA_TYPEMAP)  
	x.ReadPartition("full") 
	return x.ConvertPartitionBlockToMatrix(0) 
}

func Test_CFBRDataMatrix_ColumnLabelToIndexData(t *testing.T) {
	
	cm := SetUpMatrix1(TESTFP_1)
	fmt.Println("** ColumnLabelToIndexData **")	
	for k, _ := range NVMB_NODEDATA_TYPEMAP {
		cm.ColumnLabelToIndexData(k)
	}	
}

func Test_CFBRDataMatrix_FetchRowRangeNumerical(t *testing.T) {
	m := SetUpMatrix1(TESTFP_1)
	m.Preprocess()

	m2,_ := m.FetchRowRangeNumerical(0, 15, "int") 
	m2r,m2c := m2.Dims() 
	if m2r != 15 || m2c != 1 {
		panic("invalid dim for int retrieval")
	}

	fmt.Println("*** TCFR*** ")
	m3,_ := m.FetchRowRangeNumerical(0, 15, "float") 
	m3r,m3c := m3.Dims() 
	
	if m3r != 15 || m3c != 10 {
		panic("invalid dim for float retrieval")
	}
}

/// CAUTION: test depends on immutable data set 
func Test_CFBRDataMatrix_IndexRange(t *testing.T) {
	m := SetUpMatrix1(TESTFP_1) 
	m.Preprocess() 

	fd1, fd2 := m.convertedData["float"].Dims()	
	id1, id2 := m.convertedData["int"].Dims()	
	sd1, sd2 := m.stringData.Dims() 
	vd1, vd2 := m.vectorData.Dims() 

	if (fd1 != 349 || fd2 != 10) {
		panic("invalid dim for float data")
	}

	if (id1 != 349 || id2 != 1) {
		panic("invalid dim for int data")
	}

	if (sd1 != 349 || sd2 != 7) {
		panic("invalid dim for string data")
	}

	if (vd1 != 349 || vd2 != 7) {
		panic("invalid dim for vector data")
	}
}