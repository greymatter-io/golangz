package arrays

import (
	"fmt"
	"runtime"
	"testing"
)

type fTypeFoldl func([]int, int, func(int, int) int) (int, tTypeFoldl)
type tTypeFoldl func() (int, tTypeFoldl)

var st = new(runtime.MemStats)

func TrampolineFoldLeft(as []int, z int, f func(int, int) int) int {
	accum, p := foldL(as, z, f)
	for {
		if p == nil {
			break
		}
		accum, p = p()
	}
	return accum
}

func thunk(fn fTypeFoldl, as []int, z int, f func(int, int) int) tTypeFoldl {
	g := func() (int, tTypeFoldl) {
		a, b := fn(as, z, f)
		return a, b
	}
	return g
}
func foldL(as []int, z int, f func(int, int) int) (int, tTypeFoldl) {
	if len(as) > 1 { //Slice has a head and a tail.
		h, t := as[0], as[1:]
		zz := f(z, h)
		return zz, thunk(foldL, t, zz, sum)
	} else if len(as) == 1 { //Slice has a head and an empty tail.
		h := as[0]
		zz := f(z, h)
		return zz, thunk(foldL, Zero[int](), zz, sum)
	} else { //Causes the calling Trampoline to leave the for loop with the final accum result
		return z, nil
	}
}

var sum = func(z int, x int) int {
	return z + x
}

func TestTrampolineFoldLeft(t *testing.T) {
	//var massiveArr = make([]int, 1000000000)
	//for i := 0; i < 1000000000; i++ {
	//	massiveArr[i] = i
	//}
	massiveArr := []int{1, 2, 3, 4}
	fmt.Println("Done making big array")
	actual := TrampolineFoldLeft(massiveArr, 0, sum)
	if actual != 10 {
		t.Errorf("expected:%v, actual:%v", 10, actual)
	}

}
