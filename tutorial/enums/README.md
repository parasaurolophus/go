_Copyright &copy; Kirk Rader 2024_

# Enumerated Types in Go

Go does not have any mechanism for creating first-class enumerated types in the
style of C++ or Java `enum`. The `enums` package in this directory demonstrates
how to get as close as possible using Go's features. In particular:

- Use aliases for `int` to define application-specific scalar types.

- Use `iota` in a `const` block to declare sequential enumerated values of each
  such type.

- Implement the `fmt.Stringer` and `fmt.Scanner` interfaces for such types.

See [enums.go](./enums.go), [enums_test.go](./enums_test.go) for the demo code
and unit tests.
