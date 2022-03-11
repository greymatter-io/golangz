package arrays

//TODO make FoldLeft so you don't have to reverse the array at the end. Just like your linked list.
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

func Appender[T any](s T, as []T) []T {
	gss := append(as, s)
	return gss
}

func Zero[T any]() []T {
	return []T{}
}

func Reverse[T1 any](xs []T1) []T1 {
	f := Appender[T1]
	return FoldRight(xs, Zero[T1](), f)
}

//A structure-preserving Functor on the given array of T.
func Map[T1, T2 any](as []T1, f func(T1) T2) []T2 {
	g := func(s T1, as []T2) []T2 {
		gss := append(as, f(s))
		return gss
	}
	xs := FoldRight(as, Zero[T2](), g)
	//Put the array back in original order
	return Reverse(xs)
}

func Concat[A any](l [][]A) []A {
	g := func(s []A, as []A) []A {
		gss := append(as, s...)
		return gss
	}
	//in := Id(A[])
	return FoldRight(Reverse(l), Zero[A](), g)
}

//Similar to Map in that it takes an array of T1 and applies a function to each element.
//But FlatMap is more powerful than map. We can use flatMap to generate a collection that is either larger or smaller than the original input.
func FlatMap[T1, T2 any](as []T1, f func(T1) []T2) []T2 {
	return Concat(Map(as, f))
}

func Filter[T any](as []T, p func(T) bool) []T {
	var g = func(h T, accum []T) []T {
		if p(h) {
			return append(accum, h)
		} else {
			return accum
		}
	}
	xs := FoldRight(as, []T{}, g)

	f := Appender[T]

	//Reverse it to put the array back in original order
	return FoldRight(xs, Zero[T](), f)
}

func Append[T any](as1, as2 []T) []T {
	var g = func(h T, accum []T) []T {
		return append(accum, h)
	}
	return FoldRight(as1, as2, g)
}

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

func ContainsAllOf[T any](source []T, contains []T, equality func(l, r T) bool) bool {
	for _, v := range contains {
		if !Contains(source, v, equality) {
			return false
		}
	}
	return true
}

func SetEquality[T any](aa []T, bb []T, equality func(l, r T) bool) bool {
	return (aa == nil && bb == nil) || (len(aa) == 0 && bb == nil) || (aa == nil && len(bb) == 0) || (len(aa) == 0 && len(bb) == 0) || (ContainsAllOf(aa, bb, equality) && ContainsAllOf(bb, aa, equality))
}

//Returns the set 'a' minus set 'b'
func SetMinus[T any](a []T, b []T, equality func(l, r T) bool) []T {
	var result []T
	for _, v := range a {
		if !Contains(b, v, equality) {
			result = append(result, v)
		}
	}
	return result
}

//Returns the intersection of set 'a' and 'b'
func SetIntersection[T any](a []T, b []T, equality func(l, r T) bool) []T {
	ma := SetMinus(a, b, equality)
	mb := SetMinus(b, a, equality)
	return SetMinus(SetUnion(a, b), SetUnion(ma, mb), equality)
}

//Returns the set union of set 'a' and 'b'
func SetUnion[T any](a []T, b []T) []T {
	return Append(a, b)
}
