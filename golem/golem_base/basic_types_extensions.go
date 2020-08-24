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

	var req float64 = math.Pow(float64(10), float64(-decimalRound))

	for i, x := range s1 {
		if (math.Abs(x - s2[i]) >= req) {
			return false 
		}
	}

	return true 
}

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