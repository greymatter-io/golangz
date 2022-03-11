package function_composition

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
