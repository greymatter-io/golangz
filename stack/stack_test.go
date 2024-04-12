package stack

import (
	"fmt"
	"github.com/greymatter-io/golangz/propcheck"
	"github.com/hashicorp/go-multierror"
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
				l = l.Push(xss[i])
				p, _ := l.Peek()
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
					l = l.Push(xss[i])
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
				p, _ := l.Peek()
				l = l.Pop()
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
	s := emptyStack.Pop()
	if !s.IsEmpty() {
		t.Errorf("expected stack to be empty")
	}
}

func TestPopIsEmptyStack(t *testing.T) {
	s := NewStack[string]()
	s = s.Push("fred")
	if s.IsEmpty() {
		t.Errorf("expected stack to be empty")
	}
}
