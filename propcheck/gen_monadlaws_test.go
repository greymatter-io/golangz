package propcheck

import (
	"github.com/go-test/deep"
	"testing"
	"time"
)

func assoc(ab Pair[[]int, Pair[[]int, []int]]) Pair[Pair[[]int, []int], []int] {
	l := ab.A
	m := ab.B
	n := Pair[[]int, []int]{l, m.A}
	o := Pair[Pair[[]int, []int], []int]{n, m.B}
	return o
}

//Gen is a Functor and a Monad and also an Applicative Functor as proven by the following set of tests.
//Associative Law for Applicative Functors - It should not matter how I nest the combination of three Monadic values into one.
//Note l and r below, the way I change the nesting.
func TestMapAssociativeLaw(t *testing.T) {
	ge := ChooseInt(0, 1000)
	fa := ChooseArray(0, 10, ge)
	fb := ChooseArray(11, 20, ge)
	fc := ChooseArray(21, 30, ge)
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}
	l1 := Product(fa, fb)
	l2 := Product(l1, fc)
	l, _ := l2(rng)
	r1 := Product(fb, fc)
	r2 := Product(fa, r1)
	r, _ := Map(r2, assoc)(rng)
	if diff := deep.Equal(l, r); diff != nil {
		t.Error(diff)
	}
}

//Identity Law for Applicative Functors - A Gen composed with Id (AKA the "Identity Function")
//on the left should produce the same thing as A Gen composed with Id on the right
func TestMap2RightAndLeftIdentityLaw(t *testing.T) {
	fl := func(a int, b int) int {
		return a
	}
	ge := ChooseInt(0, 1000)
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}
	v, _ := ge(rng)
	li := Map2(Id(v), ge, fl)
	ri := Map2(ge, Id(v), fl)
	l, _ := li(rng)
	r, _ := ri(rng)
	if l != r || l != v {
		t.Errorf("l and r and v should have matched and they were instead: %v, %v, %v", l, r, v)
	}
}
