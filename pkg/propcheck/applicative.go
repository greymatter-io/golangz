package propcheck

import "fmt"

type Pair[A, B any] struct {
	A A
	B B
}

func (w Pair[A, B]) String() string {
	return fmt.Sprintf("Pair{A: %v \n, B: %v\n}\n", w.A, w.B)
}

func Product[A, B any](fa func(SimpleRNG) (A, SimpleRNG), fb func(SimpleRNG) (B, SimpleRNG)) func(SimpleRNG) (Pair[A, B], SimpleRNG) {
	f := func(a A, b B) Pair[A, B] {
		return Pair[A, B]{a, b}
	}
	g := Map2(fa, fb, f)
	return g
}

//MapN are the functions that make this an Applicative Functor. These functions allow you to compose generators without the context-sensitivity that you get with FlatMap.
//A good example of this is validation where you don't want the computation to stop because A Flatmap in the chain fails.
//This is the original version of Map2
func pMap2[A, B, C any](ra func(SimpleRNG) (A, SimpleRNG), rb func(SimpleRNG) (B, SimpleRNG), f func(a A, b B) C) func(rng SimpleRNG) (C, SimpleRNG) {
	return func(rng SimpleRNG) (C, SimpleRNG) {
		a, r1 := ra(rng)
		b, r2 := rb(r1)
		c := f(a, b)
		return c, r2
	}
}
func Map2[A, B, C any](ra func(SimpleRNG) (A, SimpleRNG), rb func(SimpleRNG) (B, SimpleRNG), f func(a A, b B) C) func(rng SimpleRNG) (C, SimpleRNG) {
	return pMap2(ra, rb, f)
}

func Map3[A, B, C, D any](ra func(SimpleRNG) (A, SimpleRNG), rb func(SimpleRNG) (B, SimpleRNG),
	rc func(SimpleRNG) (C, SimpleRNG), f func(a A, b B, c C) D) func(SimpleRNG) (D, SimpleRNG) {
	return func(rng SimpleRNG) (D, SimpleRNG) {
		fab := Product[A, B](ra, rb)
		fg := func(abd Pair[A, B], c C) D {
			return f(abd.A, abd.B, c)
		}
		g := Map2(fab, rc, fg)
		return g(rng)
	}
}

func Map4[A, B, C, D, E any](ra func(SimpleRNG) (A, SimpleRNG), rb func(SimpleRNG) (B, SimpleRNG), rc func(SimpleRNG) (C, SimpleRNG),
	rd func(SimpleRNG) (D, SimpleRNG), f func(a A, b B, c C, d D) E) func(SimpleRNG) (E, SimpleRNG) {
	return func(rng SimpleRNG) (E, SimpleRNG) {

		fab := Product(ra, rb)
		fcd := Product(rc, rd)

		fg := func(ab Pair[A, B], cd Pair[C, D]) E {
			var abd = ab
			var cdd = cd
			return f(abd.A, abd.B, cdd.A, cdd.B)
		}

		g := Map2(fab, fcd, fg)
		return g(rng)
	}
}
