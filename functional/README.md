_Copyright &copy; Kirk Rader 2024_

# Functional Programming in Go

Go is almost entirely lacking in features that have been considered the
hallmarks of "object oriented" programming languages since the early 1980's. It
has nothing analgous to C++ or Java style classes. It provides no mechanism to
share code between entities via inheritance. It does have a capability for
programmer-defined interfaces, but they are limited in many ways that make them
often more trouble to use than they are worth.

What it does have are lexical closures. It has been a maxim of Lisp programmers
for decades that "objects are a poor-man's closures." In fact, the very first
commercially-significant object-oriented programming system -- the
[Flavors](https://en.wikipedia.org/wiki/Flavors_(programming_language)) system
that was developed for the Symbolics Lisp Machine in the 1970's and evolved to
become the Common Lisp Object System (CLOS) in the 1980's -- was implemented
using lexical closures its core.

## Encapsulation

Closures provide directly one of the defining characteristics of object-oriented
programming: encapsulation. Each time a closure is created, it effectively
clones the state of the lexical environment visible to it at the point which it
is created. Lexical scoping rules allow for some variables to entirely hidden
outside of a given closure (like the "private" members of a C++ or Java class)
while other variables can be shared among a group of closures while remaining
inaccessible outside of their shared scope (like "protected" members).

Consider:

```go
package main

import "fmt"

func MakeClosures(sharedValue int) (func() int, func(int)) {

	get := func() int {
		return sharedValue
	}

	set := func(newValue int) {
		sharedValue = newValue
	}

	return get, set
}

func main() {

	get1, set1 := MakeClosures(0)
	get2, set2 := MakeClosures(42)
	fmt.Printf("%3d, %3d\n", get1(), get2())
	set1(-1)
	fmt.Printf("%3d, %3d\n", get1(), get2())
	set2(get2() + 1)
	fmt.Printf("%3d, %3d\n", get1(), get2())
}
```

which prints the following to `stdout`:

```
  0,  42
 -1,  42
 -1,  43
```

The preceding works as follows:

- Each time `MakeClosures(sharedValue)` is called, it returns two closures that
  are created in the same lexical scope. That scope includes a definition of the
  variable `sharedValue` that is invisible outside of that scope. The second
  closure defines its own lexical variable, `newValue`, visible only to it. When
  the closures are invoked, each "remembers" the particular binding of
  `sharedValue` that was in scope when it was created. In particular, when
  `main()` executes

```go
get1, set1 := MakeClosures(0)
get2, set2 := MakeClosures(42)
```

`get1` and `set1` are bound to two functions, each of which share a single
binding of `sharedValue`, initially set to 0, while `get2` and `set2` share a
distinct binding to the value 42. Invoking `get1()` and `get2()` at that point
return 0 and 42, respectively, due to their respective bindings of `sharedValue`
at the time they were created.

Further, since `get1` and `set1` share a common binding for `sharedValue` that
is separate from the binding shared by `get2` and `set2`, invoking `set1(-1)`
updates the value that will be returned by a subsequent call to `get1()` without
affecting the value that will be returned by `get2()`.

This behavior of closures is very convenient for anonymous functions defined on
the fly for things like passing as formatting functions to logging libraries.
But as can be seen from the preceding examples, they can also be used to
implement "methods" of "instances" which are actually closures created in a
common lexical environment.

## Composition

Another defining characteristic of what many people think of as "object
orientation" is inheritance. In languages like C++ and Java, when one class is
declared to extend another it automatically implements all of the inheritable
fields and methods of the parent class. This encourages code-reuse and thus
reduced redundancy (among other aspects of good code hygiene).

Go makes no provision for inheritance. While it has a very idiosyncratic version
of programmer defined types, it has no notion of "classes" in the sense of C++
or Java. It also does not support function overloading; two functions cannot be
defined with the same name but different numbers or types of arguments.

For these reasons it cannot directly implement anything very close to Flavors or
CLOS, let alone C++ or Java. You can, however, share implementation by
functional composition (i.e. calling one function in the implementation of
another). The following prints `ok` to `stdout`:

```go
package main

import (
	"fmt"
)

type (

	// Type of closure used in a demonstration of basic functional programming
	// principles.
	//
	// A Transformer accepts any type of paramter and returns the same or a
	// different value. It returns nil as a second result or else an error
	// object providing diagnostic information if a transformation fails.
	Transformer func(input any) (any, error)
)

// Return a Transformer which returns the result of adding n to its input, when
// passed an int.
//
// This demonstrates that a closure provides an encapsulation of the state of
// lexical variables that are in scope when it is created.
func IntAdder(n int) Transformer {

	return func(input any) (any, error) {

		switch v := input.(type) {

		case int:
			return v + n, nil

		default:
			return input, fmt.Errorf("unsupported value %v", v)
		}
	}
}

// Return a Transformer that returns the result of invoking the set of
// Transformers passed to it, in succession.
//
// This demonstrates the use of functional composition to share implementations
// in the absence of inheritance between instances of classes.
func CompositeTransformer(transformers ...Transformer) Transformer {

	return func(input any) (any, error) {

		var err error
		output := input

		for _, transformer := range transformers {

			output, err = transformer(output)

			if err != nil {
				return input, err
			}
		}

		return output, nil
	}
}

func main() {

	failed := 0
	oneAdder := IntAdder(1)
	twoAdder := IntAdder(2)
	composite := CompositeTransformer(twoAdder, oneAdder, twoAdder)

	output, err := composite(0)

	if err != nil {

		fmt.Println(err.Error())
		failed += 1
	}

	switch v := output.(type) {

	case int:
		if v != 5 {
			fmt.Printf("expected 5, got %d\n", v)
			failed += 1
		}

	default:
		fmt.Printf("expected 5, got %v\n", v)
		failed += 1
	}

	if failed == 0 {
		fmt.Println("ok")
	} else {
		fmt.Printf("%d failed\n", failed)
	}
}
```
