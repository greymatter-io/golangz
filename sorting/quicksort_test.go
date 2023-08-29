package sorting

import (
	"fmt"
	"github.com/greymatter-io/golangz/arrays"
	"github.com/greymatter-io/golangz/propcheck"
	"github.com/hashicorp/go-multierror"
	"sort"
	"testing"
	"time"
)

func TestQuickSortWithInts(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseArray(0, 20000, propcheck.Int())

	lessThan := func(l, r int) bool {
		if l < r {
			return true
		} else {
			return false
		}
	}
	prop := propcheck.ForAll(ge,
		"Sort an array of ints  \n",
		func(xs []int) []int {
			if len(xs) == 0 {
				return []int{}
			} else {
				return xs
			}
		},
		func(xs []int) (bool, error) {
			var expected = make([]int, len(xs))
			copy(expected, xs)
			//myStart := time.Now()
			QuickSort(xs, lessThan)
			//log.Printf(
			//	"My sort of int array of size:%v took %s", len(expected),
			//	time.Since(myStart),
			//)
			//		goStart := time.Now()
			sort.Ints(expected)
			//log.Printf(
			//	"Golang sort of int array of size:%v took %s", len(expected),
			//	time.Since(goStart),
			//)

			eq := func(l, r int) bool {
				if l == r {
					return true
				} else {
					return false
				}
			}
			var errors error
			if !arrays.ArrayEquality(xs, expected, eq) {
				errors = multierror.Append(errors, fmt.Errorf(" Actual: %v\nExpected:%v ", xs, expected))
			}
			if errors != nil {
				return false, errors
			} else {
				return true, nil
			}
		},
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]int](t, result)

}

func TestQuickSortWithStrings(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseArray(0, 20000, propcheck.String(20))

	lessThan := func(l, r string) bool {
		if l < r {
			return true
		} else {
			return false
		}
	}
	prop := propcheck.ForAll(ge,
		"Sort an array of strings  \n",
		func(xs []string) []string {
			if len(xs) == 0 {
				return []string{}
			} else {
				return xs
			}
		},
		func(xs []string) (bool, error) {
			//	myStart := time.Now()
			var expected = make([]string, len(xs))
			copy(expected, xs)
			QuickSort(xs, lessThan)
			//log.Printf(
			//	"My sort of string array of size:%v took %s", len(expected),
			//	time.Since(myStart),
			//)
			//goStart := time.Now()
			sort.Strings(expected)
			//log.Printf(
			//	"Golang sort of string array of size:%v took %s", len(expected),
			//	time.Since(goStart),
			//)

			eq := func(l, r string) bool {
				if l == r {
					return true
				} else {
					return false
				}
			}
			var errors error
			if !arrays.ArrayEquality(xs, expected, eq) {
				errors = multierror.Append(errors, fmt.Errorf(" Actual: %v\nExpected:%v ", xs, expected))
			}
			if errors != nil {
				return false, errors
			} else {
				return true, nil
			}
		},
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]string](t, result)
}

func TestSortSpecialGuyArray(t *testing.T) {
	type specialGuy struct {
		age        int
		likeCoffee bool
	}

	xs := []specialGuy{{12, true}, {122, true}, {120, true}, {3, true}}
	expected := []specialGuy{{3, true}, {12, true}, {120, true}, {122, true}}
	lessThan := func(l, r specialGuy) bool {
		if l.age < r.age {
			return true
		} else {
			return false
		}
	}
	QuickSort(xs, lessThan)
	if fmt.Sprintf("%v", xs) != fmt.Sprintf("%v", expected) {
		t.Errorf("actual:%v, expected:%v", xs, expected)
	}
}
