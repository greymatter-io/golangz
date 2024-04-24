package stack

// An Immutable stack, mostly from here: https://www.reddit.com/r/golang/comments/13hi00u/should_constructor_functions_always_return_a/
// the pointer type.  All functions are pure.
// The operators on Stack are not methods because I cannot figure out a way to work around the fact that the type definition needs
// more type parameters than the type allowed. FoldRight and FoldLeft need an additional type parameter. And who knows what
// future functions might need more.
type node[A any] struct {
	val  A
	next *node[A]
}

// the value wrapper
type Stack[A any] struct {
	root *node[A]
}

func NewStack[A any]() Stack[A] {
	return Stack[A]{}
}

func Push[A any](s Stack[A], val A) Stack[A] {
	newRoot := node[A]{val: val, next: s.root}
	return Stack[A]{root: &newRoot}
}

func Pop[A any](s Stack[A]) Stack[A] {
	if s.root == nil {
		return s
	}
	return Stack[A]{root: s.root.next}
}

func IsEmpty[A any](s Stack[A]) bool {
	if s.root == nil {
		return true
	} else {
		return false
	}
}

func Peek[A any](s Stack[A]) (val A, ok bool) {
	if s.root == nil {
		return val, false
	}
	return s.root.val, true
}

// FoldRight
/*	FoldRight[B](l Stack[A], z: B)(op: (A, B) => B): B
	  Applies a binary operator to all elements of this list and a start value, going right to left. Left is the top of the stack, and right is the bottom.
	  B -the result type of the binary operator.
	  z - the start or zero value.
	  op  - the binary operator.
	  returns - the result of inserting op between consecutive elements of the stack, going right to left with the start
		value z on the right: op(x1, op(x2, ... op(xn, z)...)) where x1, ..., xn are the elements of this list. Returns z if this list is empty.
*/
func FoldRight[A, B any](l Stack[A], z B, f func(A, B) B) B {
	if IsEmpty(l) {
		return z
	} else {
		return f(l.root.val, FoldRight(Stack[A]{root: l.root.next}, z, f))
	}
}

// FoldLeft
/* FoldLeft[B](l Stack[A], z: B)(op: (B, A) => B): B
* Applies a binary operator to all elements of this list and a start value, going left to right. Left is the top of the stack, and right is the bottom.
*   FoldLeft can be made stack-safe,  tail recursive if the compiler supports that optimization.  Note how the f function is executed
*    every time the FoldLeft is invoked without waiting for the recursive stack of function calls to reach bottom.
*    Golang's compiler does not optimize tail calls.  But you can make it do that using a thunk as shown in Arrays.FOldLeft in the arrays package of Golangz.
*	  B -the result type of the binary operator.
*	  z - the start or zero value.
*	  op  - the binary operator.
*	  returns - the result of inserting op between consecutive elements of this list, going left to right with the start
*		value z on the left: op(...op(z, x1), x2, ..., xn) where x1, ..., xn are the elements of this list. Returns z if this list is empty.
 */
func FoldLeft[A, B any](l Stack[A], z B, f func(B, A) B) B {
	if IsEmpty(l) {
		return z
	} else {
		return FoldLeft(Stack[A]{root: l.root.next}, f(z, l.root.val), f)
	}
}
