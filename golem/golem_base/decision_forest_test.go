package golem

import (
	"testing"
	"time"
	//"math"
	//"math/rand"
	"fmt"
)

/*
TODO: work on this!
*/ 
func Test_DForest__DTreeMulticlass(t *testing.T) {
	CPARTSIZE = 50000
	fp := "./datos/other_data/glass.csv"
	///fp := "./datos/other_data/creditcard.csv"
	cfbp := OneConcurrentFileBlockReader(fp)
	s := OneSampler(cfbp, 0.7)
	df := OneDForest(s, 13, []string{}, []string{"Type"})
	//df := OneDForest(s, 13, []string{"Time"}, []string{"Class"})

	start := time.Now()
	df.InitializeForest("in-order")
	duration := time.Since(start)

	fmt.Println("[1] duration ", duration)	
	df.MakeForest(3, EntropyCost)
	pm := df.PredictMC() 
	fmt.Println("DISPLAYING RESULTS")
	pm.Display()
}