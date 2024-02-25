_Copyright &copy; Kirk Rader 2024_

# Functional Programming in Go

Go's syntax and much of its semantics deliberately hearken back to the days when
languages like Pascal and C represented the main stream of programming language
design. Go has no object-oriented features. Its variables are fairly strongly
typed (though less strongly than its marketing would have you believe). Even its
exception-handling mechanisms are carefully designed to resemble, at least
cosmetically, operating-system level library functions rather than distinct
syntax like the `try ... finallly ...` special forms of other languages.

Where Go's semantics diverge from the mainstream languages of the 1970's and
1980's, it borrows just as liberally (though less visibly) from that other
ancient programming language lineage: Lisp. It has a garbage collector. But more
importantly, it has lexical closures. Closures allow a programmer to use the
same techniques that gave rise to object orientation in the first place. (It has
long been a maxim of Lisp programmers that "objects are a poor-man's closures.")

But first, some general background on Go features related to functions.

## Functions

You can only declare a named function in Go at global scope.

```go
package main

func SomeFunction() {

}

func main() {
	SomeFunction()
}
```

You can declare anonymous functions in any scope.

```go
package main

func main() {

	someFunction := func() {

	}

	someFunction()
}
```

The identity of a function is determined by its signature:

1. Number of parameters
2. Type of each parameter
3. Number of return values
4. Type of each return value

It is a syntax error to attempt to define two functions with the same name but
different signatures; i.e. Go does not allow function overloading.

Here is an example with a fully elaborated function definition with multiple
parameters, including a variadic parameter, and multiple return values:

```go
package main

import (
	"errors"
	"fmt"
	"os"
)

func SomeFunction(n int, s string, a ...any) (any, error) {

	return nil, errors.New("not yet implemnented")
}

func main() {

	v, err := SomeFunction(42, "Hello, world!", 4.2, "another string", nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println(v)
}
```

## Methods

To reiterate: Go has no object-oriented features. What it calls methods are just
a slightly different syntax for functions. There is nothing that you can or need
to do with methods that could not be done just as well with functions, other
than use them in coordination with interfaces. To be clear, methods in Go are
not dependent on interfaces, while interfaces are completely and utterly defined
by the methods they declare.

Here is an example of interface-less methods (plus one method that implements
the `fmt.String` interface) on an application-defined `struct` type:

```go
package main

import "fmt"

type MyType struct {
	member1 int
	member2 string
}

func (m *MyType) Member1() int {
	return m.member1
}

func (m *MyType) Member2() string {
	return m.member2
}

func (m *MyType) SetMember1(n int) {
	m.member1 = n
}

func (m *MyType) SetMember2(s string) {
	m.member2 = s
}

func (m MyType) String() string {
	return fmt.Sprintf("<MyType %d, \"%s\">", m.member1, m.member2)
}

// Prints <MyType 42, "Hello, world!"> to stdout
func main() {

	m := MyType{}
	m.SetMember1(42)
	m.SetMember2("Hello, world!")
	fmt.Println(m)
}
```

The only things that make `Member1()`, `SetMember1(int)` etc. methods different
from ordinary functions are that they enable the use of "." notation when
invoking them together with a little bit of (potentially confusing) syntactic
sugar related to pointers. In other words, the following program is 100%
semantically equivalent to the preceding one:

```go
package main

import "fmt"

type MyType struct {
	member1 int
	member2 string
}

func Member1(m *MyType) int {
	return m.member1
}

func Member2(m *MyType) string {
	return m.member2
}

func SetMember1(m *MyType, n int) {
	m.member1 = n
}

func SetMember2(m *MyType, s string) {
	m.member2 = s
}

func (m MyType) String() string {
	return fmt.Sprintf("<MyType %d, \"%s\">", m.member1, m.member2)
}

// Prints <MyType 42, "Hello, world!"> to stdout
func main() {

	m := MyType{}
	SetMember1(&m, 42)
	SetMember2(&m, "Hello, world!")
	fmt.Println(m)
}
```

## Function Type Aliases

Before embarking on a fuller discussion of the functional programming paradigm,
Go style, it is worth emphasizing that anonymous functions in Go are not only
lexical closures -- more later on exactly what that means -- but they are first
class data values, and so must have a type. As already noted, the type of a
function is determined by its signature. For example, the following is a syntax
error due to the fact that inferred types for `f1` and `f2` are different:

```go
f1 := func() {}
f2 := func(int) {}
// error! incompatible types
f1 = f2
```

You can give such function types names using aliases just as for any other type.

```go
package main

import "fmt"

type (
	Invariant func() int
	Transform func(int) int
)

func Composite(invariant Invariant, transform Transform) int {
	return transform(invariant())
}

// Prints 2 to stderr.
func main() {
	invariant1 := func() int { return 1 }
	increment := func(n int) int { return n + 1 }
	fmt.Println(Composite(invariant1, increment))
}
```

Such type aliases save typing when defining the signatures of functions like
`Composite()` that take closures as parameters. But they can also become
essential to crafting maintainable function-oriented programs as discussed
below.

## Lexical Closures

The term "lexical closure" (or the shorthand "closure") has been used a few
times so far. For those unfamiliar with that term, a closure is both:

- A data value that can be passed as a parameter to other functions and returned
  as a value from them.

- An executable procedure that can be invoked as a subroutine at any point
  during the execution of a program.

In addition, each closure contains its own private copy of the complete lexical
environment in which it was created. What that means is that closures can
provide both the data encapsulation and functional modularity that are two of
the core principles of object-oriented design and implementation. What they do
not provide on their own is polymorphism, but see below on how to use functional
composition as means of implementing similar software design features as is the
goal of inheritance in object-oriented programming languages.

Understanding the nature and power of closures is difficult in the abstract.
Consider the following Go code:

```go
package main

import "fmt"

type (
	Getter func() int
	Setter func(n int)
)

func MakeInstance(value int) (Getter, Setter) {

	getter := func() int { return value }
	setter := func(n int) { value = n }
	return getter, setter
}

// Prints the following to stdout:
//
//	 0  42
//	-1  43
func main() {

	getter1, setter1 := MakeInstance(0)
	getter2, setter2 := MakeInstance(42)
	fmt.Printf("%3d %3d\n", getter1(), getter2())
	setter1(-1)
	setter2(43)
	fmt.Printf("%3d %3d\n", getter1(), getter2())
}
```

The preceding works because `MakeInstance()` returns two newly created lexical
closures each time it is called. Those closures are created in a common lexical
environment with its own binding for `value`, the parameter to `MakeInstance()`.
Thus when `getter1` and `getter2` are called they each return the `value` passed
to `MakeInstance()`, even though that means there are two bindings for the
"same" variable. One binding is visible to `getter1` and the other to `getter2`.

Since `setter1` was created in the same lexical environment as `getter1`, it
sees the same binding for `value`. Ditto for `setter2` and `getter2`. I.e. the
two closures created by a single invocation of `MakeInstance()` share always
share a single binding for `value` but that binding is always different from the
binding shared by the closures created by any other invocation.

Here is the definition of `stacktraces.skipUntil()` in
[../stacktraces/stacktrace.go](../stacktraces/stacktrace.go) for a less contived
example.

```go
// Return a function that returns true when invoked for a frame with the given
// function name and all that follow it.
//
// The returned function will be the frame test used by formatStackTrace() when
// skipFrames is a string
func skipUntil(startWhen string) stackFrameTest {

	seen := false

	return func(frame *runtime.Frame) bool {
		seen = seen || startWhen == frame.Function
		return seen
	}
}
```

Using a closure in this way allows encapsulation of the `startWhen` parameter
and local `seen` variable while going with the flow of Go's overall procedural
orientation. It used to create a helper function as needed inside the
implementation of `formatStackTrace()`.

```go
switch v := skipFrames.(type) {

case int:
	if v < 0 {
		// skip past this function's caller's caller when skipFrames is
		// negative
		skip = defaultSkip
	} else {
		// skip the specified number of frames when skipFrames is
		// non-negative
		skip = v
	}

case string:
	// use a frameTest that skips all frames until a given function is seen
	frameTest = skipUntil(v)

default:
	// skip past this function's caller's caller when skipFrames is any
	// other value
	skip = defaultSkip
}
```

## Functional Programming

Just as object-oriented programming is a paradigm that represents the
fundamental building blocks of a software design as _objects_, functional
programming uses _functions_ as its fundamental organizing principle. This is
different from a procedural paradigm that simply emphasizes organizing code into
subroutines. Functions in the functional paradigm provide not only modularity
but data encapsulation and represent the fundamental model for flow of control,
to within the limits of a given language's support for it.

> As a brief historical discursion, the term Continuation Passing Style (CPS),
> which is used internally in the implementation of many compilers for many
> programming languages, was first coined to describe a fundamental aspect of
> the design of the Scheme programming language. Scheme was built from the
> ground up with the explicit goal of enabling functional programming techniques
> in the same way that C++ and Java were designed specifically to support
> object-oriented programming. In order to truly conform to a "pure" functional
> approach, a language needs to support CPS at the application programming
> level. Among real-world languages, Scheme is the only such language of which I
> am aware. (The Clojure dialect of Lisp built on the Java Virtual Machine with
> its `trampoline` construct comes close but is too limited in various ways by
> the underlying semantics of the JVM to reach the full power of Scheme.)

While Go fails to support all the semantics necessary for completely
implementing functional programming techniques at every level of software design
and implementation, its support for lexical closures provides some essential
building blocks necessary to benefit from the functional paradigm.

| Object-Oriented Concept | Functional Equivalent                                      |
|-------------------------|------------------------------------------------------------|
| Instances of  classes   | Closures                                                   |
| Methods of classes      | Closures                                                   |
| Inheritance             | Functional composition (i.e. closures calling one another) |

For purposes of this tutorial, a "functional approach" boils down to two general
recommendations:

1. When thinking about flow of control, consider passing and calling closures as
   an alternative to deeply nested chains of named subroutines.

2. When thinking about data encapsulation and modularity, consider capturing
   simple values from lexical environments in closures rather than creating
   `struct`, `map` and similar complex data types while defining methods on
   those structured types.

When implementing idioms like dependency injection and similar inversions of
control, the first principal suggests implementing them by passing closures that
conform to particular signatures rather than declaring them using `interface`
types. This has significant advantages for Go given the limitations [discussed
elsewhere](./types.md) of Go's whole approach to interfaces. A simple example of
this has already been shown in an earlier example taken from
[../stacktraces/stacktrace.go](../stacktraces/stacktrace.go). To expand that
example, the `stacktraces` library code defines the following alias for a
particular function signature:

```go
// Type of function used to test a frame for inclusion in a stack trace.
type stackFrameTest func(frame *runtime.Frame) bool
```

The `formatStackTrace()` function then creates a closure based on one of two
functions that implement that signature depending on the value passed in its
`skipFrames` parameter. I.e. either

```go
// default frame test unconditionally returns true on the assumption that
// skip will be a non-negative int
frameTest := func(*runtime.Frame) bool { return true }
```

or the result of calling `skipUntil()` when its `skipFrames` parameter turns out
to be a `string`:

```go
case string:
	// use a frameTest that skips all frames until a given function is seen
	frameTest = skipUntil(v)
```

As shown previously, `skipFrames()` is a named function that returns a
dynamically created closure which, in turn, accesses data needed to drive its
behavior in its closed-over environment.

After execution of `formatStackTrace()` has set `frameTest` to its ultimate
value, it then simply uses whichever of those two versions of `stackFrameTest`
it previously created while looping through the call stack to determine which
frames to include in the stack trace:

```go
for {

	frame, more := frames.Next()

	if frameTest(&frame) {

		if longFormatter != nil {
			longFormatter(longWriter, frameNumber, &frame)
		}

		if shortFormatter != nil {
			shortFormatter(shortWriter, frameNumber, &frame)
		}
	}

	frameNumber += 1

	if !more {
		break
	}
}
```

i.e. a simple example of dependency injection via lexical closure, using the
principles of functional programming.

Note that the loop shown above includes additional examples of inversion of
control via lexical closures in the form of the `longFormatter` and
`shortFormatter` parameters to `formatStackTrace()`. In that case, it is being
used to support a performance and memory optimization. The callers of
`formatStackTrace()` can choose to pass `nil` for either of the
`stackFrameFormatter` parameters when they only need one of the two supported
formats of stack trace or valid `strackFrameFormatter` functions for both when
both formats are requested.

Use of type aliases for particular function signatures improves both the
readability of the code and improves the compiler diagnostics when attempts are
made to pass a wrong type of closure.
