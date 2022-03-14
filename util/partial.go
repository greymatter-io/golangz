package util

//Takes a function of two arguments and turns it into a function with one argument with the first argument bound.
//This is useful replacement for a Class with instance variables in OOP.
func Partial1[T1, T2, T3 any](f func(T1, T2) T3, s1 T1) func(T2) T3 {
	return func(s2 T2) T3 {
		return f(s1, s2)
	}
}

//Takes a function of three arguments and turns it into a function with two arguments with the first argument bound.
//This is useful replacement for a Class with instance variables in OOP.
//TODO I wonder if there is a way to make the function variadic and still preserve generic type safety?  I don't think so because type parameters cannot be variadic.
func Partial2[T1, T2, T3, T4 any](f func(T1, T2, T3) T4, s1 T1) func(T2, T3) T4 {
	return func(s2 T2, s3 T3) T4 {
		return f(s1, s2, s3)
	}
}

//TODO
//def curry3[A, B, C, D](f: (A, B, C) => D): A => (B => (C => D)) = a => b => c => f(a, b, c)
//
//def partial1[A, B, C](a: A, f: (A, B) => C): B => C = (b: B) => f(a, b)
//
//def uncurry3[A, B, C, D](f: A => B => C => D): (A, B, C) => D = (a: A, b: B, c: C) => f(a)(b)(c)
