package propcheck

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestStringGenerator(t *testing.T) {
	rng := SimpleRNG{time.Now().Nanosecond()}
	stringMaxSize := 100
	var zeroLengthString = false
	var nonZeroLengthString = false
	mustBeAZeroLengthString := ForAll(String(stringMaxSize), "Must be A A zero length string. \n",
		func(x string) string {
			return x
		},
		func(s string) (bool, error) {
			if len(s) == 0 && !zeroLengthString {
				zeroLengthString = true //This is weird relying on A closure. HMMMM I should think about this.
			}
			return true, nil
		},
	)
	mustBeANonZerolengthString := ForAll(String(stringMaxSize), "Must be A  non-zero length string. \n",
		func(x string) string {
			return x
		},
		func(s string) (bool, error) {
			if len(s) > 0 && !nonZeroLengthString {
				nonZeroLengthString = true
			}
			return true, nil
		})
	_ = And[string](mustBeANonZerolengthString, mustBeAZeroLengthString).Run(RunParms{500, rng}) //Result does not matter. You have to look at the closure to verify.
	if !zeroLengthString {
		t.Errorf("There should have been A zero length string. \n")
	}
	if !nonZeroLengthString {
		t.Errorf("There should have been A non-zero length string. %v \n", nonZeroLengthString)
	}
}

func TestFloatGenerator(t *testing.T) {
	rng := SimpleRNG{time.Now().Nanosecond()}
	mustBeFloat := ForAll(Float(), "Number must be A floating point number with A fractional part. \n",
		func(y float64) float64 {
			_, z := math.Modf(y)
			return z
		},
		func(z float64) (bool, error) {
			if z > 0 {
				return true, nil
			} else {
				return false, fmt.Errorf("z was not a floating point number but was %v", z)
			}
		},
	)
	result := mustBeFloat.Run(RunParms{200, rng})
	ExpectSuccess[float64](t, result)
}

func TestChooseInt(t *testing.T) {
	var rng = SimpleRNG{time.Now().Nanosecond()}
	start := 12
	endExclusive := 130100
	mustBeInRange := ForAll(ChooseInt(start, endExclusive), fmt.Sprintf("Number must be in the range %v and %v exclusive. \n", start, endExclusive),
		func(x int) int {
			return x
		},
		func(x int) (bool, error) {
			if x >= start && x < endExclusive {
				return true, nil
			} else {
				return false, fmt.Errorf("The number was not in the proper range")
			}
		},
	)
	result := mustBeInRange.Run(RunParms{200, rng})
	ExpectSuccess[int](t, result)
}

func TestChooseDate(t *testing.T) {
	rng := SimpleRNG{time.Now().Nanosecond()}
	start := 0
	stopExclusive := 100
	mustBeBoolean := ForAll(ChooseDate(start, stopExclusive), "Date should be within 100 days of 1999-12-31. \n",
		func(x time.Time) time.Time {
			return x
		},

		func(y time.Time) (bool, error) {
			ninetynine := "1999-12-31"
			current, _ := time.Parse("2006-01-02", ninetynine)
			pastRange := current.AddDate(0, 0, -stopExclusive-start)
			futureRange := current.AddDate(0, 0, stopExclusive-start)
			if y.Before(pastRange) || y.After(futureRange) {
				return false, fmt.Errorf("date %v was out of range", y)
			} else {
				return true, nil
			}
		},
	)
	result := mustBeBoolean.Run(RunParms{200, rng})
	ExpectSuccess[time.Time](t, result)
}

func TestNonNegativeInt(t *testing.T) {
	rng := SimpleRNG{time.Now().Nanosecond()}
	mustBePositiveInt := ForAll(NonNegativeInt, "Number must be A positive integer. \n",
		func(x int) int {
			return x
		},
		func(y int) (bool, error) {
			if y < 0 {
				return false, fmt.Errorf("Number was negative")
			} else {
				return true, nil
			}
		},
	)
	result := mustBePositiveInt.Run(RunParms{200, rng})
	ExpectSuccess[int](t, result)
}

func TestThatMap2IsContextSensitive(t *testing.T) {
	rng := SimpleRNG{time.Now().Nanosecond()}
	g := func(j, k int) Pair[int, int] {
		//Note here how Map2(and Flatmap) are context-sensitive. The A value is calculated and then used to calculate the B value. That is
		//consistent with how I see the Generator working in the  FP book exercises. But it means that the A generator will not produce the same
		//B value as the A generator.
		return Pair[int, int]{j, k}
	}

	h := Map2(Int(), Int(), g)
	test := ForAll(h, "Elements from same kind of generator with same SimpleRNG should have not produced the same value. \n",
		func(jk Pair[int, int]) Pair[int, int] {
			return jk
		},
		func(jk Pair[int, int]) (bool, error) {
			jkt := jk
			if jkt.A == jkt.B {
				return false, fmt.Errorf("The same generator should have not produced the same from both generators")
			} else {
				return true, nil
			}
		},
	)
	result := test.Run(RunParms{200, rng})
	ExpectSuccess[Pair[int, int]](t, result)
}

func TestListOfNWithForAll(t *testing.T) {
	start := 1
	endExclusive := 5
	rSize := 1000
	rng := SimpleRNG{time.Now().Nanosecond()}
	u := ListOfN(rSize, ChooseInt(1, 5))
	correctLength := ForAll(u, fmt.Sprintf("Array must have A length of %v \n", rSize),
		func(xs []int) []int {
			return xs
		},
		func(xs []int) (bool, error) {
			if len(xs) != rSize {
				return false, fmt.Errorf("Sizes were not same")
			} else {
				return true, nil
			}
		},
	)
	elementsInRange := ForAll(u, "Each number in Array must be between 1 and 4 \n",
		func(xs []int) []int {
			return xs
		},
		func(ys []int) (bool, error) {
			for _, v := range ys {
				c := v
				if c < start || c > (endExclusive-1) {
					return false, fmt.Errorf("One of the numbers in array was out of range.")
				}
			}
			return true, nil
		})
	test := And[[]int](correctLength, elementsInRange)
	result := test.Run(RunParms{200, rng})
	ExpectSuccess[[]int](t, result)
}

func TestWeightedNoWeightOnFirst(t *testing.T) {
	r1 := ChooseInt(1, 5)
	r2 := ChooseInt(10, 20)
	l := []WeightedGen[int]{{Gen: r1, Weight: 0}, {Gen: r2, Weight: 10}}
	ge2 := Weighted(l)
	lengthInRange := ForAll(ge2, "Weighted should have produced an int between 10 and 20 exclusive",
		func(x int) int {
			return x
		},
		func(x int) (bool, error) {
			if x < 10 && x > 19 {
				return false, fmt.Errorf("Number was in wrong range")
			} else {
				return true, nil
			}
		},
	)
	rng := SimpleRNG{time.Now().Nanosecond()}
	result := lengthInRange.Run(RunParms{200, rng})
	ExpectSuccess[int](t, result)
}

func TestWeightedFailureLastSuccessCase(t *testing.T) {
	minListSize := 10
	maxListSize := 25
	r1 := ChooseInt(1, 5)
	r2 := ChooseInt(10, 20)

	l := []WeightedGen[int]{{Gen: r1, Weight: 200}, {Gen: r2, Weight: 10}}

	ge2 := Weighted(l)
	list := ChooseList(minListSize, maxListSize, ge2)

	lengthInRange := ForAll(list, fmt.Sprintf("Weighted should have produced A list with size between %v and %v inclusive", minListSize, maxListSize),
		func(x []int) []int {
			return x

		},
		func(x []int) (bool, error) {
			l := x
			if len(l) > maxListSize {
				return false, fmt.Errorf("List was too big")
			} else {
				return true, nil
			}
		},
	)
	rng := SimpleRNG{time.Now().Nanosecond()}
	result := lengthInRange.Run(RunParms{200, rng})
	ExpectSuccess[[]int](t, result)
}

func TestWeighted(t *testing.T) {
	var rng = SimpleRNG{time.Now().Nanosecond()}
	r1 := ChooseInt(1000, 5000)
	r2 := ChooseInt(100000, 200000)
	l := []WeightedGen[int]{{Gen: r1, Weight: 300}, {Gen: r2, Weight: 10}}
	u := Weighted(l)
	checker := func(x int) int {
		return x
	}
	assertion := func(x int) (bool, error) {
		if x < 5000 {
			if x < 1000 {
				return false, fmt.Errorf(" Number not in proper range")
			}
		} else if x < 100000 || x > 199999 {
			return false, fmt.Errorf("Number not in proper range")
		}
		return true, nil
	}
	test := ForAll(u, "Weighted should have produced A number between 1000 and 5000 exclusive or between 100000 and 200000 exclusive.", checker, assertion)
	ExpectSuccess[int](t, test.Run(RunParms{200, rng}))
}

func TestChooseListWillProduceListOfZeroElements(t *testing.T) {
	ge := ChooseInt(0, 1000)
	ge2 := ChooseList(0, 0, ge)
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}
	lengthZero := ForAll(ge2, "List must have A length of zero",
		func(xs []int) []int {
			return xs
		},
		func(xs []int) (bool, error) {
			if len(xs) != 0 {
				return false, fmt.Errorf("Length of array should be zero")
			} else {
				return true, nil
			}
		},
	)
	result := lengthZero.Run(RunParms{200, rng})
	ExpectSuccess[[]int](t, result)
}

func TestChooseListWillProduceListOfAtLeastMinElementsAndNotMoreThanHighRange(t *testing.T) {
	maxListSize := 10
	minListSize := 2
	ge := ChooseInt(0, 1000)
	ge2 := ChooseList(minListSize, maxListSize, ge)
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}
	lengthGEOne := ForAll(ge2, fmt.Sprintf("List must have at least A length of %v \n", minListSize),
		func(xs []int) []int {
			return xs
		},
		func(xs []int) (bool, error) {
			if len(xs) < minListSize {
				return false, fmt.Errorf("List too small")
			} else {
				return true, nil
			}
		},
	)
	lengthLEMax := ForAll(ge2, fmt.Sprintf("List must have A length less than %v \n", maxListSize),
		func(xs []int) []int { return xs },
		func(xs []int) (bool, error) {
			if len(xs) > maxListSize {
				return false, fmt.Errorf("List too large")
			} else {
				return true, nil
			}
		},
	)
	bigProp := And[[]int](lengthGEOne, lengthLEMax)
	result := bigProp.Run(RunParms{200, rng})
	ExpectSuccess[[]int](t, result)
}

func TestListOfNWillProduceListConsistingOfDifferentValues(t *testing.T) {
	listSize := 20
	ge2 := ListOfN(listSize, String(40))
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}
	prop := ForAll(ge2, fmt.Sprintf("List of strings contained duplicate strings %v \n", listSize),
		func(xs []string) []string {
			return xs
		},
		func(xss []string) (bool, error) {
			var last string
			for _, v := range xss {
				if last == v && len(last) > 0 {
					return false, fmt.Errorf("list of strings contained a duplicate, possible but extremmely unlikely")
				}
				last = v
			}
			return true, nil
		},
	)
	result := prop.Run(RunParms{200, rng})
	ExpectSuccess[[]string](t, result)
}
