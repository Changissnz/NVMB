package golem

import (
	"fmt"
	"math"
)

var MAX_NUMBER_OF_QUALRANGES int = 8

/*
*/
func (dt *DTree) RandomGreedyRangeFinder(xyBayes map[float64]map[string]float64, f func(float64) float64, targetY string, treeIndex int, colIndex int, r float64) (*Pair, float64, float64)  { 

	sl,sr := 0.5,0.5
	var ml,mr float64

	occupied := make([]interface{}, 0) 
	mini,maxi := MinAndMaxKeyOfBayesProbMap(xyBayes)
	floatRange := &Pair{a: mini, b: maxi}
	for i := 0; i < MAX_NUMBER_OF_QUALRANGES; i++ {
		quit := false 
		switch {
		case dt.samplr.colDC[colIndex] == "discrete": 
			br,one,two, ml_,mr_ := dt.ChooseNextBestRangeDiscrete(xyBayes, f, floatRange, occupied, targetY, treeIndex, r) 

			if one + two < sl + sr && br != math.Inf(-1) {
				occupied = append(occupied, br) 
				sl,sr = one,two 
				ml,mr = ml_,mr_
			} else {
				quit = true
			}

		case dt.samplr.colDC[colIndex] == "continuous": 
			br,one,two, ml_,mr_ := dt.ChooseNextBestRange(xyBayes, f, floatRange, occupied, targetY, treeIndex, dt.samplr.colDC[colIndex], r) 
			if one + two < sl + sr && br != nil {
				occupied = append(occupied, br) 
				sl,sr = one,two 
				ml,mr = ml_,mr_
			} else {
				quit = true
			}
		default: 
			panic(fmt.Sprintf("invalid type %s", dt.samplr.colDC[colIndex]))
		}

		if quit {
			break
		}

	}

	if len(occupied) == 0 {
		return nil, ml,mr
	}

	// make split here
	s := OneSplit(occupied, dt.samplr.colDC[colIndex])
	a,b := &Pair{a: s, b: colIndex}, &Pair{a: sl, b: sr} 
	return &Pair{a: a, b: b}, ml, mr
}

/*
*/
func (dt *DTree) ChooseNextBestRangeDiscrete(xyBayes map[float64]map[string]float64, f func(float64) float64, floatRange *Pair, occupied []interface{}, targetY string, treeIndex int, r float64) (float64, float64, float64, float64, float64) {
	qr := dt.ChooseQualifyingRange(floatRange, occupied, "discrete").([]int) 
	if len(qr) == 0 {
		return math.Inf(-1), 0.5, 0.5, 0.5, 0.5 
	}

	var bestSubRange float64 
	scoreLeft,scoreRight := 0.5,0.5
	mismatchLeft, mismatchRight := 0.5, 0.5

	occ2 := InterfaceSliceToFloatSlice(occupied) 
	for _,s := range qr {
		occ := append(occ2, float64(s)) 
		yx,x,nyx,nx := ReadBayesProbMapDiscrete(xyBayes, occ, targetY)

		////// TODO: error-check, spareColFrequency and colFrequency are the same
		/*
		q := PairSliceToFloat64Slice(dt.samplr.colFrequencyCondensed[2], 0)
		q2 := PairSliceToFloat64Slice(dt.samplr.spareColFrequencyCondensed[2], 0)
		if !EqualFloat64Slices(q,q2,4) {
			panic("not equal!")
		}
		*/ 
		one,two := SplitCostStandard(r, yx, x, nyx, nx, f)
		two = two * dt.bdts[treeIndex].outcomeClassWeight

		if one + two < scoreLeft + scoreRight {
			bestSubRange = float64(s)
			scoreLeft,scoreRight = one,two
			mismatchLeft,mismatchRight = MismatchValues(yx,x,nyx,nx) 			
		}
	}
	return bestSubRange, scoreLeft,scoreRight, mismatchLeft, mismatchRight
}

//// TODO: check this for left and right score (ordering)
/*
*/
func (dt *DTree) ChooseNextBestRange(xyBayes map[float64]map[string]float64, f func(float64) float64, floatRange *Pair, occupied []interface{}, targetY string, treeIndex int, dataType string, r float64) (*Pair, float64,float64,float64,float64) {
	qr := dt.ChooseQualifyingRange(floatRange, occupied, dataType)
	if qr == nil {
		return nil, 0.5, 0.5, 0.5, 0.5
	}

	occ2 := InterfaceSliceToPairSlice(occupied)
	sr := SampleSubrangesForRange(qr.(*Pair), nil , RANGE_SEARCH_SPACE_MAX_THRESHOLD, []*Pair{qr.(*Pair)} , 4, true)

	var bestSubRange *Pair 
	scoreLeft,scoreRight := 0.5,0.5
	mismatchLeft, mismatchRight := 0.5, 0.5
	for _,s := range sr {
		occ := append(occ2, s)
		yx,x,nyx,nx := ReadBayesProbMap(xyBayes, occ, targetY)

		////// TODO: error-check, spareColFrequency and colFrequency are the same
		/*
		q := PairSliceToFloat64Slice(dt.samplr.colFrequencyCondensed[2], 0)
		q2 := PairSliceToFloat64Slice(dt.samplr.spareColFrequencyCondensed[2], 0)
		if !EqualFloat64Slices(q,q2,4) {
			panic("not equal!")
		}
		*/ 
		one,two := SplitCostStandard(r, yx, x, nyx, nx, f)

		/// TODO: add weight here
		two = two * dt.bdts[treeIndex].outcomeClassWeight
		if one + two < scoreLeft + scoreRight {
			bestSubRange = s
			scoreLeft,scoreRight = one,two
			mismatchLeft,mismatchRight = MismatchValues(yx,x,nyx,nx) 			
		}
	}
	return bestSubRange, scoreLeft,scoreRight, mismatchLeft, mismatchRight
}

func MismatchValues(yx,x,nyx,nx float64) (float64, float64) {
	return nyx, x - yx
}

/*
*/
func (dt *DTree) ChooseQualifyingRange(floatRange *Pair, usedRanges []interface{}, dataType string) interface{} {

	switch {
	case dataType == "continuous":
		ur := InterfaceSliceToPairSlice(usedRanges)
		unoccupied := BayesianIdentityOfFloatSliceSubset(floatRange, ur)

		// choose one range in unoccupied and check for subranges of its partitions 
		if len(unoccupied) == 0 {
			return nil
		}

		ri := RandomIntS(0, len(unoccupied))
		return unoccupied[ri]
	case dataType == "discrete":
		ur := InterfaceSliceToIntSlice(usedRanges)
		floatRange.b = floatRange.b.(float64) + 1
		unoccupied := BayesianIdentityOfIndexSubset(floatRange, ur, false)
		return unoccupied
	default: 
		panic("invalid data type for choosing qualifying range") 
	}
	
}

// CAUTION: data type assumed to be `spare`
/*
*/
func (dt *DTree) IndicesOfColumnInQualRange(y []float64, targetIndices []int, colIndex int, qualRanges []*Pair) ([]int,[]int) {
	indices, nindices := make([]int,0), make([]int,0)

	for _,x := range targetIndices {
		index := dt.samplr.GroupOfValue(y[x], dt.samplr.spareColFrequencyCondensed[colIndex], dt.samplr.colDC[colIndex]) 
		l := dt.samplr.spareColFrequencyCondensed[colIndex][index].a.(float64)

		if IsValueInQualRange(l, qualRanges) {
			indices = append(indices, x)
		} else {
			nindices = append(nindices, x) 
		}
	}

	return nindices, indices
}

/*
*/
func (dt *DTree) IndicesOfColumnInQualRangeDiscrete(y []float64, targetIndices []int, colIndex int, qualRanges []float64) ([]int,[]int) {
	indices, nindices := make([]int,0), make([]int,0)

	for _,x := range targetIndices {
		index := dt.samplr.GroupOfValue(y[x], dt.samplr.spareColFrequencyCondensed[colIndex], dt.samplr.colDC[colIndex]) 
		l := dt.samplr.spareColFrequencyCondensed[colIndex][index].a.(float64)

		if Float64InSlice(qualRanges, l) {
			indices = append(indices, x)
		} else {
			nindices = append(nindices, x) 
		}
	}

	return nindices, indices
}