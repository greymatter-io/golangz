package stack

import (
	"fmt"
	"github.com/greymatter-io/golangz/option"
	"github.com/greymatter-io/golangz/propcheck"
	"github.com/hashicorp/go-multierror"
	"reflect"
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
			var l *Stack[string]
			var i int
			for {
				if len(xss) == 0 {
					break
				}
				l = Push(xss[i], l)
				p, err := Peek(l)
				if err != nil {
					fmt.Println(err)
				}
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
			makeStack := func(xss []string) *Stack[string] {
				var l *Stack[string]
				var i int
				for {
					if len(xss) == 0 {
						break
					}
					l = Push(xss[i], l)
					if i+1 == len(xss) {
						break
					} else {
						i++
					}
				}
				return l
			}

			var l = makeStack(xss)
			var p string
			var errors error
			var popErr error
			for i := len(xss) - 1; i >= 0; i-- {
				p, l, popErr = Pop(l)
				fmt.Printf("popErr:%v\n", popErr)
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
	emptyStack := &Stack[string]{}
	_, _, err := Pop(emptyStack)
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestPopOEmptyStack(t *testing.T) {
	emptyStack := &Stack[string]{}
	r := PopO(emptyStack)
	f := func(s propcheck.Pair[string, *Stack[string]]) string {
		return s.A
	}
	actual := option.Map(r, f)
	s := fmt.Sprint(reflect.TypeOf(actual))
	if s != "option.None[string]" {
		t.Errorf("Expected \"option.None[string]\" actual:%v", s)
	}
}

func TestPopNonEmptyStack(t *testing.T) {
	a := &Stack[string]{}
	b := PushO("fred", a)
	 := PopO(b)
	f := func(s propcheck.Pair[string, *Stack[string]]) string {
		return s.A
	}

	actual := option.Map(r, f)
	s := fmt.Sprint(reflect.TypeOf(actual))
	if s != "option.None[string]" {
		t.Errorf("Expected \"option.None[string]\" actual:%v", s)
	}
}
