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
func TestMap32WithALotOfFunctionComposition(t *testing.T) {
	rng := SimpleRNG{time.Now().Nanosecond()}
	ra := String(10)
	rb := String(20)
	rc := String(1)
	rd := String(2)
	re := String(15)
	rf := String(100)
	l1 := []WeightedGen[string]{
		{Gen: ra, Weight: 300}, {Gen: rb, Weight: 10}, {Gen: rc, Weight: 300}, {Gen: rd, Weight: 10}, {Gen: re, Weight: 300}, {Gen: rf, Weight: 10},
		{Gen: ra, Weight: 300}, {Gen: rb, Weight: 10}, {Gen: rc, Weight: 300}, {Gen: rd, Weight: 10}, {Gen: re, Weight: 300}, {Gen: rf, Weight: 10},
		{Gen: ra, Weight: 300}, {Gen: rb, Weight: 10}, {Gen: rc, Weight: 300}, {Gen: rd, Weight: 10}, {Gen: re, Weight: 300}, {Gen: rf, Weight: 10},
		{Gen: ra, Weight: 300}, {Gen: rb, Weight: 10}, {Gen: rc, Weight: 300}, {Gen: rd, Weight: 10}, {Gen: re, Weight: 300}, {Gen: rf, Weight: 10},
		{Gen: ra, Weight: 300}, {Gen: rb, Weight: 10}, {Gen: rc, Weight: 300}, {Gen: rd, Weight: 10}, {Gen: re, Weight: 300}, {Gen: rf, Weight: 10},
		{Gen: ra, Weight: 300}, {Gen: rb, Weight: 10}, {Gen: rc, Weight: 300}, {Gen: rd, Weight: 10}, {Gen: re, Weight: 300}, {Gen: rf, Weight: 10},
	}
	f := func(a, b, c, d, f, g, h, i, j, k, l, m, n, o, p, q, r, s, t, u, v, w, x, y, z, aa, bb, cc, dd, ee, ff, gg string) string {
		return gg
	}
	res := Map32(Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), Weighted(l1), f)

	actual := ForAll(res, "Map32 with lots of function composition",
		func(x string) string { return x },
		func(x string) (bool, error) {
			return true, nil
		},
	)
	result := actual.Run(RunParms{100, rng})
	ExpectSuccess[string](t, result)
}

func TestMap32AndBySurrogateMap16AndMap8(t *testing.T) {
	rng := SimpleRNG{time.Now().Nanosecond()}
	type goose struct {
		a  int
		b  int64
		c  string
		d  time.Time
		e  bool
		f  int
		g  int64
		h  string
		i  time.Time
		j  int
		k  int64
		l  string
		m  time.Time
		n  bool
		o  string
		p  int
		q  int64
		r  string
		s  time.Time
		t  int
		u  int64
		v  string
		w  time.Time
		x  bool
		y  int
		z  int64
		aa string
		bb time.Time
		cc string
		dd time.Time
		ee string
		ff time.Time
	}

	d := time.Now()
	i := d.Add(1)
	m := i.Add(1)
	s := m.Add(1)
	w := s.Add(1)
	bb := w.Add(1)
	dd := bb.Add(1)
	ff := dd.Add(1)

	expected := goose{
		12, int64(13), "14", d, true, 12, int64(13), "14", i, 12, int64(13), "14", m, true, "ro", 12, int64(13), "14", s, 12,
		int64(13), "14", w, false, 12, int64(13), "14", bb, "14", dd, "14", ff,
	}

	ra := Id(12)
	rb := Id(int64(13))
	rc := Id("14")
	rd := Id(d)
	re := Id(true)
	rf := Id(12)
	rg := Id(int64(13))
	rh := Id("14")
	ri := Id(i)
	rj := Id(12)
	rk := Id(int64(13))
	rl := Id("14")
	rm := Id(m)
	rn := Id(true)
	ro := Id("ro")
	rp := Id(12)
	rq := Id(int64(13))
	rr := Id("14")
	rs := Id(s)
	rt := Id(12)
	ru := Id(int64(13))
	rv := Id("14")
	rw := Id(w)
	rx := Id(false)
	ry := Id(12)
	rz := Id(int64(13))
	raa := Id("14")
	rbb := Id(bb)
	rcc := Id("14")
	rdd := Id(dd)
	ree := Id("14")
	rff := Id(ff)

	f := func(a int, b int64, c string, d time.Time, e bool, f int, g int64, h string, i time.Time, j int, k int64, l string, m time.Time, n bool, o string, p int, q int64, r string, s time.Time, t int, u int64, v string, w time.Time, x bool, y int, z int64, aa string, bb time.Time, cc string, dd time.Time, ee string, ff time.Time) goose {
		return goose{a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q, r, s, t, u, v, w, x, y, z, aa, bb, cc, dd, ee, ff}
	}
	res := Map32(ra, rb, rc, rd, re, rf, rg, rh, ri, rj, rk, rl, rm, rn, ro, rp, rq, rr, rs, rt, ru, rv, rw, rx, ry, rz, raa, rbb, rcc, rdd, ree, rff, f)
	actual, _ := res(rng)
	if actual != expected {
		t.Errorf("Map32 did not map correctly \nactual:  %v\nexpected:%v\n", actual, expected)
	}
}
