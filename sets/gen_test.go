package sets

import (
	"fmt"
	"github.com/greymatter-io/golangz/propcheck"
	"testing"
	"time"
)

func TestChooseSet(t *testing.T) {
	lt := func(l, r bool) bool {
		if l {
			return true
		} else {
			return false
		}
	}
	eq := func(l, r bool) bool {
		if l == r {
			return true
		} else {
			return false
		}
	}

	maxListSize := 3000 //Am banking on the odds being near zero that I get the same bool over 3000 rolls of the die.
	minListSize := 2
	ge := propcheck.Boolean()
	ge2 := ChooseSet(minListSize, maxListSize, ge, lt, eq)
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	lengthGEOne := propcheck.ForAll(ge2, fmt.Sprintf("List must have at least A length of %v \n", minListSize),
		func(xs []bool) []bool {
			return xs
		},
		func(xs []bool) (bool, error) {
			if len(xs) < minListSize {
				return false, fmt.Errorf("List too small")
			} else {
				return true, nil
			}
		},
	)
	lengthLEMax := propcheck.ForAll(ge2, fmt.Sprintf("List must have A length less than %v \n", maxListSize),
		func(xs []bool) []bool { return xs },
		func(xs []bool) (bool, error) {
			if len(xs) > maxListSize {
				return false, fmt.Errorf("List too large")
			} else {
				return true, nil
			}
		},
	)
	bigProp := propcheck.And[[]int](lengthGEOne, lengthLEMax)
	result := bigProp.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]bool](t, result)
}
