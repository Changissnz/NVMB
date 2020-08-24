package main

import "C"
import "./golem_base"

/*
test function 
*/
//export Add
func Add(a, b int) int { return a + b * a + b } 

//export IncorpGolem 
func IncorpGolem() int {

	s1 := []string{"one", "two", "three"} 
	s2 := []string{"one", "two", "free"} 
	res := golem.EqualStringSlices(s1, s2)

	if res {
		return 1
	}
	return 0
}

func main() {
}
