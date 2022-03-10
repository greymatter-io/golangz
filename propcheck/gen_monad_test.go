package propcheck

import (
	"github.com/go-test/deep"
	"testing"
	"time"
)

func TestSameGeneratorProducesSameRandomNumberEveryTime(t *testing.T) {
	rng1 := SimpleRNG{Seed: time.Now().Nanosecond()}
	n1, rng2 := NextInt(rng1)
	n2, rng3 := NextInt(rng1)
	if n1 != n2 {
		t.Errorf("NextInt using the same generator should have produced the same number but n1 was %v and n2 was %v. \n", n1, n2)
	}
	if rng2 != rng3 {
		t.Errorf("NextInt should not mutate its underlying generator \n")
	}
}

func TestDifferentGeneratorProducesDifferentRandomNumbersEveryTime(t *testing.T) {
	rng1 := SimpleRNG{Seed: time.Now().Nanosecond()}
	n1, rng2 := NextInt(rng1)
	n2, _ := NextInt(rng2)
	if n1 == n2 {
		t.Error("Generating two random numbers using A different generator should have produced different numbers but did not \n")
	}
}

func TestNextIntProducesDifferentGeneratorsEveryTime(t *testing.T) {
	rng1 := SimpleRNG{Seed: time.Now().Nanosecond()}
	_, rng2 := NextInt(rng1)
	if rng1 == rng2 {
		t.Error("NextInt should produce different generators every time \n")
	}
}

func TestIdIsAPureFunction(t *testing.T) {
	u := Id("hello")
	rng1 := SimpleRNG{Seed: time.Now().Nanosecond()}
	r, rng2 := u(rng1)
	r2, rng3 := u(rng1)
	n3, _ := NextInt(rng1)
	n4, _ := NextInt(rng1)
	n5, _ := NextInt(rng1)
	if rng2 != rng3 {
		t.Errorf("Generating two random numbers with the same Seed should have produced A different generator each time \n")
	}
	if r != "hello" && r2 != r {
		t.Errorf("Wrapping an arbitrary value with Id gives should give you the value back wrapped in A new Random Number Generator \n")
	}
	if !(n3 == n4 && n4 == n5) {
		t.Errorf("Id has side effects \n")
	}
}

func TestFlatMapWithInt(t *testing.T) {
	rng := SimpleRNG{time.Now().Nanosecond()}
	r := Id(12)
	g := func(x int) func(SimpleRNG) (int, SimpleRNG) {
		return Id(x + 1)
	}
	res := FlatMap(r, g)
	actual, rng2 := res(rng)
	if actual != 13 {
		t.Errorf("Map should have incremented the unit value by 1 \n")
	}
	if rng != rng2 {
		t.Error("Should have gotten the same SimpleRNG because you just flatmapped over A Id generator")
	}
}

func TestFlatMapWithStringArray(t *testing.T) {
	rng := SimpleRNG{time.Now().Nanosecond()}
	r := Id([]string{"asd", "aDS"})
	g := func(x []string) func(SimpleRNG) ([]string, SimpleRNG) {
		return Id(append(x, "dude"))
	}
	res := FlatMap(r, g)
	actual, rng2 := res(rng)
	expectedA := []string{"asd", "aDS", "dude"}
	if diff := deep.Equal(actual, expectedA); diff != nil {
		t.Error(diff)
	}
	if rng != rng2 {
		t.Error("Should have gotten the same SimpleRNG because you just flatmapped over A Id generator")
	}
}

func TestMap(t *testing.T) {
	rng := SimpleRNG{time.Now().Nanosecond()}
	r := Id(12)
	g := func(x int) int {
		return x + 1
	}
	res := Map(r, g)
	actual, rng2 := res(rng)
	if actual != 13 {
		t.Errorf("Map should have incremented the unit value by 1 \n")
	}
	if rng != rng2 {
		t.Error("Should have gotten the same SimpleRNG because you just mapped over A Id generator")
	}
}

func TestGenerateSequenceOfRandomInts(t *testing.T) {
	rSize := 1000
	rng := SimpleRNG{time.Now().Nanosecond()}
	start := 1
	endExclusive := 500
	var s []func(SimpleRNG) (int, SimpleRNG)
	for x := 0; x < rSize; x++ {
		s = append(s, ChooseInt(start, endExclusive))
	}
	u := Sequence(s)
	actual, _ := u(rng)
	if len(actual) != rSize {
		t.Errorf("Sequence should have produced an array of %v numbers but produced an array of size:%v \n", rSize, len(actual))
	}
	for _, v := range actual {
		c := v
		if c < start || c > (endExclusive-1) {
			t.Errorf("Sequence should have produced an array of numbers >= 1 and < 5 but produced %v \n", actual)
		}
	}
}

func TestGenerateSequenceOfRandomFloats(t *testing.T) {
	rSize := 1000
	rng := SimpleRNG{time.Now().Nanosecond()}
	var s []func(SimpleRNG) (float64, SimpleRNG)
	for x := 0; x < rSize; x++ {
		s = append(s, Float())
	}
	u := Sequence(s)
	actual, _ := u(rng)
	if len(actual) != rSize {
		t.Errorf("Sequence should have produced an array of %v numbers but produced an array of size:%v \n", rSize, len(actual))
	}
}

func TestGenerateSequenceOfRandomBools(t *testing.T) {
	rSize := 1000
	rng := SimpleRNG{time.Now().Nanosecond()}
	var s []func(SimpleRNG) (bool, SimpleRNG)
	for x := 0; x < rSize; x++ {
		s = append(s, Boolean())
	}
	u := Sequence(s)
	actual, _ := u(rng)
	if len(actual) != rSize {
		t.Errorf("Sequence should have produced an array of %v bools but produced an array of size:%v \n", rSize, len(actual))
	}
}
