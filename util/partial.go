package util

//Takes a function of two parameters and turns it into a function with one parameter with the first parameter bound.
//This is useful replacement for a Class with instance variables in OOP.
func Partial1[T1, T2, T3 any](f func(T1, T2) T3, s1 T1) func(T2) T3 {
	return func(s2 T2) T3 {
		return f(s1, s2)
	}
}

//TODO
//def curry3[A, B, C, D](f: (A, B, C) => D): A => (B => (C => D)) = a => b => c => f(a, b, c)
//
//def partial1[A, B, C](a: A, f: (A, B) => C): B => C = (b: B) => f(a, b)
//
//def uncurry3[A, B, C, D](f: A => B => C => D): (A, B, C) => D = (a: A, b: B, c: C) => f(a)(b)(c)
