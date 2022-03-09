package util

func qs(l, r int, partition func(l, r, pivot int) int) {
	if l < r {
		pi := partition(l, r, l+(r-l)/2)
		qs(l, pi-1, partition)
		qs(pi+1, r, partition)
	}
}

//This is a generic Quicksort.  You only need to pass in a predicate function that tells whether or not l is less than r.
func QuickSort[T any](xs []T, lessThan func(l, r T) bool) []T {
	arr := xs
	swap := func(x, y int) {
		tmp := arr[x]
		arr[x] = arr[y]
		arr[y] = tmp
	}

	partition := func(l, r, pivot int) int {
		pivotVal := arr[pivot]
		swap(pivot, r)
		j := l
		for i := l; i < r; i++ {
			if lessThan(arr[i], pivotVal) {
				swap(i, j)
				j++
			}
		}
		swap(j, r)
		return j
	}

	if len(xs) == 0 {
		return xs
	} else {
		qs(0, len(arr)-1, partition)
		return arr
	}
}
