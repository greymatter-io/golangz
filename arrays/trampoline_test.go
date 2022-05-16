package arrays

import (
	"fmt"
	"runtime"
	"testing"
)

// parameters are n, prev, and current
// return either the result (int) or a thunkType if not done
type fnType func(int64, int64, int64) (int64, thunkType)

// return either the result (int) or another thunk
type thunkType func() (int64, thunkType)

func thunk(fn fnType, n int64, prev int64, curr int64) thunkType {
	return func() (int64, thunkType) {
		return fn(n, prev, curr)
	}
}

func thunkFib(n int64, prev int64, curr int64) (int64, thunkType) {
	if n == 0 {
		return prev, nil
	}
	if n == 1 {
		return curr, nil
	}
	// since we return another thunk, the int result does not matter
	return 0 /* unused */, thunk(thunkFib, n-1, curr, curr+prev)
}

//TODO Add trampolining to FoldLeft
//TODO Investigate use heuristic the Golang standard library uses to decide the way to do a sort based upon the size of the array.  For Folds with small array sizes it is more efficient not to do the thunkk thing.
//TODO Apply heap-safe FoldLeft to FoldRight and then reverse.  Think about heuristic above for this.
//TODO Review with Rob and Ming
func trampoline(fn fnType) func(int64) int64 {
	st := new(runtime.MemStats)
	return func(n int64) int64 {
		result, f := fn(n, 0, 1) // initial values for aggregators
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
	fib := trampoline(thunkFib)
	fmt.Printf("Fibonacci(%d) = %d\n", n, fib(n))
}
