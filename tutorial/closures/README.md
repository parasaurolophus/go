_Copyright &copy; Kirk Rader 2024_

# Lexical Closures and Functional Programming in Go

- [./closures.go](./closures.go)
- [./closures_test.go](./closures_test.go)

Lexical closures (or just closures, for brevity) are a feature of any
programming language intended to support the functional programming paradigm.
Just as "objects" are the essential metaphor of "object oriented" programming,
"functions" are the fundamental organizing primicple of "functional" or
"function oriented" software design and implementation.

A closure is a function that:

1. Can be passed to and returned from functions as data.

2. Can be invoked like ordinary named functions at any point in the execution of
   other functions.

3. Encapsulate the bindings of all variables that were in lexical scope when
   they were created.

The last point is key to understanding the power and usefulness of closures.
When a function invokes a closure, the closure's code is executed not in the
calling function's lexical environment (as it would be when invoking a named
function in Go or using simple function pointers in C) but in a "closed" copy of
the lexical environment that was in effect when the closure was created.

But what exactly is a "lexical environment" and what does it mean to "close"
one? Consider:

```go
// Copyright Kirk Rader 2024

package main

import "fmt"

func ReturnInt() int {

	// Entering a given lexical scope, e.g. by calling a function, establishes a
	// lexical environment in which locally defined identifiers are visible,
	// while being invisible outside the body of that scope.

	// In this case, localVar is only visible to the body of the ReturnInt()
	// function. Effectively, ReturnInt()'s lexical scope is everything between
	// the curly braces directly above and below this comment.

	localVar := 42
	return localVar
}

func ReturnClosure() func() int {

	// Lexical scopes nest such that a scope that is defined inside another can
	// see the variables bound in the scope in which it was created even though
	// variables bound within the inner scope are invisible to the outer scope.

	// In this case, both the body of ReturnClosure() and the body of the
	// anonymous function assigned to the variable named closure can see
	// closedVar.

	// In addition, the anonymous function returned by ReturnClosure()
	// "remembers" the lexical scope in which it was created. Each time it is
	// invoked, it will see the binding for closedVar that was in effect when it
	// was created. That is what is meant by saying that such a function
	// "closes" its lexical environment (which is terminology that actually
	// refers to optimizations that occur when no lexical closure is created vs
	// when a snapshot of a given lexical environment needs to be kept and
	// referenced by one or more closures). I.e. when the function body of a
	// lexical closure is invoked, it executes inside its own closed environment
	// rather than the lexical environment from which it was invoked.

	closedVar := 42
	closure := func() int { return closedVar }
	return closure
}

func ReturnTwoClosures() (func() int, func(int)) {

	// A new set of local bindings are created each time the lexical environment
	// is entered.

	// All closures created in a given lexical environment share a common set of
	// closed bindings.

	// In this case, a single binding of anotherClosedVar will be visible to
	// both closures returned by a given invocation of ReturnTwoClosures() but
	// each invocation will return a new, unique binding.

	anotherClosedVar := 0
	getter := func() int { return anotherClosedVar }
	setter := func(n int) { anotherClosedVar = n }
	return getter, setter
}

// Prints
//
// 42
// 42
// 42,  0
//
// to stdout. The fact that the two numbers on the last line are different
// demonstrates that closures are given access to a unique "snapshot" of their
// lexical environment when they are created.
func main() {

	fmt.Println(ReturnInt())
	f := ReturnClosure()
	fmt.Println(f())
	getter1, setter1 := ReturnTwoClosures()
	getter2, _ := ReturnTwoClosures()
	setter1(42)
	fmt.Printf("%2d, %2d\n", getter1(), getter2())
}
```

See [./closures.go](./closures.go) and
[[](./closures_test.go)](./closures_test.go) for "hello world" level examples of
what can be accomplished using closures, including a version of the preceding
code.

Note in particular the implementation of `closures.PassContinuations()` and its
unit test. While implemented using the mechanisms and terminology of functional
programming and the mathematics on which it is based, it is also a perfect
example of the dependency injection pattern and similar inversions of control
that would rely on interfaces and objects in an object-oriented programming
system. It also demonstrates Go's typically idiosyncratic approach to stack
unwinding protection (`defer`) and ways to use it to `recover()` from a
`panic()` in case some injected dependency dramatically misbehaves.

## Functional vs Object-Oriented Programming

Historically, closures and the functional programming paradigm gave rise to
object orientation as a practical feature of programming languages. The world's
first commercially significant object-oriented programming system,
[Flavors](https://en.wikipedia.org/wiki/Flavors_(programming_language)), was
developed in the 1970's at MIT for what became the
[Symbolics](https://en.wikipedia.org/wiki/Symbolics) Lisp Machine.

> As an aside, _symbolics.com_ was the world's first registered commercial
> Internet domain because we used to get into just about everything, back in the
> day.

Flavors was based on closures as its underlying implementation mechanism. It is
a saying to this day among Lisp programmers that, "objects are a poor-man's
closures."

To see how concepts of the functional and object paradigms compare:

| Object Oriented Concept    | Functional Concept                                   |
|----------------------------|------------------------------------------------------|
| object / data members      | closed environment                                   |
| methods / member functions | lexical closures                                     |
| inheritance                | functional composition (one function calling others) |

In Flavors, which evolved into and survives today as the Common Lisp Object
System (CLOS), instances of classes really are, under the hood, just closed
environments captured when the lexical closures implementing their methods are
created by a given class's constructor. Inheritiance really is implemented by
one such method automatically calling others according to the rules for "method
combination" provided by Flavors / CLOS.

> The _before_ / _primary_ / _after_ method invocation pattern that is provided
> by a number of object-oriented programming systems is called "daemon method
> combination" and was invented by the designers of Flavors. In Flavors / CLOS,
> however, daemon method combination is merely the default rule. There are other
> built-in method combination rules that can be configured for a given class
> such as "_or_ method combination," "_and_ method combination," "_progn_ method
> combination" and so on. It is also possible for programmers to provide their
> own custom method combination rules on a class-by-class basis (but that way
> lies madness!)

All of that is facilitated by the underlying functional programming paradigm
from which Flavors / CLOS was built. In other words, moving from an
object-oriented mindset to the functional programming paradigm can be viewed as
simply removing one layer of abstraction from a given software design. Removing
layers of abstraction usually results in improved efficiency and expressive
power -- and with power comes responsibility.
