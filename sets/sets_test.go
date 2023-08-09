package sets

import (
	"fmt"
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

	var lt = func(l, r fancy) bool {
		if l.id < r.id {
			return true
		} else {
			return false
		}
	}
	var eq = func(l, r fancy) bool {
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
		diff := SetMinus(xs, set, eq)
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
	var lt = func(l, r fancy) bool {
		if l.id < r.id {
			return true
		} else {
			return false
		}
	}
	var eq = func(l, r fancy) bool {
		if l.id == r.id {
			return true
		} else {
			return false
		}
	}

	expected := []fancy{fancy{"a"}, fancy{"b"}}
	r := SetIntersection(arr1, arr2, lt, eq)
	if !SetEquality(r, expected, eq) {
		t.Errorf("expected:%v, actual:%v", expected, r)
	}
}

func TestSetUnion(t *testing.T) {
	type fancy struct {
		id string
	}
	var lt = func(l, r fancy) bool {
		if l.id < r.id {
			return true
		} else {
			return false
		}
	}
	var eq = func(l, r fancy) bool {
		if l.id == r.id {
			return true
		} else {
			return false
		}
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

	r := SetUnion(arr1, arr2, lt, eq)
	if !SetEquality(r, expected, equality) {
		t.Errorf("expected:%v, actual:%v", expected, r)
	}
}
