package sets

import (
	"github.com/greymatter-io/golangz/propcheck"
)

//Generates a set with A size in the indicated range using the given Gen
//lt - a predicate function that returns true if l is lexically less than r
//eq - a predicate function that returns true if l is equal to r
//Note that the returned set may not meet the minimum size requirement after de-duplication.
//This is impossible to handle in any general way(i.e. the set of bools with cardinality == 3).
func ChooseSet[T any](start, stopInclusive int, kind func(pair propcheck.SimpleRNG) (T, propcheck.SimpleRNG), lt func(l, r T) bool, eq func(l, r T) bool) func(propcheck.SimpleRNG) ([]T, propcheck.SimpleRNG) {
	return func(rng propcheck.SimpleRNG) ([]T, propcheck.SimpleRNG) {
		r := propcheck.ChooseArray(start, stopInclusive, kind)
		s, rng2 := r(rng)
		return ToSet(s, lt, eq), rng2
	}
}
