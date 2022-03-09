package linked_list

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/mikejlong60/genericpropcheck/pkg/array"
	"github.com/mikejlong60/genericpropcheck/pkg/propcheck"
	"testing"
	"time"
)

func TestPush(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 20, propcheck.String(40))

	prop := propcheck.ForAll(ge,
		"Validate Push for ConsList  \n",
		func(xs []string) []string {
			return xs
		},
		func(xss []string) (bool, error) {
			var errors error
			var l *ConsList[string]
			var i int
			for {
				if len(xss) == 0 {
					break
				}
				l = Push(xss[i], l)
				if l.head != xss[i] {
					errors = multierror.Append(errors, fmt.Errorf("Head %v  should have been %v pushed to head of ConsList", l.head, xss[i]))
				}
				if i+1 == len(xss) {
					break
				} else {
					i++
				}
			}
			if errors != nil {
				return false, errors
			} else {
				return true, nil
			}
		},
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]string](t, result)
}

func TestAddLast(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 20, propcheck.String(40))

	prop := propcheck.ForAll(ge,
		"Validate AddLast for ConsList  \n",
		func(xs []string) []string {
			return xs
		},
		func(xss []string) (bool, error) {
			var errors error
			var l *ConsList[string]
			var i int
			for {
				if len(xss) == 0 {
					break
				}
				l = AddLast(xss[i], l)
				if l.head == xss[i] && i > 0 {
					errors = multierror.Append(errors, fmt.Errorf("Head %v  should have been %v pushed to last Cons of ConsList, not the beginning", l.head, xss[i]))
				}
				if Len(l) != i+1 {
					errors = multierror.Append(errors, fmt.Errorf("Element %v did not get added to ConsList", l.head))
				}
				if i+1 == len(xss) {
					break
				} else {
					i++
				}
			}
			if errors != nil {
				return false, errors
			} else {
				return true, nil
			}
		},
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]string](t, result)
}

func TestAddWhile(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 20, propcheck.ChooseInt(1, 100))
	type resultType = []int

	prop := propcheck.ForAll(ge,
		"Validate AddWhile for ConsList  \n",
		func(xs []int) []int {
			return xs
		},
		func(xss resultType) (bool, error) { //TODO This input type has to match that of ExpectSuccess. This is a bug in your Prop mechanism you need to fix.
			p := func(x int) bool {
				if x > 50 { //only add ints > 50 to new list
					return true
				} else {
					return false
				}
			}

			xls := ToList(xss)
			xs := AddWhile(xls, p)

			var errors error
			var arrF []int
			for _, x := range xss {
				if x <= 50 {
					break
				} else {
					arrF = append(arrF, x)
				}
			}
			p2 := func(x, y int) bool {
				if x == y {
					return true
				} else {
					return false
				}
			}
			if !array.SetEquality(arrF, ToArray(xs), p2) {
				errors = multierror.Append(errors, fmt.Errorf("AddWhile did not stop adding when it hit limit. List from AddWhile:%v, expected: %v", ToArray(xs), arrF))
			}

			if errors != nil {
				//			fmt.Println(errors)
				return false, errors
			} else {
				return true, nil
			}
		},
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]int](t, result)
}

func TestTail(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 20, propcheck.String(40))

	prop := propcheck.ForAll(ge,
		"Validate Tail on ConsList  \n",
		func(xs []string) []string {
			return xs
		},
		func(xs []string) (bool, error) {
			var errors error
			var l *ConsList[string]
			_, err := Tail(l)
			if err == nil {
				errors = multierror.Append(errors, err)
			}
			if errors != nil {
				return false, errors
			} else {
				return true, nil
			}
		},
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]string](t, result)
}

func TestDropAndLen(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 20, propcheck.String(40))

	type fancyType struct {
		name   string
		number int
	}
	prop := propcheck.ForAll(ge,
		"Validate Drop and Len on ConsList  \n",
		func(xs []string) []string {
			return xs
		},
		func(xss []string) (bool, error) {
			var errors error
			var l *ConsList[fancyType]
			var i = 0
			for {
				if len(xss) == 0 {
					break
				}
				l = Push(fancyType{name: xss[i], number: i}, l)
				if i+1 == len(xss) {
					break
				} else {
					i++
				}
			}
			lenXss := Len(l)
			if lenXss != len(xss) {
				errors = multierror.Append(errors, fmt.Errorf("Len returned a length of %v but should have returned %v", lenXss, len(xss)))
			}

			l = Drop(l, len(xss)+1)
			if l != nil {
				errors = multierror.Append(errors, fmt.Errorf("Drop should have removed all elements but was %v ", l))
			}
			if errors != nil {
				return false, errors
			} else {
				return true, nil
			}
		},
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]string](t, result)
}

func TestTailUnwind(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 20, propcheck.String(40))

	prop := propcheck.ForAll(ge,
		"Validate unwinding Tail on ConsList  \n",
		func(xs []string) []string {
			return xs
		},
		func(xss []string) (bool, error) {
			var errors error
			var l *ConsList[string]
			var i int
			for {
				if len(xss) == 0 {
					break
				}
				l = Push(xss[i], l)
				if i+1 == len(xss) {
					break
				} else {
					i++
				}
			}
			for {
				if i == 0 {
					break
				} else {
					i--
				}
				l, _ = Tail(l)
				if xss[i] != Head(l) {
					errors = multierror.Append(errors, fmt.Errorf("Tail.Head %v  should have been %v ", Head(l), xss[i]))
				}
			}
			if errors != nil {
				return false, errors
			} else {
				return true, nil
			}
		},
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]string](t, result)
}

func TestHead(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 20, propcheck.String(40))

	prop := propcheck.ForAll(ge,
		"Validate Head on ConsList \n",
		func(xs []string) []string {
			return xs
		},
		func(xss []string) (bool, error) {
			var errors error
			var l *ConsList[string]
			var i int
			for {
				if len(xss) == 0 {
					break
				}
				l = Push(xss[i], l)
				if Head(l) != xss[i] {
					errors = multierror.Append(errors, fmt.Errorf("Head %v  should have been %v pushed to head of ConsList", Head(l), xss[i]))
				}
				if i+1 == len(xss) {
					break
				} else {
					i++
				}
			}
			if errors != nil {
				return false, errors
			} else {
				return true, nil
			}
		},
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]string](t, result)
}

func TestFoldRight(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 20, propcheck.String(40))

	prop := propcheck.ForAll(ge,
		"Validate FoldRight ConsList  \n",
		func(xs []string) []string {
			return xs
		},
		func(xss []string) (bool, error) {
			var errors error
			var l *ConsList[string]
			var i = 0
			for {
				if len(xss) == 0 {
					break
				}
				l = Push(xss[i], l)
				if i+1 == len(xss) {
					break
				} else {
					i++
				}
			}
			f := func(s string, accum []string) []string {
				return append(accum, s)
			}
			fConcat := FoldRight(l, []string{}, f)
			equality := func(l, r string) bool {
				if l == r {
					return true
				} else {
					return false
				}
			}

			if !array.SetEquality(xss, fConcat, equality) {
				errors = multierror.Append(errors, fmt.Errorf("FoldRight with toArray function returned %v but should have returned %v", fConcat, xss))
			}
			if errors != nil {
				return false, errors
			} else {
				return true, nil
			}
		},
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]string](t, result)
}

func TestFoldLeftAndFoldRight(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 200, propcheck.String(40))

	prop := propcheck.ForAll(ge,
		"Validate FoldLeft and FoldRight of ConsList  \n",
		func(xs []string) []string {
			return xs
		},
		func(xss []string) (bool, error) {
			var errors error
			var l *ConsList[string]
			var xssConcat string
			var i = 0
			for {
				if len(xss) == 0 {
					break
				}
				l = Push(xss[i], l)
				xssConcat = fmt.Sprintf("%v%v", xssConcat, xss[i])
				if i+1 == len(xss) {
					break
				} else {
					i++
				}
			}

			f := func(accum []string, s string) []string {
				return append([]string{s}, accum...)
			}
			flConcat := FoldLeft(l, []string{}, f)

			fr := func(s string, accum []string) []string {
				return append([]string{s}, accum...)
			}
			frConcat := FoldRight(l, []string{}, fr)
			equality := func(l, r string) bool {
				if l == r {
					return true
				} else {
					return false
				}
			}

			if !array.SetEquality(xss, flConcat, equality) {
				errors = multierror.Append(errors, fmt.Errorf("FoldLeft with toArray function returned %v but should have returned %v", flConcat, xss))
			}
			if !array.SetEquality(xss, frConcat, equality) {
				errors = multierror.Append(errors, fmt.Errorf("FoldRight with toArray function returned %v but should have returned %v", frConcat, xss))
			}

			if errors != nil {
				return false, errors
			} else {
				return true, nil
			}
		},
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]string](t, result)
}

func TestOrderingOfFoldLeftAndFoldRight(t *testing.T) {
	rng := propcheck.SimpleRNG{Seed: time.Now().Nanosecond()}
	ge := propcheck.ChooseList(0, 200, propcheck.String(40))

	prop := propcheck.ForAll(ge,
		"Validate FoldLeft and FoldRight of ConsList  \n",
		func(xs []string) []string {
			return xs
		},
		func(xss []string) (bool, error) {
			var errors error
			var l *ConsList[string]
			var xssConcat string
			var i = 0
			for {
				if len(xss) == 0 {
					break
				}
				l = Push(xss[i], l)
				xssConcat = fmt.Sprintf("%v%v", xssConcat, xss[i])
				if i+1 == len(xss) {
					break
				} else {
					i++
				}
			}

			//Verify way that FoldLeft and FoldRight order things in List
			h := func(accum string, s string) string {
				return fmt.Sprintf("%v%v", s, accum)
			}
			sConcatL := FoldLeft(l, "", h)
			if xssConcat != sConcatL {
				errors = multierror.Append(errors, fmt.Errorf("FoldLeft with string concat function returned [%v] but should have returned [%v]", sConcatL, xssConcat))
			}

			ii := func(s string, accum string) string {
				return fmt.Sprintf("%v%v", accum, s)
			}
			sConcatR := FoldRight(l, "", ii)
			if xssConcat != sConcatR {
				errors = multierror.Append(errors, fmt.Errorf("FoldRight with string concat function returned [%v] but should have returned [%v]", sConcatR, xssConcat))
			}

			if errors != nil {
				return false, errors
			} else {
				return true, nil
			}
		},
	)
	result := prop.Run(propcheck.RunParms{100, rng})
	propcheck.ExpectSuccess[[]string](t, result)
}
