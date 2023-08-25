package heap

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Heap[A any, B comparable] struct {
	hp          []*A
	heapLocator map[B]int //the type B represents the key type in the heap, providing a O(1) way of looking up the index of a key in the hp array(the heap)
}

func New[A any, B comparable]() Heap[A, B] {
	return Heap[A, B]{
		hp:          make([]*A, 0),
		heapLocator: make(map[B]int, 0),
	}
}

// i int - the index in the given heap of the parent of element i. Array indices start with the number zero.
// Performance - O(1)
func ParentIdx(i int) int {
	//Odd number
	if i%2 > 0 {
		return i / 2
	} else { // even number
		return (i / 2) - 1
	}
}

// Definition of almost-a-heap. Only one node in the tree has a value less than it's parent as per the lt function and that
// node is at the bottom rung of the heap.
// Definition of a heap.  Every node in the tree has a greater value than it's parent as per the lt function.
// This is a not pure function
// Parameters:
//
//	heap []A - the slice that is holding the heap
//	i int - the index into the heap of the element you want to move up. Array indices start with the number zero.
//	lt func(l, r A) bool - A predicate function that determines whether or not the left A element is less than the right A element.
//
// Returns - The modified heap (as a slice) that has the i'th element in its proper position in the heap
// Performance - O(log N) assuming that the array is almost-a-heap with the key: heap(i) too small.
func heapifyUp[A any, B comparable](h Heap[A, B], i int, lt func(l, r *A) bool) Heap[A, B] {
	if len(h.hp) == 0 {
		return h
	}
	if i > 0 {
		j := ParentIdx(i)
		if lt(h.hp[i], h.hp[j]) {
			//Swap elements
			temp := h.hp[i]
			temp2 := h.hp[j]
			h.hp[j] = temp
			h.hp[i] = temp2
			h = heapifyUp(h, j, lt)
		}
	}
	return h
}

// This is not a pure function because it modified the array each time.
// Parameters:
//
//	heap []A - the slice that is holding the heap
//	i int - the index into the heap of the element you want to move down. Array indices start with the number zero.TODO change
//	lt func(l, r A) bool - A predicate function that determines whether or not the left A element is less than the right A element.
//
// Returns - The original heap (as a slice) that has the i'th element in its proper position
// Performance - O(log N) assuming that the array is almost-a-heap with the key: heap(i) too big.
// func (e Heap[A, B]) heapifyDown(h []*A, i int, lt func(l, r *A) bool) []*A {
func heapifyDown[A any, B comparable](h Heap[A, B], i int, lt func(l, r *A) bool) Heap[A, B] {
	var j int
	n := len(h.hp)
	if (2*i)+1 > n {
		return h
	} else if (2*i)+1 < n {
		j = 0
		//These differ from book definition because array indices start with zero
		left := (2 * i) + 1
		right := (2 * i) + 2
		leftVal := h.hp[left]
		if right < n {
			rightVal := h.hp[right]
			if lt(leftVal, rightVal) {
				j = left
			} else {
				j = right
			}
		} else {
			j = left
		}
	} else if (2*i)+1 == n {
		j = (2 * i) + 1
	}
	if j < n && lt(h.hp[j], h.hp[i]) {
		//Swap elements
		temp := h.hp[i]
		temp2 := h.hp[j]
		h.hp[j] = temp
		h.hp[i] = temp2
		h = heapifyDown(h, j, lt)
	}
	return h
}

// This is a pure function.
// Parameters:
//
//	heap []A - the slice that is holding the heap
//
// Returns -the minimum element in the given heap without removing it. O(1)
// Performance - O(1)
func FindMin[A any, B comparable](h Heap[A, B]) (*A, error) {
	if len(h.hp) == 0 || h.hp[0] == nil {
		return nil, fmt.Errorf("heap is empty so findMin is therefore irrelevant")
	}
	return h.hp[0], nil
}

// Inserts the given element into the given heap and returns the modified heap.
//
// O(log n)
//
// Parameters:
//
//	heap []A - the slice that is holding the heap
//	a  A - the element you want to insert into the heap
//	lt func(l, r A) bool - A predicate function that determines whether or not the left A element is less than the right A element.
//
// Returns - The original heap (as a slice) that has the given element in its proper position
// Performance - O(log N)
// NOTE This function assumes that the heap slice has no empty elements. It always adds a new one.
func HeapInsert[A any, B comparable](h Heap[A, B], a *A, lt func(l, r *A) bool) Heap[A, B] {
	h.hp = append(h.hp, nil)
	l := len(h.hp) - 1
	h.hp[l] = a
	return heapifyUp(h, l, lt)
}

// Deletes an element from the given heap. This is not a pure function.
// Parameters:
//
//	heap []A - the slice that is holding the heap
//	i int - the index into the heap of the element you want to delete. Array indices start with the number zero.
//	lt func(l, r A) bool - A predicate function that determines whether or not the left A element is less than the right A element.
//
// Returns - The original heap that has the given element in its proper position
// Performance - O(log N)
func HeapDelete[A any, B comparable](h Heap[A, B], i int, lt func(l, r *A) bool) (Heap[A, B], error) {

	n := len(h.hp)
	if n == 0 {
		return h, nil
	}

	if i > len(h.hp)-1 {
		log.Errorf("element:%v you are trying to delete is longer than heap length: %v", i, len(h.hp)-1)
		return h, fmt.Errorf("element:%v you are trying to delete is longer than heap length: %v", i, len(h.hp)-1)
	}

	//Delete last and only element from heap
	if i == len(h.hp)-1 {
		h.hp = []*A{}
		return h, nil
	}
	h.hp[i] = h.hp[len(h.hp)-1]
	h.hp = h.hp[0 : len(h.hp)-1]
	if len(h.hp) == 1 {
		return h, nil
	}

	parent := ParentIdx(i)
	if parent > 0 && lt(h.hp[i], h.hp[parent]) {
		return heapifyUp(h, i, lt), nil
	} else {
		return heapifyDown(h, i, lt), nil
	}
}
