package propcheck

import (
	"strconv"
	"testing"
	"time"
)

//Naturality of Product Law for Applicative Functors- You can combine values of Gen either before or after applying them with Map2 to make a new Gen.
func TestNaturalityOfProductLaw(t *testing.T) {
	f := func(a int) int {
		return a
	}

	g := func(a int) int {
		return a
	}
	ge := ChooseInt(0, 1000)
	gf := ChooseInt(0, 1000)
	rng := SimpleRNG{Seed: time.Now().Nanosecond()}

	naturalityOfProductLaw(t, rng, ge, gf, f, g)
}

type stringDouble struct {
	a string
	b string
}

func naturalityOfProductLaw(t *testing.T, rng SimpleRNG, fa func(SimpleRNG) (int, SimpleRNG), fb func(SimpleRNG) (int, SimpleRNG), faf func(int) int, fbf func(int) int) {
	productF := func(f, g func(int) int) func(int, int) stringDouble {
		return func(x int, y int) stringDouble {
			a := f(x)
			b := g(y)
			return stringDouble{strconv.Itoa(a), strconv.Itoa(b)}
		}
	}
	l := Map2(fa, fb, productF(faf, fbf))
	r1 := Map(fa, faf)
	r2 := Map(fb, fbf)
	r := Product(r1, r2)
	rl, _ := l(rng)
	rr, _ := r(rng)
	rla := rl.a
	lar, _ := strconv.Atoi(rla)
	rar := rr.A

	rlb := rl.b
	lbr, _ := strconv.Atoi(rlb)
	rbr := rr.B

	if lar != rar {
		t.Errorf("The generated integers from left value in Pair should have been the same and were instead: %v, %v", lar, rar)
	}
	if lbr != rbr {
		t.Errorf("The generated integers from right value in Pair should have been the same and were instead: %v, %v", lbr, rbr)
	}
}
