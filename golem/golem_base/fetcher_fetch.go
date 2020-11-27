/*
this file is where Fetcher actually fetches the variables.

NOTE: below methods are designed for safety over efficiency. 
*/
package golem

import (
	"fmt"
)

func (f *Fetcher) GatherVariablesAtTimestamp() []*Variable {

	output := make([]*Variable, 0)
	output = append(output, f.FetchVar_CFGainAtTimestamp())
	output = append(output, f.FetchVar_ContractGain()) 
	output = append(output, f.FetchVar_NumberOfContracts())
	output = append(output, f.FetchVar_Transmission())
	output = append(output, f.FetchVar_CfAcceptanceRate())
	output = append(output, f.FetchVar_Competition())
	output = append(output, f.FetchVar_Greed())
	output = append(output, f.FetchVar_Negotiation())
	output = append(output, f.FetchVar_Growth())
	output = append(output, f.FetchVar_NodeRank())
	output = append(output, f.FetchVar_Currency())

	output = append(output, f.FetchVar_CompetitionMeasure())
	output = append(output, f.FetchVar_ContractMeasure())
	output = append(output, f.FetchVar_NewContractMeasure())
	output = append(output, f.FetchVar_BondAdvantageMeasure())
	output = append(output, f.FetchVar_BondDeletionAdvantageMeasure())

	output = append(output, f.FetchVar_NodeConnectivity())
	output = append(output, f.FetchVar_NumPathsRatio())
	output = append(output, f.FetchVar_Neighbors())
	
	v1, v2 := f.FetchVar__TargetsPath()
	output = append(output, v1) 
	output = append(output, v2)
	return output
}

/*
gathers initial variables into structvar<analysis> 
*/ 
func (f *Fetcher) GatherVariablesInitial() int {
	var stat bool
	var timestampVariables []*Variable
	f.analysis = make(map[int][]*Variable,0)
	startTimestamp := 0 
	f.currentRowIndex = 0

	for {
		// get timestamp
		f.currentRowIndex, stat = f.FetchTimestampData(f.currentRowIndex)
		if (!stat) {
			break
		}

		timestampVariables = f.GatherVariablesAtTimestamp()
		f.analysis[startTimestamp] = timestampVariables
		startTimestamp++
	}
	
	return startTimestamp
}

/*
*/
func (f *Fetcher) FormatVariablesForDelta(start int, end int) map[string]interface{} {///[]float64 {

	if start > end && end != -1 {
		panic(fmt.Sprintf("invalid range: %d-%d", start,end))
	}

	if start < 0 {
		return nil 
	}

	output := f.MakeOutputForFormat() 
	if end == -1 {
		end = len(f.analysis)
	}

	for i := start; i < end; i++ {
		for _, v_ := range f.analysis[i] {
			switch {
				case StringIndexInSlice(NVMB_STRING_VARS, v_.varName) != -1:
					y := append(output[v_.varName].([]string), v_.varValueS)
					output[v_.varName] = y
				default: 
					y := append(output[v_.varName].([]float64), v_.varValue) 
					output[v_.varName] = y
			}
		}
	}

	return output 
}

func (f *Fetcher) MakeOutputForFormat() map[string]interface{} {
	output := make(map[string]interface{}, 0) 
	for _, v_ := range NVMB_FETCHVAR_NAMES {
		switch {

		case StringIndexInSlice(NVMB_STRING_VARS, v_) != -1: 
			output[v_] = []string{}
		default: 
			output[v_] = []float64{}
		}
	}
	return output
}

/*
*/
func (f *Fetcher) CalculateDelta(deltaOp string, start int, end int) []*Variable {
	m := f.FormatVariablesForDelta(start,end)

	keysToIterate := make([]string,0) 
	for k,_ := range m {
		keysToIterate = append(keysToIterate, k) 
	}

	keysToIterate = SortStringSliceAscending(keysToIterate)
	output := make([]*Variable, 0)

	for _, k := range keysToIterate {
		
		switch {
		case StringIndexInSlice(NVMB_SET_SUBTRACT_VARS, k) != -1:
			vn := fmt.Sprintf("%s sub__delta__%d_%d", k, start, end)
			vn2 := fmt.Sprintf("%s add__delta__%d_%d", k, start, end)
			sub, add := DeltaOpOnSetSequence(m[k].([]string))
			variable := &Variable{varName: vn, varValue: float64(sub)} 
			variable2 := &Variable{varName: vn2, varValue: float64(add)} 
			
			output = append(output, variable)  			
			output = append(output, variable2)

		/// string variable  
		case StringIndexInSlice(NVMB_STRING_VARS, k) != -1: 
			/// m[k]
			vn := fmt.Sprintf("%s__delta__%d_%d", k, start, end)
			s := DeltaOpOnDiscreteSequence(m[k].([]string))
			variable := &Variable{varName: vn, varValueS: s} 
			output = append(output, variable)

		default: 
			vn := fmt.Sprintf("%s__delta__%d_%d", k, start, end)
			///fmt.Println("DELTA ON ", k)
			variable := &Variable{varName: vn, varValue: DeltaOpOnSequence(m[k].([]float64), deltaOp), varValueS: deltaOp}
			output = append(output, variable)
		}
	} 

	return output 
}