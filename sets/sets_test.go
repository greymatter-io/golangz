package sets

import (
	"fmt"
	"github.com/greymatter-io/golangz/arrays"
	"github.com/greymatter-io/golangz/propcheck"
	"github.com/hashicorp/go-multierror"
	"testing"
	"time"
)

func TestToSet(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}       //Generate a random seed based upon current timestamp
	ge := propcheck.ChooseArray(0, 200, propcheck.ChooseInt(0, 20)) //Produce a generator a list of length 0 - 200 of the integers 0 - 20.  You will probably get some duplication which is what we want

	type fancy struct {
		id int
	}

	lt := func(l, r fancy) bool {
		if l.id < r.id {
			return true
		} else {
			return false
		}
	}
	eq := func(l, r fancy) bool {
		if l.id == r.id {
			return true
		} else {
			return false
		}
	}

	//This is function of type "f func(A) B" that takes the generated type A and returns another type B and then passes it along to the list of assertion functions.
	arrayToFancyType := func(xs []int) []fancy {
		var r []fancy
		for _, x := range xs {
			r = append(r, fancy{x})
		}
		return r
	}
	//Assertion functions
	setCorrectLength := func(xs []fancy) (bool, error) {
		var errors error
		set := ToSet(xs, lt, eq)
		if len(xs) < len(set) {
			errors = multierror.Append(errors, fmt.Errorf("Set length was longer than underlying array"))
		}
		if errors != nil {
			return false, errors
		} else {
			return true, nil
		}
	}
	setComplete := func(xs []fancy) (bool, error) {
		var errors error
		set := ToSet(xs, lt, eq)
		diff := arrays.SetMinus(xs, set, eq)
		if len(diff) > 0 {
			errors = multierror.Append(errors, fmt.Errorf("Difference between orignal array and its set should have been zero but was %v", diff))
		}
		if errors != nil {
			return false, errors
		} else {
			return true, nil
		}
	}

	prop := propcheck.ForAll(ge,
		"Assert the creation of a set from an array of user-defined types  \n",
		arrayToFancyType,
		setCorrectLength, setComplete,
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]int](t, result)
}
