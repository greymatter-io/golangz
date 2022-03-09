package propcheck

import (
	"testing"
	"time"
)

func TestMap2WithInt(t *testing.T) {
	rng := SimpleRNG{time.Now().Nanosecond()}
	ra := Id(12)
	rb := Id(13)
	f := func(a int, b int) int {
		return a + b
	}
	res := Map2(ra, rb, f)
	actual, _ := res(rng)
	if actual != 25 {
		t.Errorf("Map2 should have summed the A and B inside the SimpleRNG: %v but resulted in %v \n", 12+13, actual)
	}
}

func TestMap2ChooseInt(t *testing.T) {
	var rng = SimpleRNG{time.Now().Nanosecond()}
	start := 1
	endExclusive := 5
	r1 := ChooseInt(start, endExclusive)
	r2 := ChooseInt(start, endExclusive)
	f := func(a int, b int) []int {
		return []int{a, b}
	}

	r3 := Map2(r1, r2, f)
	a, _ := r3(rng)
	actual := a
	if len(actual) != 2 || actual[0] < start || actual[0] > (endExclusive-1) || actual[1] < start || actual[1] > (endExclusive-1) {
		t.Errorf("Map2 should have produced an array of two numbers >= 1 and < 5 but produced %v \n", actual)
	}
}

func TestMap3WithInt(t *testing.T) {
	rng := SimpleRNG{time.Now().Nanosecond()}
	ra := Id(12)
	rb := Id(13)
	rc := Id(14)
	f := func(a int, b int, c int) int {
		return a - b + c
	}
	res := Map3(ra, rb, rc, f)
	actual, _ := res(rng)
	if actual != 13 {
		t.Errorf("Map3 should have summed A - B + c to 13 but resulted in %v \n", actual)
	}
}

func TestMap4WithInt(t *testing.T) {
	rng := SimpleRNG{time.Now().Nanosecond()}
	ra := Id(12)
	rb := Id(13)
	rc := Id(14)
	rd := Id(15)
	f := func(a int, b int, c int, d int) int {
		return a + (b - c) + d
	}
	res := Map4(ra, rb, rc, rd, f)
	actual, _ := res(rng)
	if actual != 26 {
		t.Errorf("Map4 should have summed A  + (B  - c) + d to 26 but resulted in %v \n", actual)
	}
}
