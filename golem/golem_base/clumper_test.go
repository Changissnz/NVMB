package golem 

import (
	"fmt"
	"testing"
	//"time"
	//"reflect"
	//"strconv"
	//"encoding/csv"
	//"io"
)

func SetUpClumperInput() []*InstanceCapture {

	a := OneAnalyzer_(TESTFP_1) 
	a.CaptureAllData("exact", "ends", 2, 3) 
	a.SetThresher("one", NVMB_DELTA_NEG_CORR)
	a.GatherJudgmentValues() 

	return a.instanceInfo
}

func Test_Clump_KMeansClump(t *testing.T) {
	
	ic := SetUpClumperInput() 
	fmt.Println("NUM DATAS : ", len(ic)) 

	oc := OneClumper(ic) 
	
	fmt.Println("NUM DATAS : ", oc.data.Len()) 
	oc.KMeansClump(5, -1) 
} 


