package golem

import (
	"fmt"
	"testing"
	//"math"
)

func Test_Fetcher_FetchTimestampData(t *testing.T) {

	f := SetUpFetcher1(TESTFP_1)
	f.currentRowIndex = 0
	f.FetchTimestampData(0)

	/// TEST CASE 1	
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
		
		if answers[v_.varName] != v_.varValue {
			panic(fmt.Sprintf("incorrect variable value for %s, want %f got %f", v_.varName, answers[v_.varName], v_.varValue))
		}		
	}
	
	f.FetchTimestampData(13)

	endRow, _ := f.FetchTimestampData(13) 
	if endRow != 16 {
		panic(fmt.Sprintf("end row %d, want 16", endRow))
	}

	v = f.GatherVariablesAtTimestamp() 
	answers2 := map[string]float64{"cf gain" : 0,
								"contract gain" : 0,
								"number of contracts" : 0,
								"transmission" : 1,
								"cf acceptance rate" : 0,
								"competition" : 0.3155980110168457, 
								"greed" : 0.28494301438331604,
								"negotiation" : 0.24060100317001343,
								"growth" : 0.484127014875412,
								"nodeRank" : 0.6162790060043335,
								"currency" : 28.708572387695312,
								"node connectivity" : 0.3,
								"num paths ratio" : 0.5333330035209656,
								"targets path risk ratio" : 1.6349596666666668,
								"number of paths per target ratio" : 1.6666666666666667}


	for _, v_ := range v {
		if answers2[v_.varName] != v_.varValue {
			panic(fmt.Sprintf("incorrect variable value for %s", v_.varName))
		}
	}
}

func Test_Fetcher_FormatVariablesForDelta(t *testing.T) {
	f := SetUpFetcher1(TESTFP_1)
	f.GatherVariablesInitial()
	fmt.Println()
	m := f.FormatVariablesForDelta(0, 5) 
	fmt.Println("** Gather Variables Initial **")
	fmt.Println(m)
}

func Test_Fetcher_CalculateDelta(t *testing.T) {
	f := SetUpFetcher1(TESTFP_1)
	f.GatherVariablesInitial() 
	
	fmt.Println("END")
	fmt.Println(f.CalculateDelta("ends", 0, -1)) 
	fmt.Println() 

	fmt.Println("MEAN TO END")
	fmt.Println(f.CalculateDelta("mean-to-end", 0, -1)) 
	fmt.Println() 

	fmt.Println("MEAN TO START")
	fmt.Println(f.CalculateDelta("mean-to-start", 0, -1)) 
	fmt.Println() 

	fmt.Println("LINEAR")
	fmt.Println(f.CalculateDelta("linear", 0, -1)) 
	fmt.Println()  

	fmt.Println("---------------------------------------") 

	fmt.Println("END 2")
	fmt.Println(f.CalculateDelta("ends", 0, 10)) 
	fmt.Println() 
}