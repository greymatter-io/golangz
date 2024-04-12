package stack

// An Immutable stack, mostly from here: https://www.reddit.com/r/golang/comments/13hi00u/should_constructor_functions_always_return_a/
// the pointer type
type node[T any] struct {
	val  T
	next *node[T]
}

// the value wrapper
type Stack[T any] struct {
	root *node[T]
}

func NewStack[T any]() Stack[T] {
	return Stack[T]{}
}

// returns values - hides the pointers
func (s Stack[T]) Push(val T) Stack[T] {
	newRoot := node[T]{val: val, next: s.root}
	return Stack[T]{root: &newRoot}
}

func (s Stack[T]) Pop() Stack[T] {
	if s.root == nil {
		return s
	}
	return Stack[T]{root: s.root.next}
}

func (s Stack[T]) IsEmpty() bool {
	if s.root == nil {
		return true
	} else {
		return false
	}
}

func (s Stack[T]) Peek() (val T, ok bool) {
	if s.root == nil {
		return val, false
	}
	return s.root.val, true
}
