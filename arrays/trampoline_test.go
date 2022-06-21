package arrays

import (
	"fmt"
	"runtime"
	"testing"
)

type fTypeFoldl[T1, T2 any] func([]T1, T2, func(T2, T1) T2) (T2, tTypeFoldl[T2])
type tTypeFoldl[T2 any] func() (T2, tTypeFoldl[T2])

var st = new(runtime.MemStats)

func TrampolineFoldLeft[T1, T2 any](as []T1, z T2, f func(T2, T1) T2) T2 {
	accum, p := foldL(as, z, f)
	for {
		runtime.ReadMemStats(st)
		fmt.Printf("\nMemstat HeapObjects:%v\n", st.HeapObjects)
		if p == nil {
			break
		}
		accum, p = p()
	}
	return accum
}

func thunk[T1, T2 any](fn fTypeFoldl[T1, T2], as []T1, z T2, f func(T2, T1) T2) tTypeFoldl[T2] {
	g := func() (T2, tTypeFoldl[T2]) {
		a, b := fn(as, z, f)
		return a, b
	}
	return g
}
func foldL[T1, T2 any](as []T1, z T2, f func(T2, T1) T2) (T2, tTypeFoldl[T2]) {
	if len(as) > 1 { //Slice has a head and a tail.
		h, t := as[0], as[1:]
		zz := f(z, h)
		return zz, thunk(foldL[T1, T2], t, zz, f)
	} else if len(as) == 1 { //Slice has a head and an empty tail.
		h := as[0]
		zz := f(z, h)
		return zz, thunk(foldL[T1, T2], Zero[T1](), zz, f)
	} else { //Causes the calling Trampoline to leave the for loop with the final accum result
		return z, nil
	}
}

func TestTrampolineFoldLeft(t *testing.T) {
	var massiveArr = make([]int64, 1000000000)
	for i := 0; i < 1000000000; i++ {
		massiveArr[i] = int64(i)
	}
	//massiveArr := []int{1, 2, 3, 4}
	fmt.Println("Done making big array")
	sum := func(z int64, x int64) int64 {
		return z + x
	}
	actual := TrampolineFoldLeft(massiveArr, 0, sum)
	fmt.Println(actual)
	if actual != 10 {
		t.Errorf("expected:%v, actual:%v", 10, actual)
	}
	concat := func(z string, x int64) string {
		return fmt.Sprintf("%v-%v", z, x)
	}
	actual2 := TrampolineFoldLeft(massiveArr, "", concat)
	fmt.Println(actual2)
	if actual2 != "-1-2-3-4" {
		t.Errorf("expected:%v, actual2:%v", "-1-2-3-4", actual2)
	}

}
