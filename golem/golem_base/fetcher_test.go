package golem

/*
this file contains accuracy as well as time tests for
class ConcurrentFileBlockReader
*/

import (
	"fmt"
	"testing"
	"math"
	//"time"
	//"reflect"
	
	//"strconv"
	//"encoding/csv"
	//"io"
)

func SetUpFetcher1(fp string) *Fetcher {

	m := SetUpMatrix1(fp)
	m.Preprocess()
	f := Fetcher{data: m}
	return &f
}

func Test_Fetcher_ImpactOfEventType(t *testing.T) {

	f := SetUpFetcher1(TESTFP_1) 
	f.FetchTimestampData(16) 
	impact1 := f.ImpactOfEventType("break_bond/phantom")
	if math.Abs(impact1 - (-10.179814)) > 0.01 {
		panic(fmt.Sprintf("wrong calculation, impact is %f, want %f", impact1, -10.179814)) 
	} 
}

func Test_Fetcher_DataValueAt(t *testing.T) {

	f := SetUpFetcher1(TESTFP_1)
	f.FetchTimestampData(16) 

	// case 1
	q := f.CollectValuesAtTimestamp("event_type", "summary", "transmission")
	answer := []DataValue{1.000} 
	if !EqualFloat64Slices(DataValueSliceToFloat64Slice(q), DataValueSliceToFloat64Slice(answer), 3) {
		///fmt.Println("WRONGGGGH")
		panic("[0] wrong collected values")
	}

	// case 2 
	q2 := f.CollectValuesAtTimestamp("event_type", "break_bond/phantom", "impact") 
	answer2 := []DataValue{-3.304824, -3.304824, -3.570166} 
	if !EqualFloat64Slices(DataValueSliceToFloat64Slice(q2), DataValueSliceToFloat64Slice(answer2), 3) {
		panic(fmt.Sprintf("[1] wrong collected values, want %v, got %v", q2, answer2))
	}
}

// TODO: make tests for GetTimestampRowRange