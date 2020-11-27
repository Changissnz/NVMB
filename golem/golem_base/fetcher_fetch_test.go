package golem

import (
	"fmt"
	"testing"
	//"math"
)

func Test_Fetcher_FetchTimestampData(t *testing.T) {

	f := SetUpFetcher1(TESTFP_1)
	f.currentRowIndex = 0

	/// TEST CASE 1	
	f.FetchTimestampData(0)
	answers := map[string]float64{"cf gain" : 0,
								"contract gain" : 0,
								"number of contracts" : 0,
								"transmission" : -1,
								"cf acceptance rate" : -1,
								"competition" : 0.3155980110168457, 
								"greed" : 0.28494301438331604,
								"negotiation" : 0.24060100317001343,
								"growth" : 0.484127014875412,
								"nodeRank" : 1,
								"currency" : 100,
								"node connectivity" : 1.0,
								"num paths ratio" : 0,
								"targets path risk ratio" : 0,
								"number of paths per target ratio" : 0}
	

	v := f.GatherVariablesAtTimestamp()
	for _, v_ := range v {

		_,ok := answers[v_.varName] 
		if !ok {
			continue
		}

		if answers[v_.varName] != v_.varValue {
			panic(fmt.Sprintf("incorrect variable value for %s, want %f got %f", v_.varName, answers[v_.varName], v_.varValue))
		}		
	}

	/// TEST CASE 2
	endRow, _ := f.FetchTimestampData(17) 
	if endRow != 23 {
		panic(fmt.Sprintf("end row %d, want 23", endRow))
	}

	v = f.GatherVariablesAtTimestamp() 
	
	fmt.Println("gathering done ", v == nil) 

	// TODO: check for values in v!
	///DisplayVarSlice(v)
}

func Test_Fetcher_FormatVariablesForDelta(t *testing.T) {
	f := SetUpFetcher1(TESTFP_1)
	f.GatherVariablesInitial()
	fmt.Println()
	m := f.FormatVariablesForDelta(0, 5) 
	fmt.Println("** Gather Variables Initial **")
	fmt.Println(m)
}