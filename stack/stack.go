package stack

import (
	"fmt"
	"github.com/greymatter-io/golangz/linked_list"
	"github.com/greymatter-io/golangz/option"
	"github.com/greymatter-io/golangz/propcheck"
)

// Cannot have generic alias with a type parameter as of Golang 1.22
type Stack[T any] struct {
	l *linked_list.LinkedList[T]
}

func PushO[T any](h T, s *Stack[T]) option.Option[*Stack[T]] {
	newStack := &Stack[T]{&linked_list.LinkedList[T]{}}
	if s == nil {
		newStack.l.Head = h
	} else {
		newStack.l.Head = h
		newStack.l.Tail = s.l
	}
	return option.Some[*Stack[T]]{newStack}
}

func Push[T any](h T, s *Stack[T]) *Stack[T] {
	newStack := &Stack[T]{&linked_list.LinkedList[T]{}}
	if s == nil {
		newStack.l.Head = h
		return newStack
	} else {
		newStack.l.Head = h
		newStack.l.Tail = s.l
		return newStack
	}
}

func PopO[T any](s *Stack[T]) option.Option[propcheck.Pair[T, *Stack[T]]] {
	newStack := &Stack[T]{&linked_list.LinkedList[T]{}}
	top, err := Peek(s)
	if err != nil {
		return option.None[T]{}
	} else {
		if s.l.Tail != nil {
			newStack.l.Head = s.l.Tail.Head
			newStack.l.Tail = s.l.Tail.Tail
		}
		r := option.Some[propcheck.Pair[T, *Stack[T]]]{propcheck.Pair[T, *Stack[T]]{top, newStack}}
		return r
	}
}

// TODO this should return an Option[Pair[T, Stack[T]]
func Pop[T any](s *Stack[T]) (T, *Stack[T], error) {
	newStack := &Stack[T]{&linked_list.LinkedList[T]{}}
	top, err := Peek(s)
	if err != nil {
		return top, newStack, err
	} else {
		if s.l.Tail != nil {
			newStack.l.Head = s.l.Tail.Head
			newStack.l.Tail = s.l.Tail.Tail
		}
		return top, newStack, nil
	}
}

// option.Option[propcheck.Pair[T, *Stack[T]]]
func PeekO[T any](s *Stack[T]) option.Option[T] {
	if s.l == nil {
		return option.None[T]{}
	} else {
		return option.Some[T]{s.l.Head}
	}
}

func Peek[T any](s *Stack[T]) (T, error) {
	if s.l == nil {
		return Zero[T](), fmt.Errorf("Cannot access top of an empty stack")
	} else {
		return s.l.Head, nil
	}
}

// TODO don;t need zero value if you use option
func Zero[T any]() T {
	var r T
	return r
}
