/*
basic string-set implementation
*/

package golem

import (
	"fmt"
)

// TODO : basic set will operate using string keys

/*
*/
type Placer struct{}

/*
*/
type BasicSet struct {
	datos map[string]Placer
}

/*
*/
func OneBasicSet() *BasicSet {
	return &BasicSet{datos: make(map[string]Placer)}
}

/*
*/
func (b *BasicSet) Len() int {
	return len(b.datos)
}

/*
*/
func (b *BasicSet) AddOne(s string) bool {

	if _, ok := b.datos[s]; !ok {
		b.datos[s] = Placer{}
		return true
	}
	return false
}

/*
*/
func (b *BasicSet) AddBunch(s []string) int {
	var c int = 0
	for _, s_ := range s {
		if b.AddOne(s_) {
			c++
		}
	}
	return c
}

/// TODO: complete this.
func (b *BasicSet) DeleteOne() {

}

/*
*/
func (b *BasicSet) ToSliceIndexFormat() [][]int {

	output := make([][]int, 0)
	for k := range b.datos {
		is, stat := DefaultStringToIntSlice(k, DEFAULT_DELIMITER) 

		if (!stat) {
			panic("set data has invalid int slice!") 
		}

		output = append(output, is)
	}
	return output
}

/*
*/
func (b *BasicSet) PrintOut() {

	fmt.Println("*BASIC SET")
	for k := range b.datos {
		fmt.Println(k)
	}
}
