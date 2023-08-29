package option

import (
	"fmt"
	"testing"
)

func TestMapSome(t *testing.T) {
	a := Some[int]{1}

	addOne := func(x int) string {
		return fmt.Sprintf("%v", x+1)
	}
	b := Map[int](a, addOne)

	switch v := b.(type) {
	case Some[string]:
		if v.Value != "2" {
			t.Errorf("Actual:%v, Expected:%v", v.Value, "2")
		}
	default:
		t.Errorf("Expected type of Option to be Some but was:%T", v)
	}
}

func TestMapNone(t *testing.T) {
	addOne := func(x int) string {
		return fmt.Sprintf("%v", x+1)
	}
	b := Map[int](None[int]{}, addOne)

	switch v := b.(type) {
	case None[string]:
	default:
		t.Errorf("Expected type of Option to be None but was:%T", v)
	}
}

func TestGetOrElseNoneChangeContainedType(t *testing.T) {
	addOne := func(x int) string {
		return fmt.Sprintf("%v", x+1)
	}
	b := Map[int](None[int]{}, addOne)

	actual := GetOrElse(b, "12")
	if actual != "12" {
		t.Errorf("Expected '12' actual:%v", actual)
	}
}

func TestGetOrElseSomeChangeContainedType(t *testing.T) {
	addOne := func(x int) string {
		return fmt.Sprintf("%v", x+1)
	}
	b := Map[int](Some[int]{12}, addOne)

	actual := GetOrElse(b, "12")
	if actual != "13" {
		t.Errorf("Expected '13' actual:%v", actual)
	}
}

func TestGetOrElseNone(t *testing.T) {
	addOne := func(x int) int {
		return x + 1
	}
	b := Map[int](None[int]{}, addOne)

	actual := GetOrElse(b, 12)
	if actual != 12 {
		t.Errorf("Expected '12' actual:%v", actual)
	}
}

func TestGetOrElseSome(t *testing.T) {
	addOne := func(x int) int {
		return x + 1
	}
	b := Map[int](Some[int]{12}, addOne)

	actual := GetOrElse(b, 13)
	if actual != 13 {
		t.Errorf("Expected 13 actual:%v", actual)
	}
}

func TestMapChainingWithNone(t *testing.T) {

	addOne := func(x int) string {
		return fmt.Sprintf("%v", x+1)
	}
	concat := func(x string) string {
		return fmt.Sprintf("%v%v", x, x)
	}
	b := Map(Map[int](None[int]{}, addOne), concat)

	switch v := b.(type) {
	case None[string]:
	default:
		t.Errorf("Expected type of Option to be None but was:%T", v)
	}
}

func TestMapChainingWithSome(t *testing.T) {
	a := Some[int]{1}

	addOne := func(x int) string {
		return fmt.Sprintf("%v", x+1)
	}
	concat := func(x string) string {
		return fmt.Sprintf("%v%v", x, x)
	}
	b := Map(Map[int](a, addOne), concat)

	switch v := b.(type) {
	case Some[string]:
		if v.Value != "22" {
			t.Errorf("Actual:%v, Expected:%v", v.Value, "22")
		}
	default:
		t.Errorf("Expected type of Option to be Some but was:%T", v)
	}
}

func TestFlatMapSome(t *testing.T) {
	a := Some[int]{1}

	addOne := func(x int) Option[string] {
		return Some[string]{fmt.Sprintf("%v", x+1)}
	}
	b := FlatMap[int](a, addOne)

	switch v := b.(type) {
	case Some[string]:
		if v.Value != "2" {
			t.Errorf("Actual:%v, Expected:%v", v.Value, "2")
		}
	default:
		t.Errorf("Expected type of Option to be Some but was:%T", v)
	}
}

func TestFlatMapNone(t *testing.T) {
	a := None[string]{}

	addOne := func(x int) Option[string] {
		return Some[string]{fmt.Sprintf("%v", x+1)}
	}
	b := FlatMap[int](a, addOne)

	switch v := b.(type) {
	case None[string]:
	default:
		t.Errorf("Expected type of Option to be None but was:%T", v)
	}
}

func TestFlatMapChainingWithNone(t *testing.T) {
	a := None[string]{}

	addOne := func(x int) Option[string] {
		return Some[string]{fmt.Sprintf("%v", x+1)}
	}
	concat := func(x string) Option[string] {
		return Some[string]{fmt.Sprintf("%v%v", x, x)}
	}
	b := FlatMap[string](FlatMap[int](a, addOne), concat)

	switch v := b.(type) {
	case None[string]:
	default:
		t.Errorf("Expected type of Option to be None but was:%T", v)
	}
}

func TestFlatMapChainingWithSome(t *testing.T) {
	a := Some[int]{1}

	addOne := func(x int) Option[string] {
		return Some[string]{fmt.Sprintf("%v", x+1)}
	}
	concat := func(x string) Option[string] {
		return Some[string]{fmt.Sprintf("%v%v", x, x)}
	}
	b := FlatMap[string](FlatMap[int](a, addOne), concat)

	switch v := b.(type) {
	case Some[string]:
		if v.Value != "22" {
			t.Errorf("Actual:%v, Expected:%v", v.Value, "22")
		}
	default:
		t.Errorf("Expected type of Option to be Some but was:%T", v)
	}
}
