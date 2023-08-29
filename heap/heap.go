package heap

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// A generic heap that supports O(log n) insert, delete and change key operation as well as O(1) find minimum value operations.
//
//Invariants:
//   1. Make sure your heap keys are unique, otherwise the heap locator replaceKey function will not work because it
//  relies on a Golang map.

// A generic heap containing the heap's underlying array and a corresponding map that allows
// lookup of the key's index in the underlying array with O(1) cost.
// Also contains a function that allows extraction of the A's key value for insertion into the position map.
type Heap[A any, B comparable] struct {
	hp                  []*A
	position            map[B]int  //the type B represents the key type in the heap, providing a O(1) way of looking up the index of a key in the hp array(the heap)
	elementKeyExtractor func(*A) B //A function that extracts the key of the given hp element from an A instance.
}

func New[A any, B comparable](keyExtractor func(*A) B) Heap[A, B] {
	return Heap[A, B]{
		hp:                  make([]*A, 0),
		position:            make(map[B]int, 0),
		elementKeyExtractor: keyExtractor,
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
//	h - the generic heap object containing the heap(represented as a slice) and the reverse-lookup map.
//	i int - the index into the heap of the element you want to move up. Array indices start with the number zero.
//	lt func(l, r A) bool - A predicate function that determines whether the left A element is less than the right A element.
//
// Returns - The modified heap that has the i'th element in its proper position in the heap
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
			aa := h.elementKeyExtractor(h.hp[j])
			bb := h.elementKeyExtractor(h.hp[i])
			h.position[aa] = j
			h.position[bb] = i
			h = heapifyUp(h, j, lt)
		}
	}
	return h
}

// This is not a pure function because it modified the array each time.
// Parameters:
//
//	h - the generic heap object containing the heap(represented as a slice) and the reverse-lookup map.
//	i int - the index into the heap of the element you want to move up. Array indices start with the number zero.
//	lt func(l, r A) bool - A predicate function that determines whether the left A element is less than the right A element.
//
// Returns - The modified heap that has the i'th element in its proper position in the heap
// Performance - O(log N) assuming that the array is almost-a-heap with the key: heap(i) too big.
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
		///
		aa := h.elementKeyExtractor(h.hp[j])
		bb := h.elementKeyExtractor(h.hp[i])
		h.position[aa] = j
		h.position[bb] = i
		///
		h = heapifyDown(h, j, lt)
	}
	return h
}

// This is a pure function.
// Parameters:
//
//	h - the generic heap object containing the heap(represented as a slice) and the reverse-lookup map.
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
//	h - the generic heap object containing the heap(represented as a slice) and the reverse-lookup map.
//	a  A - the element you want to insert into the heap
//	lt func(l, r A) bool - A predicate function that determines whether or not the left A element is less than the right A element.
//
// Returns - The original heap (as a slice) that has the given element in its proper position
// Performance - O(log N)
// NOTE This function assumes that the heap slice has no empty elements. It always adds a new one.
func HeapInsert[A any, B comparable](h Heap[A, B], a *A, lt func(l, r *A) bool) Heap[A, B] {
	h.hp = append(h.hp, nil) //Adds an empty element at end
	l := len(h.hp) - 1       //Get index of end of heap nd stick new element there
	h.hp[l] = a

	aa := h.elementKeyExtractor(h.hp[l])
	h.position[aa] = l

	//Now move it up as necessary until that part of tree satisfies heap property
	return heapifyUp(h, l, lt)
}

// TODO this is not complete and its API may change. And I have not tested it.  Still fumbling around thinking about what I do with ther newKey
func ChangeKey[A, B comparable](h Heap[A, B], currentKey, newKey *A, lt func(l, r *A) bool) Heap[A, B] {
	aa := h.elementKeyExtractor(currentKey)
	l := h.position[aa]

	h.position[aa] = l

	parent := ParentIdx(l)
	if parent > 0 && lt(h.hp[l], h.hp[parent]) {
		return heapifyUp(h, l, lt)
	} else {
		return heapifyDown(h, l, lt)
	}
}

// Deletes an element from the given heap. This is not a pure function.
// Parameters:
//
//	h - the generic heap object containing the heap(represented as a slice) and the reverse-lookup map.
//	i int - the index into the heap of the element you want to delete. Array indices start with the number zero.
//	lt func(l, r A) bool - A predicate function that determines whether or not the left A element is less than the right A element.
//
// Returns - The original heap that has the given element in its proper position.
//
//	If the heao is empty or the indice you are trying to delete is longer than the heap(zero indexed) then you get an error
//
// Performance - O(log N)
func HeapDelete[A any, B comparable](h Heap[A, B], i int, lt func(l, r *A) bool) (Heap[A, B], error) {
	if i > len(h.hp)-1 || len(h.hp) == 0 {
		log.Errorf("element:%v you are trying to delete is longer than heap length: %v", i, len(h.hp)-1)
		return h, fmt.Errorf("element:%v you are trying to delete is longer than heap length: %v", i, len(h.hp)-1)
	}

	//Delete last element and return. No need to move anything around.
	if i == len(h.hp)-1 {
		k := h.elementKeyExtractor(h.hp[i])
		delete(h.position, k)
		h.hp = h.hp[0 : len(h.hp)-1]
		return h, nil
	}
	k := h.elementKeyExtractor(h.hp[i])
	delete(h.position, k)
	h.hp[i] = h.hp[len(h.hp)-1]
	j := h.elementKeyExtractor(h.hp[i])
	h.position[j] = i
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
