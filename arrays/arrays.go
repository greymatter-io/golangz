package arrays

// The efficiency of this algorithm is O(N) but it reverses the list.  Use FoldLeft instead if you don't want this.
func FoldRight[T1, T2 any](as []T1, z T2, f func(T1, T2) T2) T2 {
	if len(as) > 1 { //Slice has a head and a tail.
		h, t := as[0], as[1:len(as)]
		return f(h, FoldRight(t, z, f))
	} else if len(as) == 1 { //Slice has a head and an empty tail.
		h := as[0]
		return f(h, FoldRight(Zero[T1](), z, f))
	}
	return z
}

type fTypeFoldl[T1, T2 any] func([]T1, T2, func(T2, T1) T2) (T2, tTypeFoldl[T2])
type tTypeFoldl[T2 any] func() (T2, tTypeFoldl[T2])

// This is a stack-safe, tail recursive, pure function that uses the trampoline technique to ensure that the runtime
// does not face recursion that is too deep(i.e. The garbage collector will run before the recursion gets deep enough to blow the stack).
// See https://trinetri.wordpress.com/2015/04/28/tail-call-thunks-and-trampoline-in-golang/
// A tail call is a function call that is the last action performed in a function.
// The efficiency of this algorithm is O(N) and it does not reverse the list like FoldRight does.
// Most of the other array functions in this package use FoldLeft and thus all are stack-safe.
func FoldLeft[T1, T2 any](as []T1, z T2, f func(T2, T1) T2) T2 {
	accum, p := foldL(as, z, f)
	for {
		if p == nil {
			break
		}
		accum, p = p()
	}
	return accum
}

func thunk[T1, T2 any](fn fTypeFoldl[T1, T2], as []T1, z T2, f func(T2, T1) T2) tTypeFoldl[T2] {
	g := func() (T2, tTypeFoldl[T2]) {
		a, b := fn(as, z, f)
		return a, b
	}
	return g
}
func foldL[T1, T2 any](as []T1, z T2, f func(T2, T1) T2) (T2, tTypeFoldl[T2]) {
	if len(as) > 1 { //Slice has a head and a tail.
		h, t := as[0], as[1:]
		zz := f(z, h)
		return zz, thunk(foldL[T1, T2], t, zz, f)
	} else if len(as) == 1 { //Slice has a head and an empty tail.
		h := as[0]
		zz := f(z, h)
		return zz, thunk(foldL[T1, T2], Zero[T1](), zz, f)
	} else { //Causes the stack-safe function FoldLeft that uses the trampoline technique to leave the recursion with the final accum result
		return z, nil
	}
}

func Appender[T any](s T, as []T) []T {
	gss := append(as, s)
	return gss
}

func Zero[T any]() []T {
	return []T{}
}

// The efficiency of this algorithm is O(N)
func Reverse[T1 any](xs []T1) []T1 {
	f := Appender[T1]
	return FoldRight(xs, Zero[T1](), f)
}

// A structure-preserving Functor on the given array of T.
// The efficiency of this algorithm is O(N)
func Map[T1, T2 any](as []T1, f func(T1) T2) []T2 {
	g := func(as []T2, s T1) []T2 {
		gss := append(as, f(s))
		return gss
	}
	xs := FoldLeft(as, Zero[T2](), g)
	return xs
}

// The efficiency of this algorithm is O(N)
// Collapses the given array of arrays without changing the order.
func Concat[A any](l [][]A) []A {
	g := func(s []A, as []A) []A {
		gss := append(s, as...)
		return gss
	}
	return FoldLeft(l, Zero[A](), g)
}

// Similar to Map in that it takes an array of T1 and applies a function to each element.
// But FlatMap is more powerful than map. We can use flatMap to generate a collection that is either larger or smaller than the original input.
// The efficiency of this algorithm is O(N squared)
func FlatMap[T1, T2 any](as []T1, f func(T1) []T2) []T2 {
	return Concat(Map(as, f))
}

// The efficiency of this algorithm is O(N)
func Filter[T any](as []T, p func(T) bool) []T {
	var g = func(accum []T, h T) []T {
		if p(h) {
			return append(accum, h)
		} else {
			return accum
		}
	}
	return FoldLeft(as, []T{}, g)
}

// The efficiency of this algorithm is O(N)
// Appends as2 to the end of as1
func Append[T any](as1, as2 []T) []T {
	var g = func(accum []T, h T) []T {
		return append(accum, h)
	}
	return FoldLeft(as2, as1, g)
}

// The efficiency of this algorithm is O(N)
func Contains[T any](source []T, contains T, equality func(l, r T) bool) bool {
	p := func(s T) bool {
		if equality(s, contains) {
			return true
		} else {
			return false
		}
	}
	r := Filter(source, p)
	if len(r) > 0 {
		return true
	} else {
		return false
	}
}

// The efficiency of this algorithm is O(N)
func ContainsAllOf[T any](source []T, contains []T, equality func(l, r T) bool) bool {
	for _, v := range contains {
		if !Contains(source, v, equality) {
			return false
		}
	}
	return true
}

// The efficiency of this algorithm is O(N)
func ArrayEquality[T any](aa []T, bb []T, equality func(l, r T) bool) bool {
	f := func(aa, bb []T) bool {
		for i, _ := range aa {
			if !equality(aa[i], bb[i]) {
				return false
			}
		}
		return true
	}

	return (aa == nil && bb == nil) || (len(aa) == 0 && bb == nil) || (aa == nil && len(bb) == 0) || (len(aa) == 0 && len(bb) == 0) ||
		(len(aa) == len(bb) && f(aa, bb))
}
