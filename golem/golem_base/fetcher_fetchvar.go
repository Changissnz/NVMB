package golem

import (
)

/*
*/
var NVMB_FETCHVAR_NAMES []string = []string{"cf gain", "contract gain", "number of contracts",
						"transmission", "cf acceptance rate", "competition",
						"greed", "negotiation", "growth", "nodeRank", "currency",
						"competitionMeasure", "contractMeasure", "newContractMeasure",
						"bondAdvantageMeasure", "bondDeletionAdvantageMeasure",
						"node connectivity", "num paths ratio", "neighbors", 
						"targets path risk ratio", "number of paths per target ratio"} 

/*
*/
func (f *Fetcher) FetchVar_CFGainAtTimestamp() *Variable {
	values := DataValueSliceToFloat64Slice(f.CollectValuesAtTimestamp("event_type", "cf tax", "impact"))
	output := SumIterable(values)
	return &Variable{varName: "cf gain", varValue: output}
}

/*
*/
func (f *Fetcher) FetchVar_ContractGain() *Variable{
	cav := DataValueSliceToFloat64Slice(f.CollectValuesAtTimestamp("event_type", "contract activity", "impact")) 
	cfv := DataValueSliceToFloat64Slice(f.CollectValuesAtTimestamp("event_type", "contract activity failure", "impact")) 

	cav = append(cav, cfv...) 
	output := SumIterable(cav) 
	return &Variable{varName: "contract gain", varValue: output} 
}

/*
*/
func (f *Fetcher) FetchVar_NumberOfContracts() *Variable {
	contracts := f.CollectValuesAtTimestamp("event_type", "summary", "contracts")[0] 
	iSlice,_ := DefaultStringToIntSlice(contracts.(string), ' ')
	size := len(iSlice)
	return &Variable{varName: "number of contracts", varValue: float64(size)}
}
 
/*
*/
func (f *Fetcher) FetchVar_Transmission() *Variable {
	transmission := f.CollectValuesAtTimestamp("event_type", "summary", "transmission")[0] 
	return &Variable{varName: "transmission", varValue: transmission.(float64)} 
}

/*
*/
func (f *Fetcher) FetchVar_CfAcceptanceRate() *Variable {
	cfAcceptanceRate := f.CollectValuesAtTimestamp("event_type", "summary", "cfAcceptanceRate")[0] 
	return &Variable{varName: "cf acceptance rate", varValue: cfAcceptanceRate.(float64)}
}

/*
*/
func (f *Fetcher) FetchVar_Competition() *Variable {
	competition := f.CollectValuesAtTimestamp("event_type", "summary", "competition")[0] 
	return &Variable{varName: "competition", varValue: competition.(float64)}
}

/*
*/
func (f *Fetcher) FetchVar_Greed() *Variable {
	greed := f.CollectValuesAtTimestamp("event_type", "summary", "greed")[0] 
	return &Variable{varName: "greed", varValue: greed.(float64)}
}

/*
*/
func (f *Fetcher) FetchVar_Negotiation() *Variable {
	negotiation := f.CollectValuesAtTimestamp("event_type", "summary", "negotiation")[0] 
	return &Variable{varName: "negotiation", varValue: negotiation.(float64)}
}

/*
*/
func (f *Fetcher) FetchVar_Growth() *Variable {
	growth := f.CollectValuesAtTimestamp("event_type", "summary", "growth")[0] 
	return &Variable{varName: "growth", varValue: growth.(float64)}
}

/*
*/
func (f *Fetcher) FetchVar_NodeRank() *Variable {
	nodeRank := f.CollectValuesAtTimestamp("event_type", "summary", "nodeRank")[0] 
	return &Variable{varName: "nodeRank", varValue: nodeRank.(float64)}
} 

/*
*/
func (f *Fetcher) FetchVar_Currency() *Variable {
	currency := f.CollectValuesAtTimestamp("event_type", "summary", "currency")[0].(float64) 
	return &Variable{varName: "currency", varValue: currency}
}

/*
*/
func (f *Fetcher) FetchVar_CompetitionMeasure() *Variable {
	competitionMeasure := f.CollectValuesAtTimestamp("event_type", "summary", "competitionMeasure")[0].(string)
	return &Variable{varName: "competitionMeasure", varValueS: competitionMeasure}
}

/*
*/
func (f *Fetcher) FetchVar_ContractMeasure() *Variable {
	contractMeasure := f.CollectValuesAtTimestamp("event_type", "summary", "contractMeasure")[0].(string) 
	return &Variable{varName: "contractMeasure", varValueS: contractMeasure}
}

/*
*/
func (f *Fetcher) FetchVar_NewContractMeasure() *Variable {
	newContractMeasure := f.CollectValuesAtTimestamp("event_type", "summary", "newContractMeasure")[0].(string) 
	return &Variable{varName: "newContractMeasure", varValueS: newContractMeasure}
}

/*
*/
func (f *Fetcher) FetchVar_BondAdvantageMeasure() *Variable {
	bondAdvantageMeasure := f.CollectValuesAtTimestamp("event_type", "summary", "bondAdvantageMeasure")[0].(string) 
	return &Variable{varName: "bondAdvantageMeasure", varValueS: bondAdvantageMeasure}
}

/*
*/
func (f *Fetcher) FetchVar_BondDeletionAdvantageMeasure() *Variable {
	bondDeletionAdvantageMeasure := f.CollectValuesAtTimestamp("event_type", "summary", "bondDeletionAdvantageMeasure")[0].(string) 
	return &Variable{varName: "bondDeletionAdvantageMeasure", varValueS: bondDeletionAdvantageMeasure}
}

/*
*/
func (f *Fetcher) FetchVar_NodeConnectivity() *Variable {
	neighbors := f.CollectValuesAtTimestamp("event_type", "summary", "neighbors")[0].(string)
	n, _ := DefaultStringToFloat64Slice(neighbors, ' ')

	knownNodes := f.CollectValuesAtTimestamp("event_type", "summary", "knownNodes")[0].(string)
	kn, _ := DefaultStringToFloat64Slice(knownNodes, ' ')
	nodeConn := ZeroDiv(float64(len(n)), float64(len(kn)), 0.0, 0.0)

	return &Variable{varName: "node connectivity", varValue: nodeConn}
}

/*
*/
func (f *Fetcher) FetchVar_NumPathsRatio() *Variable {
	ratio := f.CollectValuesAtTimestamp("event_type", "summary", "numPathsRatio")[0] 
	return &Variable{varName: "num paths ratio", varValue: ratio.(float64)}
}

/*
*/
func (f *Fetcher) FetchVar_Neighbors() *Variable {
	neighbors := f.CollectValuesAtTimestamp("event_type", "summary", "neighbors")[0] 
	return &Variable{varName: "neighbors", varValueS: neighbors.(string)} 
}

/*
*/
func (f *Fetcher) FetchVar__TargetsPath() (*Variable, *Variable) {
	pathInfoMap := f.PathInfoToMap()
	competitors, _ := DefaultStringToIntSlice(f.CollectValuesAtTimestamp("event_type", "summary", "competitors")[0].(string), ' ')

	var er, np float64
	var riskTotal float64 = 0.0
	var numPathRatioTotal float64 = 0.0

	for _, c := range competitors {
		pi, stat := pathInfoMap[c]
	
		if (!stat) {
			er = 0.0
			np = 0.0
		} else {
			er = pi[1]
			np = pi[0]
		}
		
		riskTotal += er
		numPathRatioTotal += np
	}

	pathRiskRatio := ZeroDiv(riskTotal, float64(len(competitors)), 0.0, 0.0) 
	numPathsRatio := ZeroDiv(numPathRatioTotal, float64(len(competitors)), 0.0, 0.0)

	v1 := &Variable{varName: "targets path risk ratio", varValue: pathRiskRatio}
	v2 := &Variable{varName: "number of paths per target ratio", varValue: numPathsRatio} 

	return v1, v2
}

/*
*/
func (f *Fetcher) PathInfoToMap() map[int][]float64 {
	pathInfo := f.CollectValuesAtTimestamp("event_type", "summary", "pathsInfo")[0].(string)
	q, _ := DefaultStringToFloat64Slice(pathInfo, ' ')
	output := make(map[int][]float64,0)

	for i := 0; i < len(q); {
		output[int(q[i])] = []float64{q[i + 1], q[i + 2]}  
		i += 3
	}
	return output 
}

