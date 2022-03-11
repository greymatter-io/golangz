package arrays

import (
	"fmt"
	"github.com/go-test/deep"
	"github.com/hashicorp/go-multierror"
	"github.com/mikejlong60/golangz/propcheck"
	"testing"
	"time"
)

func TestFoldRightForStrings(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 20, propcheck.String(40))

	prop := propcheck.ForAll(ge,
		"Filter out all strings > than 10 characters long  \n",
		func(s []string) []string {
			var p = func(s string) bool {
				if len(s) <= 10 {
					return true
				} else {
					return false
				}
			}

			return Filter(s, p)
		},
		func(xss []string) (bool, error) {
			var errors error
			for _, s := range xss {
				if len(s) > 10 {
					errors = multierror.Append(errors, fmt.Errorf("string %v was longer than 10 characters", s))
				}
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

func TestFoldRightForInts(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 20, propcheck.Int())

	prop := propcheck.ForAll(ge,
		"Filter out all ints > than 1000 \n",
		func(xs []int) []int {
			var p = func(s int) bool {
				if s <= 1000 {
					return true
				} else {
					return false
				}
			}

			return Filter(xs, p)
		},
		func(xs []int) (bool, error) {
			var errors error
			for _, s := range xs {
				if s > 1000 {
					errors = multierror.Append(errors, fmt.Errorf("int %v was larger than 1000", s))
				}
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

func TestFlatMap(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(3, 10, propcheck.Int())

	prop := propcheck.ForAll(ge,
		"Test Flatmap turns an array of ints into a larger sized array of strings because FlatMap let's you remove a layer from the given array.\n",
		func(xs []int) []int {
			return xs
		},
		func(xs []int) (bool, error) {
			var errors error
			f := func(x int) []string {
				someExtraThingsInResult := []string{"a", "b", "c", fmt.Sprint(x)}
				return append(someExtraThingsInResult)
			}

			actual := FlatMap(xs, f)
			var expected []string
			for _, x := range xs {
				expected = append(expected, []string{"a", "b", "c", fmt.Sprint(x)}...)
			}

			p := func(x, y string) bool {
				if x == y {
					return true
				} else {
					return false
				}
			}
			if !SetEquality(actual, expected, p) {
				errors = multierror.Append(errors, fmt.Errorf("Actual:%v\n    Expected:%v", actual, expected))
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

func TestConcatArrayOfArrays(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 20, propcheck.ListOfN(10, propcheck.Int()))

	prop := propcheck.ForAll(ge,
		"Test Concat that flattens an array of arrays by 1 level and preserves order\"  \n",
		func(xs [][]int) [][]int {
			return xs
		},
		func(xss [][]int) (bool, error) {
			var errors error
			actual := Concat(xss)
			var expected []int
			for _, xs := range xss {
				expected = append(expected, xs...)
			}

			p := func(x, y int) bool {
				if x == y {
					return true
				} else {
					return false
				}
			}
			if !SetEquality(actual, expected, p) {
				errors = multierror.Append(errors, fmt.Errorf("Actual:%v\n    Expected:%v", actual, expected))
			}
			if errors != nil {
				return false, errors
			} else {
				return true, nil
			}
		},
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[][]int](t, result)

}

func TestMapForStrings(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 20, propcheck.String(40))

	prop := propcheck.ForAll(ge,
		"Make the strings all uppercase  \n",
		func(xs []string) []string {
			//Make all strings <= in length the constant "DUDE".
			//Otherwise make it "MAMA"
			var p = func(xs string) string {
				if len(xs) <= 10 {
					if len(xs) == 0 {
						return "NONE"
					} else {
						return "DUDE"
					}
				} else {
					return xs
				}
			}

			return Map(xs, p)
		},
		func(xss []string) (bool, error) {
			var errors error
			for _, s := range xss {
				if len(s) <= 10 {
					if s != "DUDE" && s != "NONE" {
						errors = multierror.Append(errors, fmt.Errorf("All string values shorter than 11 characters should have been 'DUDE' and this one was[%v]", s))
					}
				}
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

func TestAppendStringArrays(t *testing.T) {
	strings := []string{"a", "b", "c"}
	bigarray := Append(Append(strings, strings), Append(strings, strings))

	if diff := deep.Equal(bigarray, []string{"a", "b", "c", "c", "b", "a", "a", "b", "c", "c", "b", "a"}); diff != nil {
		t.Error(diff)
	}
}

func TestIntArrays(t *testing.T) {
	arr := []int{1, 2, 3}
	bigarray := Append(Append(arr, arr), Append(arr, arr))
	if diff := deep.Equal(bigarray, []int{1, 2, 3, 3, 2, 1, 1, 2, 3, 3, 2, 1}); diff != nil {
		t.Error(diff)
	}
}

func TestFilterIntArray(t *testing.T) {
	arr := []int{1, 2, 3, 3, 3, 3}
	var p = func(s int) bool {
		if s == 3 {
			return true
		} else {
			return false
		}
	}

	bigarray := Filter(arr, p)
	if diff := deep.Equal(bigarray, []int{3, 3, 3, 3}); diff != nil {
		t.Error(diff)
	}
}

func TestSetEqualityForIntArray(t *testing.T) {
	arr1 := []int{1, 2, 3, 3, 3, 3, 3, 3, 3, 1, 2}
	arr2 := []int{1, 2, 3, 3, 3, 3}
	equality := func(l, r int) bool {
		if l == r {
			return true
		} else {
			return false
		}
	}

	if !SetEquality(arr1, arr2, equality) {
		t.Error("sets should have been equal but were not")
	}
}

func TestSetInequalityForIntArray(t *testing.T) {
	arr1 := []int{1, 2, 3, 12, 3, 3, 3}
	arr2 := []int{1, 2, 3, 3, 3, 3}
	equality := func(l, r int) bool {
		if l == r {
			return true
		} else {
			return false
		}
	}
	if SetEquality(arr1, arr2, equality) {
		t.Error("sets should not have been equal but were")
	}
}

func TestSetMinusForIntArray(t *testing.T) {
	arr1 := []int{1, 2, 3, 12, 3, 3, 3}
	arr2 := []int{1, 2, 3, 3, 3, 3}
	equality := func(l, r int) bool {
		if l == r {
			return true
		} else {
			return false
		}
	}

	r := SetMinus(arr1, arr2, equality)
	if !SetEquality(r, []int{12}, equality) {
		t.Errorf("expected:%v, actual:%v", []int{12}, r)
	}
}

func TestSetMinusForStringArray(t *testing.T) {
	arr1 := []string{"a", "b", "c", "d"}
	arr2 := []string{"a", "b"}
	equality := func(l, r string) bool {
		if l == r {
			return true
		} else {
			return false
		}
	}

	r := SetMinus(arr1, arr2, equality)
	if !SetEquality(r, []string{"c", "d"}, equality) {
		t.Errorf("expected:%v, actual:%v", []string{"c", "d"}, r)
	}
}

func TestSetIntersection(t *testing.T) {
	type fancy struct {
		id string
	}
	arr1 := []fancy{fancy{"a"}, fancy{"b"}, fancy{"c"}, fancy{"d"}}
	arr2 := []fancy{fancy{"a"}, fancy{"b"}}
	equality := func(l, r fancy) bool {
		if l.id == r.id {
			return true
		} else {
			return false
		}
	}

	expected := []fancy{fancy{"a"}, fancy{"b"}}
	r := SetIntersection(arr1, arr2, equality)
	if !SetEquality(r, expected, equality) {
		t.Errorf("expected:%v, actual:%v", expected, r)
	}
}

func TestSetUnion(t *testing.T) {
	type fancy struct {
		id string
	}
	arr1 := []fancy{fancy{"a"}, fancy{"b"}, fancy{"c"}, fancy{"d"}}
	arr2 := []fancy{fancy{"a"}, fancy{"b"}, fancy{"z"}}
	equality := func(l, r fancy) bool {
		if l.id == r.id {
			return true
		} else {
			return false
		}
	}

	expected := []fancy{fancy{"a"}, fancy{"z"}, fancy{"b"}, fancy{"c"}, fancy{"d"}}

	//	r := SetUnion(arr1, []string{""})// Note that this will not compile(type error) Yaaaa!!!!
	r := SetUnion(arr1, arr2)
	if !SetEquality(r, expected, equality) {
		t.Errorf("expected:%v, actual:%v", expected, r)
	}
}
