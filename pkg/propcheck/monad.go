package propcheck

func fId[A any](a A) func(SimpleRNG) (A, SimpleRNG) {
	return func(r SimpleRNG) (A, SimpleRNG) {
		return a, r
	}
}

func Id[A any](a A) func(SimpleRNG) (A, SimpleRNG) {
	return fId(a)
}

func pFlatMap[A, B any](f func(SimpleRNG) (A, SimpleRNG), g func(A) func(SimpleRNG) (B, SimpleRNG)) func(SimpleRNG) (B, SimpleRNG) {
	return func(rng SimpleRNG) (B, SimpleRNG) {
		a, r1 := f(rng)
		b, r2 := g(a)(r1)
		return b, r2
	}
}

func FlatMap[A, B any](f func(SimpleRNG) (A, SimpleRNG), g func(A) func(SimpleRNG) (B, SimpleRNG)) func(SimpleRNG) (B, SimpleRNG) {
	return pFlatMap(f, g)
}

func pMap[A, B any](s func(SimpleRNG) (A, SimpleRNG), f func(A) B) func(SimpleRNG) (B, SimpleRNG) {
	return func(rng SimpleRNG) (B, SimpleRNG) {
		fa := func(a A) func(SimpleRNG) (B, SimpleRNG) { return Id(f(a)) }
		r := FlatMap(s, fa)
		return r(rng)
	}
}

func Map[A, B any](s func(SimpleRNG) (A, SimpleRNG), f func(A) B) func(SimpleRNG) (B, SimpleRNG) {
	return pMap(s, f)
}

func Sequence[T any](rs []func(SimpleRNG) (T, SimpleRNG)) func(SimpleRNG) ([]T, SimpleRNG) {
	var f []T
	var g = Id(f)
	h := func(accum []T, b T) []T {
		return append(accum, b)
	}
	for _, r := range rs {
		g = Map2(g, r, h)
	}
	return g
}
