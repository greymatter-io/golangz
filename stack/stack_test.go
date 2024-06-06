package stack

import (
	"fmt"
	"github.com/greymatter-io/golangz/propcheck"
	"github.com/hashicorp/go-multierror"
	"math"
	"testing"
	"time"
)

func TestPush(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseArray(0, 10, propcheck.String(40))

	prop := propcheck.ForAll(ge,
		"Validate Push for Stack  \n",
		func(xs []string) []string {
			return xs
		},
		func(xss []string) (bool, error) {
			var errors error
			var l = NewStack[string]()
			var i int
			for {
				if len(xss) == 0 {
					break
				}
				l = Push(l, xss[i])
				p, _ := Peek(l)
				if p != xss[i] {
					errors = multierror.Append(errors, fmt.Errorf("string %v  should have been %v pushed to top of Stack", p, xss[i]))
				}
				if i+1 == len(xss) {
					break
				} else {
					i++
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

func TestPop(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseArray(0, 10, propcheck.String(40))

	prop := propcheck.ForAll(ge,
		"Validate Pop for Stack  \n",
		func(xs []string) []string {
			return xs
		},
		func(xss []string) (bool, error) {
			makeStack := func(xss []string) Stack[string] {
				var l = NewStack[string]()
				var i int
				for {
					if len(xss) == 0 {
						break
					}
					l = Push(l, xss[i])
					if i+1 == len(xss) {
						break
					} else {
						i++
					}
				}
				return l
			}

			var l = makeStack(xss)
			var errors error
			for i := len(xss) - 1; i >= 0; i-- {
				p, _ := Peek(l)
				l = Pop(l)
				if p != xss[i] {
					errors = multierror.Append(errors, fmt.Errorf("string %v  should have been %v popped from top of Stack", p, xss[i]))
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

func TestPopEmptyStack(t *testing.T) {
	emptyStack := NewStack[string]()
	s := Pop(emptyStack)
	if !IsEmpty(s) {
		t.Errorf("expected stack to be empty")
	}
}

func TestPopIsEmptyStack(t *testing.T) {
	s := NewStack[string]()
	s = Push(s, "fred")
	if IsEmpty(s) {
		t.Errorf("expected stack to be empty")
	}
}

func TestFoldRight(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseArray(0, 10, propcheck.ChooseInt(-10, 1000))

	prop := propcheck.ForAll(ge,
		"Validate FoldRight for Stack  \n",
		func(xs []int) []int {
			return xs
		},
		func(xss []int) (bool, error) {
			var errors error
			var l = NewStack[int]()
			var i int
			for {
				if len(xss) == 0 {
					break
				}
				l = Push(l, xss[i])
				if i+1 == len(xss) {
					break
				} else {
					i++
				}
			}

			f := func(currentMax, x int) int {
				if x < currentMax {
					return x
				} else {
					return currentMax
				}
			}
			foldMin := FoldRight[int, int](l, math.MaxInt, f)
			goMin := math.MaxInt
			for _, x := range xss {
				if x < goMin {
					goMin = x
				}
			}
			if foldMin != goMin {
				t.Errorf("actual: %v, expected: %v", foldMin, goMin)
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

func TestFoldLeft(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseArray(0, 10, propcheck.ChooseInt(-10, 1000))

	prop := propcheck.ForAll(ge,
		"Validate FoldLeft for Stack  \n",
		func(xs []int) []int {
			return xs
		},
		func(xss []int) (bool, error) {
			var errors error
			var l = NewStack[int]()
			var i int
			for {
				if len(xss) == 0 {
					break
				}
				l = Push(l, xss[i])
				if i+1 == len(xss) {
					break
				} else {
					i++
				}
			}

			f := func(toString []string, x int) []string {
				toString = append(toString, fmt.Sprintf("%v", x))
				return toString
			}
			zx := FoldLeft[int, []string](l, []string{}, f)
			if len(zx) != len(xss) {
				t.Errorf("actual: %v, expected: %v", len(zx), len(xss))
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
