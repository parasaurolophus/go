_Copyright &copy; Kirk Rader 2024_

# Go Data Types

Like all programming languages, Go supports a variety of data types natively and
provides mechanisms for programmers to extend those with application-specific
types. Those mechanisms deliberately exclude any that are explicitly designed to
support object-oriented programming in the style of languages like C++ or Java.
In particular, Go's semantics include no provision for polymorphism, other than
of a few very specific and limited kinds (as discussed below). This has a huge
impact on software design and implementation, especially when combined with the
fact that variables are strictly typed in Go and it provides very few implicit
C-style type conversions nor function overloading.

> As a bit of foreshadowing, the [functional programming
> paradigm](./functional.md) is generally the way to Go (pun intended).

## Scalar Types

Go defines a bewildering variety of numeric data types. See
<https://go.dev/ref/spec#Numeric_types> for the complete list. While that link
reveals a fairly predictible, if unusually large, variety of numeric types, as
with evrything related to Go, there are some idiosyncracies of which to be
aware.

### Integers

Like C of a bygone era, Go defines `int` and `uint` types whose precision varies
depending on the architecture of the underlying hardware on which a Go program
was compiled. I.e. `int` and `uint` will consume 32 bits on 32-bit hardware, 64
bits otherwise. For maximum efficiency of integer calculations where it makes no
difference whether an integer value is 32 or 64 bits, use `int` or `uint`.

In contexts where it makes a difference exactly how many bits there are in a
given integer value, you can use `int8`, `int16`, `int32` or `int64` types
explicitly and their unsigned equivalents, `unint8`, `uint16`, `uint32` and
`uint64`.

> Note that signed integers are represented as twos-complement binary numbers as
> is the case for most programming languages. While that is rarely of much
> interest to anyone not writing a Go compiler, it does have some implications
> for edge-case computations where the numeric magnitude reaches the limit of
> precision of the given integer type. See any college-level programming text
> book for details (Knuth's venerable series on _The Art of Computer
> Programming_ remains a personal favorite, which may reveal something about my
> age).

### Non-Numeric "Integer" Types

Go defines the `byte` type as a synonym for `uint8`. Its intended use is in the
representation of arbitrary binary data rather than numbers suitable for use in
mathematical calculations.

Similarly, the `rune` type is a synonym (oddly enough) for `int32` (rather than,
as one might have assumed, `uint32`). It is used to represent Unicode "code
points," i.e. more or less equivalent to the `char` type in many other languages.

Finally, there is the `uintptr` type which, like `int` and `uint`, varies in
size according to the target hardware architecture. It is used to represent
[pointers](#pointers) (about which much is said later).

### Floats

The `float32` and `float64` types represent
[IEEE-754](https://standards.ieee.org/ieee/754/6210/) single- and
double-precision floating point numbers, respectively.

### Imaginary Numbers

The `complex64` and `complex128` tyes represent 64- and 128-bit imaginary
numbers, respectively.

> As a reminder from your college math classes, an imaginary number is
> represented as a sum of some ordinary number and some multiple of _i_, the
> square root of -1. Such values are "complex" because of the summing of two
> numeric values intrinsic to their definition. The mathematical entities they
> represent are called "imaginary" because, of course, there is in reality no
> such thing as the square root of any negative number. But, when you get right
> to the heart of things, there really is no such thing as a negative number to
> start with, so what's your point?

```go
// the inferred type of c is complex128
c := (4.2 + -1.1i)

// prints (4.2-1.1i) to stdout
fmt.Println(c)
```

### Strings

Unicode made everything to do with representing text vastly more complicated
than the good old days when if it couldn't be represented in ASCII, it didn't
need representing. Go deals with these complexities by representing a `string`
value as a sequence of `byte` values (not `rune` values!) which represent a
sequence of Unicode code points. Go does not guarantee such sequences to be
normalized, so there is no predictable relationship between the number of bytes,
let alone runes, and the length of the string when considered as a sequence of
"characters."

As a result of all this, most of Go's standard library functions that manipulate
text operate on sequences (really `byte` [slices](#arrays-and-slices)) and you
turn such a sequence into a `string` only at the end, when you are ready to
display, transmit, store or otherwise use the byte sequence as an actual string
of human-readable text. For example, the standard library's `json` package
operates on byte slices which can then be turned into actual strings only when
needed:

```go
value := SomeStruct{
    SomeMember:      1,
    SomeOtherMember: "two",
}

// the inferred type for b is []byte, i.e. a slice of bytes
b, _ := json.Marshal(value)

// the inferred type of s is string
//
// it contains the JSON representation of the value passed
// to json.Marshal()
//
// note that the only prediciable relationship between len(b)
// and len(s) is that len(b) >= len(s), i.e. it can take
// different numbers of bytes to represent a single Unicode
// "character" as determined by a given encoding
//
// and, yes, that is true of UTF-8, as well; remember that
// the whole point of UTF-8 is to allow Unicode strings that
// only contain ASCII code points be interchangeable with
// ASCII strings while still supporting Unicode code points
// that are outside of ASCII's range
//
// the fact that Go's standard library makes no attempt
// to normalize Unicode strings makes the relationship between
// "number of characters" and "size in memory" that much more
// tenuous
s := string(b)

// this prints {"SomeMember":1,"SomeOtherMember":"two"} to stdout
fmt.Println(s)
```

## Arrays and Slices

Like most programming languages, Go possesses a built-in `array` type that
represents ordered collections of elements of some specific type. Unlike most
languages, programmers rarely operate on arrays directly in Go. Instead, most
built-in syntax and standard library operations use the `slice` type
corresponding to a given `array` type.

"They use a what, now?" I hear you ask.

A Go `slice` is a mutable view of a contiguous segment of some underlying
`array`. Anywhere in Go that you see syntax that looks as if it were operating
on an array is probably actually operating on a `slice`, instead. Even the
built-in `make()` function, used to create arrays, actually returns a `slice` to
the otherwise-hidden array. The following code (taken from
[../stacktraces/function_info.go](../stacktraces/function_info.go)) illustrates
a typical pattern for implicitly creating an array using `make()` and
manipulating the state of the slice as needed.

```go
pc := make([]uintptr, maxDepth)
n := runtime.Callers(skip, pc)

if n < 1 {
    return 0, "", "", 0, false
}

pc = pc[:n]
frames := runtime.CallersFrames(pc)
```

In the preceding example, the inferred type of `pc` is `[]intptr`, i.e. the same
as the type parameter to `make()`. The second parameter is the size of the
underlying array to be created. This dynamically allocates an array of `uintptr`
elements and creates a slice that is the same as if you had used something like:

```go
a := [maxDepth]uintptr{}
pc := a[:]
```

To extend the preceding example:

```go
// a's inferred type is "array of 3 ints"
a := [3]int{1, 2, 3}

// prints 3 to stdout
fmt.Println(len(a))

// prints 2, i.e. array indexes are 0-based
fmt.Println(a[1])

// s's inferred type is "slice of 2 ints"
s := a[1:]

// prints 2 since s starts at the second element in a
fmt.Println(len(s))

// prints 2 since the first element of s is the second element of a
fmt.Println(s[0])

// drop the first element from s
s = s[1:]

// prints 1 since we removed one element from the beginning of the slice
fmt.Println(len(s))

// prints 3 since s now starts at the third element of a
fmt.Println(s[0])
```

Note that the preceding example is shown for completeness. The preferred way to
create and manipulate arrays is through functions that implicitly or explicitly
call `new()` or `make()`, returning slices that hide the underlying arrays
entirely. This reduces the risk of multiple slices referring to overlapping
regions in the same array which can create havoc in your code, especially if
[goroutines](./goroutines.md) are thrown into the mix.

## Structs

Go provides C-style `struct` types. A struct combines a fixed number of values,
each of a specific type into a single data object. Go's `struct` syntax is
suggestive of JavaScript objects but their semantics are actually pure C in that
every object of a given `struct` type has exactly the same number of members,
each of which are of a fixed type. Nor does each instance of a `struct` type
need to allocate store for the "keys," which actually just become
programmer-defined identifiers, even though they are used with "." notation
designed to suggest key / value pairs. The compiler does keep track of such
identifiers such that they can be used at runtime by the `json` standard library
package.

```go
type MyStruct struct {
    Member1 int
    Member2 float32
    Member3 string
}

s := MyStruct{
    Member1: 40,
    Member2: 2.0,
    Member3: "Hello, world!",
}

// prints Hellow, world!
fmt.Println(s.Member3)

b, _ := json.Marshal(s)

// prints {"Member1":40,"Member2":2,"Member3":"Hello, world!"}
fmt.Println(string(b))
```

Because Go's syntax imposes constraints on such "keys," it also makes provision
for annotating `struct` declarations with "tags" that act as hints to the `json`
library functions. Note that in the preceding example, the `struct` type and its
member names are all declared using capitalized identifiers. This is Go's syntax
for making an identifier visible outside of the package in which it is defined.
By default, the JSON marshaler will use the member names as keys, resulting in
capitalized keys in the JSON string. It is a common convention to not capitalize
key names in JSON. Here is the same example, enhanced with JSON tags to produce
more idiomatic JSON output:

```go
type MyStruct struct {
    Member1 int     `json:"member1"`
    Member2 float32 `json:"member2"`
    Member3 string  `json:"member3"`
}

s := MyStruct{
    Member1: 40,
    Member2: 2.0,
    Member3: "Hello, world!",
}

b, _ := json.Marshal(s)

// this version prints {"member1":40,"member2":2,"member3":"Hello, world!"}
fmt.Println(string(b))
```

## Maps

In addition to structs, Go provides a general purpose type called `map`. Maps
are true dictionaries in that each map object can associate an arbitrary number
of keys with a corresponding number of values. They are still bound by Go's
overall highly constrained type system such that, unlike hash tables in many
other languages, all of the values in a given map must be of the type specified
when the map type was declared.

```go
m := map[string]int{
    "one":   1,
    "two":   2,
    "three": 3,
}

// prints 1
fmt.Println(m["one"])

b, _ := json.Marshal(m)

// prints {"one":1,"three":3,"two":2}
fmt.Println(string(b))
```

Maps can be specified to use any type as keys, not just strings.

```go
type Key int

const (
    One Key = iota
    Two
    Three
)

m := map[Key]string{
    One:   "one",
    Two:   "two",
    Three: "three",
}

// prints two
fmt.Println(m[Two])

b, _ := json.Marshal(m)

// prints {"0":"one","1":"two","2":"three"}
fmt.Println(string(b))
```

See [Type Aliases](#type-aliases), below, for more on what is happening in the
latest example to cause the possibly somewhat unexpected key strings to be
emitted in the JSON output.

## Pointers

Go passes (nearly) all values to and from functions by value, not reference. Its
support for passing by reference is modeled on that of C (not C++). For every
type, `T`, there is a corresponding type `*T` which denotes "pointer to T."
There are corresponding "address of" and "referenced" operators, all inspired by
C's traditional pointer syntax.

```go
var p *int
n := 0
p = &n
*p = 42

// prints 42
fmt.Println(n)
```

Such pointer types are of most benefit when used as parameters to [methods and
functions](./functional.md). For example, here is how to implement a method with
a side-effect on its receiver:

```go
package main

import "fmt"

// Define an alias for the built-in int type.
type MyInt int

// Define a method for MyInt with a side-effect thanks
// to operating on a pointer receiver.
func (p *MyInt) Increment() MyInt {
	*p += 1
	return *p
}

// Prints the following to stdout.
//
// 0
// 1
// 1
func main() {

	n := MyInt(0)
	fmt.Println(n)
	fmt.Println(n.Increment())
	fmt.Println(n)
}
```

Go's pointers are far more limited in their semantics compared to C. In
particular, Go does not support pointer arithmetic, which is one of the defining
characteristics of C and one of the main reasons it remains popular for
system-level programming to this day. In part this is a design choice, but it is
inevitable due to Go being burdened by a garbabe collector. When coupled with
Go's rather inconsistent rules regarding implicit conversions, Go's pointer type
definitely takes getting used to, whether or not you are familiar with how
pointers work in C.

## Type Aliases and Constraints

A number of the examples on this page have already used type aliases, e.g.

```go
type MyInt int
```

Such aliases are used for a variety of purposes. Very frequently they are used
merely as a work-around for the rule that [methods](./functional.md) can only be
declared for types defined in the same package.

```go
package main

// Syntax error, because the built-in type int is defined in a different package
// than the current one.
//
// func (p *int) Increment() int {
// 	*p += 1
// 	return *p
// }

// Define an alias for int in the current package.
type MyInt int

// Define a method on MyInt.
func (p *MyInt) Increment() MyInt {
	*p += 1
	return *p
}
```

Aliases are also used to provide more succinct syntax and to express aspects of
software design. For example, [../logging/logger.go](../logging/logger.go),
[../stacktraces/stacktrace.go](../stacktraces/stacktrace.go) and
[../stacktraces/function_info.go](../stacktraces/function_info.go) each defines
a number of aliases for functions with specific signatures, e.g.

```go
// Type of function passed to logging methods for lazy evaluation of message
// formatting.
//
// The returned string becomes the value of the log entry's msg attribute.
//
// Such a function is invoked only if a given verbosity is enabled for a
// given logger.
MessageBuilder func() string

// Type of function passed as first argument to Logger.Defer() and
// Logger.DeferContext().
Finally func()

// Type of function passed to Logger.Defer() and Logger.DeferContext() to
// allow for including the value returned by recover() in the log entry.
RecoverHandler func(recovered any) string
```

These named function signatures not only reduce typing where they are used, they
help elucidate their intended usage.

In addition to defining type aliases, Go's syntax overloads the `type` keyword
for another purpose. While Go's version of generic functions stop far short of
providing true parameterized types, they do support creating
application-specific named combinations of types to make generic functions more
useful.

```go
package main

import (
	"fmt"
	"strconv"
)

type (
	Integer interface {
		int | int8 | int16 | int32 | int64
	}

	Unsigned interface {
		uint | uint8 | uint16 | uint32 | uint64
	}

	Float interface {
		float32 | float64
	}

	Complex interface {
		complex64 | complex128
	}

	Number interface {
		byte | rune | Integer | Unsigned | Float | Complex
	}
)

func Sum[N Number](numbers ...N) N {

	var n N = N(0)

	for _, v := range numbers {
		n += v
	}

	return n
}

type Key int

const (
	One Key = iota
	Two
	Three
)

func (k Key) String() string {

	switch k {

	case One:
		return "one"

	case Two:
		return "two"

	case Three:
		return "three"

	default:
		return strconv.Itoa(int(k))
	}
}

func (k *Key) Scan(state fmt.ScanState, verb rune) error {

	b, err := state.Token(true, nil)

	if err != nil {
		return err
	}

	token := string(b)

	switch token {

	case One.String():
		*k = One

	case Two.String():
		*k = Two

	case Three.String():
		*k = Three

	default:
		n, err := strconv.Atoi(token)
		if err != nil {
			return err
		}
		*k = Key(n)
	}

	return nil
}

// Prints (5.1+5.2i) to stdout.
func main() {
	complexParameters := []complex64{1i, 2, 3.1, 4.2i}
	complexResult := Sum(complexParameters...)
	fmt.Println(complexResult)
}
```

In the preceding example, the definitions of `Integer`, `Float`, `Number` etc.
are reminiscent of type aliases, but they are actually type constraints for use
only in declaring generic methods like `Sum()` from the same example. You cannot
use such type constraints as if they were actual types, as you can with type
aliases.

```go
// Syntax error because Integer is a constraint, not an alias.
//
// var i Integer = 0
```

## Interfaces and Nilable Types

Anyone familiar with Go might have already objected on reading elsewhere on this
page my several assertions that Go lacks any object-oriented features. What
about interfaces? Go does support those, in a typically idiosyncratic fashion
rendering them quite different from the feature of the same name in languages
like Java. For one thing, Go provides no semantics for declaring "is a"
relationships among types, including interfaces. So one of the primary uses of
interfaces in object-oriented programming, to support polymorphism, is
eliminated. A consequence of this design choice makes any Go program that uses
interfaces more difficult to maintain in ways that interfaces in other languages
help to prevent.

Specifically, you cannot declare that a given type must implement a particular
set of interfaces. This is very much in keeping with Go's overall "no
polymorphism" policy, but makes the specification of which types implement which
interfaces rather tenuous and fragile. As discussed in [Functions and
Methods](./functional.md), functional programming techniques can be used to
implement anything for which you might want to use object-oriented programming
features such as interfaces in a fashion that is more consistent with Go's
overall design.

All of that said, it is important to understand some aspects of how Go's
interfaces work, because they are subtly pervasive throughout Go's design. As
already mentioned, you cannot explicitly declare that a given type is intended
to implement a given interface. Interfaces are implemented in Go simply by
implementing all of the methods they declare.

```go
package main

import "fmt"

// Declare an interface with two methods.
type Adder interface {
	Add(n int) int
	Subtract(m int) int
}

// Declare a function that uses the Adder interface.
func Add(m Adder, n int) int {
	return m.Add(n)
}

// Declare another function that uses the Adder interface.
func Subtract(m Adder, n int) int {
	return m.Subtract(n)
}

// Declare an alias for int that will implement Adder.
//
// Note that you can declare that int, itself, declare Adder because int is
// defined by a different package.
type MyInt int

// Implement Adder.Add().
func (m MyInt) Add(n int) int {
	return int(m) + n
}

// Implement Adder.Subtrace().
func (m MyInt) Subtract(n int) int {
	return m.Add(-n)
}

// Prints 42 to stdout
func main() {
	// Note all the type-casting required.
	m := MyInt(40)
	m = MyInt(Add(m, 4))
	fmt.Println(Subtract(m, 2))
}
```

The fact that `MyInt` implements the `Adder` interface is implicit in the fact
that it implements all of `Adder`'s methods, and by nothing else. In the given
example, it would be easy to make a number of mistakes initially when writing
this code or later when trying to tidy the code. For example, had you declared
the `MyInt` methods to take `MyInt` parameters and return `MyInt` values so as
to avoid frequent casting between `int` and `MyInt` then `MyInt` would no longer
implement `Adder` due to that difference in their method signatures. Of course,
the `Adder` interface could in this contrived example simply have been declared
to operate on `MyInt` values rather than `int` to solve the same problem. But
that solution quickly breaks down in multi-package builds.

Another non-obvious consequence of the way that the Go compiler determines
whether or not a given type implements a given interface is that every type
implements at least one interface, the empty one: `interface{}`. Go even has a
name for the empty interface, `any`. That name denotes the fact that a variable
of type `any` can hold any type of value, because every type of value implements
`interface{}`.

Let's talk that through a bit more slowly. The identity of an interface like
`Adder` in the preceding example is determined solely by its list of method
signatures. The fact that a type like `MyInt` implements `Adder` is determined
solely by the fact that there is no method declared by `Adder` that is not
implemnted by `MyInt`. The empty interface, `interface{}`, declares no methods
so there is no method of the empty interface not implemented by every type. Go
provides the type name `any` for the empty interface, so every type is
convertible to  `any`.

```go
package main

import "fmt"

func Print(v any) {
	fmt.Printf("%v\n", v)
}

// Prints the following to stdout
//
// 42
// Hello, world!
func main() {
	Print(42)
	Print("Hello, world!")
}
```

And so much for Go's vaunted type safety. Note the extensive use of the `any`
type in the logging and stacktrace library code in this repository. This is
common in these kinds of utility libraries whose whole purpose is to implement
cross-cutting concerns independent of particular types. And before any fingers
are wagged, this is just as true of standard library packages like `fmt` as it
is in mine, as shown by the preceding example.

Interfaces not only allow an escape hatch for cross-cutting concerns, they also
magically confer nilability on non-nilable types. That is a feature exploited by
Go's error reporting idiom by way of its built-in `error` interface.

```go
package main

import "fmt"

// Declare an alias for int.
type MyInt int

// Implement the error interface.
func (m MyInt) Error() string {
	return fmt.Sprintf("Error %d", m)
}

func Fail() error {
	return MyInt(42)
}

// Prints Error 42 to stdout.
func main() {
	err := Fail()

	if err != nil {
		fmt.Println(err.Error())
	}
}
```

As can be seen in the preceding example, non-nilable types like `MyInt` become
nilable when accessed as objects implementing an interface (the built-in `error`
interface, in this case). This is how Go's standard library actually implements
commonly used `error` objects. The built-in `string` type implements the `error`
interface so all that functions like `errors.New("...")` and `fmt.Errorf("...",
...)` need do is return the specified `string`, et voila!

The implicit conversion from some concrete type to some interface type,
including `any`, is not just a one-way street. While you cannot directly cast an
object whose type is an interface to some other type, Go does provide a
mechanism for what is often called "down casting" in object-oriented programming
languages.

```go
package main

import "fmt"

// Declare an alias for int.
type MyInt int

// Implement the error interface.
func (m MyInt) Error() string {
	return fmt.Sprintf("Error %d", m)
}

// Signal error using a MyInt.
func Fail() error {
	return MyInt(42)
}

// Prints the following to stdout:
//
// A MyInt was used to report Error 42
func main() {
	err := Fail()

	switch v := err.(type) {

	case MyInt:
		// v's inferred type here is MyInt
		fmt.Printf("A MyInt with value of %d was used to report %s", v, v.Error())

	case nil:
		// v's inferred type here is nil, meaning that that v == nil is true
		fmt.Println("Success")

	default:
		// v's type is not known, so just fall back on the type that is known
		// (the error interface, in this case)
		fmt.Println(err.Error())
	}
}
```

And more generally:

```go
package main

import (
	"fmt"
	"strconv"
)

func String(v any) (string, error) {
	switch v := v.(type) {

	case int:
		return strconv.Itoa(v), nil

	case string:
		return v, nil

	default:
		return "", fmt.Errorf("unsupported value: %v", v)
	}
}

// Prints the following to stdout
//
// 42
// Hello, world!
// unsupported value: 4.2
func main() {
	s, _ := String(42)
	fmt.Println(s)
	s, _ = String("Hello, world!")
	fmt.Println(s)
	_, err := String(4.2)
	fmt.Println(err.Error())
}
```
