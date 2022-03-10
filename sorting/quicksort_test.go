package sorting

import (
	"fmt"
	"github.com/go-test/deep"
	"github.com/hashicorp/go-multierror"
	"github.com/mikejlong60/golangz/propcheck"
	"sort"
	"testing"
	"time"
)

func TestQuickSortWithInts(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 20000, propcheck.Int())

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
			return xs
		},
		func(expected []int) (bool, error) {
			//myStart := time.Now()
			actual := QuickSort(expected, lessThan)
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

			var errors error
			if diff := deep.Equal(actual, expected); diff != nil {
				errors = multierror.Append(errors, fmt.Errorf("Arrays were not equal %v ", diff))
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
	ge := propcheck.ChooseList(0, 20000, propcheck.String(20))

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
			return xs
		},
		func(expected []string) (bool, error) {
			//			myStart := time.Now()
			actual := QuickSort(expected, lessThan)
			//log.Printf(
			//	"My sort of string array of size:%v took %s", len(expected),
			//	time.Since(myStart),
			//)
			//			goStart := time.Now()
			sort.Strings(expected)
			//log.Printf(
			//	"Golang sort of string array of size:%v took %s", len(expected),
			//	time.Since(goStart),
			//)

			var errors error
			if diff := deep.Equal(actual, expected); diff != nil {
				errors = multierror.Append(errors, fmt.Errorf("Arrays were not equal %v ", diff))
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
	actual := QuickSort(xs, lessThan)
	if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", expected) {
		t.Errorf("actual:%v, expected:%v", actual, expected)
	}
}
