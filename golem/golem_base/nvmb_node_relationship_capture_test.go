package golem 

import (
	"testing" 
)

func Test_ParseDeltaString(t *testing.T) {

	target := "sadfasdfsd__delta__351_3435"
	d,_,_ := ParseDeltaString(target) 
	
	if d != "sadfasdfsd__delta" {
		panic("invalid parse of delta string") 
	}
}

