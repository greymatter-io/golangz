package propcheck

import (
	"fmt"
	"strings"
	"time"
)

type SimpleRNG struct {
	Seed int
}

func (w SimpleRNG) String() string {
	return fmt.Sprintf("SimpleRMG{Seed: %v}", w.Seed)
}

func NextInt(r SimpleRNG) (int, SimpleRNG) {
	newSeed := (r.Seed*0x5DEECE66D + 0xB) & 0xFFFFFFFFFFFF
	nextRNG := SimpleRNG{newSeed}
	n := newSeed >> 16
	return n, nextRNG
}

// All the subsequent functions return A function which takes A SimpleRNG and returns an (A(or B or C), SimpleRNG) pair.

// Generate A random Int.
func Int() func(SimpleRNG) (int, SimpleRNG) {
	return func(r SimpleRNG) (int, SimpleRNG) {
		return NextInt(r)
	}
}

type WeightedGen[A any] struct {
	Gen    func(rng SimpleRNG) (A, SimpleRNG)
	Weight int
}

// Generates A random value from A set of generators in proportion to an individual generator's weight in the list.
func Weighted[A any](wgen []WeightedGen[A]) func(SimpleRNG) (A, SimpleRNG) {
	return func(rng SimpleRNG) (A, SimpleRNG) {
		var r []func(SimpleRNG) (A, SimpleRNG)
		for _, p := range wgen {
			for i := 0; i < p.Weight; i++ {
				r = append(r, p.Gen)
			}
		}
		a := ChooseInt(0, len(r))
		b, _ := a(rng)
		d := b
		g, rng2 := r[d](rng)
		return g, rng2
	}
}

// Generates A non-negative integer
var NonNegativeInt = func(rng SimpleRNG) (int, SimpleRNG) {
	i, r := NextInt(rng)
	if i < 0 {
		return -(i + 1), r
	} else {
		return i, r
	}
}

// Generates A float64 floating point number
var Float = func() func(SimpleRNG) (float64, SimpleRNG) {
	fa := func(a int) float64 {
		aa := a
		aaa := float64(aa)
		if aaa > 0 {
			return 1.0 / aaa
		} else {
			return float64(0)
		}
	}
	return Map(NonNegativeInt, fa)
}

var EmptyString = func() func(SimpleRNG) (string, SimpleRNG) {
	return Id("")
}

// See https://www.joelonsoftware.com/2003/10/08/the-absolute-minimum-every-software-developer-absolutely-positively-must-know-about-unicode-and-character-sets-no-excuses/
// You must understand these things about Unicode and character sets. The jist is that A Unicode codepoint('\u04E00' for example) can result in A string of 1 to 4 characters
// depending on the character set.

// Be careful about specifying the stringMaxSize because if you make it too large you will probably never end up with an empty string. A rule of
// thumb is to make the stringMaxSize 1/3 of the number of test cases you are running.
func String(unicodeMaxSize int) func(SimpleRNG) (string, SimpleRNG) {
	f := func(numOfCharactersInSet int, startingRune rune) []string {
		var unicodeStrings []string
		var currentUnicodeString string
		var currentRune = startingRune
		for x := 0; x <= numOfCharactersInSet; x++ {
			currentUnicodeString = fmt.Sprintf("%v", string(currentRune))
			unicodeStrings = append(unicodeStrings, currentUnicodeString)
			currentRune = currentRune + 0x01
		}
		return unicodeStrings
	}

	//latin := f(128, '\u0000')
	arabicNumbers := f(10, '\u0030')
	asciiLowercaseLetters := f(26, '\u0061')
	asciiUppercaseLetters := f(26, '\u0041')
	//hanchinesesubset := f(1000, '\u4E00')
	//cyrillic := f(646-456+60, '\u0400')
	//greek := f(1023-880, '\u0370')
	//bigUnicodeList := append(append(cyrillic, greek...), append(latin, hanchinesesubset...)...)
	bigUnicodeList := append(arabicNumbers, append(asciiLowercaseLetters, asciiUppercaseLetters...)...)

	start := 0
	stopInclusive := len(bigUnicodeList)

	return func(rng SimpleRNG) (string, SimpleRNG) {
		var i int                                            //The index into the big array of Unicode codepoints.
		var lr = rng                                         //The ever-changing random number generator inside the loop below.
		var res []string                                     //The growing list of unicode codepoints for making A single string at the end.
		var randomMaxSize int                                //The max size of this string measured by the number of unicode code points, not necessarily the size of the resulting string.
		randomMaxSize, lr = ChooseInt(0, unicodeMaxSize)(lr) //Randomly choose A value for the number of Unicode code points.
		for x := 0; x < randomMaxSize; x++ {
			_, lr = NextInt(lr)
			i, lr = ChooseInt(start, stopInclusive)(lr)
			res = append(res, bigUnicodeList[i])
		}
		return strings.Join(res, ""), lr
	}
}

// Generates a random date (stopExclusive - start) days from or preceding 1999-12-31.
func ChooseDate(start int, stopExclusive int) func(SimpleRNG) (time.Time, SimpleRNG) {
	g := func(days int, past bool) time.Time {
		ninetynine := "1999-12-31"
		current, _ := time.Parse("2006-01-02", ninetynine)
		if past {
			return current.AddDate(0, 0, -days).Round(time.Hour)
		} else {
			return current.AddDate(0, 0, days).Round(time.Hour)
		}
	}
	return Map2(ChooseInt(start, stopExclusive), Boolean(), g)
}

// Generates an integer between start and stop exclusive
func ChooseInt(start int, stopExclusive int) func(SimpleRNG) (int, SimpleRNG) {
	fa := func(a int) int {
		var divisor = stopExclusive - start
		if divisor <= 0 {
			divisor = 1
		}
		aa := a //.(int)
		r := start + aa%(divisor)
		return r
	}
	return Map(NonNegativeInt, fa)
}

// Generates A random boolean
func Boolean() func(SimpleRNG) (bool, SimpleRNG) {
	fa := func(a int) bool {
		aa := a
		return aa%2 == 0
	}
	return Map(NonNegativeInt, fa)
}

// Generates an array of N elements from the given generator
func ArrayOfN[T any](n int, g func(SimpleRNG) (T, SimpleRNG)) func(SimpleRNG) ([]T, SimpleRNG) {
	var s []func(SimpleRNG) (T, SimpleRNG)
	for x := 0; x < n; x++ {
		s = append(s, g)
	}
	u := Sequence(s)
	return u
}

// Generates an array with A size in the indicated range using the given Gen
func ChooseArray[T any](start, stopInclusive int, kind func(SimpleRNG) (T, SimpleRNG)) func(SimpleRNG) ([]T, SimpleRNG) {
	return func(rng SimpleRNG) ([]T, SimpleRNG) {
		if start < 0 || start > stopInclusive {
			panic(fmt.Sprintf("Low range[%v] was < 0 or exceeded the high range[%v]", start, stopInclusive))
		}
		i, _ := ChooseInt(start, stopInclusive)(rng)
		r, rng2 := ArrayOfN(i, kind)(rng)
		return r, rng2
	}
}
