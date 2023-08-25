package heap

import (
	"fmt"
	"github.com/greymatter-io/golangz/propcheck"
	"github.com/greymatter-io/golangz/sorting"
	"github.com/hashicorp/go-multierror"
	"testing"
	"time"
)

type Cache struct {
	key  int
	data string
}

func lt(l, r *Cache) bool {
	if l != nil && r != nil && l.key < r.key {
		return true
	} else {
		return false
	}
}

func keyGT(l, r *Cache) bool {
	if l.key > r.key {
		return true
	} else {
		return false
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

func parentIsLessThanOrEqual[A any, B comparable](h Heap[A, B], lastIdx int, parentGT func(l, r *A) bool) error {
	var pIdx = ParentIdx(lastIdx)
	var cIdx = lastIdx
	var errors error
	for pIdx >= 0 {
		if parentGT(h.hp[pIdx], h.hp[cIdx]) {
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

func insertIntoHeap(xss []int) Heap[Cache, string] {
	var r = New[Cache, string]()
	for _, x := range xss {
		r = HeapInsert(r, &Cache{x, fmt.Sprintf("key:%v", x)}, lt)
	}
	return r
}

func validateIsAHeap(p Heap[Cache, string]) (bool, error) {
	var errors error
	for idx := range p.hp {
		errors = parentIsLessThanOrEqual(p, idx, keyGT)
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
	g := propcheck.ChooseArray(0, 5, propcheck.ChooseInt(0, 10000))
	rng := propcheck.SimpleRNG{time.Now().Nanosecond()}

	prop := propcheck.ForAll(g,
		"Validate heapifyUp  \n",
		insert,
		validateIsAHeap, validateHeapMin,
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]int](t, result)
}

func TestHeapInsertWithNonEmptyHeapHeap(t *testing.T) {
	g := propcheck.ChooseArray(10, 1000, propcheck.ChooseInt(0, 10000))
	rng := propcheck.SimpleRNG{time.Now().Nanosecond()}

	prop := propcheck.ForAll(g,
		"Validate heapifyUp  \n",
		insert,
		validateIsAHeap, validateHeapMin,
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]int](t, result)
}

func TestHeapDeleteSpecificElements(t *testing.T) {
	var delete6ElementsFromHeapOfAtLeast6 = func(xss []int) Heap[Cache, string] {
		var h = insertIntoHeap(xss)
		h, _ = HeapDelete(h, 5, lt)
		h, _ = HeapDelete(h, 4, lt)
		h, _ = HeapDelete(h, 3, lt)
		h, _ = HeapDelete(h, 2, lt)
		h, _ = HeapDelete(h, 1, lt)
		h, _ = HeapDelete(h, 0, lt)
		return h
	}

	g0 := propcheck.ChooseArray(6, 15, propcheck.ChooseInt(1, 2000))
	rng := propcheck.SimpleRNG{time.Now().Nanosecond()}
	prop := propcheck.ForAll(g0,
		"Validate HeapDelete  \n",
		delete6ElementsFromHeapOfAtLeast6,
		validateIsAHeap, validateHeapMin,
	)
	result := prop.Run(propcheck.RunParms{100, rng}) //The 3rd iteration paniced with array out or bounds
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

	g0 := propcheck.ChooseArray(0, 1000, propcheck.ChooseInt(1, 200000))
	rng := propcheck.SimpleRNG{time.Now().Nanosecond()}
	prop := propcheck.ForAll(g0,
		"Validate HeapDelete  \n",
		deleteAllFromHeap,
		validateIsAHeap, heapWrong,
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]int](t, result)
}
