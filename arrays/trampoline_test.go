package arrays

import (
	"fmt"
	"runtime"
	"testing"
)

// parameters are n, prev, and current
// return either the result (int) or a thunkType if not done
//func FoldLeft[T1, T2 any](as []T1, z T2, f func(T2, T1) T2) T2 {

type fnType[T1, T2 any] func(T1, T2) (T1, thunkType[T1, T2])

// return either the result (int) or another thunk
type thunkType[T1, T2 any] func() (T1, thunkType[T1, T2])

func thunk[T1, T2 any](fn fnType[T1, T2], n T1, u T2) thunkType[T1, T2] {
	return func() (T1, thunkType[T1, T2]) {
		a, b := fn(n, u)
		fmt.Printf("a in thunk funtion:%v\n", a)
		fmt.Printf("b in thunk funtion:%v\n", b)
		return a, b //fn(n, u)
	}
}

func thunkFib[T1, T2 any](n T1, u T2) (T1, thunkType[T1, T2]) {
	// since we return another thunk, the int result does not matter
	fmt.Printf("in thunkFib:n:%v u:%v\n", n, u)
	return n /* unused */, thunk[T1, T2](thunkFib[T1, T2], n, u)
}

//TODO Add trampolining to FoldLeft
//TODO Investigate use heuristic the Golang standard library uses to decide the way to do a sort based upon the size of the array.  For Folds with small array sizes it is more efficient not to do the thunkk thing.
//TODO Apply heap-safe FoldLeft to FoldRight and then reverse.  Think about heuristic above for this.
//TODO Review with Rob and Ming
func trampoline[T1, T2 any](fn fnType[T1, T2]) func(T1, T2) T1 {
	st := new(runtime.MemStats)
	return func(n T1, u T2) T1 {
		result, f := fn(n, u) // initial values for aggregators
		for {
			if f == nil {
				break
			}
			result, f = f()
			runtime.ReadMemStats(st)
			fmt.Printf("\nMemstat HeapObjects:%v\n", st.HeapObjects)
		}
		return result
	}
}

func TestTrampoline(t *testing.T) {
	n := int64(200) // fib(10) == 55
	u := []int{1, 2, 3}
	fib := trampoline(thunkFib[int64, []int])
	fmt.Printf("Fibonacci(%d) = %d\n", n, fib(n, u))
}
