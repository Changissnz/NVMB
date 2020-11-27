/*
where type conversion occurs, for 
*/ 

package golem

import (
	"fmt"
	"strconv"
	"strings"
	"math"
	"math/rand"
	"gonum.org/v1/gonum/mat"
	"sort"
	"sync"
)

const DEFAULT_DELIMITER rune = ','
const DATA_DEFAULT_DELIMITER rune = ' ' 
const NULL_VALUE_REPR string = "NULLOS" 

/*
outputs random boolean 
*/ 
func RandomBoolean() float32 {

	if (rand.Float32() >= 0.5) {
		return 1.0 
	} else {
		return 0.0 
	}
}

/*
outputs random integer  
*/ 
func RandomInt(minumum int32, maximum int32) int32 {
	
	if (maximum < minumum) {
		panic("max cannot be less than min")
	}

	difos := maximum - minumum 
	d := rand.Int31n(difos)
	return minumum + d 
}

func RandomIntS(minimum int, maximum int) int {
	return int(RandomInt(int32(minimum), int32(maximum))) 
} 

/*
random integer slice 
*/ 
func RandomIntSlice(minumum int32, maximum int32, size int) []int32 {
	f := make([]int32,0)

	for i := 0; i < size; i++ {
		f = append(f,RandomInt(minumum,maximum)) 
	}
	return f 
}

/*
random integer slice w/ no duplicate elements

CAUTION: no argument error-checking
*/ 
func RandomUniqueIntSlice(minumum int32, maximum int32, size int, exclude []int32) []int32 {

	s := make([]int32,0)
	for {

		if (len(s) == size) {
			break 
		}

		q := RandomInt(minumum,maximum)
		r := append(s, exclude...) // not efficient 

		if (!IntInSlice(r, q)) {
			s = append(s,q) 
		}
	}
	return s 
}

/*
output random float 
*/ 
func RandomFloat(minumum float32, maximum float32) float32 {

	if (minumum > maximum) {
		panic("min cannot be greater than max")
	}

	var diff float32 = maximum - minumum 
	var output float32 = minumum
	var coeff float32 = rand.Float32() 
	
	return output + (diff * coeff) 
}

func RandomFloat64(minumum float64, maximum float64) float64 {

	if (minumum > maximum) {
		panic("min cannot be greater than max")
	}

	var diff float64 = maximum - minumum 
	var output float64 = minumum
	var coeff float64 = rand.Float64() 
	
	return output + (diff * coeff) 
}


/*
random float slice
*/ 
func RandomFloatSlice(minumum float32, maximum float32, size int) []float32 {
	f := make([]float32,0)

	for i := 0; i < size; i++ {
		f = append(f,RandomFloat(minumum,maximum)) 
	}
	return f 
}


/*
*/ 
func UnionFloat64ProbabilityMaps(m1 map[int]float64, m2 map[int]float64) map[int]float64 {
	if m1 == nil {
		return m2
	} 

	if m2 == nil {
		return m1
	}

	for k, v := range m2 {
		m1[k] += v
	}

	return m1
}

/*
*/
func RandomProbabilityMap(minumum float32, maximum float32, numColumns int) map[int]float32 {
	f := make(map[int]float32, 0)  

	for i := 0; i < numColumns; i++ {
		f[i] = RandomFloat(minumum,maximum)
	}
	return f 
} 

/*
*/
func DisplayProbabilityMap(m map[int]float32) {
	for k,v := range m {
		fmt.Println(fmt.Sprintf("%d\t%f", k,v)) 
	}
}

/*
outputs an integer in [minumum,maximum]
*/ 
func RandomIntInRange(minumum int, maximum int) int {

	diff := maximum - minumum
	if diff == 0 {
		return -1 
	}

	toAdd := rand.Intn(diff) 
	return minumum + toAdd 
}

/*
outputs a random slice element and its index 
*/
func RandomElementInIntSlice(input []int) (int,int) {
	r:= RandomIntInRange(0, len(input))
	return input[r], r 
}

/*
*/
func RandomElementInStringSlice(s []string) (string,int) {

	qi := RandomIntInRange(0, len(s))

	if qi == -1 {
		return "", -1 
	}
	return s[qi], qi 
}

/*
*/
func RandomElementInSlice(s []interface{}) (interface{}, int) {
	qi := RandomIntInRange(0, len(s))
	return s[qi], qi
}

/*
mode := r(replace)|nr(non-replace) 
*/ 
func RandomElementsInIntSlice(input []int, numElements int, mode string) []int {
	
	if (numElements > len(input)) {
		panic("number of wanted elements greater than input size")
	}

	if (mode != "r" && mode != "nr") {
		s := fmt.Sprintf("invalid mode %s, want `r` or `nr`", mode)
		panic(s)
	}

	output := make([]int, 0) 
	for i := 0; i < numElements; i++ {
		x, j := RandomElementInIntSlice(input)
		output = append(output, x) 

		if (mode == "nr") {
			input = append(input[:j], input[j+1:]...) 
		}
	}

	return output 
}

/*
default formats float to 5 places 
*/ 
func DefaultFloatToString(x float32) string {
	return strconv.FormatFloat(float64(x), 'E', 5, 32) 
}

/*
*/
func DefaultFloat64ToString(x float64) string {
	return strconv.FormatFloat(float64(x), 'E', 5, 64) 
}

/*
outputs a float32 representation of string x 
*/ 
func DefaultStringToFloat(x string) float32 {

	f, err := strconv.ParseFloat(x, 32)  
	if (err != nil) {
		s := fmt.Sprintf("invalid float as string : %s", x) 
		panic(s)
	}
	return float32(f)  
}

/*
*/
func DefaultStringToFloat64(x string) (float64, bool) {

	f, err := strconv.ParseFloat(x, 64)  
	if (err != nil) {
		return -1.0, false
	}
	return float64(f), true
}

/*
*/ 
func BoolToFloat64(b bool) float64 {
	if b {
		return float64(1)
	} else {
		return float64(0)
	}
}

/*
outputs string repr. of integer slice 
*/ 
func DefaultIntSliceToString(x []int, delimiter rune) string {
	if (len(x) == 0) {
		return ""
	}
	var y string = ""
	for _, x_ := range x {
		y += strconv.Itoa(x_) + string(DEFAULT_DELIMITER)
	}

	return y[:len(y) -1] 
}

/*
converts int32 slice to int slice 
*/ 
func Int32SliceToIntSlice(x []int32) []int{

	output := make([]int, 0)
	for _, y := range x {
		output = append(output, int(y))
	}
	return output 
}

/*
*/
func StripString(x string, whitespace rune) string {
	// strip forward
	index := -1 
	for i, x_ := range x {
		if x_ != whitespace {
			index = i
			break 
		}
	}

	if index != -1 {
		x = x[index:]
	}

	// strip backward
	index = -1
	l := len(x) 
	for i := l - 1; i >= 0;i-- {
		if rune(x[i]) != whitespace {
			index = i
			break  
		}
	}

	if index != -1 {
		x = x[:index + 1]
	}
	
	return x
} 

/*
*/
func ParseOneElementFromString(x string, delimiter rune) (string,string) {
	l := len(x) 
	for i := 0; i < l; i++ {
		if (rune(x[i]) == delimiter) {
			return x[:i], x[i+1:]
		} 
	}
	return x, ""
}

/*
*/
func DefaultStringToIntSlice(x string, delimiter rune) ([]int, bool) {

	if (len(x) == 0) {
		return make([]int,0), true
	}
	
	output := make([]int, 0) 
	x = StripString(x, delimiter)

	for {
		s1,s2 := ParseOneElementFromString(x, delimiter) 
		q1,q2 := strconv.Atoi(s1)


		if (q2 != nil) {
			return make([]int,0), false 
		}

		output = append(output, q1)

		// done 
		x = s2 
		if (x == "") {
			break
		}

	}

	return output, true  
}

/*
*/
func DefaultStringToFloat64Slice(x string, delimiter rune) ([]float64, bool) {

	if (len(x) == 0) {
		return make([]float64,0), true
	}

	x = StripString(x, delimiter)
	output := make([]float64, 0) 
	for {
		s1,s2 := ParseOneElementFromString(x, delimiter) 
		q1, q2 := DefaultStringToFloat64(s1) // TODO: STRIP 

		if (!q2) {
			return make([]float64,0), false 
		}

		output = append(output, q1)

		// done 
		x = s2 
		if (x == "") {
			break
		}

	}

	return output, true 
}

func DefaultFloat64SliceToString(x []float64, delimiter rune) string {

	if (len(x) == 0) {
		return ""
	}

	output := "" 
	for _, x_ := range x {

		s := DefaultFloat64ToString(x_)
		output += s + ","
	}

	output = output[:len(output) -1] 
	return output
} 

/*
CAUTION: performs no argument error-checking
*/
func DefaultStringSliceToFloat64Slice(x []string) []float64 {
	output := make([]float64,0) 
	for _,x_ := range x {
		c, succ := DefaultStringToFloat64(x_)
		if (!succ) {
			panic("could not convert")
		}
		output = append(output, c)
	}
	return output 
} 

/*
*/
func DefaultStringToStringSlice(x string, delimiter rune) ([]string, bool) {

	if (len(x) == 0) {
		return make([]string,0), true
	}

	x = StripString(x, delimiter)
	output := make([]string, 0) 
	for {
		s1,s2 := ParseOneElementFromString(x, delimiter) 		
		output = append(output, s1)

		// done 
		x = s2 
		if (x == "") {
			break
		}
	}

	return output, true 
}

/*
*/ 
func DefaultStringSliceToString(x []string, delimiter rune) string {

	output := "" 

	if len(x) == 0 {
		return output
	}

	for _,x_ := range x {
		output += x_ + string(delimiter) 
	}

	return output[:len(output) -1] 
} 

/*
*/
func DataValueSliceToFloat64Slice(x []DataValue) []float64 {
	output := make([]float64, 0) 

	for _, y := range x {
		output = append(output, y.(float64))
	}
	return output 
}

/*
*/
func DataValueSliceToStringSlice(x []DataValue) []string{
	output := make([]string, 0) 

	for _, y := range x {
		output = append(output, y.(string))
	}
	return output 
}

/*
return: string,int,float,vector, or ?
*/ 
 func DeduceBasicStringType(s string) string{
	_, err := strconv.ParseInt(s, 10, 64)
	if (err == nil) {
		return "int" 
	}

	_, err2 := strconv.ParseFloat(s, 64)
	if (err2 == nil) {
		return "float"
	}

	if (s == "NULLOS" || StripString(s, ' ') == "") {
		return "?"
	}
	_, stat := DefaultStringToFloat64Slice(s, ' ')
	if (stat) {
		return "vector" 
	}

	return "string"
}

/*
determines if two string slices are identical 
*/ 
func EqualStringSlices(s1 []string, s2 []string) bool { 

	if (len(s1) != len(s2)) {
		return false 
	}

	for i, x := range s1 {
		if (x != s2[i]) {
			return false 
		}
	}

	return true 
}

// TODO : refactor with above
/*
determines if two float64 slices are identical 
*/ 
func EqualFloat64Slices(s1 []float64, s2 []float64, decimalRound int) bool {
	if (len(s1) != len(s2)) {
		return false 
	}

	if decimalRound < 0 {
		panic("decimal round must be non-neg.")
	}

	for i, x := range s1 {
		if !EqualFloat64(x,s2[i], decimalRound) {
			return false 
		}
	}

	return true 
}

func EqualFloat64(f1 float64, f2 float64, decimalRound int) bool {
	var req float64 = math.Pow(float64(10), float64(-decimalRound))

	if (math.Abs(f1 - f2) >= req) {
		return false 
	} 

	return true 
}

/// TODO: rename as Int32...
/*
determines if int in slice 
*/ 
func IntInSlice(s []int32, n int32) bool {
	for _, q := range s {
		if (q == n) {
			return true 
		}
	}
	return false 
}

func IntSInSlice(s []int, n int) bool {
	for _, q := range s {
		if (q == n) {
			return true 
		}
	}
	return false
}

func Float64InSlice(s []float64, n float64) bool {
	for _, q := range s {
		if (q == n) {
			return true 
		}
	}
	return false
}

/*
determines if int in slice 
*/ 
func StringIndexInSlice(s []string, n string) int {
	for i, q := range s {
		if (q == n) {
			return i 
		}
	}
	return -1 
}

/*
Assume s and s2 are slice repr. of sets. 
Outputs the index of each string in s2 in s. 
*/ 
func StringSliceIndicesInStringSlice(s1 []string, s2 []string) []int {
	indices := make([]int,0)
	for _,s := range s2 { 
		indices = append(indices,StringIndexInSlice(s1, s))
	}
	return indices 
}

/*
outputs starting indices of all occurrences of substring in target
*/ 
func SubstringInStringAll(target string, substring string) []int {
	output := make([]int, 0)
	l := len(substring)

	targ := target 
	offset := 0 
	for {
		index := strings.Index(targ, substring)
		
		if index == -1 {
			break 
		}

		output = append(output, index + offset) 
		offset += index + l
		targ = targ[index + l:]
	}

	return output
}

func StringSliceDifference(x []string, diff []string) []string {
	output := make([]string,0)
	for _,x_ := range x {
		if StringIndexInSlice(diff,x_) != -1 {
			continue
		}
		output = append(output, x_) 
	}

	return output
}

/*
outputs a 2-d "matrix" (slice of string slices)
*/ 
func TwoDimStringMatrix(i, j int) [][]string {

	r := make([]string, j) 
	output := make([][]string, i)
	for p := 0; p < i; p++ {
		output[p] = r 
	}
	return output 
}

/*
CAUTION no error-handling 
*/ 
func GetTwoDStringMatrixChunk(x [][]string, si int, ei int) [][]string {

	m := make([][]string,0) 

	for i := si; i < ei; i++ {
		m = append(m, x[i])
	}

	return m 
}

/*
subtracts f2 - f1 
*/ 
func SubtractFloat64Slices(f1,f2 []float64) []float64 {

	if len(f1) == 0 || len(f1) != len(f2) {
		panic("slice sizes have to match for subtraction")
	}

	m1 := mat.NewDense(1, len(f1), f1)
	m2 := mat.NewDense(1, len(f2), f2)

	output := mat.NewDense(1, len(f1), nil)
	output.Sub(m2, m1) 
	
	return MatrixTo1DFloat(output)
}

/*
boolean function applied on f1 and f2 
*/ 
func FBApplyFloat64Slices(f1, f2 []float64, q func(float64, float64) bool) []float64 {

	if len(f1) != len(f2) || len(f1) == 0 {
		return nil 
	}

	output := make([]float64,0) 
	// TODO : concurrency here 
	for i, x := range f1 {
		o := q(x, f2[i])
		output = append(output,BoolToFloat64(o))   
	}
	return output 
}

func Float64SliceCollectByIndices(x []float64, indices []int, exclude bool) []float64 {
	var output []float64

	switch {
		case exclude == true: 
			for i,x_ := range x {
				if IntSInSlice(indices,i) {
					continue
				}
				output = append(output, x_) 
			}

		default: 
			output = make([]float64,len(indices))
			for k, i := range indices {
				output[k] = x[i]
			}
	}

	return output
}

/*
*/ 
func IntRange(start int, end int) []int {

	output := make([]int,0)  
	for i := start; i < end; i++ {
		output = append(output, i)
	}
	return output 
}

func IntSliceCollectByIndices(x []int, indices []int, exclude bool) []int {
	var output []int

	switch {
		case exclude == true: 
			for i,x_ := range x {
				if IntSInSlice(indices,i) {
					continue
				}
				output = append(output, x_) 
			}

		default: 
			output = make([]int,len(indices))
			for k, i := range indices {
				output[k] = x[i]
			}
	}

	return output
}

func IntSliceDifference(x []int, diff []int) []int {
	output := make([]int,0)
	for _,x_ := range x {
		if IntSInSlice(diff,x_) {
			continue
		}
		output = append(output, x_) 
	}

	return output
}


/*
*/ 
func StringSliceCollectByIndices(x []string, indices []int, exclude bool) []string {
	var output []string

	switch {
		case exclude == true: 
			for i,x_ := range x {
				if IntSInSlice(indices,i) {
					continue
				}
				output = append(output, x_) 
			}

		default: 
			output = make([]string,len(indices))
			for k, i := range indices {
				output[k] = x[i]
			}
	}

	return output
}

/*
numerical function applied to two rows f1,f2 
*/ 
func FNApplyFloat64Slices(f1, f2 []float64, q func(float64, float64) float64) []float64 {
	if len(f1) != len(f2) || len(f1) == 0 {
		return nil 
	}

	output := make([]float64,0) 
	// TODO : concurrency here 
	for i, x := range f1 {
		o := q(x, f2[i])
		output = append(output,o)   
	}
	return output 
}

func FNApplyAccumulateFloat64Slices(f1, f2 []float64, q func(float64, float64) float64, accumulator func(float64, float64) float64) float64 {
	if len(f1) != len(f2) || len(f1) == 0 {
		panic("cannot apply-accumulate on empty slices")
	}

	output := float64(0)
	for i, x := range f1 {
		o := q(x, f2[i])
		output = accumulator(output, o) 
	}

	return output
}

func FNApplyFloat64Slice(f1 []float64, f2 float64, q func(float64,float64) float64) []float64 {
	if len(f1) == 0 {
		return nil 
	}

	output := make([]float64,0) 
	// TODO : concurrency here 
	for _, x := range f1 {
		o := q(x, f2)
		output = append(output,o)   
	}
	return output
}

func SumIterable(r1 []float64) float64 {
	x := float64(0)

	for _, r := range r1 {
		x += r
	}

	return x
}

func MeanIterable(r1 []float64) float64 {
	l := len(r1) 

	// case: 0-div
	if (l == 0) {
		return 0 
	}

	s := SumIterable(r1) 
	return s / float64(l) 
}

func SortStringSliceAscending(s []string) []string {

	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})

	return s
}

/// CAUTION: inefficient algorithm
/// CAUTION: does not perform argument-checking on ranges
/// CAUTION: type-casting used b/t float 32/64
/*
/// description 
if x2 in [ranges[x1], ranges[x1 + 1]), choose it
*/
func DrawUniqueIndicesFromRanges(ranges []float64, numSamples int) []int {

	b := OneBasicSet() 
	var targetSize int
	var exclude bool = false 

	// determine if choose indices to exclude or include 
	if float64(numSamples) > float64(float64(len(ranges)) / 2.0) {
		targetSize = len(ranges) - numSamples
		exclude = true 
	} else {
		targetSize = numSamples 
	}

	// choose indices
	minimum, maximum := float32(0), float32(ranges[len(ranges) -1])
	for {
		if b.Len() == targetSize {
			break 
		}
		
		rf := RandomFloat(minimum, maximum)
		index := IndexOfSampleGivenValue(ranges, rf)
		b.AddOne(strconv.Itoa(index))
		fmt.Println("LEN B ", b.Len())
	}

	// output result based on exclude or include 
	output := make([]int,0) 
	switch {
		case exclude:
			x := len(ranges) 
			for i := 0; i < x; i++ {
				a := strconv.Itoa(i)
				if !b.DoesExist(a) {
					output = append(output, i)
				} 
			}

		default: 
			for k,_ := range b.datos {
				a,_ := strconv.Atoi(k)
				output = append(output, a)   
			}
	}
	return output

}

///////////////////////////////// concurrent version of above
var DEFAULT_NUMBER_OF_THREADS int = 10

func DrawIndicesFromRangesConcurrent(ranges []float64, numSamples int, unique bool) []int {

	b := OneBasicSet()
	b2 := []int{} 

	minimum, maximum := float32(0), float32(ranges[len(ranges) -1])
	partialThreadSize := int(math.Mod(float64(numSamples), float64(DEFAULT_NUMBER_OF_THREADS))) 
	wholeThreadSize := numSamples / DEFAULT_NUMBER_OF_THREADS
	numThreads := DEFAULT_NUMBER_OF_THREADS 

	if partialThreadSize >= 0 {
		numThreads++ 	
	}

	wg := sync.WaitGroup{}
	mux := sync.Mutex{}
	for i := 0; i < numThreads; i++ {
				
		wg.Add(1)
		go func(i int) {
			var sz int
			if i < numThreads - 1 {
				sz = wholeThreadSize 
			} else {
				sz = partialThreadSize
			}

			for sz > 0 {
				rf := RandomFloat(minimum, maximum)
				index := IndexOfSampleGivenValue(ranges, rf)
	
				mux.Lock()
				switch {
				case unique: 
					if !b.DoesExist(strconv.Itoa(index)) { 
						b.AddOne(strconv.Itoa(index))
						sz--
					}
				default: 
					b2 = append(b2, index)
					sz--
				}
	
				mux.Unlock() 
			}
			wg.Done() 
		}(i) 	
	}
	wg.Wait() 

	output := make([]int,0)
	switch {
		case unique: 
			for k,_ := range b.datos {
				a,_ := strconv.Atoi(k)
				output = append(output, a)   
			}
		default: 
			output = b2
	}

	return output 
}

/// TODO: look into refactoring
/*
*/ 
func DrawIntegersGivenWeights(valueRange int, numSamples int, weights []float64) []int { 

	if len(weights) != valueRange {
		panic("number of weights wrong.")
	}

	// construct weight range
	valueWeightSlice := make([]float64, 0)
	weight := float64(0) 

	for i := 0; i < valueRange; i++ {
		weight += weights[i]
		valueWeightSlice = append(valueWeightSlice, weight) 
	}

	return DrawIndicesFromRangesConcurrent(valueWeightSlice, numSamples, false)
} 


///////////////////////////////// end: concurrent version of above 

/// CAUTION: type-casting used.
/*
iterates through ranges and chooses the index that v is in
*/ 
func IndexOfSampleGivenValue(ranges []float64, v float32) int {
	prev := float32(0) 
	for i,r := range ranges {
		if v >= prev && v <= float32(r) {
			return i
		}
		prev = float32(r)
	}
	return -1
}


//////////////////////////////// start: map functions 

	//////////////////////////////// TODO: refactor these methods
/*
outputs minimum value, maximum value, minimum key, maximum key 
*/
func MinAndMaxOfMapStringToFloat64(x map[string]float64) (float64,float64, string, string) {
	
	if len(x) == 0 {
		panic("no min and max on empty map")
	}
	
	minimum, maximum := math.Inf(1), math.Inf(-1)
	var minKey, maxKey string 
	for k, v := range x {
		if v < minimum {
			minimum = v
			minKey = k 
		}

		if v > maximum {
			maximum = v
			maxKey = k
		}
	}

	return minimum, maximum, minKey, maxKey
}

func MinAndMaxOfMapFloat64ToFloat64(x map[float64]float64, keyOrVal string) (float64, float64) {
	
	if len(x) == 0 {
		panic("no min and max on empty map")
	}
	
	minimum, maximum := math.Inf(1), math.Inf(-1)
	for k, v := range x {
		switch {
		case keyOrVal == "key": 
			if k < minimum {
				minimum = k
			}

			if k > maximum {
				maximum = k
			}
		
		case keyOrVal == "value": 
			if v < minimum {
				minimum = v
			}

			if v > maximum {
				maximum = v
			}
		}
	}

	return minimum, maximum
}

func MinAndMaxKeyByValueOfMapFloat64ToFloat64(x map[float64]float64) (float64, float64) { 
	mini, maxi := math.Inf(1), math.Inf(-1)
	var minKey, maxKey float64
	for k,v := range x {
		if v < mini {
			mini = v
			minKey = k
		}

		if v > maxi {
			maxi = v
			maxKey = k
		}
	}

	return minKey, maxKey

}

//// TODO: in the case of ties, max and min will be the last element
func MinAndMaxKeyByValueOfMapStringToFloat64(x map[string]float64) (string, string) { 
	mini, maxi := math.Inf(1), math.Inf(-1)
	var minKey, maxKey string
	for k,v := range x {
		if v < mini {
			mini = v
			minKey = k
		}

		if v > maxi {
			maxi = v
			maxKey = k
		}
	}

	return minKey, maxKey
}

	///////////////////////////////////////////////////// END TODO: refactor these methods 


/*
*/ 
func CollectValuesOfMapIntToPairSlice(x map[int][]*Pair) []*Pair { 
	output := make([]*Pair,0) 
	for _,v := range x {
		output = append(output, v...) 
	}
	return output
}

func CollectKVOfMapFloat64ToFloat64(x map[float64]float64, kv string) []float64 {
	output := make([]float64,0)
	switch {
	case kv == "key": 
		for v,_ := range x {
			output = append(output, v)
		}
	case kv == "value": 
		for _,v := range x {
			output = append(output, v)
		}
	default: 
		panic("must be key or value")
	} 

	return output
} 

func PairSliceToFloat64Slice(p []*Pair, index int) []float64 {
	if index != 1 && index != 0 {
		panic("invalid pair index")
	}
	
	output := make([]float64,0) 
	for _, p_ := range p {
		switch {
			case index == 0: 
				output = append(output, p_.a.(float64)) 
			default: 
				output = append(output, p_.b.(float64)) 
		}
	}

	return output
}

func IndexOfFloatPairInPairSlice(p []*Pair, p2 *Pair) int {

	for i, k := range p {
		if EqualFloat64(k.a.(float64), p2.a.(float64), 5) && EqualFloat64(k.b.(float64), p2.b.(float64), 5) {
			return i
		}
	}

	return -1
}

/*
complement of a subset given chosen indices
*/
func BayesianIdentityOfIndexSubset(indexRange *Pair, wantedIndexes []int, sortedSlice bool) []int{
	
	if !sortedSlice {
		sort.Ints(wantedIndexes) 
	}
	
	output := make([]int,0) 

	// type deduction for float and int
	var mini, maxi int 
	switch indexRange.a.(type){
	case int: 
		mini, maxi = indexRange.a.(int), indexRange.b.(int) 
	case float64: 
		mini, maxi = int(indexRange.a.(float64)), int(indexRange.b.(float64))
	}


	for i := mini; i < maxi; i++ {
		if len(wantedIndexes) == 0 {
			output = append(output, i)
			continue
		}

		if i == wantedIndexes[0] {
			wantedIndexes = wantedIndexes[1:] 
			continue
		}
		output = append(output, i) 
	}
	return output
}

func BayesianIdentityOfFloatSliceSubset(floatRange *Pair, usedIndexRanges []*Pair) []*Pair {

	output := make([]*Pair,0)
	switch {
		case len(usedIndexRanges) == 0:
			output = append(output, floatRange) 

		default:
			head := floatRange.a.(float64)
			for i, k := range usedIndexRanges {

				if k.b.(float64) < k.a.(float64) {
					panic(fmt.Sprintf("invalid index range at %d", i)) 
				}

				if head < k.a.(float64) {
					p := &Pair{a: head, b: k.a.(float64)} 
					output = append(output, p)
				}
				head = k.b.(float64) 
			}

			if head < floatRange.b.(float64) {
				p := &Pair{a: head, b: floatRange.b.(float64)} 
				output = append(output, p) 
			}
	}
	return output
}

/*
return: 
- y|x, x, y|!x, !x
*/ 
func ReadBayesProbMap(bayesProbMap map[float64]map[string]float64, qualRange []*Pair, wantedOutcome string) (float64,float64,float64,float64) {
	yx, x, notX, notYX := 0.0,0.0,0.0,0.0

	for k,v := range bayesProbMap {
		if IsValueInQualRange(k, qualRange) {
			for k2,v2 := range v {
				if k2 == wantedOutcome {
					yx += v2
				}
				x += v2
			}
		} else {
			for k2,v2 := range v {
				if k2 == wantedOutcome {
					notYX += v2
				}
				notX += v2
			}
		}
	}
	return yx,x,notYX,notX
}

func ReadBayesProbMapDiscrete(bayesProbMap map[float64]map[string]float64, qualRange []float64, wantedOutcome string) (float64, float64, float64, float64) {
	yx, x, notX, notYX := 0.0,0.0,0.0,0.0

	for k,v := range bayesProbMap {
		if Float64InSlice(qualRange, k) {
			for k2,v2 := range v {
				if k2 == wantedOutcome {
					yx += v2
				}
				x += v2
			}
		} else {
			for k2,v2 := range v {
				if k2 == wantedOutcome {
					notYX += v2
				}
				notX += v2
			}
		}
	}
	return yx,x,notYX,notX
}

/// CAUTION: check inequality statement 
func IsValueInQualRange(v float64, qualRange []*Pair) bool {
	for _,q := range qualRange {
		if v >= q.a.(float64) && v <= q.b.(float64) {
			return true
		}
	}
	return false 
}

func SortKeysMapFloat64ToFloat64(x map[float64]float64) []float64 {
	k := make([]float64,0) 
	for x1,_ := range x {
		k = append(k, x1)
	}

	sort.Float64s(k) 
	return k
}

func SortKeysByValueMapStringToFloat64(x map[string]float64, order string) []*Pair {
	k := make([]*Pair, 0) 
	for x1,x2 := range x {
		k = append(k, &Pair{x1,x2}) 
	}

	switch {
	case order == "increasing": 
		sort.Slice(k, func (i,j int) bool {
			return k[i].b.(float64) < k[j].b.(float64) 
		})
	case order == "decreasing": 
		sort.Slice(k, func (i,j int) bool {
			return k[i].b.(float64) > k[j].b.(float64) 
		})
	default: 
		panic("argument is not valid order")
	}

	return k
}

/*
compares the elements of s1 and s2 at each index by a function
*/
func CountOfMatchTwoFloat64Slices(s1 []float64, s2 []float64) {
}

/*
TODO: refactor with FNApply... method above
*/ 
func CountOfMatchTwoStringSlices(s1, s2 []string) int {

	if len(s1) != len(s2) { 
		panic("slices must be equal length") 
	}

	c := 0 
	for i, s := range s1 {
		if s == s2[i] {
			c++ 
		}
	}
	return c 
}

//////////////////////////////// end: map functions 

//////////////////////////////// start: interface slice conversion functions 

func PairSliceToInterfaceSlice(p []*Pair) []interface{}{  
	output := make([]interface{},0) 

	for _,p_ := range p {
		output = append(output, p_) 
	}
	return output
}

func InterfaceSliceToPairSlice(p []interface{}) []*Pair {  
	output := make([]*Pair,0) 

	for _,p_ := range p {
		output = append(output, p_.(*Pair)) 
	}
	return output
}

func InterfaceSliceToIntSlice(p []interface{}) []int {
	output := make([]int,0) 
	for _,p_ := range p {
		switch px := p_.(type) {

		case int: 
			output = append(output, px)

		case float64: 
			output = append(output, int(px))
		default: 
			panic("invalid interface type for int slice conversion")
		}
	}
	return output
}

//////////////////////////////// end: interface slice conversion functions 
