package golem

import (
	"gonum.org/v1/gonum/mat"
	"fmt"
	"math/rand"
	"sync"
)

var MAX_DEPTH_FOR_TREE int = 18

type DForest struct {
	dts []*DTree
	size int
	excludeCols []string
	outputCols []string
	keyOrder string
	samplr *Sampler 
}

func OneDForest(s *Sampler, numTrees int, excludeCols []string, outputCols []string) *DForest {
	
	if numTrees < 1 {
		panic("forest cannot be empty")
	}
	
	df := &DForest{samplr: s, size: numTrees, excludeCols: excludeCols, outputCols: outputCols}
	return df
}

/// TODO: check to make sure data matrix for each tree is the same. 
/*
*/
func (df *DForest) InitializeForest(keyOrder string) {

	// make the first tree
	df.keyOrder = keyOrder 
	dt := OneDTree(df.samplr, df.excludeCols, df.outputCols)
	dt.Preprocess(df.keyOrder) 
	df.dts = append(df.dts,dt)

	// make remaining trees
	for i := 1; i < df.size; i++ {
		dt2 := OneDTree(df.samplr, df.excludeCols, df.outputCols)
		dt2.data = dt.data
		dt2.classOrder = dt.classOrder
		dt2.AssignOutcomeClassesToBDT()
		df.dts = append(df.dts, dt2) 
	}
}

/*
*/ 
func (df *DForest) MakeForest(maxDepth int, f func(float64) float64) { 
	for i, dt := range df.dts {
		// update sampling data and make tree
		fmt.Println("making tree ", i)
		dt.FetchSamplingData(true)
		dt.Make(maxDepth, f) 
	}
}

/*
*/ 
func (df *DForest) PredictMC() *PredictionMetric {

	// update sampling data
	df.dts[0].FetchSamplingData(true) 
	
	r,_ := df.dts[0].data.Dims() 
	actual,_ := MatrixColToFloat64Slice(df.dts[0].data, df.samplr.outputColumnIndices[0])
	results := make([]float64, r)

	wg := sync.WaitGroup{}
	for i := 0; i < r; i++ {

		wg.Add(1) 
		go func(i int) {
			row,_ := MatrixRowToFloat64Slice(df.dts[0].data, i)
			results[i] = df.VoteOnSample(row)
		}(i)
		wg.Done()
	}
	wg.Wait()

	return RunPredictionMetricMC(actual, results)
}

func (df *DForest) VoteOnSample(x []float64) float64 {

	outputCounts := make(map[float64]float64,0) 
	for _, dt := range df.dts {
		out := dt.PredictOneMC(x)
		outputCounts[out]++
	}

	// get key w/ max value
	_,maxKey := MinAndMaxKeyByValueOfMapFloat64ToFloat64(outputCounts)
	return maxKey
}

/*
return: 
- prediction, number of positives
*/
func (df *DForest) PredictSampleByClass(outcomeClassValue float64) []float64 {
	
	//// TODO: refactor this chunk 
	var data *mat.Dense
	switch {
	case df.samplr.imbal2:
		df.samplr.MakeDataForImbalanced(false, true)
		data = df.samplr.imbal2SamplingData 
	default: 
		/// TODO: needs to be checked.
		data = df.samplr.samplingData
	}

	///
	r,_ := data.Dims() 
	output := make([]float64, r)
	wg := sync.WaitGroup{}

	median := (outcomeClassValue - 1) / 2.0
	for j := 0; j < r; j++ {
		x,_ := MatrixRowToFloat64Slice(data, j)
		wg.Add(1) 
		go func(j int) {
			y := float64(0)
			for _, dt := range df.dts {
				index := dt.IndexOfTreeOfClass(outcomeClassValue)
				/// TODO: lock?? 
				y += dt.bdts[index].PredictOne(dt.bdts[index].ruut, x)
			}

			y = ZeroDivStandard(y, float64(len(df.dts)))

			// make random guess
			switch {
			case y == median: 
				if rand.Float64() >= 0.5 {
					output[j] = 1
				} else {
					output[j] = -1
				}

			case y > median: 
				output[j] = outcomeClassValue

			default: 
				output[j] = -1
			}
			wg.Done()
		}(j)
	}

	wg.Wait()
	return output
}