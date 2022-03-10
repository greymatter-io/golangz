package arraylist

func FoldRight[T1, T2 any](as []T1, z T2, f func(T1, T2) T2) T2 {
	if len(as) > 1 { //Slice has a head and a tail.
		h, t := as[0], as[1:len(as)]
		return f(h, FoldRight(t, z, f))
	} else if len(as) == 1 { //Slice has a head and an empty tail.
		h := as[0]
		return f(h, FoldRight([]T1{}, z, f))
	}
	return z
}

func Id[T any](s T, as []T) []T {
	gss := append(as, s)
	return gss
}

func Reverse[T1 any](xs []T1) []T1 {
	id := Id[T1]
	return FoldRight(xs, []T1{}, id)
}

//A structure-preserving Functor on the given array of T.
func Map[T1, T2 any](as []T1, f func(T1) T2) []T2 {
	g := func(s T1, as []T2) []T2 {
		gss := append(as, f(s))
		return gss
	}
	xs := FoldRight(as, []T2{}, g)
	//Put the array back in original order
	return Reverse(xs)
}

func Concat[A any](l [][]A) []A {
	g := func(s []A, as []A) []A {
		gss := append(as, s...)
		return gss
	}
	return FoldRight(Reverse(l), []A{}, g)
}

func FlatMap[T1, T2 any](as []T1, f func(T1) []T2) []T2 {
	return Concat(Map(as, f))
}

//g after f in order of fs
func ComposeAll[T any](fs []func(s T) T) func(s T) T {
	var g = func(s T) T {
		return s
	}

	for _, f := range fs {
		g = Compose(f, g)
	}
	return g
}

//f after g
func Compose[T any](f, g func(s T) T) func(s T) T {
	return func(s T) T {
		return f(g(s))
	}
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

	id := Id[T]

	//Reverse it to put the array back in original order
	return FoldRight(xs, []T{}, id)
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
		//if fmt.Sprint(s) == fmt.Sprint(contains) {
		//	return true
		//} else {
		//		return false
		//	}
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
