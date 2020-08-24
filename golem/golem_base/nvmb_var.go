/*
this file contains relevant methods and vars. for NVMB node data analysis.
*/ 

package golem 

import (
	"fmt"
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
	"pathsInfo": "vector"} 

/// TODO: do this below
type DataVariable struct {
	cfGain float64
	contractGain float64 
	currency float64 
	negotiation float64 
	nodeConnectivity float64 
	nodeRank float64 
	numPathsRatio float64 
	numberOfContracts float64 
	numberOfPathsPerTargetRatio float64 
	targetsPathRiskRatio float64 
	transmission float64 
}

type DataVariableSlice struct {
	data []*DataVariable
}

func (dv *DataVariable) Len() int { return 12 }

func (dvs *DataVariableSlice) Len() int { return len(dvs.data)} 

func (dvs *DataVariableSlice) Values(i int) []float64 {

	output := make([]float64,0) 
	dv := dvs.data[i]

	output = append(output, dv.cfGain) 
	output = append(output, dv.contractGain) 
	output = append(output, dv.currency) 
	output = append(output, dv.negotiation) 
	output = append(output, dv.nodeConnectivity) 
	output = append(output, dv.nodeRank) 
	output = append(output, dv.numPathsRatio)  
	output = append(output, dv.numberOfContracts) 
	output = append(output, dv.numberOfPathsPerTargetRatio) 
	output = append(output, dv.targetsPathRiskRatio) 
	output = append(output, dv.transmission)  

	return output
}

type Variable struct {
	varName string
	varValue float64 // TODO: maybe DataValue is better?
	spareValue string
	varValue2 float64 // can be used as loss
	associatedValues []float64
}

func (v *Variable) MakeCopy() *Variable {
	v2 := &Variable{varName: v.varName, varValue: v.varValue,
		spareValue: v.spareValue, varValue2: v.varValue2,
		associatedValues: v.associatedValues} 
	return v2 
}

func DisplayVarSlice(v []*Variable) {
	for _, v_ := range v {
		fmt.Println(fmt.Sprintf("%s %f", v_.varName, v_.varValue))
	}
}