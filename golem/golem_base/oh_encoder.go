package golem 

import (
	"sync"
	"gonum.org/v1/gonum/mat"
)

/*
code contains implementation of one hot encoder
*/

func UniqueElementsOfStringSlice(x []string) *BasicSet {

	output := OneBasicSet()
	for _, x_ := range x {
		output.AddOne(x_) 
	}

	return output 
} 

/*
CAUTION: assumes encoding starts at value 0.0
*/
func UniqueElementsOfStringSlice_UsesPreviousEncoding(x []string, previousEncoding map[string]float64) *BasicSet { 
	output := OneBasicSet() 

	for _, x_ := range x {	
		_, ok := previousEncoding[x_] 
		if !ok {
			output.AddOne(x_)
		}
	}

	return output
}

func OHEncodeColumn_Dumb(x []string) []float64 { 

	encoding := make(map[string]int, 0) 
	ue := UniqueElementsOfStringSlice(x) 

	i := 0 
	for k,_  := range ue.datos {
		encoding[k] = i 
		i++ 
	}

	output := make([]float64,0)
	for _, x_ := range x {
		output = append(output, float64(encoding[x_]))  
	}
	return output 
}

func OHEncodeBlock_Dumb(b *Block) *mat.Dense {
	if (!b.transposeOn) {
		b.Transpose()
	}

	l := b.Length()

	encodedData := mat.NewDense(l, b.Width(), nil)	
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for i := 0; i < l; i++ { 
			x := OHEncodeColumn_Dumb(b.datos[i])
			encodedData.SetRow(i, x)
		}
		wg.Done() 
	}()

	wg.Wait()

	encodedData = TransposeMatrix(encodedData) 
	if (b.transposeOn) {
		b.Transpose()
	}

	return encodedData 
} 

/*
*/
func OHEncodeColumn_Dumb_UsesPreviousEncoding(x []string, previousEncoding map[string]float64) []float64 { 
	ue := UniqueElementsOfStringSlice_UsesPreviousEncoding(x, previousEncoding)
	i := len(previousEncoding) 
	for k,_  := range ue.datos {
		previousEncoding[k] = float64(i)
		i++ 
	}

	output := make([]float64,0)
	for _, x_ := range x {
		output = append(output, float64(previousEncoding[x_]))  
	}

	return output 
}

/*
*/
func OHEncodeBlock_Dumb_UsesPreviousEncoding(b *Block, previousEncoding map[int]map[string]float64) *mat.Dense {

	if (!b.transposeOn) {
		b.Transpose()
	}

	l := b.Length()
	encodedData := mat.NewDense(l, b.Width(), nil)	
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for i := 0; i < l; i++ {
			
			if previousEncoding[i] == nil {
				previousEncoding[i] = make(map[string]float64,0) 
			}

			x := OHEncodeColumn_Dumb_UsesPreviousEncoding(b.datos[i], previousEncoding[i])
			encodedData.SetRow(i, x)
		}
		wg.Done() 
	}()

	wg.Wait()

	encodedData = TransposeMatrix(encodedData) 
	if (b.transposeOn) {
		b.Transpose()
	}

	return encodedData 
}