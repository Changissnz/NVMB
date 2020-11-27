package golem

import (
	
	"testing"
	"fmt"
	"gonum.org/v1/gonum/mat"
)

func SetUpSampler() *Sampler {
	CPARTSIZE = 5000

	fp := "./datos/other_data/creditcard.csv"
	cfbp := OneConcurrentFileBlockReader(fp)
	s := OneSampler(cfbp, 0.1)
	return s
}

/*
*/
func Test_DTree__FetchSamplingData(t *testing.T) { 

	fp := "./datos/other_data/glass.csv"
	cfbp := OneConcurrentFileBlockReader(fp)
	s := OneSampler(cfbp, 0.7)
	dt := OneDTree(s, []string{}, []string{"Type"})
	dt.Preprocess("decreasing")

	dt.FetchSamplingData(true)
	x := dt.data
	dt.FetchSamplingData(true)
	x2 := dt.data 

	if MatricesAreEqual(x,x2) {
		panic("matrices cannot be equal!")
	}

	dt.FetchSamplingData(false)
	x = dt.data
	if !MatricesAreEqual(x,x2) {
		panic("matrices have to be equal!")
	}
}

/*
*/
func Test_DTree__SplitNodeChild(t *testing.T) {
	s := SetUpSampler() 
	dt := OneDTree(s, []string{"Time"}, []string{"Class"})
	dt.Preprocess("in-order")
	dt.SplitNodeChild(dt.bdts[0].ruut, "l", EntropyCost, 0, 2)
}

/*
*/ 
func Test_DTree__MakeTree(t *testing.T) {
	s := SetUpSampler() 
	dt := OneDTree(s, []string{"Time"}, []string{"Class"})
	//dt.weightedOutput = true
	fmt.Println("before")
	dt.Preprocess("in-order")
	fmt.Println("collect freq and mean")

	/////// on new data 
	dt.MakeTree(dt.data, 0, 2, EntropyCost)
	dt.FetchSamplingData(true)
	encoded := dt.samplr.MatrixToGroupEncoding(dt.data, "spare") 
	pred,_ := dt.bdts[0].Predict(encoded, "full") 
	
	fmt.Println("AFTER PREDICTION")
	actual,_ := MatrixColToFloat64Slice(dt.data, dt.samplr.outputColumnIndices[0])
	pm := RunPredictionMetric(actual,pred, dt.bdts[0].outcomeClassValue) 
	pm.Display()
	
	fmt.Println("DISPLAYING NODE DATA") 
	dt.bdts[0].OperatePreOrderTraversal_(dt.bdts[0].DisplayNodeData)
	
		///// MC
	/*
	dt.Make(2, EntropyCost)
	dt.FetchSamplingData(true)
	encoded := dt.samplr.MatrixToGroupEncoding(dt.data, "spare") 
	pred,_ := dt.bdts[0].Predict(encoded, "full") 
	
	fmt.Println("AFTER PREDICTION")
	actual,_ := MatrixColToFloat64Slice(encoded, dt.samplr.outputColumnIndices[0])
	pm := RunPredictionMetric(actual,pred, dt.bdts[0].outcomeClassValue) 
	pm.Display()
	*/
}

func Test_DTree__RandomGreedyRangeFinder(t *testing.T) {
	s := SetUpSampler() 
	dt := OneDTree(s, []string{"Time"}, []string{"Class"})
	dt.Preprocess("in-order")

	bpm := dt.samplr.BayesProbOverInputAllColumns(dt.data)
	r,_ := dt.data.Dims()
	p, x1, x2 := dt.RandomGreedyRangeFinder(bpm[0], EntropyCost, DefaultFloat64ToString(dt.bdts[0].outcomeClassValue), 0, 0, float64(r)) 

	fmt.Println("BEST INFO")
	fmt.Println(p.a, "    ", p.b)
	fmt.Println("mismatch ", x1, x2)
	
	fmt.Println("best range for split") 
	p.a.(*Pair).a.(*Split).Display()

}

func Test_DTree__MakeEntire__GlassData(t *testing.T) {
	fp := "./datos/other_data/glass.csv"
	cfbp := OneConcurrentFileBlockReader(fp)
	s := OneSampler(cfbp, 0.7)

	dt := OneDTree(s, []string{}, []string{"Type"})
	dt.weightedOutput = true
	dt.Preprocess("decreasing")
	
	///// checking one error on mc
	dt.Make(1, StandardCost)
	dt.FetchSamplingData(true) 
	encoded := dt.samplr.MatrixToGroupEncoding(dt.data, "spare") 
	for i,_ := range dt.bdts {
		DisplayOneTreePrediction(dt,encoded,i)
	}
} 

func DisplayOneTreePrediction(dt *DTree, encoded *mat.Dense, index int) {
	fmt.Println("AFTER PREDICTION ", index)
	pred,_ := dt.bdts[index].Predict(encoded, "full") 

	actual,_ := MatrixColToFloat64Slice(encoded, dt.samplr.outputColumnIndices[0])
	pm := RunPredictionMetric(actual,pred, dt.bdts[index].outcomeClassValue) 
	pm.Display()

	/// uncomment below to display tree contents
	/*
	dt.bdts[1].InOrderTraversalForData_(encoded, true)
	fmt.Println("DISPLAYING NODE DATA") 
	dt.bdts[1].OperatePreOrderTraversal_(dt.bdts[index].DisplayNodeData)
	*/
}

/*
*/
func Test_DTree__Make__DataGlass(t *testing.T) {
	fmt.Println("MAKING TREE ON GLASS DATA")
	fp := "./datos/other_data/glass.csv"

	numTrees := 10
	for i := 0; i < numTrees; i++ {
		fmt.Println("making one tree ", i)
		cfbp := OneConcurrentFileBlockReader(fp)
		s := OneSampler(cfbp, 0.7)
		dt := OneDTree(s, []string{}, []string{"Type"})
		//dt.weightedOutput = true
		dt.Preprocess("increasing")
	
		dt.Make(10, StandardCost)
		dt.FetchSamplingData(true) 
		pm := dt.PredictMC()
		pm.Display()
	}
}

/// BM
func Test_DTree_Make__DataDrug200(t *testing.T) {
	fp := "./datos/other_data/drug200.csv" 
	cfbp := OneConcurrentFileBlockReader(fp)
	s := OneSampler(cfbp, 1.0)
	///s.RunSampling2([]string{}, []string{"Drug"})  

	dt := OneDTree(s, []string{}, []string{"Drug"})
	//dt.weightedOutput = true
	dt.Preprocess("increasing")// "in-order")

	dt.Make(5, StandardCost)
	dt.FetchSamplingData(true) 
	pm := dt.PredictMC()
	pm.Display()
}

/// Iris
func Test_DTree_Make__DataIris(t *testing.T) {

	//// analysis of data here
	fp := "./datos/other_data/IRIS.csv"
	cfbp := OneConcurrentFileBlockReader(fp)
	s := OneSampler(cfbp, 0.7)

	//// TODO: refactor to method `s.DisplaySamplingInfo()`
	/*
	s.RunSampling2([]string{}, []string{"species"})//, true, "outcome probability")
	fmt.Println("imbalanced: ", s.imbal2)
	fmt.Println("minority class keys: ", s.minorityClassKeys) 
	fmt.Println("column types") 
	fmt.Println(s.colDC)
	*/
	////----------
	dt := OneDTree(s, []string{}, []string{"species"})
	dt.Preprocess("increasing")
	dt.Make(11, StandardCost) 
	dt.FetchSamplingData(true) 
	pm := dt.PredictMC()
	pm.Display()
	
}