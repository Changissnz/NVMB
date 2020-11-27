package golem

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
)

type DataValue interface{}

/*
The Fetcher struct is used to fetch related variables pertaining to data 
found in `CFBRDataMatrix`. 

NOTE: Only NVMBNodeData fetcher type is currently supported
*/
type Fetcher struct {
	data *CFBRDataMatrix
	
	// load current timestamp
	tInt *mat.Dense // timestamp info numerical 
	tFloat *mat.Dense
	tString [][]string // timestamp info string
	tVector [][]string
	currentRowIndex int /// TODO: unused 

	analysis map[int][]*Variable
}

func OneFetcher(d *CFBRDataMatrix) *Fetcher {
	return &Fetcher{data: d, currentRowIndex: 0} 
}

/*
*/
func (f *Fetcher) GetTimestampRowRange(rowIndex int) (int,int) {

	colIndex, dataType := f.data.ColumnLabelToIndexData("time")
	ts, exists := f.data.FetchValueNumerical(rowIndex, colIndex, dataType)
	if !exists {
		return -1, -1
	}

	cr := rowIndex + 1
	for {
		ts2, _ := f.data.FetchValueNumerical(cr, colIndex, dataType)
		if ts2 != ts {
			break
		}
		cr++
	}

	return rowIndex,cr
}

/*
*/
func (f *Fetcher) FetchTimestampData(rowIndex int) (int,bool) {

	f.ClearTimestampData()
	start, end := f.GetTimestampRowRange(rowIndex) 
	if (start == -1) {
		return -1, false
	}

	f.tInt,_ = f.data.FetchRowRangeNumerical(start,end,"int")
	f.tFloat,_ = f.data.FetchRowRangeNumerical(start,end,"float") 
	f.tString,_ = f.data.FetchRowRangeString(start,end, "string")
	f.tVector,_ = f.data.FetchRowRangeString(start,end, "vector") 
	return end, true
}

/*
*/
func (f *Fetcher) GetColumnValuesAtTimestampString(column string) []string {
	ind,t := f.data.ColumnLabelToIndexData(column) 
	
	if t != "string" && t != "vector" {
		panic(fmt.Sprintf("invalid data type %s", t))
	}
	
	x := make([]string, 0)
	switch {
		case (t == "string"):
			for _, q := range f.tString {
				x = append(x, q[ind])
			}

		default:
			for _, q := range f.tVector {
				x = append(x, q[ind])
			}

	}
	
	return x
}

/*
*/ 
func (f *Fetcher) GetColumnValuesAtTimestampFloat(column string) []float64 {
	ind,t := f.data.ColumnLabelToIndexData(column) 
	
	if t != "float" && t != "int" {
		panic(fmt.Sprintf("invalid data type %s", t))
	}
	
	x := make([]float64, 0)
	switch {
	case (t == "float"):

		r,_ := f.tFloat.Dims() 

		for i := 0; i < r; i++ {
			row, _ := MatrixRowToFloat64Slice(f.tFloat, i)
			x = append(x, row[ind])
		}

	default:
		r,_ := f.tInt.Dims() 

		for i := 0; i < r; i++ {
			row, _ := MatrixRowToFloat64Slice(f.tInt, i)
			x = append(x, row[ind])
		}
	}

	return x
}

/*
*/
func (f *Fetcher) ClearTimestampData() {
	f.tInt = nil 
	f.tFloat = nil  
	f.tString = nil  
	f.tVector = nil 
}

/*
dataType := int|float|string|vector
index := 2-d index 

CAUTION:  does not perform index checking 
*/
func (f *Fetcher) DataValueAt(dataType string, index []int) DataValue {

	if dataType == "int" {
		return f.tInt.At(index[0], index[1])
	}

	if dataType == "float" {
		return f.tFloat.At(index[0], index[1]) 
	}

	if dataType == "string" {
		return f.tString[index[0]][index[1]]
	}

	if dataType == "vector" {
		return f.tVector[index[0]][index[1]]
	}

	panic(fmt.Sprintf("data type %s invalid", dataType))
}

/// TODO: below needs to be deleted. Re-factor associated test 
/*
sums up impact by event_type
*/ 
func (f *Fetcher) ImpactOfEventType(eventType string) float64 {
	var impact float64 = 0

	// get impact column and eventType column
	index1, _ := f.data.ColumnLabelToIndexData("event_type")
	index2, _ := f.data.ColumnLabelToIndexData("impact")

	// iterate through and consider all rows with eventType
	r,_ :=  f.tInt.Dims()
	for i := 0; i < r; i++ {
		if f.tString[i][index1] == eventType {
			impact += f.tFloat.At(i, index2) 
		}
	}
	return impact 
}

/*
referenceCol := column label used to filter out timestamp data 
referenceColVal := value of column label to equal
columnToCollect := column label in which values will be collected
*/
func (f *Fetcher) CollectValuesAtTimestamp(referenceColLabel string, referenceColVal DataValue, columnToCollect string) []DataValue {

	// get impact column and eventType column
	index1, dt1 := f.data.ColumnLabelToIndexData(referenceColLabel)
	index2, dt2 := f.data.ColumnLabelToIndexData(columnToCollect) 

	// declare output collector
	var output []DataValue
	r,_ := f.tInt.Dims() 
	for i := 0; i < r; i++ {
		switch {
			case (dt1 == "int" || dt1 == "float"):
				if f.DataValueAt(dt1, []int{i,index1}).(float64) == referenceColVal.(float64) {
					output = append(output, f.DataValueAt(dt2, []int{i,index2}))
				}
			case (dt1 == "string" || dt1 == "vector"):				
				if f.DataValueAt(dt1, []int{i, index1}).(string) == referenceColVal.(string) {
					output = append(output, f.DataValueAt(dt2, []int{i,index2}))
				}
			default: 
				panic("could not proceed with collecting values at timestamp")
		}		
	}

	return output 
}

/*
eventType := key for event_type column
deltaCol := column to sum delta up on
deltaOperation := 

CAUTION: no error-checking performed
*/
func (f *Fetcher) DeltaOfEvent(eventType string, deltaCol string, deltaOperation string) float64 {

	collected := DataValueSliceToFloat64Slice(f.CollectValuesAtTimestamp("event_type", eventType, deltaCol))
	return DeltaOpOnSequence(collected, deltaOperation)
}

/*
Outputs the ratio of active contracts over known nodes.
*/
func (f *Fetcher) ActiveContractKnownRatio() float64 {
	numContractsString := f.CollectValuesAtTimestamp("event_type", "summary", "contracts")[0]
	iSlice, _ := DefaultStringToIntSlice(numContractsString.(string), ' ')
	var size int = len(iSlice)

	numKnownString := f.CollectValuesAtTimestamp("event_type", "summary", "knownNodes")[0]
	iSlice2, _ := DefaultStringToIntSlice(numKnownString.(string), ' ')
	var size2 int = len(iSlice2) 
	return ZeroDiv(float64(size), float64(size2), 0.0, 0.0)
}

/*
Outputs the ratio of neighbors over known nodes.
*/
func (f *Fetcher) NodeConnectivityKnownRatio() float64 {
	neighborsString := f.CollectValuesAtTimestamp("event_type", "summary", "neighbors")[0]
	iSlice, _ := DefaultStringToIntSlice(neighborsString.(string), ' ') 
	var size int = len(iSlice) 

	numKnownString := f.CollectValuesAtTimestamp("event_type", "summary", "knownNodes")[0]
	iSlice2, _ := DefaultStringToIntSlice(numKnownString.(string), ' ')
	var size2 int = len(iSlice2) 

	return ZeroDiv(float64(size), float64(size2), 0.0, 0.0)
}