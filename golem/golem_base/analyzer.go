package golem

import (
	"fmt"
	//"math"
)

type Analyzer struct {
	fetchr *Fetcher
	readr *ConcurrentFileBlockReader
	readrBlockIndex int

	leftovers *CFBRDataMatrix // use for inter-partition processing.
	priorH int 
	postH int 

	instanceInfo []*InstanceCapture

	thresher *Thresher /// TODO: this needs to be implemented 
	clumpr *Clumper
}

/// TODO: delete below constructor 
func OneAnalyzer(d *CFBRDataMatrix) *Analyzer {
	d.Preprocess() 
	f := OneFetcher(d)
	t := OneThresher() 
	return &Analyzer{fetchr: f, thresher: t, readrBlockIndex: 0} 
}

// TODO: this constructor will replace above
func OneAnalyzer_(fp string) *Analyzer {
	r := OneConcurrentFileBlockReader(fp) 
	r.ManualColumnTypeSet(NVMB_NODEDATA_TYPEMAP)
	t := OneThresher() 
	a := &Analyzer{readr: r, readrBlockIndex: 0, thresher: t} 
	return a 
}

/// TODO: this needs to be error checked. 
/*
loads one partition from reader onto fetchr. If no data is left, 
outputs false.  
*/ 
func (a *Analyzer) LoadOneRead(readType string, postHop int) bool {

	// read partition 
	x := a.readr.ReadPartition(readType) 
	if x == 0 {
		return false 
	}

	// read partition into matrix and add leftovers from last partition 
	matrix := a.readr.ConvertPartitionBlockToMatrix(0)
	
	if a.leftovers != nil {
		a.leftovers.StackMatrix(matrix)
		matrix = a.leftovers 
		a.leftovers = nil  
	}

	a.readrBlockIndex = 0 
	matrix.Preprocess()
	a.fetchr = OneFetcher(matrix) 
	return true 
}

/*
slides partition block forward to account for `postHop`. 
Outputs true if done with 
*/ 
func (a *Analyzer) SlideFetcher(postHop int) bool {
	var r int 
	// get leftovers
	
	a.leftovers = nil 
	if a.fetchr != nil {
		r,_ = a.fetchr.data.stringData.Dims()
		if postHop >= r {
			fmt.Println(fmt.Sprintf("hop length %d must be less than data length %d", postHop, r)) 
			return true
		}
		a.leftovers = a.fetchr.data.IndexRange(r - postHop, r) 
	} else {
		panic("cannot slide data with starting null data")
	}

	// end of partition read
	if len(a.readr.blockData) <= a.readrBlockIndex {
		return true 
	}

	// new block
	matrix := a.readr.ConvertPartitionBlockToMatrix(a.readrBlockIndex)
	matrix.Preprocess()

	// get leftovers 
	a.leftovers.StackMatrix(matrix)
	matrix = nil 
	matrix = a.leftovers 
	a.leftovers = nil 

	// make new fetcher 
	a.fetchr = OneFetcher(matrix) 
	a.readrBlockIndex++ 
	return false 
}

////////// `main` method: process dataset. 
/*
*/
func (a *Analyzer) CaptureAllData(readType string, deltaOp string, priorHop int, postHop int) int {
	a.priorH = priorHop 
	a.postH = postHop
	a.instanceInfo = make([]*InstanceCapture, 0) 
	c := 0 

	for {
		if !a.CollectDataOnePartition(readType, deltaOp) {
			break 
		}
		c++ 
	}
	return c 
}

/*
*/ 
func (a *Analyzer) CollectDataOnePartition(readType string, deltaOp string) bool {

	if !a.LoadOneRead(readType, a.postH) {
		return false 
	}

	// format variables into analysis
	a.fetchr.GatherVariablesInitial() 

	// capture data
	c := 0 
	for {
		a.CollectBlockData(deltaOp) 
		c++ 

		if !a.SlideFetcher(a.postH) {
			break 
		}
	}

	return true
} 

/*
*/
func (a *Analyzer) CollectBlockData(deltaOp string) {
	l := len(a.fetchr.analysis) - a.postH
	for i := a.priorH; i < l; i++ {
		ic := a.CaptureAtTimestamp(i, deltaOp, a.priorH, a.postH) 
		a.instanceInfo = append(a.instanceInfo, ic) 
	}
}

/*
collects delta info given a range.
*/ 
func (a *Analyzer) CollectData(deltaOp string, start int, end int) []*Variable {
	a.fetchr.GatherVariablesInitial()
	deltaV := a.fetchr.CalculateDelta(deltaOp, start, end)  	
	return deltaV
}

/*
*/
func (a *Analyzer) FormatCapture_PriorAtAfter(bef []*Variable, at []*Variable, aft []*Variable) *InstanceCapture {
	ic := OneInstanceCapture("NVMB")
	ic.CaptureInput(bef) 
	ic.CaptureControl(at) 
	ic.CaptureOutput(aft)
	return ic 
}

/*
captures data at timestamp and outputs an InstanceCapture. 
*/ 
func (a *Analyzer) CaptureAtTimestamp(timestamp int, deltaOp string, priorHop int, postHop int) *InstanceCapture {///([]*Variable, []*Variable, []*Variable) {
	beforeRangeStart, beforeRangeEnd := timestamp - priorHop, timestamp 
	afterRangeStart, afterRangeEnd :=  timestamp + 1, timestamp + postHop + 1 
	varBef := a.CollectData(deltaOp, beforeRangeStart, beforeRangeEnd) 
	varAft := a.CollectData(deltaOp, afterRangeStart, afterRangeEnd)
	return a.FormatCapture_PriorAtAfter(varBef, a.fetchr.analysis[timestamp], varAft)
}

/*
description:
criteria 1 takes a look at the following values: 
	OUTPUT_VARS - INPUT_VARS
*/
func (a *Analyzer) ObtainDifferencePrePostOp(is *InstanceCapture, negVars []string) []float64 {

	// obtain difference b/t pre and post
	diff := make([]float64,0) 

	for i,v := range is.inputVariables {
		v_ := is.outputVariables[i]

		vn := ParseDeltaString(v_.varName) 
		if vn != ParseDeltaString(v.varName) {
			panic("ERROR")
		}

		mul := float64(1) 
		if StringIndexInSlice(negVars, vn) != -1 {
			mul = -1
		}
		diff = append(diff, mul * (v_.varValue - v.varValue))  
	}
	
	return diff
}

func (a *Analyzer) JudgeInstanceCapture_Criteria1(diff []float64, threshold float64) bool { 

	c := 0 
	for _, d := range diff {
		if d > threshold {
			c++ 
		}
	}

	r := ZeroDiv(float64(c), float64(len(diff)), 0.0, 0.0)
	if r >= 0.5 {
		return true 
	}
	return false 

}

func (a *Analyzer) SetThresher(criteria string, negVars []string) {
	a.thresher.thresholdType = criteria
	a.thresher.directionality = negVars
} 

func (a *Analyzer) GatherJudgmentValues() {
	if a.thresher.thresholdType == "" {
		panic("threshold type must be set!")
	}

	switch {
		case a.thresher.thresholdType == "one": 
			for _, c := range a.instanceInfo {
				// calculate difference
				c.judgmentValues = a.ObtainDifferencePrePostOp(c, a.thresher.directionality)
			}

		default: 
			panic(fmt.Sprintf("criteria %s does not exist", a.thresher.thresholdType)) 
	}
}

/*
*/
func (a *Analyzer) ConvertInstanceInfoToClumperFormat() {

}

func (a *Analyzer) StartThresher() { 
	// get number of variables 
	if len(a.instanceInfo) == 0 {
		return 
	}

	a.thresher.SetThreshold(len(a.instanceInfo[0].inputVariables), "zero", "boolean") 
	return 
}

func (a *Analyzer) RunThresher(limit float64, increment float64) {
	
	// TODO: code this below 
	posCount := 0 
	for _, c  := range a.instanceInfo {
		c.judgment = a.JudgeInstanceCapture_Criteria1(c.judgmentValues, a.thresher.thresholdTypeOne + increment)
		posCount += 1
	}

	a.thresher.thresholdTypeOne += increment 
	a.thresher.positive = ZeroDiv(float64(posCount), float64(len(a.instanceInfo)), 0.0, 0.0)
} 

