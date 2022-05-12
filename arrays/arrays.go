package arrays

//The efficiency of this algorithm is O(N) but it reverses the list.  Use FoldLeft instead if you don't want this.
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

//The efficiency of this algorithm is O(N) and it does not reverse the list like FoldRight does.
//TODO Make this function tail-recursive
func FoldLeft[T1, T2 any](as []T1, z T2, f func(T2, T1) T2) T2 {
	if len(as) > 1 { //Slice has a head and a tail.
		h, t := as[0], as[1:len(as)]
		return FoldLeft(t, f(z, h), f)
	} else if len(as) == 1 { //Slice has a head and an empty tail.
		h := as[0]
		return FoldLeft(Zero[T1](), f(z, h), f)
	}
	return z
}

func Appender[T any](s T, as []T) []T {
	gss := append(as, s)
	return gss
}

func Zero[T any]() []T {
	return []T{}
}

//The efficiency of this algorithm is O(N)
func Reverse[T1 any](xs []T1) []T1 {
	f := Appender[T1]
	return FoldRight(xs, Zero[T1](), f)
}

// A structure-preserving Functor on the given array of T.
//The efficiency of this algorithm is O(N)
func Map[T1, T2 any](as []T1, f func(T1) T2) []T2 {
	g := func(as []T2, s T1) []T2 {
		gss := append(as, f(s))
		return gss
	}
	xs := FoldLeft(as, Zero[T2](), g)
	return xs
}

//The efficiency of this algorithm is O(N)
func Concat[A any](l [][]A) []A {
	g := func(s []A, as []A) []A {
		gss := append(as, s...)
		return gss
	}
	return FoldLeft(l, Zero[A](), g)
}

// Similar to Map in that it takes an array of T1 and applies a function to each element.
// But FlatMap is more powerful than map. We can use flatMap to generate a collection that is either larger or smaller than the original input.
//The efficiency of this algorithm is O(N squared)
func FlatMap[T1, T2 any](as []T1, f func(T1) []T2) []T2 {
	return Concat(Map(as, f))
}

//The efficiency of this algorithm is O(N)
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

//The efficiency of this algorithm is O(N)
func Append[T any](as1, as2 []T) []T {
	var g = func(h T, accum []T) []T {
		return append(accum, h)
	}
	return FoldRight(as1, as2, g)
}

//The efficiency of this algorithm is O(N)
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

//The efficiency of this algorithm is O(N)
func ContainsAllOf[T any](source []T, contains []T, equality func(l, r T) bool) bool {
	for _, v := range contains {
		if !Contains(source, v, equality) {
			return false
		}
	}
	return true
}

//The efficiency of this algorithm is O(N)
func SetEquality[T any](aa []T, bb []T, equality func(l, r T) bool) bool {
	return (aa == nil && bb == nil) || (len(aa) == 0 && bb == nil) || (aa == nil && len(bb) == 0) || (len(aa) == 0 && len(bb) == 0) || (ContainsAllOf(aa, bb, equality) && ContainsAllOf(bb, aa, equality))
}

// Returns the set 'a' minus set 'b'
//The efficiency of this algorithm is O(N)
func SetMinus[T any](a []T, b []T, equality func(l, r T) bool) []T {
	var result []T
	for _, v := range a {
		if !Contains(b, v, equality) {
			result = append(result, v)
		}
	}
	return result
}

// Returns the intersection of set 'a' and 'b'
//The efficiency of this algorithm is O(5 * N)
func SetIntersection[T any](a []T, b []T, equality func(l, r T) bool) []T {
	ma := SetMinus(a, b, equality)
	mb := SetMinus(b, a, equality)
	return SetMinus(SetUnion(a, b), SetUnion(ma, mb), equality)
}

// Returns the set union of set 'a' and 'b'
//The efficiency of this algorithm is O(N)
func SetUnion[T any](a []T, b []T) []T {
	return Append(a, b)
}
