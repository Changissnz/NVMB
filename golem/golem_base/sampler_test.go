package golem

import (
	"testing"
	"fmt"
	"time"
	"math/rand" 
)

func SetUpTestSamplerReader() *ConcurrentFileBlockReader {
	fp := "./datos/other_data/creditcard.csv"
	cfbp := OneConcurrentFileBlockReader(fp)
	return cfbp
}

/*
*/
func Test_Sampler_PreprocessReader(t *testing.T) {
	readr := SetUpTestSamplerReader()
	
	start := time.Now() 
	s := OneSampler(readr, 0.526672448)
	duration := time.Since(start)

	fmt.Println("[0] duration ", duration)

	// time test 
	start = time.Now()
	s.PreprocessReader("regular", []string{}, []string{"Class"}) 
	duration = time.Since(start)
	fmt.Println("[1] duration ", duration) 

	// check column 31
	q,_ := MatrixColToFloat64Slice(s.sampleMAT, 30)
	fmt.Println("len of q ", len(q))
}

/*
*/
func Test_Sampler_AnalyzeDataForTypeDC(t *testing.T) {

	DISCRETE_COLUMN_DISCOVERY_DATA_SIZE_THRESHOLD = 0.0175557
	readr := SetUpTestSamplerReader()
	s := OneSampler(readr, 0.526672448)
	s.PreprocessReader("regular", []string{}, []string{"Class"}) 
	
	s.AnalyzeDataForTypeDC()
	if s.colFrequency[30][float64(0)] != 4997 {
		panic(fmt.Sprintf("incorrect frequency collected for column %d, got %f", 31, s.colFrequency[30][float64(0)]))
	}

	discreteCols := []int32{30}
	for k,v := range s.colDC {
		if IntInSlice(discreteCols,int32(k)) && v != "discrete" {
			panic(fmt.Sprintf("incorrect d/c type for column %d", k))
		}
	}
}

/*
*/
func Test_Sampler__AnalyzeDataForTypeDC__TimeTest(t *testing.T) {
	readr := SetUpTestSamplerReader()
	s := OneSampler(readr, 0.526672448)
	s.PreprocessReader("regular", []string{}, []string{"Class"}) 

	start := time.Now() 
	s.AnalyzeDataForTypeDC()
	duration := time.Since(start)
	fmt.Println("\tduration: ", duration)
}

/*
*/
func Test_Sampler_FrequencyCondenseColumns(t *testing.T) {
	readr := SetUpTestSamplerReader()
	s := OneSampler(readr, 0.526672448)
	s.PreprocessReader("regular", []string{}, []string{"Class"}) 
	s.AnalyzeDataForTypeDC()

	fmt.Println("COL FREQ")
	fmt.Println(s.colFrequency[30])

	start := time.Now() 
	s.FrequencyCondenseColumns(nil, "sample")
	duration := time.Since(start) 

	fmt.Println("duration: ", duration)

	//// uncomment below for displaying column frequencies
	/*
	for k, v := range s.colFrequencyCondensed {
		fmt.Println("key ", k) 
		fmt.Println("value ")
		DisplayPairSlice(v)
	}
	*/
}

/*
*/
func Test_Sampler_AnalyzeForOutcomeClassFrequencies(t *testing.T) {
	readr := SetUpTestSamplerReader()
	s := OneSampler(readr, 0.526672448)
	s.PreprocessReader("regular", []string{}, []string{"Class"}) 
	s.AnalyzeDataForTypeDC()

	start := time.Now() 
	s.FrequencyCondenseColumns(nil, "sample")

	duration := time.Since(start) 

	fmt.Println("\tduration: ", duration)
	fmt.Println("possible outcome classes ", len(s.colFrequencyCondensed[2]), " ", len(s.colFrequencyCondensed[3]), " ", len(s.colFrequencyCondensed[30]))
	s.AnalyzeForOutcomeClassFrequencies()
	fmt.Println("number of outcome classes ", len(s.probabilityFrequencies))

	////// display thresholds here 
	for k,v := range s.colFrequencyCondensed {
		fmt.Println("key ", k) 
		fmt.Println("value") 
		DisplayPairSlice(v) 
		fmt.Println("_________________________________")

	}
}

/*
*/
func Test_Sampler_PreAnalysis(t *testing.T) {
	readr := SetUpTestSamplerReader()
	s := OneSampler(readr, 0.526672448)

	s.PreAnalysis([]string{"Time"}, []string{"Class"})

	fmt.Println("ANALYSIS") 
	fmt.Println("\t prob weights") 
	fmt.Println(s.probabilityWeights)

	fmt.Println("\t prob frequencies") 
	fmt.Println(s.probabilityFrequencies)
} 

/*
*/
func Test_Sampler_GatherRowsWithKeyEntireData(t *testing.T) {
	CPARTSIZE = 50000
	readr := SetUpTestSamplerReader()
	s := OneSampler(readr, 0.526672448)

	s.PreAnalysis([]string{"Time"}, []string{"Class"})

	key := []string{DefaultFloat64ToString(1)}
	start := time.Now()
	s.GatherRowsWithKeysEntireData(key, false)
	duration := time.Since(start) 
	fmt.Println("duration ", duration) 

	r,_ := s.outcomeClassData[DefaultFloat64ToString(1)].Dims() 
	if r != 492 {
		panic("incorrect number of rows drawn")
	}
}

/*
*/
func Test_Sampler__RunSampling2(t *testing.T) {
	CPARTSIZE = 50000
	readr := SetUpTestSamplerReader()
	s := OneSampler(readr, 0.526672448)

	now := time.Now() 
	stat := s.RunSampling2([]string{"Time"}, []string{"Class"})
	if !stat {
		panic("minority keys not caught") 
	}
	duration := time.Since(now)
	fmt.Println("first draw time: ", duration)

	key := DefaultFloat64ToString(1) 
	r,c := s.outcomeClassData[key].Dims()

	if r != 492 || c != 30 {
		panic("incorrect outcome class dim.")
	}

	if len(s.outcomeClassData) != 1 {
		panic("incorrect number of outcome classes") 
	}

	r,c = s.samplingData.Dims()
	fmt.Println("dim. of sampling data: ", r, " ", c)
	
	/// try sampling again
	now = time.Now() 
	stat = s.RunSampling2([]string{"Time"}, []string{"Class"})
	duration = time.Since(now)
	fmt.Println("second draw time: ", duration)

	r2,c2 := s.samplingData.Dims()

	if r != r2 || c != c2 {
		panic("invalid second sampling") 
	}

	fmt.Println("outcome class data size: ", len(s.outcomeClassData))
}

func Test_Sampler_DrawSamplesConcurrent(t *testing.T) { 
	CPARTSIZE = 50000
	readr := SetUpTestSamplerReader()
	s := OneSampler(readr, 0.526672448)
	s.PreprocessReader("regular", []string{}, []string{"Class"}) 

	fmt.Println("full ", s.fullSize, " part ", s.leftoverSize)
	s.AnalyzeDataForTypeDC()
	s.FrequencyCondenseColumns(nil, "sample") 
	s.AnalyzeForOutcomeClassFrequencies()

	fmt.Println("time-test : draw samples concurrent") 
	start := time.Now() 
	s.DrawSamplesConcurrent(true, "outcome probability", s.sampleSize) 
	duration := time.Since(start)
	fmt.Println("\tduration: ", duration)
}

func Test_Sampler_DrawSamples_NotConcurrent(t *testing.T) {
	CPARTSIZE = 50000
	readr := SetUpTestSamplerReader()
	s := OneSampler(readr, 0.526672448)
	s.PreprocessReader("regular", []string{}, []string{"Class"}) 

	fmt.Println("full ", s.fullSize, " part ", s.leftoverSize) 

	s.AnalyzeDataForTypeDC()
	s.FrequencyCondenseColumns(nil, "sample") 
	s.AnalyzeForOutcomeClassFrequencies()

	fmt.Println("time-test : draw samples not concurrent") 
	start := time.Now() 
	s.DrawSamples(true, "outcome probability", 50000) 
	duration := time.Since(start)
	fmt.Println("\tduration: ", duration)
} 

/// TODO: make assertions for sampling data frequencies
func Test_Sampler___AnalyzeSamplingData(t *testing.T) {
	/// prelim
	rand.Seed(12)
	readr := SetUpTestSamplerReader()
	s := OneSampler(readr, 0.526672448)
	s.PreprocessReader("regular", []string{}, []string{"Class"}) 

	fmt.Println("full ", s.fullSize, " part ", s.leftoverSize)
	s.AnalyzeDataForTypeDC()
	s.FrequencyCondenseColumns(nil, "sample") 
	s.AnalyzeForOutcomeClassFrequencies()
	s.ModifyWeightsAccordingToSpec("inverse")
	s.DrawSamplesConcurrent(true, "outcome probability", s.sampleSize) 
	// analyzing sampling data frequencies 
	s.AnalyzeSamplingData()
		
	fmt.Println("sampling data frequencies") 
	fmt.Println(len(s.outcomeDataFrequencies)) 
	for k,v := range s.outcomeDataFrequencies {
		fmt.Println("key ", k, " || value ", v)
	}
}

/// X
func Test_Sampler__FixImbalancedData(t *testing.T) {
	/// prelim
	rand.Seed(12)
	CPARTSIZE = 50000

	readr := SetUpTestSamplerReader()
	s := OneSampler(readr, 0.526672448)
	s.PreprocessReader("regular", []string{}, []string{"Class"}) 

	fmt.Println("full ", s.fullSize, " part ", s.leftoverSize)
	s.AnalyzeDataForTypeDC()
	s.FrequencyCondenseColumns(nil, "sample")

	s.AnalyzeForOutcomeClassFrequencies()
	s.ModifyWeightsAccordingToSpec("inverse")
	s.DrawSamplesConcurrent(true, "outcome probability", s.sampleSize) 

	// analyzing sampling data frequencies 
	s.AnalyzeSamplingData()
	s.FixImbalancedData() 

	// check outcome class ratio 
	fmt.Println("OUTCOME DATA FREQUENCIES")
	fmt.Println(s.outcomeDataFrequencies)

		
	rat := s.outcomeDataFrequencies[DefaultFloat64ToString(float64(1))] / s.outcomeDataFrequencies[DefaultFloat64ToString(float64(0))]
	if rat < RATIO_BALANCE_THRESHOLD {
		panic("ratio was not fixed") 
	}
		
	CPARTSIZE = 1000
}

/*
*/
func Test_Sampler__CollectFrequencyMeanOnInputMatrix(t* testing.T) {
	CPARTSIZE = 50000
	readr := SetUpTestSamplerReader()
	s := OneSampler(readr, 0.526672448)
	s.RunSampling([]string{"Time"}, []string{"Class"}, true, "outcome probability")
	s.CollectFrequencyMeanOnInputMatrix(s.samplingData)
}

/*
//
*/
func Test_Sampler_MakeDataForImbalanced(t *testing.T) {
	CPARTSIZE = 50000
	readr := SetUpTestSamplerReader()
	s := OneSampler(readr, 0.526672448)
	s.RunSampling2([]string{"Time"}, []string{"Class"}) 	
	s.MakeDataForImbalanced(false,true) 

	// check that size of sampling data is 
	r2,_ := s.samplingData.Dims() 
	r,_ := s.outcomeClassData[DefaultFloat64ToString(1)].Dims() 
	if !EqualFloat64(float64(r) / float64(r2), 0.35, 2) {
		panic("sampling data is incorrect size")
	}

	// check that sampling data is not equal after each draw
	d1 := s.samplingData
	s.MakeDataForImbalanced(false,true) 	
	d2 := s.samplingData

	if (MatricesAreEqual(d1,d2)) {
		panic("matrices from two draws cannot be equal!")
	}
}

/////////////////////////////////////////////////////////////

/// TODO: make check for below
/// TODO: check how to merge imbal2SamplingData and samplingData
func Test_Sampler___DatasetGlass(t *testing.T) {
	fp := "./datos/other_data/glass.csv"
	cfbp := OneConcurrentFileBlockReader(fp)
	s := OneSampler(cfbp, 1.0) 
	s.RunSampling2([]string{}, []string{"Type"}) 

	fmt.Println("column types ") 
	fmt.Println(cfbp.columnTypes)

	fmt.Println("D OR C")
	fmt.Println(s.colDC)

	r,_ := s.samplingData.Dims() 
	fmt.Println("number of rows for samplingData ", r)

	fmt.Println("frequency mean on sample")
}