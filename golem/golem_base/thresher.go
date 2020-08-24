package golem 

import (
	"fmt"
)

type Thresher struct {
	thresholdType string

	thresholdTypeOne float64 
	thresholdTypeTwo []float64

	/// TODO: this variable will need to be more closely inspected.
	directionality []string // variables with negative correlation
	positive float64 
	numRefinement int
}

func OneThresher() *Thresher {
	return &Thresher{thresholdType: "", numRefinement: 0, positive: 0}
}

/// TODO: below needs work. 
/*
placementType := random|zero
dataType := boolean|continuous
*/ 
func (t *Thresher) SetThreshold(numVars int, placementType string, dataType string) {

	if placementType == "zero" {
		if dataType == "boolean" {
			switch {
			case t.thresholdType == "one": 
				t.thresholdTypeOne = float64(0)
			case t.thresholdType == "two": 
				t.thresholdTypeTwo = make([]float64,0) 
				for i := 0; i < numVars; i++ {
					t.thresholdTypeTwo[i] = 0 
				}
			default: 
				panic("invalid threshold type")
			}
			return
		}
		panic(fmt.Sprintf("data type %s not yet coded", dataType ))
	}

	panic(fmt.Sprintf("placement type %s not yet coded", placementType))
}