package golem

import (
	"github.com/biogo/cluster/kmeans"
	"fmt"
)

type Clumper struct {
	dataCache []*DataVariable
	data DataVariableSlice
}

func OneClumper(x []*InstanceCapture) *Clumper { 
	c := &Clumper{} 
	c.LoadData(x) 
	return c
}

/// TODO: code missing from previous function 
func (c *Clumper) LoadData(x []*InstanceCapture) {
	q := make([][]float64,0) 

	for i,x_ := range x {
		fmt.Println("YES ", i)
		dv := InstanceInfoToDataVariable(x_) 
		q = append(q, dv) 
	}

	c.data = DataVariableSlice{data: q}
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
	centers := km.Centers()

	for _, c := range centers {
		fmt.Println("center is ", c)
	}


	fmt.Println("NUM VALUES : ", len(values))
	for _, v := range values {
		q := fmt.Sprintf("cluster of %v is %v", v.V(), v.Cluster()) 
		fmt.Println(q) 
		fmt.Println() 
	}
}


func InstanceInfoToDataVariable(instInf *InstanceCapture) []float64 {
	dv := make([]float64, 0)

	for _, x := range instInf.judgmentValues {
		dv = append(dv, x) 
	}
	return dv 
}