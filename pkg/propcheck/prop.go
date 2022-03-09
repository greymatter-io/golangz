package propcheck

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"testing"
)

type TestCases = int //The number of test cases to run.
type PropName = string
type Run = func(RunParms) Result

type RunParms struct {
	TestCases TestCases
	Rng       SimpleRNG
}
type Result interface {
	IsFalsified() bool
}

type Prop struct {
	Run  Run
	Name PropName
}

func (w Prop) String() string {
	return fmt.Sprintf("Prop{Run: %T, Name: %v}", w.Run, w.Name)
}

type Falsified[A any] struct {
	Name            string
	FailedCase      A
	Successes       int
	LastSuccessCase A
	Errors          error
}

func (w Falsified[A]) String() string {
	return fmt.Sprintf("\u001B[31m Falsified{Name: %v, FailedCase: %v, Successes: %v, LastSuccessCase: %v, Errors: %v \u001B[30m}", w.Name, w.FailedCase, w.Successes, w.LastSuccessCase, w.Errors)
}

type Passed[A any] struct{}

func (w Passed[A]) String() string {
	return fmt.Sprintf("Passed{}")
}

func (f Falsified[A]) IsFalsified() bool {
	return true
}

func (f Passed[A]) IsFalsified() bool {
	return false
}

//This is a lazily evaluated And that composes two properties.
func And[A any](p1, p2 Prop) Prop {
	run := func(n RunParms) Result {
		r := p1.Run(n)
		if !r.IsFalsified() {
			return p2.Run(n)
		} else {
			return r
		}
	}
	return Prop{run, p1.Name}
}

//This is a lazily evaluated Or that composes two properties.
func Or[A any](p1, p2 Prop) Prop {
	run := func(n RunParms) Result {
		r := p1.Run(n)
		if !r.IsFalsified() {
			return r
		} else {
			return p2.Run(n)
		}
	}
	return Prop{run, p1.Name}
}

/**
Given a Generator(ge), a generated-value transformation function(f), and a variadic list of predicate functions(assertions),
ForAll produces a function(of type Prop) that will run a set number of test cases with a given generator.

Parameters:
	ge - a generator of type "func(SimpleRNG) (A, SimpleRNG)"
	name - a name to assign the Prop
	f - a function of type "f func(A) B" that takes the generated type A and returns another type B and then passes it along to the list of assertion functions.
	assertions - a variadic list of assertion functions of type "func(B) (bool, error)", each returning a pair consisting of a boolean success and a possible list of errors.
Returns:
    Prop - a data structure consisting of a descriptive name for the property and a function of type func(n RunParms) Result. Result is a sum type that can be
		either Falsified or Passed. The FailedCase and LastSuccessCase attributes of the Falsified type(type parameter A)
        contain the value that caused the test failure and the last successful value for the test.
*/
func ForAll[A, B any](ge func(SimpleRNG) (A, SimpleRNG), name string, f func(A) B, assertions ...func(B) (bool, error)) Prop {
	run := func(n RunParms) Result {
		var rng = n.Rng
		var failedCases []Falsified[A]
		var successCases []Result
		var lastSuccessCase A
		var testData A
		for x := 0; x < n.TestCases; x++ {
			testData, rng = ge(rng)
			b := f(testData)
			var errors error
			for _, s := range assertions {
				success, err := s(b)
				if !success {
					if err != nil {
						errors = multierror.Append(errors, err)
					}
					break
				}
			}
			if errors == nil {
				successCases = append(successCases, Passed[A]{})
				lastSuccessCase = testData
			} else {
				f := Falsified[A]{
					Name:            name,
					FailedCase:      testData,
					Successes:       x,
					LastSuccessCase: lastSuccessCase,
					Errors:          errors,
				}
				failedCases = append(failedCases, f)
			}
			_, rng = NextInt(rng)
		}
		if len(failedCases) > 0 {
			return failedCases[0]
		} else {
			return Passed[A]{}
		}

	}
	return Prop{run, name}
}

func ExpectSuccess[A any](t *testing.T, result Result) {
	switch v := result.(type) {
	case Falsified[A]:
		t.Errorf("\033[31m Test Falsified with: %v  \u001B[30m \n", v)
	case Passed[A]:
	default:
		panic(fmt.Sprintf("Expected type of Result to be:%T which is the type of the generator.", v))
	}
}

func ExpectFailure[A any](t *testing.T, result Result) {
	switch v := result.(type) {
	case Passed[A]:
		t.Errorf("\u001B[31m Expected test to be Falsified but it was: %v \u001B[30m \n", v)
	case Falsified[A]:
	default:
		panic(fmt.Sprintf("Expected type of Result to be:%T which is the type of the generator.", v))
	}
}

//Combines a list of assertion functions of type "func(A) (bool, error)" into a single new function that returns their logical OR.
//Note that like Prop.Or above it evaluates lazily.  As soon as a true function is encountered it returns.
//Otherwise it returns all the accumulated errors.
func AssertionOr[A any](assertions ...func(A) (bool, error)) func(A) (bool, error) {
	return func(b A) (bool, error) {
		var errors error
		for _, s := range assertions {
			success, err := s(b)
			if success {
				return true, nil
			} else {
				if err != nil {
					errors = multierror.Append(errors, err)
				}
			}
		}
		return false, errors
	}
}

//Combines a list of assertion functions of type "func(A) (bool, error)" into a new function that returns their logical AND
//Note that like Prop.And above it evaluates lazily.   Branches are evaluated only until one fails.
func AssertionAnd[A any](assertions ...func(A) (bool, error)) func(A) (bool, error) {
	return func(b A) (bool, error) {
		var errors error
		for _, s := range assertions {
			success, err := s(b)
			if !success {
				if err != nil {
					errors = multierror.Append(errors, err)
				}
				return false, errors
			}
		}
		return true, nil
	}
}
