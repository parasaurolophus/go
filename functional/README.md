_Copyright &copy; Kirk Rader 2024_

# Functional Programming in Go

> The documentation and code examples here are contrived for clarity of
> exposition and focus on the relevant topics. See
> [../logging/logger.go](../logging/logger.go) and
> [../stacktraces/stacktrace.go](../stacktraces/stacktrace.go) for real-world
> applications of these principles. Note in particular
> `logging.newAttrReplacer()`, which uses its `oldReplacer` parameter for
> "inheritance" via functional composition. Also see both
> `stacktraces.longFrameFormatter()` and `stacktraces.shortFrameFormatter()`
> which use locally declared bindings to encapsulate data values that persist
> between invocations and are used in an implemention of the dependency
> injection pattern using functional composition, among other examples in the
> logging and stacktraces utility libraries that are included in this
> repository.

Go is almost entirely lacking in features that have been considered the
hallmarks of "object oriented" programming languages since the early 1980's. It
has nothing analgous to C++ or Java classes. It provides no mechanism to share
code between entities via inheritance. It does have a capability for
programmer-defined interfaces, but they are limited in many ways that make them
often more trouble to use than they are worth from the point of view of software
architecture.

> One of the primary advantages of user-defined intefaces as first-class data
> types is for use in idioms like dependency injection and similar inversions of
> control. These uses are somewhat underminded in Go due to its lack of support
> for corollary features such as function overloading and "is a" relationships
> between interfaces. Inversion of control can still be implemented in Go in a
> flexible and efficient manner using functional composition, as described
> [below](#composition).

What Go does have are lexical closures. It has been a maxim of Lisp programmers
for decades that "objects are a poor-man's closures." In fact, the world's first
commercially-significant object-oriented programming system -- the
[Flavors](https://en.wikipedia.org/wiki/Flavors_(programming_language)) system
that was developed for the Symbolics Lisp Machine in the 1970's and evolved to
become the Common Lisp Object System (CLOS) by the late 1980's -- was
implemented using lexical closures at its core.

## Encapsulation

Closures provide directly one of the defining characteristics of object-oriented
programming: encapsulation. Each time a closure is created, it effectively
clones the state of the lexical environment visible to it at the point which it
is created. Lexical scoping rules allow for some bindings to be entirely hidden
outside of a single closure (like the "private" members of a C++ or Java class)
while other variables can be shared among a group of closures while remaining
inaccessible outside of their shared scope (like "protected" members).

Consider:

```go
package main

import "fmt"

// Declare three types of functions to be used as lexical closures.
type (

	// Return the value of a lexically-scoped binding.
	Getter func() int

	// Modify the value of a lexically-scoped binding.
	Setter func(newValue int)

	// Access two lexically-scoped bindings.
	Restorer func() int
)

// A single binding of each lexically scoped variable is visible to all of the
// closures returned by this function. Such bindings are unique to each
// invocation of this function and persist between invocations of the closures.
func MakeClosures(sharedValue int) (Getter, Setter, Restorer) {

	// Other local bindings are also shared by closures, not just those declared
	// as parameters to functions.
	previousValue := sharedValue

	// Return a closure that returns the current value of sharedValue when
	// invoked.
	get := func() int {
		return sharedValue
	}

	// Return a closure that modifies the value of previousValue and sharedValue
	// when invoked.
	set := func(newValue int) {
		previousValue = sharedValue
		sharedValue = newValue
	}

	// Return a closure that updates sharedValue to previousValue when invoked.
	restore := func() int {
		n := sharedValue
		sharedValue = previousValue
		return n
	}

	return get, set, restore
}

// Demonstrate local bindings in lexical closures.
func main() {

	// Create two sets of three closures and print the values of their initial
	// local bindings. Note that each closure retains a distinct binding for
	// sharedValue based on the original parameter value passed to the function
	// that created it.
	get1, set1, restore1 := MakeClosures(0)
	get2, set2, restore2 := MakeClosures(42)
	fmt.Printf("%3d, %3d\n", get1(), get2())

	// Use one of the closures to modify its local bindings and print their
	// current state. Note that only one of the two bindings of sharedValue will
	// have been affected.
	set1(-1)
	fmt.Printf("%3d, %3d\n", get1(), get2())

	// Use another of the closures to modify its local bindings independtly of
	// the first's.
	set2(get2() + 1)
	fmt.Printf("%3d, %3d\n", get1(), get2())

	// Note that the same rules for local bindings apply to variables declared
	// inside a lexical environment, not just those appearing as parameters to
	// functions as demonstrated by a Restorer closure's use of previousValue.
	fmt.Printf("%3d, %3d\n", restore1(), restore2())
	fmt.Printf("%3d, %3d\n", get1(), get2())
}
```

which prints the following to `stdout`:

```
  0,  42
 -1,  42
 -1,  43
 -1,  43
  0,  42
```

The preceding works as follows:

- Each time `MakeClosures(sharedValue)` is called, it returns three closures that
  are created in the same lexical scope.
  
- That scope includes a bindings of `sharedValue` `previousValue` that are
  invisible outside of that scope.
  
- The second closure defines its own lexical
  variable, `newValue`, visible only to it.
  
- When the closures are invoked, each "remembers" the particular bindings of
  `sharedValue` and `previousValue` that were in scope when it was created.
  
In particular, when `main()` executes:

```go
get1, set1, restore1 := MakeClosures(0)
get2, set2, restore2 := MakeClosures(42)
```

`get1`, `set1` and `restore1` are bound to functions, each of which share a
single binding each of `sharedValue` and `previousValue`, initially set to 0,
while `get2`, `set2` and `restore2` share a distinct binding to the value 42.
Invoking `get1()` and `get2()` at that point return 0 and 42, respectively, due
to their respective bindings of `sharedValue` at the time they were created.

Further, since `get1` and `set1` share a common binding for `sharedValue` that
is separate from the binding shared by `get2` and `set2`, invoking `set1(-1)`
updates the value that will be returned by a subsequent call to `get1()` without
affecting the value that will be returned by `get2()`.

```go
set1(-1)
```

This behavior of closures is very convenient for anonymous functions defined on
the fly for things like passing as formatting functions to the logging library
in this same repository. But as can be seen from the preceding examples, they
can also be used to implement "methods" of "instances" which are actually just
groups of closures created in a common lexical environment.

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

Functional composition can not only be used to implement a simple kind of
"inheritance," but more generally for many other idioms commonly associated with
object-oriented implementations such as dependency injection.
