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

// MapN are the functions that make this an Applicative Functor. These functions allow you to compose generators without the context-sensitivity that you get with FlatMap.
// A good example of this is validation where you don't want the computation to stop because A Flatmap in the chain fails.
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

func Map8[A, B, C, D, E, F, G, H, I any](ra func(SimpleRNG) (A, SimpleRNG), rb func(SimpleRNG) (B, SimpleRNG), rc func(SimpleRNG) (C, SimpleRNG),
	rd func(SimpleRNG) (D, SimpleRNG), re func(SimpleRNG) (E, SimpleRNG), rf func(SimpleRNG) (F, SimpleRNG),
	rg func(SimpleRNG) (G, SimpleRNG), rh func(SimpleRNG) (H, SimpleRNG), f func(a A, b B, c C, d D, e E, f F, g G, h H) I) func(SimpleRNG) (I, SimpleRNG) {
	return func(rng SimpleRNG) (I, SimpleRNG) {

		fab := Product(ra, rb)
		fcd := Product(rc, rd)
		fef := Product(re, rf)
		fgh := Product(rg, rh)

		fg := func(ab Pair[A, B], cd Pair[C, D], ef Pair[E, F], gh Pair[G, H]) I {
			return f(ab.A, ab.B, cd.A, cd.B, ef.A, ef.B, gh.A, gh.B)
		}

		g := Map4(fab, fcd, fef, fgh, fg)

		return g(rng)
	}
}

func Map16[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q any](ra func(SimpleRNG) (A, SimpleRNG), rb func(SimpleRNG) (B, SimpleRNG), rc func(SimpleRNG) (C, SimpleRNG),
	rd func(SimpleRNG) (D, SimpleRNG), re func(SimpleRNG) (E, SimpleRNG), rf func(SimpleRNG) (F, SimpleRNG),
	rg func(SimpleRNG) (G, SimpleRNG), rh func(SimpleRNG) (H, SimpleRNG), ri func(SimpleRNG) (I, SimpleRNG), rj func(SimpleRNG) (J, SimpleRNG), rk func(SimpleRNG) (K, SimpleRNG), rl func(SimpleRNG) (L, SimpleRNG),
	rm func(SimpleRNG) (M, SimpleRNG), rn func(SimpleRNG) (N, SimpleRNG), ro func(SimpleRNG) (O, SimpleRNG), rp func(SimpleRNG) (P, SimpleRNG),
	f func(a A, b B, c C, d D, e E, f F, g G, h H, i I, j J, k K, l L, m M, n N, o O, p P) Q) func(SimpleRNG) (Q, SimpleRNG) {
	return func(rng SimpleRNG) (Q, SimpleRNG) {

		fab := Product(ra, rb)
		fcd := Product(rc, rd)
		fef := Product(re, rf)
		fgh := Product(rg, rh)
		fij := Product(ri, rj)
		fkl := Product(rk, rl)
		fmn := Product(rm, rn)
		fop := Product(ro, rp)

		fg := func(ab Pair[A, B], cd Pair[C, D], ef Pair[E, F], gh Pair[G, H], ij Pair[I, J], kl Pair[K, L], mn Pair[M, N], op Pair[O, P]) Q {
			return f(ab.A, ab.B, cd.A, cd.B, ef.A, ef.B, gh.A, gh.B, ij.A, ij.B, kl.A, kl.B, mn.A, mn.B, op.A, op.B)
		}

		g := Map8(fab, fcd, fef, fgh, fij, fkl, fmn, fop, fg)

		return g(rng)
	}
}

func Map32[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, AA, BB, CC, DD, EE, FF, GG, HH, II any](ra func(SimpleRNG) (A, SimpleRNG), rb func(SimpleRNG) (B, SimpleRNG), rc func(SimpleRNG) (C, SimpleRNG),
	rd func(SimpleRNG) (D, SimpleRNG), re func(SimpleRNG) (E, SimpleRNG), rf func(SimpleRNG) (F, SimpleRNG),
	rg func(SimpleRNG) (G, SimpleRNG), rh func(SimpleRNG) (H, SimpleRNG), ri func(SimpleRNG) (I, SimpleRNG), rj func(SimpleRNG) (J, SimpleRNG), rk func(SimpleRNG) (K, SimpleRNG), rl func(SimpleRNG) (L, SimpleRNG),
	rm func(SimpleRNG) (M, SimpleRNG), rn func(SimpleRNG) (N, SimpleRNG), ro func(SimpleRNG) (O, SimpleRNG), rp func(SimpleRNG) (P, SimpleRNG), rq func(SimpleRNG) (Q, SimpleRNG), rr func(SimpleRNG) (R, SimpleRNG),
	rs func(SimpleRNG) (S, SimpleRNG), rt func(SimpleRNG) (T, SimpleRNG), ru func(SimpleRNG) (U, SimpleRNG), rv func(SimpleRNG) (V, SimpleRNG), rw func(SimpleRNG) (W, SimpleRNG), rx func(SimpleRNG) (X, SimpleRNG),
	raa func(SimpleRNG) (AA, SimpleRNG), rbb func(SimpleRNG) (BB, SimpleRNG), rcc func(SimpleRNG) (CC, SimpleRNG), rdd func(SimpleRNG) (DD, SimpleRNG), ree func(SimpleRNG) (EE, SimpleRNG), rff func(SimpleRNG) (FF, SimpleRNG),
	rgg func(SimpleRNG) (GG, SimpleRNG), rhh func(SimpleRNG) (HH, SimpleRNG),
	f func(a A, b B, c C, d D, e E, f F, g G, h H, i I, j J, k K, l L, m M, n N, o O, p P, q Q, r R, s S, t T, u U, v V, w W, x X, aa AA, bb BB, cc CC, dd DD, ee EE, ff FF, gg GG, hh HH) II) func(SimpleRNG) (II, SimpleRNG) {
	return func(rng SimpleRNG) (II, SimpleRNG) {
		fab := Product(ra, rb)
		fcd := Product(rc, rd)
		fef := Product(re, rf)
		fgh := Product(rg, rh)
		fij := Product(ri, rj)
		fkl := Product(rk, rl)
		fmn := Product(rm, rn)
		fop := Product(ro, rp)
		fqr := Product(rq, rr)
		fst := Product(rs, rt)
		fuv := Product(ru, rv)
		fxy := Product(rw, rx)
		faabb := Product(raa, rbb)
		fccdd := Product(rcc, rdd)
		feeff := Product(ree, rff)
		fgghh := Product(rgg, rhh)
		fg := func(ab Pair[A, B], cd Pair[C, D], ef Pair[E, F], gh Pair[G, H], ij Pair[I, J], kl Pair[K, L], mn Pair[M, N],
			op Pair[O, P], qr Pair[Q, R], st Pair[S, T],
			uv Pair[U, V], wx Pair[W, X], aabb Pair[AA, BB], ccdd Pair[CC, DD], eeff Pair[EE, FF], gghh Pair[GG, HH]) II {
			return f(ab.A, ab.B, cd.A, cd.B, ef.A, ef.B, gh.A, gh.B, ij.A, ij.B, kl.A, kl.B, mn.A, mn.B, op.A, op.B, qr.A, qr.B,
				st.A, st.B, uv.A, uv.B, wx.A, wx.B, aabb.A, aabb.B, ccdd.A, ccdd.B, eeff.A, eeff.B, gghh.A, gghh.B)
		}
		g := Map16(fab, fcd, fef, fgh, fij, fkl, fmn, fop, fqr, fst, fuv, fxy, faabb, fccdd, feeff, fgghh, fg)
		return g(rng)
	}
}
