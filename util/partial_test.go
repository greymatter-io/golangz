package util

import (
	"testing"
)

func TestPartial1(t *testing.T) {
	p := func(s1 string, s2 string) string {
		return s1 + s2
	}

	var f = Partial1(p, "Hello")
	expected := "Hello World"
	var actual = f(" World")
	if actual != expected {
		t.Errorf("actual:%v, expected:%v", actual, expected)

	}
}
