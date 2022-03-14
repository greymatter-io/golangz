package util

import (
	"fmt"
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

func TestPartial2(t *testing.T) {
	p := func(s1 string, s2 string, s3 int) string {
		return s1 + s2 + fmt.Sprintf("%v", s3)
	}

	var f = Partial2(p, "Hello")
	expected := "Hello World-12"
	var actual = f(" World-", 12)
	if actual != expected {
		t.Errorf("actual:%v, expected:%v", actual, expected)

	}
}
