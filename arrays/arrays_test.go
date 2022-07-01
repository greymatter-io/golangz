package arrays

import (
	"fmt"
	"github.com/go-test/deep"
	"github.com/greymatter-io/golangz/propcheck"
	"github.com/hashicorp/go-multierror"
	"testing"
	"time"
)

func TestFoldLeftMatrixSum(t *testing.T) {
	matrixSumWithFoldLeft := func(source []int) [][]int64 {
		var inner [][]int64
		var append = func(xs [][]int64, x int) [][]int64 {
			var result = make([]int64, 2)
			currentX := int64(x)
			result[0] = currentX
			currentAccumLen := len(xs)
			if currentAccumLen == 0 { //Set first sum to first element value
				result[1] = currentX
			} else { //Just grab previous sum
				result[1] = xs[currentAccumLen-1][1] + currentX
			}
			xs = append(xs, result)
			return xs
		}
		result := FoldLeft(source, inner, append)
		return result
	}

	g0 := propcheck.ChooseInt(1, 3000)
	g1 := propcheck.ChooseArray(0, 10000, g0)
	now := time.Now().Nanosecond()
	rng := propcheck.SimpleRNG{now}
	prop := propcheck.ForAll(g1,
		"Validate FoldLeft works and does not change order of resulting array. \n",
		func(xs []int) propcheck.Pair[[][]int64, []int] {
			r := matrixSumWithFoldLeft(xs)
			return propcheck.Pair[[][]int64, []int]{r, xs}
		},
		func(p propcheck.Pair[[][]int64, []int]) (bool, error) {
			xss := p.A
			xs := p.B
			var errors error
			for i := 1; i < len(xss); i++ {
				last := xss[i-1][1]
				if xss[i][1] < last || xss[i][0] != int64(xs[i]) {
					errors = multierror.Append(errors, fmt.Errorf("Array element sum[%v] should not have been less than previous accumulated value", xss[i][1]))
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

func TestTrampolineFoldLeft(t *testing.T) {
	massiveArr := []int64{1, 2, 3, 4}
	sum := func(z int64, x int64) int64 {
		return z + x
	}
	actual := FoldLeft(massiveArr, 0, sum)
	if actual != 10 {
		t.Errorf("expected:%v, actual:%v", 10, actual)
	}
	concat := func(z string, x int64) string {
		return fmt.Sprintf("%v-%v", z, x)
	}
	actual2 := FoldLeft(massiveArr, "", concat)
	if actual2 != "-1-2-3-4" {
		t.Errorf("expected:%v, actual2:%v", "-1-2-3-4", actual2)
	}

}

func TestFoldRightForStrings(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseArray(0, 20, propcheck.String(40))

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
	ge := propcheck.ChooseArray(0, 20, propcheck.Int())

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
	ge := propcheck.ChooseArray(3, 10, propcheck.Int())

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
			if !ArrayEquality(actual, expected, p) {
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
	ge := propcheck.ChooseArray(0, 20, propcheck.ArrayOfN(10, propcheck.Int()))

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
			if !ArrayEquality(actual, expected, p) {
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
	ge := propcheck.ChooseArray[string](0, 20, propcheck.String(40))

	prop := propcheck.ForAll[[]string](ge,
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

	eq := func(l, r string) bool {
		if l == r {
			return true
		} else {
			return false
		}
	}
	if !ArrayEquality(bigarray, []string{"a", "b", "c", "a", "b", "c", "a", "b", "c", "a", "b", "c"}, eq) {
		t.Error("Arrays not equal")
	}
}

func TestAppend(t *testing.T) {
	a1 := []string{"a", "b", "c"}
	a2 := []string{"d", "e", "f"}
	bigarray := Append(a1, a2)

	eq := func(l, r string) bool {
		if l == r {
			return true
		} else {
			return false
		}
	}
	if !ArrayEquality(bigarray, []string{"a", "b", "c", "d", "e", "f"}, eq) {
		t.Error("Arrays not equal")
	}
}

func TestIntArrays(t *testing.T) {
	eq := func(l, r int) bool {
		if l == r {
			return true
		} else {
			return false
		}
	}
	arr := []int{1, 2, 3}
	actual := Append(Append(arr, arr), Append(arr, arr))
	expected := []int{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}
	if !ArrayEquality(actual, expected, eq) {
		t.Errorf("Actual:%v, Expected;%v", actual, expected)
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
