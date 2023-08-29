package either

// Represents a value of one of two possible types (a disjoint union).
// An instance of Either is an instance of either Left or Right.
// A common use of Either is as an alternative to Option for dealing with possibly missing values.
// In this usage, None is replaced with a Left which can contain useful information.
// Right takes the place of Some. Convention dictates that Left is used for failure and Right is used for success.
// So left is often an error type but it does not have to be.
//
// For example, you could use Either[String, Int] to indicate whether a received input is a String or an Int.
// Either is right-biased, which means that Right is assumed to be the default case to operate on.
// If it is Left, operations like Map and FlatMap return the Left value unchanged:
type Either[A, B any] interface{}

type Left[A any] struct {
	Value A
}

type Right[B any] struct {
	Value B
}

func Map[A, B any](a Either[A, B], f func(B) B) Either[A, B] {
	switch v := a.(type) {
	case Right[B]:
		return Right[B]{f(v.Value)} //s
	default:
		return v
	}
}

func FlatMap[A, B any](a Either[A, B], f func(B) Either[A, B]) Either[A, B] {
	switch v := a.(type) {
	case Right[B]:
		return f(v.Value)
	default:
		return v
	}
}

func GetOrElse[A, B any](a Either[A, B], def A) Either[A, B] {
	switch v := a.(type) {
	case Right[B]:
		return v.Value
	default:
		return def
	}
}
