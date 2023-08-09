package option

// Represents optional values. Instances of Option are either an instance of Some or a type-safe empty structure None.
// The most idiomatic way to use an Option instance is to treat it as a collection or monad and use Map or FlatMap.
// If None is returned from any operation in the chain, the entire expression results in None[T].
// This allows for sophisticated chaining of Option values without having to check for the existence of a value.
type Option[A any] interface{}

type None[A any] struct{}

type Some[A any] struct {
	Value A
}

func Map[A, B any](a Option[A], f func(A) B) Option[B] {
	switch v := a.(type) {
	case Some[A]:
		return Some[B]{f(v.Value)}
	default:
		return None[B]{}
	}
}

func FlatMap[A, B any](a Option[A], f func(A) Option[B]) Option[B] {
	switch v := a.(type) {
	case Some[A]:
		return f(v.Value)
	default:
		return None[B]{}
	}
}

func GetOrElse[A any](a Option[A], def A) A {
	switch v := a.(type) {
	case Some[A]:
		return v.Value
	default:
		return def
	}
}
