package golem

import (
	"fmt"
	"math"
	"sync"
	"sort"
	"gonum.org/v1/gonum/mat"
)

/*
NOTE: 
sampler is not specialized for time-series data
*/

//--------------------------------------------------------------------------------
/// ISSUE: sampler will be used only for regular-type matrices 
///			cfbr-type matrices will not work!! 
/// ISSUE: discrete-valued columns have classification assumed to be dumb-encoded
/// [!!] ISSUE: pre-analysis should run counts on entire file instead of the first n samples.

//--------------------------------------------------------------------------------

/// make method to automatically alter prop. to data file size
var DISCRETE_COLUMN_DISCOVERY_DATA_SIZE_THRESHOLD float64 = 0.15

// used for automatic detection of discrete data based on class-to-data ratio
var CLASS_TO_DATA_RATIO_MAX_THRESHOLD float64 = float64(10) / float64(100) 
var MAX_OUTCOMESET_CLASSES = 100000 // set this variable to limit the number of outcome-set classes to consider
var DEFAULT_CONDENSE_SIZE int = 50

/////// TODO: make zoom on colFrequencyCondensed
var NUM_PARTITIONS_PER_CHUNK int = 100 
var RATIO_IMBALANCE_THRESHOLD float64 = float64(0.1) 
var RATIO_BALANCE_THRESHOLD float64 = 0.35

type Sampler struct {
	readr *ConcurrentFileBlockReader
	sampleSize int
	sampleSizeRatio float64 
	initialRead bool

	sampleCFBR *CFBRDataMatrix
	sampleMAT *mat.Dense

	samplingData *mat.Dense
	///colMinMax map[int]*Pair
	imbal2SamplingData *mat.Dense
	outcomeDataFrequencies map[string]float64

	/// spare variables are used to load a sub-matrix of sampling data  
	spareData *mat.Dense
	spareColMeans map[int]float64
	spareColFrequency map[int]map[float64]float64
	spareColFrequencyCondensed map[int][]*Pair

	leftoverSize int32 
	fullSize int32 

	colDC map[int]string
	colMeans map[int]float64
	colFrequency map[int]map[float64]float64
	colFrequencyCondensed map[int][]*Pair // Pair(float64::range, float64::score)
	condenseSize int
	condenseZoom map[int]int

	// TODO: find a way to refactor [0] and [1] below
	outcomeClasses []string 
	probabilityFrequencies map[string]float64 // for outcome class 
	probabilityWeights []float64 // equal to length of partition

	/// concurrency [1] 
	sampleChunk *mat.Dense
	chunkIndex int64

	// bayesian stats
	/// 	map[colIndex]map[inputClass]map[outcomeClass]frequency
	bayesProbMap map[int]map[float64]map[string]float64

	excludedColumnIndices []int
	outputColumnIndices []int

	/// var. for sampling method #2
	pre2 bool
	imbal2 bool // imbalanced data
	numberOfPresampled int

	valid bool 
	outcomeClassData map[string]*mat.Dense
	minorityClassKeys []string
}

/*
*/
func OneSampler(r *ConcurrentFileBlockReader, ssr float64) *Sampler {
	x := &Sampler{readr: r, sampleSizeRatio: ssr, initialRead: false, colMeans: make(map[int]float64,0), colFrequency: make(map[int]map[float64]float64,0), condenseZoom: make(map[int]int,0), pre2: false, numberOfPresampled: 0}  
	x.fullSize, x.leftoverSize = x.CalculateSamplesPerPartition() 
	return x
} 

//////////////////////////////// start: index functions 

/*
*/
func (s *Sampler) InputColumnIndices() []int {
	_,c := s.samplingData.Dims()
	input := make([]int,0) 

	for i := 0; i < c; i++ {
		if !IntSInSlice(s.outputColumnIndices, i) {
			input = append(input, i) 
		}
	}
	return input 
}

/*
*/
func (s *Sampler) AllIndices() []int {
	_,c := s.samplingData.Dims()
	input := make([]int,0) 

	for i := 0; i < c; i++ {
		input = append(input, i) 
	}
	return input 
}

//////////////////////////////// end: index functions 

//////////////////////////////// start: read partition 

/*
collect's one partition's worth of data starting from the line marker
*/ 
func (s *Sampler) CollectOnePartition(sampleType string, encodeType string) {///, numSamples int) {
	if len(s.readr.blockData) == 0 {
		s.readr.ReadPartition("exact")
	}

	if sampleType == "regular" {
		s.PartitionToMatrix(sampleType, encodeType)
	} else {
		panic("[FUTURE] CFBR data collection not supported")
	}
}

/*
converts partition data to a matrix
*/
func (s *Sampler) PartitionToMatrix(sampleType string, encodeType string) {

	var d1 *CFBRDataMatrix 
	var d2 *mat.Dense
	for i := 0; i < len(s.readr.blockData); i++ {
		switch {
			case sampleType == "cfbr": 
				x := s.readr.ConvertPartitionBlockToMatrix(i)
				x.Preprocess() 

				if d1 == nil {
					d1 = x
				} else {
					d1.StackMatrix(x) 
				}

			case sampleType == "regular": 
				x2 := s.readr.ConvertPartitionBlockToRegularMatrix(i, encodeType) 
				if d2 == nil {
					d2 = x2 
				} else {
					d2 = StackTwoMatrices(d2, x2)  
				}

			default: 
				panic("invalid sample type")
		}
	}  

	s.sampleCFBR = d1 
	s.sampleMAT = d2 
}

//////////////////////////////// end: read partition 

/////////////////////////////// start: preprocessing

/*
preprocesses read for valid column types, output column indices, and column class type
*/
func (s *Sampler) PreprocessReader(wantedDataType string, excludedColumns []string, outputColumns []string) bool {
	s.excludedColumnIndices = StringSliceIndicesInStringSlice(s.readr.columns, excludedColumns)
	s.outputColumnIndices = StringSliceIndicesInStringSlice(s.readr.columns, outputColumns)
	s.outputColumnIndices = s.readr.SetColumnsForRemoval(s.outputColumnIndices, s.excludedColumnIndices)
	// get column types
	s.readr.DeduceColumnTypes() 
	s.CollectOnePartition(wantedDataType, "dumb")
	s.valid = s.CheckValidColumnTypes(wantedDataType)

	if !s.valid {
		return false
	}

	return true
}

/*
checks column types for validity
*/
func (s *Sampler) CheckValidColumnTypes(wantedDataType string) bool {

	bs := s.readr.GetColumnTypes()
	
	if (bs.DoesExist("?")) {
		return false
	}
	
	switch {

		case wantedDataType == "regular":  
			if (bs.DoesExist("vector")) {
				return false
			}

		case wantedDataType == "cfbr":  
			return true
		default: 
			panic("[X] invalid data type")
	}

	return true
}

/////////////////////////////// end: preprocessing

/////////////////////////////// start: deduce column types

/// TODO: risky. may want to do random sampling for entire data instead.
/*
analyzes columns for discrete or continuous types

CAUTION: use only with regular matrices
*/
func (s *Sampler) AnalyzeDataForTypeDC() int {
	s.readr.SetReaderToBeginning() 

	numConsidered := int(math.Ceil(float64(s.readr.fileLength) * DISCRETE_COLUMN_DISCOVERY_DATA_SIZE_THRESHOLD))
	distSize := int(0)
	for {
		if distSize >= numConsidered {
			break
		}

		s.CollectOnePartition("regular", "dumb")
		distSize = s.CollectFrequencyMeanOnData("regular", "partition", distSize)
		s.readr.blockData = make([]*Block, 0) 
	}

	s.AnalyzeColumnsForType(distSize)
	return distSize
}

/*
*/
func (s *Sampler) AnalyzeColumnsForType(distSize int) {
	s.colDC = make(map[int]string,0)
	for k, v := range s.colFrequency {
		x := s.FrequencyDiscreteDataRatio(v, distSize)
		if x {
			s.colDC[k] = "discrete"
		} else {
			s.colDC[k] = "continuous"
		}
	}
}

/*
*/
func (s *Sampler) FrequencyDiscreteDataRatio(freq map[float64]float64, distSize int) bool {
	results := float64(len(freq)) / float64(distSize)
	if results <= CLASS_TO_DATA_RATIO_MAX_THRESHOLD {
		return true 
	}
	return false
}

/*
*/
func (s *Sampler) CollectFrequencyMeanOnData(wantedDataType string, dataType string, distSize int) int {

	////// no conc
		// get r and c 
	var r,c int 
	switch {
		case dataType == "spare": 
			r,c = s.spareData.Dims() 
		case dataType == "sample": 
			r,c = s.samplingData.Dims() 
		case dataType == "partition": 
			r,c = s.sampleMAT.Dims()
		default: 
			panic("invalidaatic")
	}
	
	for i := 0; i < c; i++ {
		switch {
			case dataType != "spare": 
				s.colMeans[i], s.colFrequency[i] = s.CollectFrequencyMeanOnColumn(wantedDataType, dataType, i, s.colMeans[i], s.colFrequency[i], distSize)

			default:
				s.spareColMeans[i], s.spareColFrequency[i] = s.CollectFrequencyMeanOnColumn(wantedDataType, dataType, i, s.spareColMeans[i], s.spareColFrequency[i], distSize) 
		}
	}

	////// conc patt 2 
	/*
	wg := sync.WaitGroup{}
	mut := sync.Mutex{} 
	r,c := s.sampleMAT.Dims() 
	
	for i := 0; i < c; i++ {
		
		wg.Add(1)
		go func(i int) {
			mut.Lock() 
			s.colMeans[i], s.colFrequency[i] = s.CollectFrequencyMeanOnColumn(wantedDataType, i, s.colMeans[i], s.colFrequency[i], distSize)
			mut.Unlock() 
			wg.Done() 
		}(i)  
	} 

	wg.Wait() 
	*/
	///////////////////////////////////////////////////////////////

	return distSize + r 
}

/*
*/ 
// CAUTION: no argument check 
func (s *Sampler) CollectFrequencyMeanOnInputMatrix(m *mat.Dense) { 
	s.spareData = m
	s.spareColMeans = make(map[int]float64,0)
	s.spareColFrequency = make(map[int]map[float64]float64,0) 
	s.CollectFrequencyMeanOnData("regular", "spare", 0)
} 

/*
*/
func (s *Sampler) CollectFrequencyMeanOnColumn(wantedDataType string, dataType string, column int, mean float64, freq map[float64]float64, distSize int) (float64, map[float64]float64) { 

	var r float64
	var mean2 float64

	if freq == nil {
		freq = make(map[float64]float64, 0)
	}

	switch {

		case wantedDataType == "cfbr": 
			panic("[FUTURE] output stat measures for cfbr")		

		case wantedDataType == "regular":
			// collect average
			var col []float64 
			col = s.CollectOneRCByDataType(dataType, "c", column)
			mean2 = MeanIterable(col)
			var r_ int

			switch {
				case dataType == "partition": 
					r_,_ = s.sampleMAT.Dims()
					for i := 0; i < r_; i++ {
						v := s.sampleMAT.At(i,column)
						freq[v] += float64(1)
					}
				case dataType == "sample": 
					r_,_ = s.samplingData.Dims() 
					for i := 0; i < r_; i++ {
						v := s.samplingData.At(i,column)
						freq[v] += float64(1)
					}
				case dataType == "spare": 
					r_,_ = s.spareData.Dims() 
					for i := 0; i < r_; i++ {
						v := s.spareData.At(i,column)
						freq[v] += float64(1)
					}
				default: 
					panic("YES")
				
			}
			r = float64(r_)
		}

	mean = (mean2 * float64(r) + mean * float64(distSize)) / (float64(r) + float64(distSize)) 	
	return mean, freq
}

/////////////////////////////////////////////////////////////// sampling

		/////////////////////////// start: concurrency experiment //////////////////

/*
draw samples concurrent
*/ 
func (s *Sampler) DrawSamplesConcurrent(unique bool, weightType string, sampleSize int) {
	s.readr.SetReaderToBeginning() 
	s.chunkIndex = 0 
	s.samplingData = nil 
	chunksz := 0
	j := 0
	for chunksz < int(s.readr.fileLength) {
		m, sz := s.DrawFromChunk(unique, weightType, sampleSize)
		
		if sz == 0 {
			break
		}

		if s.samplingData == nil {
			s.samplingData = m
		} else {
			s.samplingData = StackTwoMatrices(s.samplingData, m) 
		}

		chunksz += sz
		j++ 
	}
}

/// TODO: refactor this with `DrawFromPartition`
/*
*/
func (s *Sampler) DrawFromChunk(unique bool, weightType string, sampleSize int) (*mat.Dense, int) {
	sz := s.MakeChunk() 	
	
	if sz == 0 {
		fmt.Println("\tchunk is size 0!")
		return nil, sz
	}
	
	numSamples := int(math.Floor( float64(sz) / float64(s.readr.fileLength) * float64(sampleSize)))
	s.MakeWeightsForData("chunk", weightType)
	drawnIndices := DrawIndicesFromRangesConcurrent(s.probabilityWeights, numSamples, unique)
	return MatrixCollectByIndices(s.sampleChunk, drawnIndices), sz 
}


/*
*/
func (s *Sampler) MakeChunk() int {
	s.sampleChunk = nil 
	chunksz := 0

	for i := 0; i < NUM_PARTITIONS_PER_CHUNK; i++ {
		if s.chunkIndex == s.readr.partitionLength + 1 {
			break 
		}

		s.CollectOnePartition("regular", "dumb") 

		if s.sampleChunk == nil {
			s.sampleChunk = s.sampleMAT 
		} else {
			s.sampleChunk = StackTwoMatrices(s.sampleChunk, s.sampleMAT) 
		}

		r,_ := s.sampleMAT.Dims() 
		chunksz += r

		s.readr.blockData = make([]*Block,0) 
		s.chunkIndex++ 
	}

	return chunksz
}

		/////////////////////////// end: concurrency experiment //////////////////

/*
*/
func (s *Sampler) DrawSamples(uniqueDraw bool, weightType string, sampleSize int) { 
	s.readr.SetReaderToBeginning() 
	s.samplingData = nil

	for i := 0; i < int(s.readr.partitionLength); i++ {
		s.CollectOnePartition("regular", "dumb") 
		r,_ := s.sampleMAT.Dims()
		numSamples := int(math.Floor( float64(r) / float64(s.readr.fileLength) * float64(sampleSize)))

		newMat := s.DrawFromPartition(numSamples, uniqueDraw, weightType) 
		if s.samplingData != nil {
			s.samplingData = StackTwoMatrices(s.samplingData, newMat) 
		} else {
			s.samplingData = newMat
		}

		s.readr.blockData = make([]*Block,0) 
	}

}

/*
Draws required number of samples from partition by randomly choosing 
weighted indices.

Sampling can be one of unique|non-unique
*/
func (s *Sampler) DrawFromPartition(numSamples int, uniqueDraw bool, weightType string) *mat.Dense {
	
	// construct the weights
	s.MakeWeightsForData("partition", weightType)

	// drawn random indices 
	drawnIndices := DrawIndicesFromRangesConcurrent(s.probabilityWeights, numSamples, uniqueDraw)
	return MatrixCollectByIndices(s.sampleMAT, drawnIndices) 
}

/////
/*
makes weights for matrix rows in proportion to outcome class 

CAUTION: only regular matrix type supported at this time. 
*/
func (s *Sampler) MakeWeightsForData(datasizeType string, weightType string) {
	s.probabilityWeights = make([]float64,0) 
	weight := float64(0) 
	var r int 

		/// NOTE: delete here after experiment 
	switch {
		case datasizeType == "chunk": 
			r,_ = s.sampleChunk.Dims() 
		case datasizeType == "partition": 
			r,_ = s.sampleMAT.Dims() 
		default: 
			panic("invalid datasize type") 
	}

	switch {
	case weightType == "outcome probability": 
		var q []float64
		for i := 0; i < r; i++ {
			q = s.CollectOneRCByDataType(datasizeType, "r", i)
			class := s.DetermineOutputClassOfRow(q, "sample") /// CAUTION 
			weight += s.probabilityFrequencies[class] 
			s.probabilityWeights = append(s.probabilityWeights, weight) 
		}

	case weightType == "uniform": 
		x := float64(1) / float64(r) 
		for i := 1; i <= r; i++ {
			s.probabilityWeights = append(s.probabilityWeights, x * float64(i)) 
		}

	default: 
		panic(fmt.Sprintf("invalid weight type %s", weightType))
	}
}

/*
calculate samples per partition based on total size of dataset
*/ 
func (s *Sampler) CalculateSamplesPerPartition() (int32,int32) {
	// get total size
	s.readr.PrelimRead()
	s.sampleSize = int(math.Ceil(float64(s.readr.fileLength) * s.sampleSizeRatio))

	// get samples for last partition
	leftovers := math.Mod(float64(s.readr.fileLength), float64(CPARTSIZE))
	leftovers = leftovers / float64(s.readr.fileLength) * float64(s.sampleSize) // leftover sample size

	// get samples for full partition
	var samplesPerFullPart float64 = 0.0 
	if s.readr.partitionLength > 1 {
		samplesPerFullPart = math.Floor((float64(s.sampleSize) - leftovers) / float64(s.readr.partitionLength - 1))
	}

	return int32(samplesPerFullPart), int32(leftovers) 
}

/////////////////////////////////////////////////////////////// start: probability weight modification 

/*
modify outcome-class weights according to frequency 
*/ 
func (s *Sampler) ModifyWeightsAccordingToSpec(spec string) {
	switch {
		case spec == "direct": 
			return 
		case spec == "inverse": 
			// find the maximum and minimum frequency
			minimumOutcomeFrequency, maximumOutcomeFrequency,_,_ := MinAndMaxOfMapStringToFloat64(s.probabilityFrequencies)

			// iterate through invert the values 
			for k,v := range s.probabilityFrequencies {
				s.probabilityFrequencies[k] = s.InverseFrequencyOfValue(v, minimumOutcomeFrequency, maximumOutcomeFrequency)
			}
			 
			return

		default: 
			panic("invalid spec for weight modification")
	}
}

/*
*/
func (s *Sampler) InverseFrequencyOfValue(v float64, minimum float64, maximum float64) float64 {
	return maximum - (v - minimum) 
} 

/////////////////////////////////////////////////////////////// end: probability weight modification 


//////////////////////////////////////////////////////////// start: outcome class frequency collection

/*
collects the frequencies of outcome classes by iterating through rows
*/
func (s *Sampler) AnalyzeForOutcomeClassFrequencies() {

	s.readr.SetReaderToBeginning() 
	/// TODO: caution, clearing variables 
	s.probabilityWeights = make([]float64,0) 
	s.probabilityFrequencies = make(map[string]float64, 0)
	s.outcomeDataFrequencies = make(map[string]float64, 0)

	distSize := int(0) 
	numConsidered := int(math.Ceil(float64(s.readr.fileLength) * DISCRETE_COLUMN_DISCOVERY_DATA_SIZE_THRESHOLD))
	for {
		if distSize >= numConsidered {
			break
		}

		s.CollectOnePartition("regular", "dumb")
		distSize += s.DetermineOutcomeClassesOfData("partition") 
		s.readr.blockData = make([]*Block,0) 
	}
}

/*
CAUTION: only regular matrices supported at this time
*/
func (s *Sampler) DetermineOutcomeClassesOfData(dataType string) int {

	if len(s.outputColumnIndices) == 0 {
		panic("cannot determine partition classes with 0 columns to consider!")
	}

	mux := sync.Mutex{} 
	wg := sync.WaitGroup{}

	var r int 
	switch {
		case dataType == "partition": 
			r,_ = s.sampleMAT.Dims()

		case dataType == "sample":  
			r,_ = s.samplingData.Dims() 
			s.outcomeDataFrequencies = make(map[string]float64,0)	
		case dataType == "imbal": 
			r,_ = s.imbal2SamplingData.Dims() 
			s.outcomeDataFrequencies = make(map[string]float64,0)	

		default: 
			panic("invalid data type")
	}

	for i := 0; i < r; i++ {
		wg.Add(1) 
		go func(i int) {
			
			var q []float64 
			switch {
				case dataType == "partition": 
					q,_ = MatrixRowToFloat64Slice(s.sampleMAT, i) 
					class := s.DetermineOutputClassOfRow(q, "sample") /// CAUTION 
					mux.Lock() 
					s.probabilityFrequencies[class] += float64(1) 
					mux.Unlock() 
				case dataType == "sample": 
					q,_ = MatrixRowToFloat64Slice(s.samplingData, i) 
					class := s.DetermineOutputClassOfRow(q, "sample") /// CAUTION 
					mux.Lock() 
					s.outcomeDataFrequencies[class] += float64(1) 
					mux.Unlock()

				case dataType == "imbal": 
					q,_ = MatrixRowToFloat64Slice(s.imbal2SamplingData, i) 
					class := s.DetermineOutputClassOfRow(q, "sample") /// CAUTION
					mux.Lock() 
					s.outcomeDataFrequencies[class] += float64(1)
					mux.Unlock()
			}
			wg.Done()
		}(i)
	} 
	wg.Wait()
	
	return r 
}

/////////////////////////////// start: outcome class ratio imbalance solution 

/*
analyzes sampling data based on outcome classes
*/
func (s *Sampler) AnalyzeSamplingData()  {
	s.outcomeDataFrequencies = make(map[string]float64,0)
	s.DetermineOutcomeClassesOfData("sample") 
}

/*
fixes data with imbalanced outcome class frequencies by substitute-duplication of those samples 
*/
func (s *Sampler) FixImbalancedData() {
	mc := s.ObtainMinorityClasses()
	for _, c := range mc {
		s.FixImbalancedOutcomeClass(c) 
	} 
}

/// TODO: not tested on non-binary outcome class data
/*
*/ 
func (s *Sampler) FixImbalancedOutcomeClass(key string) {

	// choose an outcome with excessive data to replace
	_,_,_, maxKey := MinAndMaxOfMapStringToFloat64(s.outcomeDataFrequencies)
	s.ReplaceMajorityOutcomeClassWithMinority(maxKey, key) 
}

/*
Replaces outcome class and updates probability frequencies
Used when data is disproportionately skewed towards 1+ outcome classes.
*/ 
func (s *Sampler) ReplaceMajorityOutcomeClassWithMinority(majorityClassKey string, minorityClassKey string) { 
	
	// gather all minority samples
	m,_ := s.GatherRowsWithKeys("sample", []string{minorityClassKey}, "data", false)

	// get indices of majority 
	_,indices := s.GatherRowsWithKeys("sample", []string{majorityClassKey}, "indices", false)

	// calculate the number of variables that must be replaced
	numToReplace := s.CalculateNumberOfMajorityToReplace(majorityClassKey, minorityClassKey)
	if len(indices) < numToReplace {
		panic("cannot replace more samples than its total size") 
	}

	// get the random majority indices to replace
	indices = RandomElementsInIntSlice(indices, numToReplace, "nr")

	// iterate through the indices and replace them 
	r,_ := m.Dims() 
	for _,i := range indices {
		// choose a random 
		r_ := RandomIntInRange(0, r) 
		row,_ := MatrixRowToFloat64Slice(m, r_) 
		s.samplingData.SetRow(i, row) 
	}

	s.outcomeDataFrequencies[majorityClassKey] -= float64(numToReplace) 
	s.outcomeDataFrequencies[minorityClassKey] += float64(numToReplace) 
}  

/*
*/ 
func (s *Sampler) CalculateNumberOfMajorityToReplace(majorityClassKey string, minorityClassKey string) int { 
	m1,m2 :=  s.outcomeDataFrequencies[minorityClassKey], s.outcomeDataFrequencies[majorityClassKey]
	toRep := (m1 - RATIO_BALANCE_THRESHOLD * m2) / (-1.0 - RATIO_BALANCE_THRESHOLD) 
	return int(math.Ceil(toRep))
}

/// TODO: refactor for classes other than outcome classes
/// TODO: implement concurrency
/*
gather samples with outcome class equal to key, samples will be in the same order 
as they exist in the input matrix

return: 
- matrix of rows w/ outcome class 
- indices of rows w/ outcome class
*/ 
func (s *Sampler) GatherRowsWithKeys(dataType string, keys []string, outputType string, collectFrequency bool) (*mat.Dense, []int) { 
	
	sampleIndices := make([]int,0)

	if collectFrequency {
		s.outcomeDataFrequencies = make(map[string]float64,0)
	}

	var r,c int
	switch {
		case dataType == "partition": 
			r,c = s.sampleMAT.Dims()

		case dataType == "sample": 
			r,c = s.samplingData.Dims() 

		case dataType == "spare": 
			r,c = s.spareData.Dims() 

		case dataType == "chunk": 
			r,c = s.sampleChunk.Dims() 

		default: 
			panic("invalid data type")
	}

	r_ := 0
	mux := sync.Mutex{} 
	wg := sync.WaitGroup{}
	var condenseDataType string

	/// TODO: caution here
	if dataType == "spare" {
		condenseDataType = dataType
	} else {
		condenseDataType = "sample"
	}

	for i := 0; i < r; i++ {
		wg.Add(1) 
		go func(i int) {
			row := s.CollectOneRCByDataType(dataType, "r", i)
			class := s.DetermineOutputClassOfRow(row, condenseDataType)

			if collectFrequency {
				mux.Lock() 
				s.outcomeDataFrequencies[class]++
				mux.Unlock() 
			}

			if StringIndexInSlice(keys, class) != -1 {
				mux.Lock() 
				sampleIndices = append(sampleIndices, i) 
				r_++ 
				mux.Unlock() 
			}
			wg.Done()

		}(i) 
	}
	wg.Wait() 

	sort.Ints(sampleIndices) 
	switch {
		case outputType == "indices": 
			return nil, sampleIndices
		default:
			data := mat.NewDense(r_, c, nil) 
			
			for i, ind := range sampleIndices {
				wg.Add(1) 
				go func(ind int, i int) {
					row := s.CollectOneRCByDataType(dataType, "r", ind)  
					data.SetRow(i, row)
					wg.Done() 
				}(ind, i)
			}
			wg.Wait()
			
			return data, nil
	}
}

//// TODO: refactor this method and above. 
/*
same as above, but for entire dataset
*/ 
func (s *Sampler) GatherRowsWithKeysEntireData(keys []string, collectFrequency bool) {///(*mat.Dense, []int) {
	if s.outcomeClassData == nil {
		s.outcomeClassData = make(map[string]*mat.Dense,0) 
	}
	
	s.readr.SetReaderToBeginning() 
	var x *mat.Dense
	chunksz := 0
	s.chunkIndex = 0
	for chunksz < int(s.readr.fileLength) {
		chunksz += s.MakeChunk()
		m,_ := s.GatherRowsWithKeys("chunk", keys, "data", collectFrequency) 
		switch {
		case x == nil: 
			x = m 
		default: 
			x = StackTwoMatrices(x, m) 
		}
	}
	s.sampleChunk = nil 

	s.outcomeClassData[DefaultStringSliceToString(keys, DEFAULT_DELIMITER)] = x
}

/*
*/
func (s *Sampler) CollectOneRCByDataType(dataType string, rc string, index int) []float64 {
	
	if rc != "r" && rc != "c" {
		panic(fmt.Sprintf("invalid, row (r) or column (c), have %s", rc))
	}

	var q []float64
	var stat bool
	switch {
		case dataType == "partition":
			if rc == "r" {
				q,stat = MatrixRowToFloat64Slice(s.sampleMAT, index)
			} else {
				q,stat = MatrixColToFloat64Slice(s.sampleMAT, index)
			}

		case dataType == "sample": 
			if rc == "r" {
				q,stat = MatrixRowToFloat64Slice(s.samplingData, index)
			} else {
				q,stat = MatrixColToFloat64Slice(s.samplingData, index)
			}

		case dataType == "spare": 
			if rc == "r" {
				q,stat = MatrixRowToFloat64Slice(s.spareData, index)
			} else {
				q,stat = MatrixColToFloat64Slice(s.spareData, index)
			}
			
		case dataType == "chunk": 
			if rc == "r" {
				q,stat = MatrixRowToFloat64Slice(s.sampleChunk, index)

			} else {
				q,stat = MatrixColToFloat64Slice(s.sampleChunk, index)
			}
			

		default: 
			panic("invalid data type")
	} 

	if !stat {
		panic("could not collect row")
	}
	return q

}


/*
minority outcome classes are those whose frequency divided by max is below default
*/
func (s *Sampler) ObtainMinorityClasses() []string {
	// obtain max value 
	_,maximum,_,_ := MinAndMaxOfMapStringToFloat64(s.outcomeDataFrequencies)
	// get minority classes 
	s.minorityClassKeys = make([]string,0) 
	for k,v := range s.outcomeDataFrequencies {
		q := v /maximum 
		if q <= RATIO_IMBALANCE_THRESHOLD {
			s.minorityClassKeys = append(s.minorityClassKeys,k)
		}
	}

	return s.minorityClassKeys
}

/////////////////////////////// end: outcome class ratio imbalance solution 

//////////////////////////////////////////// start: bayesian probability

/// TODO: this method will replace above. 
func (s *Sampler) BayesProbOverInputAllColumns(m *mat.Dense) map[int]map[float64]map[string]float64 { 
	// check dim 
	r,c := m.Dims()
	_,c2 := s.samplingData.Dims() 
	if c != c2 {
		panic("input data has incorrect dim.") 
	}

	/// run frequency analysis 
	s.CollectFrequencyMeanOnInputMatrix(m)  /// MODHERE
	s.FrequencyCondenseColumns(nil, "spare")
	bayesProb := make(map[int]map[float64]map[string]float64, 0)
	inputColIndices := s.InputColumnIndices()
	for i := 0; i < r; i++ {
		x,_ := MatrixRowToFloat64Slice(m, i) 
		s.BayesProbOverRow(x, bayesProb, inputColIndices, false)
	}
	return bayesProb
}


func (s *Sampler) IndicesOfColumnInQualRange(y []float64, colIndex int, dataType string, qualRange []*Pair) ([]int,[]int) {

	indices, nindices := make([]int,0), make([]int,0)
	var l float64

	for i,y_ := range y {

		switch {
		case dataType == "sample": 
			index := s.GroupOfValue(y_, s.colFrequencyCondensed[colIndex], s.colDC[colIndex]) 
			l = s.colFrequencyCondensed[colIndex][index].a.(float64)
		
		case dataType == "spare": 
			index := s.GroupOfValue(y_, s.spareColFrequencyCondensed[colIndex], s.colDC[colIndex]) 
			l = s.spareColFrequencyCondensed[colIndex][index].a.(float64)
		
		default: 
			panic("invaid data type")
		}

		if IsValueInQualRange(l, qualRange) {
			indices = append(indices, i)
		} else {
			nindices = append(nindices,i) 
		}
	}

	return nindices, indices
}

	///////////////////////////////////////////////////////////

/*
WARNING: uses spare's condense intervals
*/ 
func (s *Sampler) BayesProbOverRow(x []float64, bayesProb map[int]map[float64]map[string]float64, inputColIndices []int, encoded bool) {

	var class []float64
	switch {
		case encoded: 
			class = x
		default: 
			class,_ = s.DetermineClassOfRow(x, s.AllIndices(), "slice", "spare")
	}

	_, outputClass := s.DetermineClassOfRow(x, s.outputColumnIndices, "string", "spare")  	
	for _,y := range inputColIndices {

		_,ok := bayesProb[y] 

		if !ok {
			bayesProb[y] = make(map[float64]map[string]float64,0) 
		}

		_,ok = bayesProb[y][class[y]] 
		if !ok {
			bayesProb[y][class[y]] = make(map[string]float64,0) 
		}

		bayesProb[y][class[y]][outputClass]++
	}
}

//// //// //// //// //// //// //// //// //// //// //// //// 

/*
sums the frequency of classes that fall in the range p.a and p.b  
*/ 
func (s *Sampler) FrequencyOfClassRange(condenseType string, colIndex int, p *Pair) float64 {
	var l int 
	switch {
		case condenseType == "sample": 
			l = len(s.colFrequencyCondensed[colIndex])
		case condenseType == "spare": 
			l = len(s.spareColFrequencyCondensed[colIndex])
		default: 
			panic("panicla panicla")
	}
	
	sum := float64(0) 
	for i := 0; i < l; i++ {
		switch {
			case condenseType == "sample": 
				if s.colFrequencyCondensed[colIndex][i].a.(float64) >= p.a.(float64) && s.colFrequencyCondensed[colIndex][i].a.(float64) <= p.b.(float64) {
					sum += s.colFrequencyCondensed[colIndex][i].b.(float64) 
				}

			default: 
				if s.spareColFrequencyCondensed[colIndex][i].a.(float64) >= p.a.(float64) && s.spareColFrequencyCondensed[colIndex][i].a.(float64) <= p.b.(float64) {
					sum += s.spareColFrequencyCondensed[colIndex][i].b.(float64) 
				}
		}
	}
	return sum
}

/////////////////////////////// end: bayesian probability

//////////////////////////////  start: outcome class functions

/// TODO: this needs to be error-checked
func (s *Sampler) OutcomeClassSize() int {
	x := 1 

	for _, q := range s.outputColumnIndices {
		x = x * len(s.colFrequencyCondensed[q]) 
	}

	return x
}

/// TODO: work on condenseZoom feature for discrete-valued columns 
/// 	  columns cannot be
/// TODO 2: this is incorrect for multi-column outcomes.
///			there will be n_0 x n_1 x n_m classes,
///			n_x denotes number of classes for col. x
func (s *Sampler) CollectOutcomeClasses() []string {
	switch {
		case s.outcomeClasses == nil: 
			s.outcomeClasses = make([]string,0) 
		default: 
			return s.outcomeClasses
	}
		
	for _, q := range s.outputColumnIndices {
		x := s.colFrequencyCondensed[q] 
		for _,x_ := range x {
			s.outcomeClasses = append(s.outcomeClasses,  DefaultFloat64ToString(x_.a.(float64)))  
		}
	}

	return s.outcomeClasses 
}

/*
*/
func (s *Sampler) CheckOutcomeClassSize() bool {
	x := s.OutcomeClassSize() 
	if x >= MAX_OUTCOMESET_CLASSES {
		return false 
	}
	return true 
} 

////////////////////////////// end: outcome class functions

//////////////////////////////  start: class determination of row elements 

/// WARNING: uses condenseDataType `sample`
/*
targetThresholds : slice of pairs in which each pair corresponds to the class range specified by .a and .b
*/ 
func (s *Sampler) IsRowAClassMember(x []float64, targetIndices []int, targetThresholds []*Pair) bool {

	if len(targetIndices) == 0 || len(targetIndices) != len(targetThresholds) {
		panic("invalid target index or threshold argument") 
	}
	
	class,_ := s.DetermineClassOfRow(x, targetIndices, "slice", "sample") 

	for i,c := range class {

		if c >= targetThresholds[i].a.(float64) && c <= targetThresholds[i].b.(float64) {
			continue 
		}
		return false 
	}
	return true
}

/*
*/
func (s *Sampler) DetermineOutputClassOfRow(x []float64, condenseDataType string) string { 
	_, out := s.DetermineClassOfRow(x, s.outputColumnIndices, "string", condenseDataType)  
	return out
}

/// CAUTION: class will be the same as value if type is discrete
//// TODO: error here! uses `colFrequencyCondensed` instead
/*
*/
func (s *Sampler) DetermineClassOfRow(x []float64, targetIndices []int, outputType string, condenseDataType string) ([]float64, string) { 
	class := make([]float64,0)
	var index int 
	for _,ti := range targetIndices {
		switch {
		case condenseDataType == "sample": 
			index = s.GroupOfValue(x[ti], s.colFrequencyCondensed[ti], s.colDC[ti]) 

		case condenseDataType == "spare":
			index = s.GroupOfValue(x[ti], s.spareColFrequencyCondensed[ti], s.colDC[ti]) 
		
		default: 
			panic("invalid condense data type")
		}
		class = append(class, s.colFrequencyCondensed[ti][index].a.(float64))  
	}

	switch {
		case outputType == "string": 
			return nil, DefaultFloat64SliceToString(class, DEFAULT_DELIMITER)
		case outputType == "slice": 
			return class, ""
		default: 
			panic("invalid outputType")
	}
}

//////////////////////////////  end: class determination of row elements 

//////////////////////////////  start: encoding of matrix data values to their group (partition) encodings

/*
*/ 
func (s *Sampler) SamplingDataToGroupEncoding() *mat.Dense {
	return s.MatrixToGroupEncoding(s.samplingData, "sample")
}

/*
WARNING: uses sample's condensed intervals 
*/ 
func (s *Sampler) MatrixToGroupEncoding(m *mat.Dense, dataType string) *mat.Dense {

	r,c := m.Dims()
	m_ := mat.NewDense(r, c, nil) 
	x := IntRange(0, c) 
	var slice []float64
	wg := sync.WaitGroup{}

	for i := 0; i < r; i++ {
		wg.Add(1)
		go func(i int) {
			row,_ := MatrixRowToFloat64Slice(m, i)
			slice,_ = s.DetermineClassOfRow(row, x, "slice", dataType)
			m_.SetRow(i, slice)
			wg.Done()
		}(i) 
	}
	wg.Wait() 
	return m_ 
}

/*
*/ 
func (s *Sampler) CollectOutputForGroupEncodedMatrix(m *mat.Dense) []string {
	r,_ := m.Dims()
	output := make([]string,0)

	for i := 0; i < r; i++ {
		v,_ := MatrixRowToFloat64Slice(m, i) 
		fs := DefaultFloat64SliceToString(Float64SliceCollectByIndices(v, s.outputColumnIndices, false), DEFAULT_DELIMITER)
		output = append(output, fs) 	
	} 
	return output
}

//////////////////////////////  end: encoding of matrix data values to their group (partition) encodings

/////////////////////// start: frequency condensation

//// TODO: make condenseZoom
/*
condenseSizes : if set to nil, then uses standard
*/
func (s *Sampler) FrequencyCondenseColumns(condenseSizes map[int]int, condenseDataType string) { 

	keys := make([]int,0)
	switch {
		case condenseDataType == "sample": 
			s.colFrequencyCondensed = make(map[int][]*Pair,0)
			for k,_ := range s.colFrequency {
				keys = append(keys,k)
			}
		default:
			s.spareColFrequencyCondensed = make(map[int][]*Pair,0)
			for k,_ := range s.spareColFrequency {
				keys = append(keys,k)
			}
	}

	// gather keys
	for _,k := range keys {
		switch {
			case condenseSizes == nil:
				/// TODO: make condense zoom for specific cases
				/*
				_,ok := s.condenseZoom[k] 

				if !ok {
					s.condenseZoom[k] = DEFAULT_CONDENSE_SIZE
				} else {
					s.condenseZoom[k] *= DEFAULT_CONDENSE_SIZE
					fmt.Println("NEWCON")
				}
				s.FrequencyCondenseColumnSpecifySize(k, s.condenseZoom[k], condenseDataType)
				*/
				s.FrequencyCondenseColumnSpecifySize(k, DEFAULT_CONDENSE_SIZE, condenseDataType)

			default:
				s.FrequencyCondenseColumnSpecifySize(k, condenseSizes[k], condenseDataType) 				
		}
	}
}

/// TODO: delete or merge w/ future work
/*
*/
func (s *Sampler) FrequencyCondenseColumnAutoSize(colIndex int, condenseDataType string) { 
	_,ok := s.condenseZoom[colIndex] 

	if !ok {
		s.condenseZoom[colIndex] = DEFAULT_CONDENSE_SIZE
	} else {
		s.condenseZoom[colIndex] *= DEFAULT_CONDENSE_SIZE
	}

	s.FrequencyCondenseColumnSpecifySize(colIndex, s.condenseZoom[colIndex], condenseDataType)
}

/*
*/
func (s *Sampler) FrequencyCondenseColumnSpecifySize(colIndex int, condenseSize int, condenseDataType string) {
	switch {
		case condenseDataType == "sample": 
			s.colFrequencyCondensed[colIndex] = s.FrequencyCondense(colIndex, s.colFrequency[colIndex], s.colDC[colIndex], condenseDataType, condenseSize)
		case condenseDataType == "spare":
			s.spareColFrequencyCondensed[colIndex] = s.FrequencyCondense(colIndex, s.spareColFrequency[colIndex], s.colDC[colIndex], condenseDataType, condenseSize)
		default: 
			panic("panicla panicla")
	}
}

/*
splits up frequency counts into `condenseSize` intervals, and assigns each of those intervals a weight
assumes class encoding starts at 0
*/
func (s *Sampler) FrequencyCondense(colIndex int, freq map[float64]float64, dataType string, condenseType string, condenseSize int) []*Pair {///map[float64]float64 { 
	var partThreshold []*Pair 
	
	if condenseType == "spare" {
		// duplicate the sample data and zero out data 
		partThreshold = s.colFrequencyCondensed[colIndex]
		for i,_ := range partThreshold {
			partThreshold[i].b = 0.0 
		}
	} else {
		partThreshold = s.PartitionsForCondense(freq, dataType, condenseSize)
	}

	// find density of each partition
	var index int 
	for k,v := range freq {
		index = s.GroupOfValue(k, partThreshold, dataType)
		partThreshold[index].b = partThreshold[index].b.(float64) + v
	}

	return partThreshold
}

/*
splits up frequency counts into `condenseSize` intervals

if data type is discrete, then number of intervals equal to number of classes

return: 
- slice: element is Pair[a: interval(float), b: freq=0] 
*/
func (s *Sampler) PartitionsForCondense(freq map[float64]float64, dataType string, condenseSize int) []*Pair {

	partThreshold := make([]*Pair, 0) 
	switch {

		case dataType == "continuous":

			// iterate through keys and find min and max
			minimum, maximum := math.Inf(1), math.Inf(-1) 
			for k,_ := range freq {
				if minimum > k {
					minimum = k
				}

				if maximum < k {
					maximum = k
				}
			}

			partSize := (maximum - minimum) / float64(condenseSize)

			// make partition keys
			for i := 1; i <= condenseSize; i++ {
				partThreshold = append(partThreshold, &Pair{a: minimum + (float64(i) * partSize), b: float64(0)})
			}

		case dataType == "discrete":
			keys := CollectKVOfMapFloat64ToFloat64(freq, "key")
			sort.Float64s(keys) 
			for _,x := range keys {
				partThreshold = append(partThreshold, &Pair{a: float64(x), b: float64(0)}) 
			}

		default: 
			panic("invalid data type")
	}

	return partThreshold
}

/*
determines the index of group that v belongs to. 
assumes Pair slice b-values are in increasing order.

arguments: 
- groupThresholds: slice of pairs of thresholds, contiguous values

CAUTION: if v is greater than all thresholds, assigns v to the last group 
*/
/// TODO: GroupOfValue needs to separate b/t discrete and cont. 
func (s *Sampler) GroupOfValue(v float64, groupThresholds []*Pair, dataType string) int {
	
	switch {

	case dataType == "continuous": 
		prev := math.Inf(-1)
		for i,v_ := range groupThresholds { 
			if v >= prev && v < v_.a.(float64) {
				return i  
			}
			prev = v_.a.(float64)
		}
		return len(groupThresholds) - 1

	case dataType == "discrete": 
		for i, gt := range groupThresholds {
			if gt.a.(float64) == v {
				return i
			}
		}
		return -1

	default: 
		panic("invalid data type for group search")
	}


}

////////////////////////////// end: frequency condensation

////////////////////////////// start: main function 

func (s *Sampler) PreAnalysis(excludeClasses []string, outcomeClasses []string) {
	s.PreprocessReader("regular", excludeClasses, outcomeClasses)
	if !s.valid {
		panic("dataset is not valid to process") 
	}
	s.AnalyzeDataForTypeDC()
	s.FrequencyCondenseColumns(nil, "sample") 
	s.AnalyzeForOutcomeClassFrequencies()
	///// TODO: uncomment below to modify weights
	///s.ModifyWeightsAccordingToSpec("inverse") /// CAUTION
}



// CAUTION: method will not work in pipeline, need to revise
// CAUTION: cfbr method still needs to be removed
/// TODO: refer to Sampling #2. var needs to be stored for `pre-analysis`
///	TODO: test this method for successive runs. Runs should output consistent data
func (s *Sampler) RunSampling(excludeClasses []string, outcomeClasses []string, uniqueDraw bool, weightType string) { 
	s.PreAnalysis(excludeClasses, outcomeClasses) 
	s.DrawSamplesConcurrent(uniqueDraw, weightType, s.sampleSize) 
	s.AnalyzeSamplingData() 
	s.FixImbalancedData() /// CAUTION
	s.ClearData()
	s.CollectFrequencyMeanOnData("regular", "sample", 0)
	s.FrequencyCondenseColumns(nil, "sample")
}

/*
runs sampling in the case of minority classes
*/ 
func (s *Sampler) Presample() bool {

	if s.pre2 {
		return true 
	}

	// check for minority classes here
	mc := s.ObtainMinorityClasses()
	
	if len(mc) == 0 {
		return false 
	}

	s.GatherRowsWithKeysEntireData(mc, false)

	// calculate the number of majority key samples for each
	x := len(s.probabilityFrequencies) - len(s.minorityClassKeys)
	r := make([]int,0)
	
	for _,v := range s.outcomeClassData {
		r_,_ := v.Dims()
		r = append(r, r_) 
	}
	sort.Ints(r) 
	y := float64(r[0])
	s.numberOfPresampled = int(math.Ceil(y / RATIO_BALANCE_THRESHOLD * float64(x))) 
	return true
}

/// TODO: make sure that sampling works with different input in successive rounds
/*
updated version of RunSampling, used for datasets with minority classes, with the sample
data size less than `s.sampleSize` if `imbal2` is detected. 

return: 
- bool, true indicates reduced sampling size 
*/
func (s *Sampler) RunSampling2(excludeClasses []string, outcomeClasses []string) bool { 
	// conduct pre-analysis
	if !s.pre2 {
		s.PreAnalysis(excludeClasses, outcomeClasses)
		s.DrawSamplesConcurrent(false, "uniform", s.sampleSize)  
		//s.DrawSamples(false, "uniform", s.sampleSize)

		s.AnalyzeSamplingData()

		s.imbal2 = s.Presample()
		if !s.imbal2 {
			s.ClearData()
			return false
		}
	}
	s.pre2 = true

	// sample
	switch {
		case s.imbal2:
			s.DrawSamplesConcurrent(false, "uniform", s.numberOfPresampled)
			//s.DrawSamples(false, "uniform", s.numberOfPresampled)
		default:
			s.DrawSamplesConcurrent(false, "uniform", s.sampleSize)
			//s.DrawSamples(false, "uniform", s.sampleSize)
			
	}

	s.ClearData() 
	return true
}

/// TODO: this needs to be tested
/*
used for the imbalanced outcome class case
*/
func (s *Sampler) MakeDataForImbalanced(shuffle bool, refresh bool) {
	if !s.imbal2 {
		panic("data is balanced!")
	}

	if shuffle {
		panic("shuffle has not been implemented") 
	}

	if refresh {
		s.DrawSamplesConcurrent(false, "uniform", s.numberOfPresampled)
	}

	s.imbal2SamplingData = nil
	s.imbal2SamplingData = s.samplingData
	for _,v := range s.outcomeClassData {
		s.imbal2SamplingData = StackTwoMatrices(s.imbal2SamplingData, v)
	}

}

/*
*/ 
func (s *Sampler) ClearData() {
	s.chunkIndex = 0
	s.probabilityWeights = nil 
	s.sampleChunk = nil
	s.colMeans = make(map[int]float64,0)
	s.colFrequency = make(map[int]map[float64]float64,0)
	s.condenseZoom = make(map[int]int,0)
}

////////////////////////////// end: main function