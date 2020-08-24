package golem

/*
this file contains accuracy as well as time tests for
class ConcurrentFileBlockReader
*/

import (
	"fmt"
	"testing"
	//"time"
	//"reflect"
	//"strconv"
	//"encoding/csv"
	//"io"
)

func SetUpAnalyzer(fp string) *Analyzer {
	m := SetUpMatrix1(fp) 
	a := OneAnalyzer(m) 
	//a.Preprocess()
	return a
}

func Test_Analyzer_LoadOneRead(t *testing.T) {
	a := OneAnalyzer_(TESTFP_1) 
	CPARTSIZE = 3

	stat := a.LoadOneRead("exact",5) 
	if !stat {
		panic("[0] incorrect stat") 
	}
	
	fmt.Println("number of blocks:\t", len(a.readr.blockData), stat) 
	
	if len(a.readr.blockData) != 10 {
		panic("[0] incorrect number of blocks")
	}

	stat = a.LoadOneRead("exact", 5) 
	if stat {
		panic("[1] incorrect stat") 
	} 
	
	stat = a.LoadOneRead("exact", 5) 
	if stat {
		panic("[1] incorrect stat") 
	}

	CPARTSIZE = 1000
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
		
		if c == 8 && a.fetchr.data.Length() != 3 {
			panic("block 8 incorrect length")
		} 
		
		if c != 8 && a.fetchr.data.Length() != 5 {
			panic(fmt.Sprintf("block %d incorrect length", c))
		}
		
	}

	if c != numBlocks {
		panic(fmt.Sprintf("incorrect number of reads %d, want 8", c)) 
	}

	CPARTSIZE = 1000 
}


/*
displays all delta calculations delta calculations 
*/ 
func Test_Analyzer_DeltaDisplay(t *testing.T) {
	a := SetUpAnalyzer(TESTFP_1) 
	a.fetchr.GatherVariablesInitial() 

	fmt.Println("END")
	fmt.Println(a.fetchr.CalculateDelta("ends", 0, -1)) 
	fmt.Println() 

	fmt.Println("MEAN TO END")
	fmt.Println(a.fetchr.CalculateDelta("mean-to-end", 0, -1)) 
	fmt.Println() 

	fmt.Println("MEAN TO START")
	fmt.Println(a.fetchr.CalculateDelta("mean-to-start", 0, -1)) 
	fmt.Println() 

	fmt.Println("LINEAR")
	fmt.Println(a.fetchr.CalculateDelta("linear", 0, -1)) 
	fmt.Println()
}


func Test_Analyzer_CollectData(t *testing.T) {

	a := SetUpAnalyzer(TESTFP_1)
	x := a.CollectData("ends", 0, -1)
	DisplayVarSlice(x)
	fmt.Println("----------------------------------------------")
	x2 := a.CollectData("ends", 0, 5) 
	DisplayVarSlice(x2)

	fmt.Println("----------------------------------------------")
	for k, v := range a.fetchr.analysis {
		fmt.Println(k) 
		DisplayVarSlice(v) 
		fmt.Println() 
	}
}

func Test_Analyzer_CaptureAllData(t *testing.T) {
	a := OneAnalyzer_(TESTFP_1) 
	numPart := a.CaptureAllData("exact", "ends", 2, 3) 
	l := len(a.instanceInfo)

	if numPart != 1 {
		panic(fmt.Sprintf("incorrect # of partitions %d, want %d", numPart, 1))
	}

	if l != 17 - 2 -3 {
		panic(fmt.Sprintf("incorrect # of instance captures %d, want %d", l, 17 - 2 -3))
	}
}

func Test_Analyzer_CollectDataOnePartition(t *testing.T) {
	a := OneAnalyzer_(TESTFP_1)
	a.priorH = 2
	a.postH = 2

	stat := a.CollectDataOnePartition("full","ends")
	if (!stat) {
		panic("invalid collect data one partition")
	}

	if len(a.instanceInfo) != 17 - 2 - 2 {
		panic("incorrect length for instance capture")
	}
} 

/// IN PROGRESS 
func Test_Analyzer_Judgment__Display(t *testing.T) {
	a := OneAnalyzer_(TESTFP_1) 
	a.CaptureAllData("full", "ends", 2, 3) 
	a.SetThresher("one", NVMB_DELTA_NEG_CORR)
	a.GatherJudgmentValues() 

	for _, c := range a.instanceInfo {
		c.DisplayJudgmentFull() 
		fmt.Println("#--------------------#")
	}

	q := a.instanceInfo[0]

	for _, q_ := range q.inputVariables {
		ds := ParseDeltaString(q_.varName) 
		fmt.Println("delta string: ", ds) 
	}

	/*
	a.StartThresher()
	for i := 0; i < 10; i++ {
		a.RunThresher(0.5, 0.1)
		fmt.Println("positive: ", a.thresher.positive) 
	}
	*/

	// iterate through and display judgment
	/*
	for _, c := range a.instanceInfo {
		c.DisplayJudgment() 
		fmt.Println("#--------------------#")
	}
	*/
} 