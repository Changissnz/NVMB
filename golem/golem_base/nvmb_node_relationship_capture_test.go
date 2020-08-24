package golem 

import (
	//"fmt"
	"testing" 
)

func Test_ParseDeltaString(t *testing.T) {

	target := "sadfasdfsd__delta__3513435"
	d := ParseDeltaString(target) 
	
	if d != "sadfasdfsd__delta" {
		panic("invalid parse of delta string") 
	}
}

