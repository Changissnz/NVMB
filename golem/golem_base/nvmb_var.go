/*
this file contains relevant methods and vars. for NVMB node data analysis.
*/ 

package golem 

import (
	"fmt"
	"strconv"
)

var NVMB_NODEDATA_TYPEMAP = map[string]string{
	"time": "int", 
	"event_type": "string", // 
	"impact": "float",
	"role": "string",
	"nodes": "vector",
	"responses": "vector", 
	"currency": "float", 
	"neighbors": "vector",
	//"neighborsRisk": "float", // tbc
	"contracts": "vector",
	"transmission": "float",
	"cfAcceptanceRate": "float",
	"competition": "float",
	"greed": "float",
	"negotiation": "float",
	"growth": "float",
	"nodeRank": "float",
	"knownNodes": "vector",
	"competitors": "vector",
	"numPathsRatio": "float",
	"pathsInfo": "vector",
	"competitionMeasure": "string",
	"contractMeasure": "string",
	"newContractMeasure": "string",
	"bondAdvantageMeasure": "string",
	"bondDeletionAdvantageMeasure": "string"}  

/////////////////////////////////////////////////

type InputOutputDataVariable struct {
	input []float64
	output []float64
}

func SliceToInputOutputDataVariable(data []float64, inputIndices *BasicSet, outputIndices *BasicSet) *InputOutputDataVariable {
	x := InputOutputDataVariable{} 

	for i, d := range data {
		q := strconv.Itoa(i)

		switch {
		
		case inputIndices.DoesExist(q): 
			x.input = append(x.input, d) 
			
		case outputIndices.DoesExist(q): 
			x.output = append(x.output, d) 

		default: 
			panic(fmt.Sprintf("index %d is not input or output", i)) 
		}
	}

	return &x 
}

type DataVariable []float64 

type DataVariableSlice struct {
	data [][]float64
}

func (dvs DataVariableSlice) Len() int { return len(dvs.data)} 

func (dvs DataVariableSlice) Values(i int) []float64 {
	return dvs.data[i] 
}

////////////////////////////////////

/// TODO: variables need to be observed.
type Variable struct {
	varName string
	varValue float64 // TODO: maybe DataValue is better?
	///varValue interface{}
	varValueS string
	varValue2 float64 // can be used as loss
	associatedValues []float64
}

func (v *Variable) MakeCopy() *Variable {
	v2 := &Variable{varName: v.varName, varValue: v.varValue,
		varValueS: v.varValueS, varValue2: v.varValue2,
		associatedValues: v.associatedValues} 
	return v2 
}

func DisplayVarSlice(v []*Variable) {
	for _, v_ := range v {
		fmt.Println(fmt.Sprintf("name: %s\tvalue: %f\tvalueS: %s", v_.varName, v_.varValue, v_.varValueS))
	}
}