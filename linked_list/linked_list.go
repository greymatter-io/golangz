package linked_list

import "fmt"

type LinkedList[T any] struct {
	head T
	tail *LinkedList[T]
}

func (w *LinkedList[T]) String() string {
	return fmt.Sprintf("LinkedList {head: %v, tail: %v}", w.head, w.tail)
}

func Push[T any](h T, l *LinkedList[T]) *LinkedList[T] {
	if l == nil {
		return &LinkedList[T]{
			head: h,
			tail: nil,
		}
	} else {
		return &LinkedList[T]{
			head: h,
			tail: l,
		}
	}
}

func AddLast[T any](h T, l *LinkedList[T]) *LinkedList[T] {
	newTail := &LinkedList[T]{
		head: h,
		tail: nil,
	}
	if l == nil {
		return newTail
	} else {
		current := l
		for current.tail != nil {
			current = current.tail
		}
		current.tail = newTail
		return l
	}
}

func Tail[T any](l *LinkedList[T]) (*LinkedList[T], error) {
	if Len(l) == 0 {
		return nil, fmt.Errorf("Cannot Tail an empty list")
	} else {
		return l.tail, nil //r.Tail, nil
	}
}

// TODO Head of Nil list will throw a NPE. Fix this
func Head[T any](l *LinkedList[T]) T {
	return l.head
}

func Drop[T any](l *LinkedList[T], n int) *LinkedList[T] {
	if n <= 0 {
		return l
	} else {
		if l == nil {
			return nil
		} else {
			return Drop(l.tail, n-1)
		}
	}
}

func Zero[T any]() *LinkedList[T] {
	return nil
}

func internalAddWhile[T any](l *LinkedList[T], r *LinkedList[T], p func(T) bool) *LinkedList[T] {
	if l == nil || !p(l.head) {
		return r
	} else {
		return internalAddWhile(l.tail, AddLast(l.head, r), p)
	}
}

// Evaluates elements of given list, adding elements to head of a new list until predicate returns false, returning the new list and preserving ordering of original list.
// Note that this is different than filter. The algorithm stops appending to the resulting list when the predicate returns false.
func AddWhile[T any](l *LinkedList[T], p func(T) bool) *LinkedList[T] {
	return internalAddWhile(l, Zero[T](), p)
}

// Evaluates elements of given list and filters out those elements which fail predicate, preserving order of original list.
func Filter[T any](l *LinkedList[T], p func(T) bool) *LinkedList[T] {
	var g = func(h T, accum *LinkedList[T]) *LinkedList[T] {
		if p(h) {
			return Push(h, accum)
		} else {
			return accum
		}
	}
	xs := FoldRight(l, Zero[T](), g)
	return xs
}

func ToList[T any](xs []T) *LinkedList[T] {
	var r = Zero[T]()
	for _, x := range xs {
		r = AddLast(x, r)
	}
	return r
}

func internalLen[T any](l *LinkedList[T], n int) int {
	if l == nil {
		return n
	} else {
		return internalLen(l.tail, n+1)
	}
}

func Len[T any](l *LinkedList[T]) int {
	return internalLen(l, 0)
}

func FoldRight[A, B any](l *LinkedList[A], z B, f func(A, B) B) B {
	if l == nil {
		return z
	} else {
		return f(l.head, FoldRight(l.tail, z, f))
	}
}

func FoldLeft[A, B any](l *LinkedList[A], z B, f func(B, A) B) B {
	if l == nil {
		return z
	} else {
		return FoldLeft(l.tail, f(z, l.head), f)
	}
}

func ToArray[A any](l *LinkedList[A]) []A {
	f2 := func(accum []A, s A) []A {
		return append(accum, s)
	}
	arr := []A{}
	fConcat := FoldLeft(l, arr, f2)
	return fConcat
}
