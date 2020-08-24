/*
A map collector is a data structure that stores a map of maps,
and runs statistical calculations on it.
*/
package golem

import (
	"fmt"
)

const DEFAULT_DOMINANCE_RATIO float64 = 0.5

type MapCollector struct {
	datos          map[int]map[string]int64
	predictedTypes map[int]string
	predictMode    string

	dr float64 // dominance ratio
}

/*
*/
func OneMapCollector() *MapCollector {
	return &MapCollector{datos: make(map[int]map[string]int64),
		dr: DEFAULT_DOMINANCE_RATIO}
}

/*
*/
func (m *MapCollector) DominanceRatio(f float64, threshold float64) bool {
	if (f >= threshold) {
		return true 
	}
	return false 
}


/// TODO: not used.
/*
*/
func (m *MapCollector) DominanceRatioThresholdTest(f float64) bool {

	if f >= m.dr {
		return true
	}

	return false
}

/*
adds one map of data to instance
*/
func (m *MapCollector) AddOne(x map[int]string, allowNewType bool, allowNewTypeKey bool) {

	for k, v := range x {
		// allows new datos types
		_, ok := m.datos[k]

		if !ok && !allowNewType {
			panic("cannot accept new type")
		}

		if !ok {
			nu := make(map[string]int64)
			m.datos[k] = nu
		}

		if allowNewType {
			// allows datos new keys
			_, ok2 := m.datos[k][v]

			if !ok2 && !allowNewTypeKey {
				panic("cannot accept new type key!")
			}

			m.datos[k][v]++
		}
	}
}

/*
runs predictions on keys of types
*/
func (m *MapCollector) Predict(mode string) {

	if mode != "threshold" && mode != "split" {
		panic(fmt.Sprintf("invalid mode %s", mode))
	}

	if mode == "split" {
		panic("still not finished")
	}

	m.predictedTypes = make(map[int]string, 0)
	m.predictMode = mode

	for k := range m.datos {
		t := m.PredictKeyOfTypeByRequiredThreshold(k, DEFAULT_DOMINANCE_RATIO)
		m.predictedTypes[k] = t
	}
}

/*
predicts key by dominance ratio
*/
func (m *MapCollector) PredictKeyOfType(x int) (string, float64) {
	r, ok := m.datos[x]

	if !ok {
		panic("key does not exist")
	}

	// sum up total
	var s int64 = 0 
	for _, v := range r {
		s += v  
	}

	var highestRatio float64 = 0
	var predictedKey string = "?" 
	for k,v := range r {

		if k == "?" {
			continue
		}

		ratio := float64(v) / float64(s)
		if ratio > highestRatio {
			highestRatio = ratio
			predictedKey = k
		}
	}

	return predictedKey, highestRatio
}

/*
*/
func (m *MapCollector) PredictKeyOfTypeByRequiredThreshold(x int, threshold float64) string {
	s, f := m.PredictKeyOfType(x)

	if m.DominanceRatio(f, threshold) {
		return s 
	}

	return "?"
}

/// TODO: incomplete
/*
use this for multi-output of predicted type
*/
func (m *MapCollector) PredictKeyOfTypeSplitDecision(numSplitDecisions int) {

}
