_Copyright &copy; Kirk Rader 2024_

# Interfaces, Pointers and the `any` Type in Go

> N.B. this README.md barely scratches the surface of all the complexities
> surrounding interfaces in Go. It makes no reference to the many options for
> "embedding" interfaces, defining them using type sets similar to type
> constraints for generics, nor any of the special rules that apply to such
> scenarios. Brew yourself a hot beverage, settle into a comfortable seat, and
> then see [Interface Types](https://go.dev/ref/spec#Interface_types) in the Go
> language specification for all of the many, messy details.

The demo code shows the most common pattern for defining and using interfaces,
with reference in comments to some subtle pitfalls of which to be wary:

- [./interfaces.go](./interfaces.go)
- [./interfaces_test.go](./interfaces_test.go)
- [./main/main.go](./main/main.go)

Since Go deliberately eschews any language features that are conventionally
understood to support object-oriented programming techniques, its version of
interfaces works quite differently from the corresponding feature in many other
languages (notably Java, which pioneered the concept of `interface` as a
first-class data type).

> Said it before and will say it again: For most purposes that you might use an
> interface in an object-oriented paradigm, you may find that [functional
> programming](../closures/) is a better fit with Go's overall design.

For this reason, Go interfaces are mainly useful for two things:

1. Cross-cutting concerns like the built-in `fmt.Stringer` and `fmt.Scanner`
   interfaces.

2. Enablement of the built-in `any` type and similar Go-specific idioms.

## The `any` Type

In Go's documentation, its designers make much of the strictness of its type
checking, intentional lack of support for (most, but far from all) implicit type
conversions and its overall antipathy for polymorphism (including function
overloading, which Go forbids). But then they included interfaces. Since
interfaces exist for no other purpose than polymorphism, they introduce many
confusing syntactical patterns and usage idioms as noted in comments in the
bodies of a number of functions defined in this package's _.go_ files. One of
the most confusing to use and explain -- and yet quite ubiquitous in practice --
is the `any` type.

Go's `any` is a synonym for the type more properly known as `interface{}`, i.e.
"the empty interface." To understand why `interface{}` is treated as a type that
can hold any value, you must understand one of the many oddities of interfaces
in Go. There is no mechanism for declaring that a given type, `T`, must
implement a given interface, `I`.

> To do so, Go would have to embrace the concept of "is-a" relationships which
> would be the first step down the slippery slope toward inheritance and true
> polymorphism, which would be antithetical to the fundamental goals of Go's
> designers.

Instead, the compiler infers that `T` implements `I` if and only if there is no
method declared by `I` that `T` fails to implement. Since `interface{}` has no
methods, there is no method in `interface{}` that any type fails to implement,
so every type implicitly implements `interface{}`.

> The preceding is probably less confusing to those who have studied Set Theory
> enough to have been introduced to the notion of "the empty set" and the
> reasoning for why it is considered a sub-set of every set, including itself.

In any case, take it for granted that `any` is an interface and that every type
is understood to implement it. On the one hand, this allows gross violations of
Go's supposed type strictness. On the other, interfaces in general and `any` in
particular are made less useful than you might imagine in Go due to the many
constraints and syntactic oddities introduced by trying to retain some degree of
type safety while avoiding true polymorphism.

> As noted above, for these among many other reasons Go is actually better
> suited to the functional paradigm where lexical closures take the place of
> values of types that implement interfaces.

## More On Interfaces and Type Conversions

Go's syntax is largely inspired by old-school K&R C.

> While its syntax is C-like, Go's semantics are inspired mostly by Lisp
> dialects of a similar vintage to K&R C, but that is a discussion better had in
> the context [closures and the functional programming paradigm](../closures/)

On the surface, Go's general preference for passing by value along with its
"pointer" construct (and the `*` and `&` syntax for using them) is taken
directly from C that is so 70's, Go's logo might as well be a pair of well worn,
paisley-vented bell-bottom jeans. But under the surface Go's data representation
and memory management mechanisms are the polar opposite of C or C++. This leads
to oddities and out-right inconsistencies in syntax and data type conversion
rules in various aspects of Go's design.

If you first learned to program in languages like Java where passing by
reference, rather than value, is the norm it may never have occured to you to
wonder whe the following pattern is possible, let alone standard, in Go:

```go
var err error

err = someFunc()

if err != nil {
	// handle the error
}
```

As distinctively Go as that is, there must be something special about the
`error` type that allows it to be nilable without it being referenced through a
pointer, as would be required in C. It turns out that this oddity is not unique
to `error`. It is a feature of every interface type in Go.

> A truism that dates back to at least the 70's is that "a feature is a bug as
> described by the marketing department."

Specifically, every type, nilable or not, becomes nilable when accessed through
any interface it implements. Since every type implements `any` and `any` is a
synonym for `interace{}`, every type can be coerced to become nilable. That is
why the `string` returned by `errors.New("...")` becomes nilable by virtue of
being wrapped by the `error` interface. Because that is what interfaces actually
do, under the covers. When a value of type `T` is assigned to a variable of some
interface, `I`, implemented by `T` the compiler performs an implicit type
conversion from `T` to `I`. Inside the Go runtime, values of interfaces store a
reference to the particular implementing value (the value of type `T`, in this
case). That reference may be erased such that the value of `I` no longer
references any underlying value. The Go syntax for such "disconnected" interface
values is `nil`.

> That is why if you look under the hood, there are actually multiple
> runtime-level types which end up being represented as `nil` in Go's syntax. A
> "disconnected interface value" is actually different from a  "pointer to
> nothing" but the compiler sorts out the different underlying values to use for
> `nil` depending on context.

All of these subtleties have knock-on consequences for Go's syntax and
semantics. For example, if you start defining interfaces like
`interfaces.Counter` from [./interfaces.go](./interfaces.go) whose semantics
require methods with pointer receivers, expect a lot of ongoing confusion around
when you must and when you must not use `&` when creating and using values that
implement them, with less than entirely helpful diagnostics from the compiler to
figure out the right syntax for what you are attempting at any given spot in the
code. Again, see the usage and associated comments in
[./interfaces.go](interfaces.go) and [./main/main.go](./main/main.go) for
examples.
