package golem

import (
	"math"
	"math/rand"
)

func MinAndMaxKeyOfBayesProbMap(bpm map[float64]map[string]float64) (float64,float64) {
	mini,maxi := math.Inf(1), math.Inf(-1)

	for k,_ := range bpm {
		if k < mini {
			mini = k
		}

		if k > maxi {
			maxi = k
		}
	}

	return mini, maxi
}

/*
*/ 
func SampleSubrangesForRange(targetRange *Pair, subRange *Pair, numberOfSubranges int, output []*Pair, ratio int, initial bool) []*Pair {

	if len(output) >= numberOfSubranges {
		return output
	}

	switch {
	case initial: 
		for i := 0; i < ratio; i++ {
			sz := (targetRange.b.(float64) - targetRange.a.(float64)) / float64(ratio)
			p := ChooseValidRangeGivenSize(targetRange, sz)
			output = append(output, p) 
		}

		nextChoice := RandomIntS(0,ratio)
		return SampleSubrangesForRange(targetRange, output[nextChoice], numberOfSubranges, output, ratio, false)
	
	default:
			// left of partition 
			stl := subRange.a.(float64) - targetRange.a.(float64)
			sz := stl / float64(ratio) 
			p := &Pair{a: subRange.a.(float64) - sz, b: subRange.a.(float64)}
			output = append(output,p)
			nextChoice := p 
			
			p2 := &Pair{a: subRange.a.(float64) -sz, b: subRange.b.(float64)} 
			output = append(output,p2)
			if rand.Float32() < 0.5 {
				nextChoice = p2
			}

			// right of partition
			stl = targetRange.b.(float64) - subRange.b.(float64)
			sz = stl / float64(ratio) 
			p3 := &Pair{a: subRange.b.(float64), b: subRange.b.(float64) + sz}
			output = append(output,p3)
			if rand.Float32() < 0.5 {
				nextChoice = p3
			}

			p4 := &Pair{a: subRange.a.(float64), b: subRange.b.(float64) + sz} 
			output = append(output,p4)
			if rand.Float32() < 0.5 {
				nextChoice = p4
			}

			return SampleSubrangesForRange(targetRange, nextChoice, numberOfSubranges, output, ratio, false) 
	}

	return output

}

// CAUTION: does not perform error-checking
func ChooseValidRangeGivenSize(newRange *Pair, size float64) *Pair {
	for {
		start := RandomFloat64(newRange.a.(float64), newRange.b.(float64))
		if start + size <= newRange.b.(float64) {
			return &Pair{a: start, b: start + size}
		}
	}
}