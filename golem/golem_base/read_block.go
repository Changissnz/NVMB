/*
Block struct is a 2-d matrix of strings
*/

package golem

//////////////////////Start: class //////////////////////////

import (
	"fmt"
	"sync"
	///"gonum.org/v1/gonum/mat"
)

type Block struct {
	datos [][]string
	keys  []int /// TODO: review this variable. may be unused
}

func OneBlock() *Block {
	return &Block{datos: make([][]string, 0)}
}

func (b *Block) Dims() (int,int) {
	if len(b.datos) == 0 {
		return 0,0
	}

	return len(b.datos), len(b.datos[0])
}

/*
adds a row to block
*/
func (b *Block) AddOne(s []string) {
	if b.Length() != 0 {
		if len(b.GetAt(0)) != len(s) {
			panic("block cannot accept row of unequal size")
		}
	}
	b.datos = append(b.datos, s)
}

/*
CAUTION: assumes b2 keys are identical to instance
*/ 
func (b *Block) StackBlock(b2 *Block) {
	for _,d := range b2.datos { 
		b.AddOne(d) 
	}

} 

/*
outputs a row
*/
func (b *Block) GetAt(i int) []string {
	if i >= len(b.datos) {
		panic(fmt.Sprintf("invalid index %d", i))
	}

	return b.datos[i]
}

/*
fetches element at 2-index
*/
func (b *Block) GetAtOne(index []int) string {
	if len(index) != 2 {
		panic("invalid index shape")
	}

	return b.datos[index[0]][index[1]]
}

/*
sets er at row i 
*/ 
func (b *Block) SetAt(i int, er []string) {
	if i >= len(b.datos) {
		panic(fmt.Sprintf("invalid index %d", i))
	}

	if len(b.datos[i]) != len(er) {
		panic("input has incorrect length")
	}
	b.datos[i] = er 
}

/*
sets v at 2-index 
*/
func (b *Block) SetAtOne(v string, index []int) {
	if len(index) != 2 {
		panic("invalid index shape")
	}

	b.datos[index[0]][index[1]] = v 
}


/*
displays all 
*/ 
func (b *Block) Display() {
	b.DisplayByRow(0, b.Length())
} 

/*
displays rows sri-eri
*/ 
func (b *Block) DisplayByRow(sri int, eri int) {

	if (sri >= eri) {
		panic("start index has to be less than end index")
	}

	for i := sri; i < eri; i++ {
		q := b.GetAt(i) 
		fmt.Println(q) 
	}
}

/*
displays rows sri-eri at columns ci 
*/ 
func (b *Block) DisplayByRowSelectedColumns(sri int, eri int, ci []int)  {

	if (sri >= eri) {
		panic("start index has to be less than end index")
	}

	for i := sri; i < eri; i++ {
		r, stat := b.FetchElementsAtRow(i, ci)
		if (!stat) {
			return 
		}
		fmt.Println(r) 
	}
}

/*
fetches elements at row r by indices 
if r is invalid, outputs [1] false 
*/ 
func (b *Block) FetchElementsAtRow(r int, indices []int) ([]string, bool) {
	r_ := make([]string,0)

	if b.Length() <= r {
		return nil, false 
	}

	if indices != nil {
		for _, c := range indices {
			r_ = append(r_, b.GetAtOne([]int{r, c})) 
		}
	} else {
		w := b.Width() 
		for i := 0; i < w; i++ {
			r_ = append(r_, b.GetAtOne([]int{r, i})) 
		} 
	}

	return r_, true  
}

func (b *Block) FetchElementsAtRowRange(si int, ei int, indices []int) ([][]string, bool) {
	
	if (si >= ei) {
		panic("invalid index for row-range fetch")
	}

	output := make([][]string, 0)
	for i := si; i < ei; i++ {
		
		r, stat :=  b.FetchElementsAtRow(i, indices)
		if (!stat) {
			return nil, false
		}

		output = append(output, r) 
	}

	return output, true
}

// TODO : implement concurrency
/*
fetches all string columns in indices
*/ 
func (b *Block) FetchColumns(indices []int) *Block {

	l, w := b.Length(), len(indices)
	if (l == 0 || w == 0) {
		return nil 
	}

	q := TwoDimStringMatrix(l, w) 

	b2 := OneBlock() 
	b2.datos = q 

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < l; i++ {
			er, stat := b.FetchElementsAtRow(i, indices)
			
			if stat {
				b2.SetAt(i, er)
			}
		}
		wg.Done() 
	}()

	wg.Wait() 

	return b2 
}

/*
outputs number of rows in block
*/
func (b *Block) Length() int {
	return len(b.datos)
}

/*
outputs number of columns 
*/ 
func (b *Block) Width() int {
	if b.Length() == 0 {
		return 0
	}
	return len(b.datos[0])
}

/*
outputs number of elements in block
*/
func (b *Block) Size() int {
	l := b.Length()

	if l == 0 {
		return 0
	}
	return b.Length() * len(b.GetAt(0))
}

func (b *Block) Clear() {
	b.datos = make([][]string,0) 
	b.keys = make([]int,0)
}

/*
compares two blocks for equality
*/ 
func BlocksEqual(b1 *Block, b2 *Block) bool {

	if len(b1.datos) != len(b2.datos) {
		return false
	}

	for i, x := range b1.datos {
		if !EqualStringSlices(x, b2.datos[i]) {
			return false
		}
	}
	return true
}