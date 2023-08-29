package sets

import (
	"github.com/greymatter-io/golangz/arrays"
	"github.com/greymatter-io/golangz/sorting"
)

// Makes a "real" set from an array. A "real" set is an array without duplicates based upon the given equality predicate.
// Note the set operations in arrays package do not depend on the array being "real".
// Making a set is expensive because it requires sorting and de-duplication.
// lt - a predicate function that returns true if l is lexically less than r
// eq - a predicate function that returns true if l is equal to r
// NOTE -- this is not a pure function. It mutates the input array a.
func ToSet[T any](a []T, lt, eq func(l, r T) bool) []T {
	sorting.QuickSort(a, lt)
	var r = []T{}
	var previous T
	for i, x := range a {
		if i == 0 { //First element is always added
			r = append(r, x)
		} else if !eq(previous, x) {
			r = append(r, x)
		}
		previous = x
	}
	return r
}

// The efficiency of this algorithm is O(N)
func SetEquality[T any](aa []T, bb []T, equality func(l, r T) bool) bool {
	return (aa == nil && bb == nil) || (len(aa) == 0 && bb == nil) || (aa == nil && len(bb) == 0) || (len(aa) == 0 && len(bb) == 0) || (arrays.ContainsAllOf(aa, bb, equality) && arrays.ContainsAllOf(bb, aa, equality))
}

// Returns the set 'a' minus set 'b'
// The efficiency of this algorithm is O(N)
func SetMinus[T any](a []T, b []T, equality func(l, r T) bool) []T {
	var result []T
	for _, v := range a {
		if !arrays.Contains(b, v, equality) {
			result = append(result, v)
		}
	}
	return result
}

// Returns the intersection of set 'a' and 'b'
// The efficiency of this algorithm is O(5 * N)
func SetIntersection[T any](a []T, b []T, lt, eq func(l, r T) bool) []T {
	ma := SetMinus(a, b, eq)
	mb := SetMinus(b, a, eq)
	return SetMinus(SetUnion(a, b, lt, eq), SetUnion(ma, mb, lt, eq), eq)
}

// Returns the set union of set 'a' and 'b'
func SetUnion[T any](a []T, b []T, lt, eq func(l, r T) bool) []T {
	xs := arrays.Append(a, b)
	return ToSet(xs, lt, eq)
}
