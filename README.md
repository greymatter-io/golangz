# Functional Programming tools in Golang that use Go 1.18 Generics

General Purpose tools for Functional Programming

- Generic set operations on Golang arrays: Filtering, Union, Intersection, Minus
- Generic Linked List
- Generic function composition
- Generic sort
- Generic function currying and partial application

Type-Checked Properties-based testing that is based upon ScalaCheck and Haskell Quickcheck

- Everything has strong typing thanks to Golang Generics.
- Many kinds of test data generators
- Every generator is composable. Lots of generators are included already and new ones are easy to make as compositions
  of existing generators.
- A property is composable with other properties
- Assertions are composable with And and Or logic
- Test Failures include the specific generated values that caused test failure as well as the last successful case.
- Programmers can better cover the scope of all possible inputs to a test (i.e. zero values, empty things, etc).
- Tests outcomes are reproducible
- Programmers can eliminate a lot of code duplication and get better tests at the same time because properties-based
  testing uses random test data.
- Properties-based testing is useful for all sorts of tests: unit, integration, course-grained
  functional/system/black-box.

A few properties-based testing library exist in Golang. [Gopter](https://github.com/leanovate/gopter/) is an example.
This library has much less code than Gopter and provides a very important feature that Gopter does not, namely that all
abstractions are fully composable.

## Two Key Abstractions

- Generators - Generators are functions that produce random test data.
    - They are composable. You can combine them to make other generators.
    - They obey algebraic laws. You can guarantee the safety of their compositions.
    - They are pure functions, freely shareable between Go Routines.
    - Generators allow you to reproduce the exact same test data by passing in the same integer seed value into a
      SimpleRNG. You should never need to save files of test data again.
- Properties - Properties are functions that execute a predicate-like function over a set of test data generated using a
  given Generator.
    - They are composable - You can combine them in arbitrary ways to make new properties.
    - They obey algebraic laws so that you can guarantee the safety of their compositions.
    - They are pure functions, freely shareable between Go Routines.
    - A property is a function that takes two function arguments:
        - one transforming a generated value into a potentially different type
        - and one evaluating the preceding value for correctness.
        - These two characteristics are a very important feature for test reuse.

## Example property for checking the correctness of the function that makes a set from an array of user-defined types:

```
func TestMakeSet(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}      //Generate a random seed based upon current timestamp
	ge := propcheck.ChooseList(0, 200, propcheck.ChooseInt(0, 20)) //Make a generator that will  produce a list of length 0 - 200 of the integers 0 - 20.  You will probably get some duplication which is what we want.

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
		set := MakeSet(xs, lt, eq)
		if len(xs) < len(set) {
			errors = multierror.Append(errors, fmt.Errorf("Set length was longer than underlying array"))
		}
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
	setComplete := func(xs []fancy) (bool, error) {
		var errors error
		set := MakeSet(xs, lt, eq)
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
	propcheck.ExpectSuccess[[]int](t, result)//The type here is the kind of generator
}
```

## Initializing Project
    Project requires go 1.18.
    From root of project.

- `go mod init`
- `go mod tidy`
- `go test -v -count=1 ./...`
