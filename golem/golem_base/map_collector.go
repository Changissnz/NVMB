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

	m.predictedTypes = make(map[int]string, 0)
	m.predictMode = mode

	switch {

	case mode == "threshold" || mode == "majority": 
		for k := range m.datos {
			m.predictedTypes[k] = m.PredictKeyOfTypeByMode(mode, k, DEFAULT_DOMINANCE_RATIO)
		}

	case mode == "split": 
		panic("split mode not finished")

	default: 
		panic(fmt.Sprintf("invalid mode %s", mode))
}

}

//// TODO: int and float types may be confused w/ one another. 
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
	
	//// TODO: below will confuse int and float types 
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
		/////// below will correct for confusing int and float
	if predictedKey == "int" && r["float"] > 0 {
		predictedKey = "float" 
		highestRatio += (float64(r["float"]) / float64(s))
	}


	return predictedKey, highestRatio
}

/*
*/
func (m *MapCollector) PredictKeyOfTypeByMode(mode string, x int, threshold float64) string {
	s, f := m.PredictKeyOfType(x)

	switch {
	case mode == "threshold": 
		if m.DominanceRatio(f, threshold) {
			return s 
		}

		return "?"
	
	case mode == "majority":
		return s

	default: 
		panic("invalid mode for predicting key")

	}

}
