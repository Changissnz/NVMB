/*
this file will automatically label the input variables of a NVMB timestamp as good 1, bad -1, or neutral 0
*/ 

package golem

/*
auto-labeller for NVMB data file 
*/ 
type AutoLabelNVMB struct {
	games []string
}

///////////////// TODO: below

/*
*/ 
func (aln *AutoLabelNVMB) TernaryLabellerSelf() {
}

var DEFAULT_MULTICLASS_LABELS []string = []string{""}

/// TODO:
/*
used to consider just one file
*/
func (aln *AutoLabelNVMB) MulticlassLabellerSelf(ic *InstanceCapture) { 
	// iterate through judgment values with weights 
}

/// TODO:
/*
use to consider entire game 
*/ 
func (aln *AutoLabelNVMB) MulticlassLabellerSet(ref []string, targets []string) {

}

/// TODO: environment is float[node variables] 
/*
*/
func (aln *AutoLabelNVMB) EnvironmentToLabel() {

}



//////


type TernaryRubricOperator struct {
	columnLabels []string
	weights []float64
}

/*
*/ 
func OneTernaryRubricOperator(cl []string) *TernaryRubricOperator {
	return &TernaryRubricOperator{columnLabels: cl} 
}

/*
if key not found in m, assigns weight of 1
*/
func (tro *TernaryRubricOperator) AssignWeights(m map[string]float64) {
	tro.weights = make([]float64,0)
	
	for _, cl := range tro.columnLabels {
		v,ok := m[cl] 
		if !ok {
			tro.weights = append(tro.weights, 1)
		} else {
			tro.weights = append(tro.weights, v)
		}
	}
}

//// TODO: make different modes for prior, prior + post, post
/*
use for stand-alone cases
*/ 
func (tro *TernaryRubricOperator) SingleInstance_RubricOne(ic *InstanceCapture) *Ternary { 
	mul := func(x1 float64, x2 float64) float64 {
		return x1 * x2
	}
	cum := func(x1 float64, x2 float64) float64 {
		return x1 + x2
	}

	x := FNApplyAccumulateFloat64Slices(ic.judgmentValues, tro.weights, mul, cum)

	/// TODO: add threshold value instead of 0?
	///			OR epsilon-radius?  
	var ternary *Ternary
	
	switch {
	case x > 0:
		ternary.value = 1
	case x == 0: 
		ternary.value = 0
	default: 
		ternary.value = -1
	}

	return ternary
}


///////////////// END TODO: below


