package sorting

func qs(l, r int, partition func(l, r, pivot int) int) {
	if l < r {
		pi := partition(l, r, l+(r-l)/2)
		qs(l, pi-1, partition)
		qs(pi+1, r, partition)
	}
}

// This is a generic Quicksort.  You only need to pass in a predicate function that tells whether or not l is less than r.
// This is NOT a pure function. It mutates the underlying xs array.
func QuickSort[T any](xs []T, lessThan func(l, r T) bool) {
	swap := func(x, y int) {
		tmp := xs[x]
		xs[x] = xs[y]
		xs[y] = tmp
	}

	partition := func(l, r, pivot int) int {
		pivotVal := xs[pivot]
		swap(pivot, r)
		j := l
		for i := l; i < r; i++ {
			if lessThan(xs[i], pivotVal) {
				swap(i, j)
				j++
			}
		}
		swap(j, r)
		return j
	}

	if len(xs) == 0 {
		return
	} else {
		qs(0, len(xs)-1, partition)
		return
	}
}
