package sets

import (
	"github.com/greymatter-io/golangz/sorting"
)

//Makes a "real" set from an array. A "real" set is an array without duplicates based upon the given equality predicate.
//Note the set operations in arrays package do not depend on the array being "real".
//Making a set is expensive because it requires sorting and de-duplication.
//lt - a predicate function that returns true if l is lexically less than r
//eq - a predicate function that returns true if l is equal to r
func ToSet[T any](a []T, lt func(l, r T) bool, eq func(l, r T) bool) []T {
	sorted := sorting.QuickSort(a, lt)
	var r = []T{}
	var previous T
	for i, x := range sorted {
		if i == 0 { //First element is always added
			r = append(r, x)
		} else if !eq(previous, x) {
			r = append(r, x)
		}
		previous = x
	}
	return r
}
