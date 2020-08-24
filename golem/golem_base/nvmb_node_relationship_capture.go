/*
This file performs instance captures on data
*/ 

/// TODO: first, have to perfect file reader. 

package golem 

import (
	"fmt" 
)

type InstanceCapture struct {
	inputVariables []*Variable
	controlVariables []*Variable 
	outputVariables []*Variable
	captureType string // TODO: only NVMB type supported for now.

	judgmentValues []float64
	judgment bool
}

func OneInstanceCapture(ct string) *InstanceCapture {
	if ct != "NVMB" {
		panic("only NVMB type supported!") 
	}
	return &InstanceCapture{captureType: ct}
}

func (ic *InstanceCapture) DisplayJudgmentFull() { 
	for i, v := range ic.outputVariables {
		fmt.Println("variable: ", v.varName) 
		fmt.Println("\tjudgment value: ", ic.judgmentValues[i]) 
		fmt.Println() 
	}
	fmt.Println("* judgment: ", ic.judgment) 
}

func (ic *InstanceCapture) DisplayJudgment() {
	fmt.Println(ic.judgmentValues)
	fmt.Println(ic.judgment)
	fmt.Println() 
}

func (ic *InstanceCapture) CaptureInput(v []*Variable) {
	ic.inputVariables = make([]*Variable, 0)
	if ic.captureType == "NVMB" {
		for _, v_ := range v {
			vs := ParseDeltaString(v_.varName)
			if StringIndexInSlice(NVMB_INPUT_OUTPUT_VARS, vs) != -1 {
				ic.inputVariables = append(ic.inputVariables, v_) 
			}
		}
	} else {
		panic("invalid capture type")
	}
}

func (ic *InstanceCapture) CaptureControl(v []*Variable) {
	ic.controlVariables = make([]*Variable, 0)

	if ic.captureType == "NVMB" {
		for _, v_ := range v {
			vs := ParseDeltaString(v_.varName)
			if StringIndexInSlice(NVMB_CONTROL_VARS, vs) != -1 {
				ic.controlVariables = append(ic.controlVariables, v_) 
			}
		}
	} else {
		panic("invalid capture type")
	}
}

func (ic *InstanceCapture) CaptureOutput(v []*Variable) {
	ic.outputVariables = make([]*Variable, 0)
	if ic.captureType == "NVMB" {
		for _, v_ := range v {
			vs := ParseDeltaString(v_.varName)
			if StringIndexInSlice(NVMB_INPUT_OUTPUT_VARS, vs) != -1 {
				ic.outputVariables = append(ic.outputVariables, v_) 
			}
		}
	} else {
		panic("invalid capture type")
	}
}

func (ic *InstanceCapture) CompareInstanceCapture(ic2 *InstanceCapture) float64{
	return -1
} 


func ParseDeltaString(ds string) string { 
	substr := "delta" 
	l := len(substr) 
	indices := SubstringInStringAll(ds, substr)
	
	if len(indices) == 0 {
		return ds
	}

	if len(indices) > 1 {
		panic("invalid delta string") 
	}

	ds_ := ds[:indices[0] + l] 
	return ds_ 
}

////////////////////// NVMB input/output/control variables 

// observational notes 
/*
How input/output variables should be viewed: 
- nodeRank : ^ -> good 
- transmission: ^ -> good 
- currency: ^ -> good 
- cf_acceptance_rate: ^ -> g
- negotiation_delta: ^ -> g 
- num paths ratio_delta: ^ -> g 
- target path risk ratio: ^ -> b   
- number of contracts: ^ -> g 
- contract gain: ^ -> g
- node connectivity: ^ -> g? 
- number of paths per target: ^ -> g 
- cf gain: ^ -> g 
*/
var NVMB_INPUT_OUTPUT_VARS []string = []string{"nodeRank__delta",
						"transmission__delta",
						"currency__delta",
						"cf_acceptance_rate__delta",
						"negotiation__delta",
						"num paths ratio__delta",
						"targets path risk ratio__delta",
						"number of contracts__delta",
						"contract gain__delta",
						"node connectivity__delta",
						"number of paths per target ratio__delta",
						"cf gain__delta"}


var NVMB_CONTROL_VARS []string = []string{"competition", "greed", "negotiation", "growth"}

// OF IMPORTANCE
/*
"competition__delta",
"greed__delta",
"growth__delta",
"negotiation__delta"
*/
var NVMB_DELTA_NEG_CORR = []string{}
//////////////////////////////////////////////////
/*
control keys:
	- competition 
	- greed 
	- negotiation 
	- growth 

input/output keys: 
	- nodeRank__delta
	- transmission__delta
	- currency__delta
	- cf_acceptance_rate__delta
	- negotiation__delta
	- num paths ratio__delta
	- targets path risk ratio__delta
	- number of contracts__delta
	- contract gain__delta
	- node connectivity__delta
	- number of paths per target ratio__delta
	- cf gain__delta

-----------------------------------------------
input keys: 
	- competition__delta
	- greed__delta
	- growth__delta
	- negotiation__delta
*/