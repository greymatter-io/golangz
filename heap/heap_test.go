package heap

import (
	"fmt"
	"github.com/greymatter-io/golangz/propcheck"
	"github.com/greymatter-io/golangz/sets"
	"github.com/greymatter-io/golangz/sorting"
	"github.com/hashicorp/go-multierror"
	"testing"
	"time"
)

type Cache struct {
	key   int
	value string
}

func ltInt(l, r int) bool {
	if l < r {
		return true
	} else {
		return false
	}
}
func eqInt(l, r int) bool {
	if l == r {
		return true
	} else {
		return false
	}
}

func lt(l, r *Cache) bool {
	if l != nil && r != nil && l.key < r.key {
		return true
	} else {
		return false
	}
}

func parentLT(parent, child *Cache) bool { //If parent is greater this is an error
	if parent.key > child.key {
		return false
	} else {
		return true
	}
}

func eqTs(l, r *Cache) bool {
	if l.key == r.key {
		return true
	} else {
		return false
	}
}

// Find the minimum and compares it to the actual min in the initial array.
// If array is empty or filled with nil pointers that is OK and FindMin should not fail
// but return a Golang error, not panic.
func minimumCorrectValue[A any, B comparable](h Heap[A, B], sorted []*A, eq func(l, r *A) bool) bool {
	key, err := FindMin(h)
	if len(h.hp) > 0 && err == nil {
		return eq(key, sorted[0])
	} else {
		return true
	}
}

func parentIsLessThanOrEqual[A any, B comparable](h Heap[A, B], lastIdx int, parentLT func(l, r *A) bool) error {
	var pIdx = ParentIdx(lastIdx)
	var cIdx = lastIdx
	var errors error
	for pIdx >= 0 {
		if !parentLT(h.hp[pIdx], h.hp[cIdx]) {
			errors = multierror.Append(errors, fmt.Errorf("parent:%v value was not less than or equal to child's:%v\n", h.hp[pIdx], h.hp[cIdx]))
		}
		cIdx = pIdx
		pIdx = ParentIdx(cIdx)
	}
	return errors
}
func insert(p []int) Heap[Cache, string] {
	xss := insertIntoHeap(p)
	return xss
}

var elementKeyExtractor = func(c *Cache) string {
	return c.value
}

func insertIntoHeap(xss []int) Heap[Cache, string] {
	var h = New[Cache, string](elementKeyExtractor)
	for _, x := range xss {
		h = HeapInsert(h, &Cache{x, fmt.Sprintf("key:%v", x)}, lt)
	}
	return h
}

func validateIsAHeap(p Heap[Cache, string]) (bool, error) {
	var errors error
	for idx, c := range p.hp {
		errors = parentIsLessThanOrEqual(p, idx, parentLT)
		k := p.elementKeyExtractor(p.hp[idx])
		if c.value != k {
			errors = multierror.Append(errors, fmt.Errorf("Expected heap locator key value:%v using heap locator to equal heap key value:%v", k, c.value))
		}
	}
	if len(p.hp) != len(p.position) {
		errors = multierror.Append(errors, fmt.Errorf("Heap locator map:%v should have been same length as heap:%v", len(p.position), len(p.hp)))
	}
	if errors != nil {
		return false, errors
	} else {
		return true, nil
	}
}
func validateHeapMin(p Heap[Cache, string]) (bool, error) {
	var errors error
	var sorted = make([]*Cache, len(p.hp))
	copy(sorted, p.hp)
	sorting.QuickSort(sorted, lt)
	if !minimumCorrectValue(p, sorted, eqTs) {
		errors = multierror.Append(errors, fmt.Errorf("FindMin should have returned:%v", sorted[0]))
	}
	if errors != nil {
		return false, errors
	} else {
		return true, nil
	}
}

func TestHeapInsertWithEmptyHeap(t *testing.T) {
	ge := propcheck.ChooseInt(0, 10000)
	g := sets.ChooseSet(0, 5, ge, ltInt, eqInt)
	rng := propcheck.SimpleRNG{time.Now().Nanosecond()}

	prop := propcheck.ForAll(g,
		"Validate heapifyUp  \n",
		insert,
		validateIsAHeap, validateHeapMin,
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]int](t, result)
}

func TestHeapInsertWithNonEmptyHeap(t *testing.T) {
	ge := propcheck.ChooseInt(0, 1000000)
	g := sets.ChooseSet(10, 1000, ge, ltInt, eqInt)
	rng := propcheck.SimpleRNG{time.Now().Nanosecond()}

	prop := propcheck.ForAll(g,
		"Validate heapifyUp  \n",
		insert,
		validateIsAHeap, validateHeapMin,
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]int](t, result)
}

func TestHeapDeleteEveryElementStartingFromLast(t *testing.T) {
	var delete6ElementsFromHeapOf6 = func(xss []int) Heap[Cache, string] {
		var h = insertIntoHeap(xss)
		h, _ = HeapDelete(h, 5, lt)
		h, _ = HeapDelete(h, 4, lt)
		h, _ = HeapDelete(h, 3, lt)
		h, _ = HeapDelete(h, 2, lt)
		h, _ = HeapDelete(h, 1, lt)
		h, _ = HeapDelete(h, 0, lt)
		return h
	}

	ge := propcheck.ChooseInt(0, 1000000)
	g0 := sets.ChooseSet(6, 6, ge, ltInt, eqInt)
	rng := propcheck.SimpleRNG{time.Now().Nanosecond()}
	prop := propcheck.ForAll(g0,
		"Validate HeapDelete  \n",
		delete6ElementsFromHeapOf6,
		validateIsAHeap, validateHeapMin,
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]int](t, result)
}

func TestHeapDeleteMinElement(t *testing.T) {
	var errors error
	correctHeapMin := func(p Heap[Cache, string]) bool {
		var sorted = make([]*Cache, len(p.hp))
		copy(sorted, p.hp)
		sorting.QuickSort(sorted, lt)
		min, err := FindMin(p)
		if len(p.hp) == 0 {
			return true
		} else if err != nil {
			return false
		} else if sorted[0].key != min.key {
			return false
		} else {
			return true
		}
	}

	var deleteAllFromHeap = func(xss []int) Heap[Cache, string] {
		var h = insertIntoHeap(xss)
		for range h.hp {
			h, _ = HeapDelete(h, 0, lt)
			if !correctHeapMin(h) {
				errors = multierror.Append(errors, fmt.Errorf("Heap property violated"))
			}
		}
		return h
	}

	heapWrong := func(p Heap[Cache, string]) (bool, error) {
		if errors != nil {
			return false, errors
		} else {
			return true, nil
		}
	}

	ge := propcheck.ChooseInt(0, 1000000)
	g0 := sets.ChooseSet(0, 1000, ge, ltInt, eqInt)
	rng := propcheck.SimpleRNG{time.Now().Nanosecond()}
	prop := propcheck.ForAll(g0,
		"Validate HeapDelete  \n",
		deleteAllFromHeap,
		validateIsAHeap, heapWrong,
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]int](t, result)
}

func TestDeleteFromEmptyHeap(t *testing.T) {
	h, err := HeapDelete(New[Cache, string](elementKeyExtractor), 0, lt)

	if err == nil {
		t.Errorf("Should have gotten an error trying to delete from an empty heap")
	}
	if len(h.hp) != 0 {
		t.Errorf("heap should have been empty")
	}
	if len(h.position) != 0 {
		t.Errorf("heaplocator map should have been empty")
	}
}

func TestDeletePastLastElement(t *testing.T) {
	var h = New[Cache, string](elementKeyExtractor)
	h = HeapInsert(h, &Cache{12, fmt.Sprintf("key:%v", 12)}, lt)
	h, err := HeapDelete(h, 2, lt)

	if err == nil {
		t.Errorf("Should have gotten an error trying to delete past the end of the heap")
	}
}
