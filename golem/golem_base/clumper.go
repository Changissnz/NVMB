package golem

import (
	"github.com/biogo/cluster/kmeans"
	"fmt"
)


type Clumper struct {
	dataCache []*DataVariable
	data *DataVariableSlice
}

func OneClumper(x []*InstanceCapture) *Clumper { 
	c := &Clumper{} 
	c.LoadData(x) 
	return c
}

func (c *Clumper) LoadData(x []*InstanceCapture) {
	c.dataCache = make([]*DataVariable, 0)
	for _, x_ := range x {
		dv := InstanceInfoToDataVariable(x_)
		///fmt.Println("got data variable for ", i, "  ", dv)
		c.dataCache = append(c.dataCache, dv) 
	}

	c.data = &DataVariableSlice{data: c.dataCache}
	c.dataCache = nil
} 

/// TODO: numIterations needs to be coded. 
func (c *Clumper) KMeansClump(numCenters int, numIterations int) {
	km, err := kmeans.New(c.data) 

	fmt.Println("JAH JAH")
	if err != nil {
		panic("could not initialize k-means algorithm.") 
	}

	/////
	km.Seed(numCenters)

	values := km.Values()
	///centers := km.Centers()

	fmt.Println("NUM VALUES : ", len(values))
	for _, v := range values {
		q := fmt.Sprintf("cluster of %v is %v", v.V(), v.Cluster()) 
		fmt.Println(q) 
		fmt.Println() 
	}
}

func InstanceInfoToDataVariable(instInf *InstanceCapture) *DataVariable {
	dv := &DataVariable{} 

	for i, ic := range instInf.inputVariables {
		vn := ParseDeltaString(ic.varName) 
		SetDataVariableByValue(dv, vn, instInf.judgmentValues[i]) 
	}
	return dv 
} 

func SetDataVariableByValue(dv *DataVariable, varName string, value float64) {

	switch {

		case (varName == "cf gain__delta"):
			dv.cfGain = value

		case (varName == "contract gain__delta"): 
			dv.contractGain = value

		case (varName == "currency__delta"): 
			dv.currency = value

		case (varName == "negotiation__delta"): 
			dv.negotiation = value

		case (varName == "node connectivity__delta"): 
			dv.nodeConnectivity = value

		case (varName == "nodeRank__delta"): 
			dv.nodeRank = value

		case (varName == "num paths ratio__delta"): 
			dv.numPathsRatio = value

		case (varName == "number of contracts__delta"): 
			dv.numberOfContracts = value

		case (varName == "number of paths per target ratio__delta"): 
			dv.numberOfPathsPerTargetRatio = value

		case (varName == "targets path risk ratio__delta"): 
			dv.targetsPathRiskRatio = value

		case (varName == "transmission__delta"): 
			dv.transmission = value
	}
} 
