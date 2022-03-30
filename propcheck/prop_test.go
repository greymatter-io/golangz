package propcheck

import (
	"fmt"
	"testing"
	"time"
)

func TestAnd(t *testing.T) {
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}
	f1 := func(RunParms) Result {
		return Passed[string]{}
	}
	f2 := func(RunParms) Result {
		return Falsified[string]{
			Name:       "I failed",
			FailedCase: "first one",
			Successes:  0,
		}
	}
	p1 := Prop{
		Run:  f1,
		Name: "first properties test",
	}
	p2 := Prop{
		Run:  f2,
		Name: "first properties test",
	}
	actual := And[string](p1, p2).Run(RunParms{200, rng})
	switch v := actual.(type) {
	case Passed[string]:
		t.Errorf("Invoking And with one Falsified and one Passed Result should have been Falsified and was %v \n", v)
	default:
	}
}

func TestOr(t *testing.T) {
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}
	f1 := func(parms RunParms) Result {
		return Passed[string]{}
	}
	f2 := func(RunParms) Result {
		return Falsified[string]{
			Name:       "I failed",
			FailedCase: "first one",
			Successes:  0,
		}
	}
	p1 := Prop{
		Run:  f1,
		Name: "first properties test",
	}
	p2 := Prop{
		Run:  f2,
		Name: "first properties test",
	}
	actual := Or[string](p1, p2).Run(RunParms{200, rng})
	switch v := actual.(type) {
	case Falsified[string]:
		t.Errorf("Invoking Or with one Falsified and one Passed Result should have been Passed and was %v \n", v)
	default:
	}
}

func TestForAllWithInts(t *testing.T) {
	ge := ChooseInt(1, 501)
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}
	actual := ForAll(ge, "first forall",
		func(x int) int { return x },
		func(x int) (bool, error) {
			if x > 500 {
				return false, fmt.Errorf("Number was too large")
			} else {
				return true, nil
			}
		},
	)
	result := actual.Run(RunParms{200, rng})
	ExpectSuccess[int](t, result)
}

func TestForAllWithArrayOfArraysThatFails(t *testing.T) {
	ge := ChooseInt(1, 1000)
	ge2 := ArrayOfN(100, ge)
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}
	actual := ForAll(ge2, "Array should be longer than 9900 elements.",
		func(xs []int) []int { return xs },
		func(xs []int) (bool, error) {
			if len(xs) < 9900 {
				return false, fmt.Errorf("Expected failure")
			} else {
				return true, nil
			}
		},
	)
	result := actual.Run(RunParms{200, rng})
	ExpectFailure[[]int](t, result)
}

func TestExpectFailure(t *testing.T) {
	ge := ChooseInt(1, 1000)
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}
	actual := ForAll(ge, "Verify the type of ExpectFailure matches the type of the generator.",
		func(x int) string { return fmt.Sprintf("%v", x) },
		func(xs string) (bool, error) {
			return false, fmt.Errorf("a test failure: %v", xs)
		},
	)
	result := actual.Run(RunParms{200, rng})
	ExpectFailure[int](t, result)
}

func TestExpectSucces(t *testing.T) {
	ge := ChooseInt(1, 1000)
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}
	actual := ForAll(ge, "Verify the type of ExpectSuccess matches the type of the generator.",
		func(x int) string { return fmt.Sprintf("%v", x) },
		func(xs string) (bool, error) {
			return true, nil
		},
	)
	result := actual.Run(RunParms{200, rng})
	ExpectSuccess[int](t, result)
}

func TestAssertionOr(t *testing.T) {
	maxListSize := 10
	minListSize := 2
	ge := ChooseInt(0, 1000)
	ge2 := ChooseArray(minListSize, maxListSize, ge)
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}
	assertion1 := func(xs []int) (bool, error) {
		if len(xs) > minListSize {
			return false, fmt.Errorf("List too small")
		} else {
			return true, nil
		}
	}
	assertion2 := func(xs []int) (bool, error) {
		if len(xs) < minListSize {
			return false, fmt.Errorf("List too small")
		} else {
			return true, nil
		}
	}
	lengthGEOne := ForAll(ge2, fmt.Sprintf("List must have at least A length of %v \n", minListSize),
		func(xs []int) []int {
			return xs
		}, AssertionOr(assertion1, assertion2))
	result := lengthGEOne.Run(RunParms{200, rng})
	ExpectSuccess[[]int](t, result)
}

func TestAssertionAndFailure(t *testing.T) {
	maxListSize := 10
	minListSize := 2
	ge := ChooseInt(0, 1000)
	ge2 := ChooseArray(minListSize, maxListSize, ge)
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}
	assertion1 := func(xs []int) (bool, error) {
		if len(xs) > minListSize {
			return false, fmt.Errorf("List too small")
		} else {
			return true, nil
		}
	}
	assertion2 := func(xs []int) (bool, error) {
		if len(xs) < minListSize {
			return false, fmt.Errorf("List too small")
		} else {
			return true, nil
		}
	}
	lengthGEOne := ForAll(ge2, fmt.Sprintf("List must have at least A length of %v \n", minListSize),
		func(xs []int) []int {
			return xs
		}, AssertionAnd(assertion1, assertion2))
	result := lengthGEOne.Run(RunParms{200, rng})
	ExpectFailure[[]int](t, result)
}

func TestAssertionAndSuccess(t *testing.T) {
	maxListSize := 10
	minListSize := 2
	ge := ChooseInt(0, 1000)
	ge2 := ChooseArray(minListSize, maxListSize, ge)
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}
	assertion1 := func(xs []int) (bool, error) {
		if len(xs) < minListSize {
			return false, fmt.Errorf("List too small")
		} else {
			return true, nil
		}
	}
	assertion2 := func(xs []int) (bool, error) {
		if len(xs) > maxListSize {
			return false, fmt.Errorf("List too large")
		} else {
			return true, nil
		}
	}
	prop := ForAll(ge2, fmt.Sprintf("List must have at least A length of %v \n", minListSize),
		func(xs []int) []int {
			return xs
		}, AssertionAnd(assertion1, assertion2))
	result := prop.Run(RunParms{200, rng})
	ExpectSuccess[[]int](t, result)
}
