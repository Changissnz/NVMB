/*
This file performs instance captures on data
*/ 
package golem 

import (
	"fmt" 
)

type Ternary struct {
	value int
}

func OneTernary(i int) *Ternary {

	if i != -1 && i != 0 && i != 1 {
		panic("invalid ternary value")
	}
	return &Ternary{value: i}

}


type InstanceCapture struct {
	inputVariables []*Variable
	controlVariables []*Variable 
	outputVariables []*Variable
	captureType string // TODO: only NVMB type supported for now.

	judgmentValues []float64
	judgment bool

	judgment3 *Ternary
	judgmentMC string
}

func OneInstanceCapture(ct string) *InstanceCapture {
	if ct != "NVMB" {
		panic("only NVMB type supported!") 
	}
	return &InstanceCapture{captureType: ct}
}

func (ic *InstanceCapture) Display(display []string) {
	fmt.Println("\tINPUT VARIABLES ", len(ic.inputVariables))
	if StringIndexInSlice(display, "input") != -1 {
		DisplayVarSlice(ic.inputVariables) 
	}

	fmt.Println("/==/==/==/") 
	fmt.Println("\tCONTROL VARIABLES", len(ic.controlVariables))
	if StringIndexInSlice(display, "control") != -1 {
		DisplayVarSlice(ic.controlVariables) 
	}	

	fmt.Println("/==/==/==/") 
	fmt.Println("\tOUTPUT VARIABLES", len(ic.outputVariables))
	if StringIndexInSlice(display, "output") != -1 {
		DisplayVarSlice(ic.outputVariables) 
	}
}

/*
*/ 
func (ic *InstanceCapture) DisplayJudgment(judgmentType []string) {
	
	if StringIndexInSlice(judgmentType, "variable") != -1 {
		fmt.Println("\tVAR. VAL.")
		for i, v := range ic.outputVariables {
			name, start, end := ParseDeltaString(v.varName) 
			fmt.Println("variable: ", name)
			_,startp,endp := ParseDeltaString(ic.inputVariables[i].varName)
			fmt.Println("start prior ", startp, " end prior ", endp, " = ", ic.inputVariables[i].varValue) 
			fmt.Println("start post ", start, " end post ", end, " = ", v.varValue) 
			fmt.Println("judgment ", ic.judgmentValues[i]) 		
			fmt.Println("--------------------------------")
		}
	}

	if StringIndexInSlice(judgmentType, "binary") != -1 {
		fmt.Println("\tBINARY")
		fmt.Println("* judgment: ", ic.judgment) 
	}

	if StringIndexInSlice(judgmentType, "ternary") != -1 {
		fmt.Println("\tTERNARY")
		fmt.Println("* judgment: ", ic.judgment3)//.value) 
	}

	if StringIndexInSlice(judgmentType, "mc") != -1 {
		fmt.Println("\tMULTI-CLASS")
		fmt.Println("* judgment: ", ic.judgmentMC)  
	}

	fmt.Println("\n=$===$===$===$===\n=$===$===$===$===")
}

/*
*/
func (ic *InstanceCapture) CaptureByType(v []*Variable, varType string) {
	var criteria []string

	switch {
	case varType == "input": 
		ic.inputVariables = make([]*Variable, 0)
		criteria = NVMB_INPUT_OUTPUT_VARS
	case varType == "control": 
		ic.controlVariables = make([]*Variable, 0)
		criteria = NVMB_CONTROL_VARS
	case varType == "output": 
		ic.outputVariables = make([]*Variable, 0)
		criteria = NVMB_INPUT_OUTPUT_VARS
	default: 
		panic("invalid variable type")
	}

	if ic.captureType != "NVMB" {
		panic("invalid capture type")
	}

	for _, v_ := range v {
		vs,_,_ := ParseDeltaString(v_.varName)
		
		if varType == "control" {
			fmt.Println("DELTA STRING ", vs)
		}
		
		if StringIndexInSlice(criteria, vs) != -1 {
			switch {
			case varType == "input": 
				ic.inputVariables = append(ic.inputVariables, v_) 
			
			case varType == "output": 
				ic.outputVariables = append(ic.outputVariables, v_) 

			default: 
				ic.controlVariables = append(ic.controlVariables, v_) 
			}
		}
	}
}

/// TODO:
func (ic *InstanceCapture) CompareInstanceCapture(ic2 *InstanceCapture) float64{
	return -1
} 

/*
return:
- before timestamp range
- after timestamp range
*/
func ParseDeltaString(ds string) (string, string, string) { 

	substr := "delta" 
	l := len(substr) 
	indices := SubstringInStringAll(ds, substr)
	
	if len(indices) == 0 {
		return ds,"",""
	}

	if len(indices) > 1 {
		panic("invalid delta string") 
	}

	ds_ := ds[:indices[0] + l]
	x := ds[indices[0] + l:]

	q := SubstringInStringAll(x, "__") 

	if len(indices) == 0 {
		panic("not valid timestamp range")
	}

	x = x[q[0] + len("__"):] 
	
	q2 := SubstringInStringAll(x, "_") 

	if len(indices) == 0 {
		panic("not valid timestamp range")
	}

	start,end := x[:q2[0]], x[q2[0] + 1:]
	return ds_, start, end 
}

////////////////////// NVMB input/output/control variables 

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
						"cf gain__delta",
						"neighbors add__delta", 
						"neighbors sub__delta", 
						/// add control variables
						"competition__delta",
						"greed__delta",
						"negotiation__delta",
						"growth__delta",
						/// TODO: work on now. 
						"competitionMeasure__delta",
						"contractMeasure__delta",
						"newContractMeasure__delta",
						"bondAdvantageMeasure__delta",
						"bondDeletionAdvantageMeasure__delta",
					}

var NVMB_STRING_VARS []string = []string{"neighbors", "competitionMeasure", "contractMeasure", 
	"newContractMeasure", "bondAdvantageMeasure", "bondDeletionAdvantageMeasure"}

var NVMB_SET_SUBTRACT_VARS []string = []string{"neighbors"} 

var NVMB_CONTROL_VARS []string = []string{"competition", "greed", "negotiation", "growth", 
	"competitionMeasure", "contractMeasure", "newContractMeasure", "bondAdvantageMeasure",
	"bondDeletionAdvantageMeasure"}

var NVMB_DELTA_NEG_CORR = []string{}