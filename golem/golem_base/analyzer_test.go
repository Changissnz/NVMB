package golem

/*
this file contains accuracy as well as time tests for
class ConcurrentFileBlockReader
*/

import (
	"fmt"
	"testing"
)

/*
*/
func Test_Analyzer_LoadOneRead(t *testing.T) {
	a := OneAnalyzer_(TESTFP_1) 
	CPARTSIZE = 3

	stat := a.LoadOneRead("exact",5) 
	fmt.Println("STAT ", stat)
	
	if !stat {
		panic("[0] incorrect stat") 
	}
	
	if len(a.readr.blockData) != 117 {
		panic("[0] incorrect number of blocks")
	}

	CPARTSIZE = 1000
}

/*
*/ 
func Test_Analyzer_CaptureAtTimestamp(t *testing.T) {
	a := OneAnalyzer_(TESTFP_1) 

	stat := a.LoadOneRead("exact",5) 
	fmt.Println("loading one ", stat)
	a.fetchr.GatherVariablesInitial() 

	fmt.Println("X")

	ic := a.CaptureAtTimestamp(4, "ends", 2, 3)

	fmt.Println("AFTER")

	ic.Display([]string{"input", "control", "output"})
}

func Test_Analyzer_SlideFetcher(t *testing.T) {
	// iterate through blocks and check for dimensions. 
	CPARTSIZE = 4
	a := OneAnalyzer_(TESTFP_1)
	stat := a.LoadOneRead("exact", 1)	
	c := 0 
	numBlocks := len(a.readr.blockData)

	for {
		stat = a.SlideFetcher(1)
		if stat {
			break 
		}

		c+= 1 		
		if c == 88 && a.fetchr.data.Length() != 2 {
			panic("block 88 incorrect length")
		} 
		
		if c != 88 && a.fetchr.data.Length() != 5 {
			panic(fmt.Sprintf("block %d incorrect length", c))
		}
		 	
	}

	
	if c != numBlocks {
		panic(fmt.Sprintf("incorrect number of reads %d, want 88", c)) 
	}
	 
	CPARTSIZE = 1000 
}


/*
displays all delta calculations delta calculations 
*/
func Test_Analyzer__FetcherCalculateDelta__SizeCorrect(t *testing.T) {

	a := OneAnalyzer_(TESTFP_1)
	a.LoadOneRead("exact", 1)
	a.fetchr.GatherVariablesInitial() 

	//// for displaying the results 
	fmt.Println("END")
	x := a.fetchr.CalculateDelta("ends", 0, -1)
	if len(x) != 22 {
		panic("incorrect number of variables")
	}
	fmt.Println(x) 
	fmt.Println() 

	fmt.Println("MEAN TO END")
	x2 := a.fetchr.CalculateDelta("mean-to-end", 0, -1)
	if len(x2) != 22 {
		panic("incorrect number of variables")
	}
	fmt.Println(x2) 
	fmt.Println() 

	fmt.Println("MEAN TO START")
	x3 := a.fetchr.CalculateDelta("mean-to-start", 0, -1)
	if len(x3) != 22 {
		panic("incorrect number of variables")
	}
	fmt.Println(x3) 
	fmt.Println() 

	fmt.Println("LINEAR")
	x4 := a.fetchr.CalculateDelta("linear", 0, -1)
	if len(x4) != 22 {
		panic("incorrect number of variables")
	}
	fmt.Println(x4) 
	fmt.Println()
	
}

/// TODO: run testing on values collected!
func Test_Analyzer_CaptureAllData(t *testing.T) {
	fmt.Println(TESTFP_1)
	a := OneAnalyzer_(TESTFP_1) 
	numPart := a.CaptureAllData("exact", "ends", 2, 3) 
	l := len(a.instanceInfo)

	if numPart != 1 {
		panic(fmt.Sprintf("incorrect # of partitions %d, want %d", numPart, 1))
	}

	if l != 91 - 2 -3 {
		panic(fmt.Sprintf("incorrect # of instance captures %d, want %d", l, 17 - 2 -3))
	}

	
	for _, c := range a.instanceInfo {
		c.Display([]string{"input", "control", "output"})
		fmt.Println("\n------------------------------")
	}
}

/// IN PROGRESS
/*
do node neighbor change, etc. 
*/ 
func Test_Analyzer_Judgment__Display(t *testing.T) {
	a := OneAnalyzer_(TESTFP_1) 
	a.CaptureAllData("full", "ends", 2, 3)
	a.SetThresher("one", NVMB_DELTA_NEG_CORR)
	a.GatherJudgmentValues() 

	fmt.Println("JUDGMENT DISPLAY")
	for _, c := range a.instanceInfo {
		c.Display([]string{"input", "control", "output"})
		c.DisplayJudgment([]string{"binary", "ternary", "mc"})
		fmt.Println("\n------------------------------")
	}
}

func Test_Analyzer_InstanceInfoToFile(t *testing.T) {
	a := OneAnalyzer_(TESTFP_1) 
	a.CaptureAllData("full", "ends", 2, 3)

	a.SetThresher("one", NVMB_DELTA_NEG_CORR)
	a.GatherJudgmentValues()	
	a.InstanceInfoToFile("./datos/testFile.csv") 
}