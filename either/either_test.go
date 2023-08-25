package either

import (
	"fmt"
	"testing"
)

func TestGetOrElseLeft(t *testing.T) {
	actual := GetOrElse[int, int](Left[int]{1}, 11)
	if actual != 11 {
		t.Errorf("Expected 11 actual:%v", actual)
	}
}

func TestGetOrElseRight(t *testing.T) {
	actual := GetOrElse[int, int](Right[int]{1}, 11)
	if actual != 1 {
		t.Errorf("Expected 1 actual:%v", actual)
	}
}

func TestGetOrElseLeftDifferentTypes(t *testing.T) {
	actual := GetOrElse[string, int](Left[string]{"1"}, "11")
	if actual != "11" {
		t.Errorf("Expected '11' actual:%v", actual)
	}
}

func TestGetOrElseRightDifferentTypes(t *testing.T) {
	actual := GetOrElse[string, int](Right[int]{1}, "11")
	if actual != 1 {
		t.Errorf("Expected 1 actual:%v", actual)
	}
}

func TestMapRight(t *testing.T) {
	a := Right[string]{"1"}

	concat := func(x string) string {
		return fmt.Sprintf("%v%v", x, x)
	}
	b := Map[int, string](a, concat)

	switch v := b.(type) {
	case Right[string]:
		if v.Value != "11" {
			t.Errorf("Actual:%v, Expected:%v", v.Value, "11")
		}
	default:
		t.Errorf("Expected type of Either to be Right but was:%T", v)
	}
}

func TestMapLeft(t *testing.T) {
	a := Left[int]{1}

	concat := func(x string) string {
		return fmt.Sprintf("%v%v", x, x)
	}
	b := Map[int, string](a, concat)

	switch v := b.(type) {
	case Left[int]:
	default:
		t.Errorf("Expected type of Either to be Left but was:%T", v)
	}
}

func TestMapLeftWithErrorLeft(t *testing.T) {
	a := Left[error]{fmt.Errorf("what the heck happened?")}

	concat := func(x string) string {
		return fmt.Sprintf("%v%v", x, x)
	}
	b := Map[error, string](a, concat)

	switch v := b.(type) {
	case Left[error]:
	default:
		t.Errorf("Expected type of Either to be Left[error] but was:%T", v)
	}
}

func TestMapChainingWithLeft(t *testing.T) {
	a := Left[int]{1}

	hello := func(x string) string {
		return fmt.Sprintf("hello:%v", x)
	}
	concat := func(x string) string {
		return fmt.Sprintf("%v%v", x, x)
	}
	b := Map[int, string](Map[int, string](a, hello), concat)

	switch v := b.(type) {
	case Left[int]:
	default:
		t.Errorf("Expected type of Either to be Left but was:%T", v)
	}
}

func TestMapChainingWithErrorLeft(t *testing.T) {
	a := Left[error]{fmt.Errorf("what the heck happened?")}

	hello := func(x string) string {
		return fmt.Sprintf("hello:%v", x)
	}
	concat := func(x string) string {
		return fmt.Sprintf("%v%v", x, x)
	}
	b := Map[error, string](Map[error, string](a, hello), concat)

	switch v := b.(type) {
	case Left[error]:
	default:
		t.Errorf("Expected type of Either to be Left but was:%T", v)
	}
}

func TestMapChainingWithRight(t *testing.T) {
	a := Right[string]{"1"}

	hello := func(x string) string {
		return fmt.Sprintf("hello:%v", x)
	}
	concat := func(x string) string {
		return fmt.Sprintf("%v%v", x, x)
	}
	b := Map[int, string](Map[int, string](a, hello), concat) //You don't have to put a type parameter on outer one but it documents it better,

	switch v := b.(type) {
	case Right[string]:
		if v.Value != "hello:1hello:1" {
			t.Errorf("Actual:%v, Expected:%v", v.Value, "hello:1hello:1")
		}
	default:
		t.Errorf("Expected type of Either to be Rightbut was:%T", v)
	}
}

func TestFlatMapRight(t *testing.T) {
	a := Right[string]{"1"}

	concat := func(x string) Either[int, string] {
		return Right[string]{fmt.Sprintf("%v%v", x, x)}
	}
	b := FlatMap[int, string](a, concat)

	switch v := b.(type) {
	case Right[string]:
		if v.Value != "11" {
			t.Errorf("Actual:%v, Expected:%v", v.Value, "11")
		}
	default:
		t.Errorf("Expected type of Either to be Right but was:%T", v)
	}
}

func TestFlatMapLeft(t *testing.T) {
	a := Left[int]{1}

	concat := func(x string) Either[int, string] {
		return Right[string]{fmt.Sprintf("%v%v", x, x)}
	}
	b := FlatMap[int, string](a, concat)

	switch v := b.(type) {
	case Left[int]:
	default:
		t.Errorf("Expected type of Either to be Left but was:%T", v)
	}
}

func TestFlatMapChainingWithLeft(t *testing.T) {
	a := Left[int]{1}

	hello := func(x string) Either[int, string] {
		return Right[string]{fmt.Sprintf("hello:%v", x)}
	}
	concat := func(x string) Either[int, string] {
		return Right[string]{fmt.Sprintf("%v%v", x, x)}
	}
	b := FlatMap(FlatMap[int, string](a, hello), concat)

	switch v := b.(type) {
	case Left[int]:
	default:
		t.Errorf("Expected type of Either to be Left but was:%T", v)
	}
}

func TestFlatMapChainingWithRight(t *testing.T) {
	a := Right[string]{"1"}

	hello := func(x string) Either[int, string] {
		return Right[string]{fmt.Sprintf("hello:%v", x)}
	}
	concat := func(x string) Either[int, string] {
		return Right[string]{fmt.Sprintf("%v%v", x, x)}
	}
	b := FlatMap[int, string](FlatMap[int, string](a, hello), concat) //You don't have to put a type parameter on outer one but it documents it better,

	switch v := b.(type) {
	case Right[string]:
		if v.Value != "hello:1hello:1" {
			t.Errorf("Actual:%v, Expected:%v", v.Value, "hello:1hello:1")
		}
	default:
		t.Errorf("Expected type of Either to be Right but was:%T", v)
	}
}

func TestFlatMapChainingWithErrorLeft(t *testing.T) {
	a := Right[string]{"1"}

	hello := func(x string) Either[error, string] {
		return Right[string]{fmt.Sprintf("hello:%v", x)}
	}
	concat := func(x string) Either[error, string] {
		return Left[error]{fmt.Errorf("what the heck happened?")}
	}
	b := FlatMap(FlatMap[error, string](a, hello), concat)

	switch v := b.(type) {
	case Left[error]:
	default:
		t.Errorf("Expected type of Either to be Left but was:%T", v)
	}
}
