_Copyright &copy; Kirk Rader 2024_

# Go Examples

- [Enumerated Values in Go](./enums/)
- [Generic Functions in Go](./generics/)
- [Interfaces, Pointers and `any`](./interfaces/)
- [Closures and Functional Programming in Go](./closures/)
- [Concurrency in Go](./concurrency/)

## Modules and Packages

In addition to the specific topics listed above, note that the directory
structure of this repository is itself a demonstration of key features of Go's
packaging mechanisms:

```
- tutorial/                 root of the module containing Go demo code 
  |
  +- go.mod                 declares the go/ directory to be a Go module
  |                         with specific meta-data such as the module
  |                         name and any required dependencies
  |
  +- closures/              source code directory for closures package
  |  |
  |  +- closures.go         demo code for using lexical closures in Go
  |  |
  |  +- closures_test.go    unit tests for closures.go
  |
  +- concurrency/           source code directory for concurrency package
  |  |
  |  +- concurrency.go      stand-alone program demonstrating basics of
  |                         goroutines and channels
  |
  +- enums/                 source code directory for enums package
  |  |
  |  +- enums.go            demo code for defining enumerated types in Go
  |  |
  |  +- enums_test.go       unit tests for enums.go
  |
  +- generics/              source code directory for generics package
  |  |
  |  +- generics.go         demo code for generic functions
  |  |
  |  +- generics_test.go    unit tests for generics.go
  |
  +- interfaces/            source code directory fo interfaces package
     |
     +- interfaces.go       demo code for defining and implementing interfaces
     |
     +- interfaces_test.go  unit tests for interfaces.go
     |
     +- main/               a compilation unit with a program entry point
        |                   function
        |
        +- main.go          source file with entry-pint function demonstrating
                            the use of the type named any
```

The _go.mod_ file at the root of the _go/_ directory defines the root of a
package hierarchy within a Go module named `parasaurolophus/tutorial`:

```
module parasaurolophus/tutorial

go 1.22.0
```

A brief summary of Go's modules and packages:

- A _package_ is a namespace

  - Symbols are scoped by the package in which they are defined

  - Visibility

    - Symbols that start with an upper-case letter are visible to other packages

    - Symbols that start with a lower-case letter are private

    - One package references public symbols from another by:

      - Specifying the fully qualified name of the other package in an `import`
        statement

      - Prefixing symbols with the simple name of the package in which they are
        defined

        ```go
        package pkg1

        import "some/other/pkg2"

        var x = pkg2.SomeType{}
        ```

      - See [Package Names](#package-names) for more information on what
        determines the fully qualified name of a package

- A _module_ is

  - The root of a package hierarchy

  - Shared configuration of the packages in the module

    - Base package name

    - Declaration of dependencies

### Package Names

Every _.go_ file must start with a `package` declaration, e.g.

```go
package enums
```

All of the symbols defined in a _.go_ file are scoped to the package declared by
that file. More than one _.go_ file can appear in a given directory, but they
all must declare the same package. The set of symbols defined by a package is
the union of all those defined in all the _.go_ files in that package's
directory.

Packages exist in a hierarchy determined by

1. The `package` declarations at the start of each _.go_ file.

2. The sub-directory nesting of the source tree under the directory containing _go.mod_.

The fully-qualified name of any given package:

- Begins with the module name from the _go.mod_ file that is in the name
  directory as the package or its nearest ancestor with such a file.

- Includes the names of any packages whose directories are ancestors of the
  given one, separated by forward slashes.

- Ends with the name used in the `package` declaration in that directory's _.go_
  files.

While the package name at each level is determined by the `package` declarations
in the _.go_ source files, their grouping and ordering is determined by the
sub-directory relationships in the file system. For readability and
maintainability you should choose directory names that map naturally to the
resulting fully-qualified package names.

In the case of this demo code, _go/go.mod_ declares a module named
`parasaurolophus/tutorial`. The _go/_ directory also includes an _enums/_
sub-directory. _go/enums/_ contains _.go_ files which begin with `package enums`
declarations. One of those is _enums.go_. The symbols defined in
_go/enums/enums.go_ thus are scoped to the `enums` package and the `enums`
package's fully-qualified name is `parasaurolophus/tutorial/enums`. The name of
the _go/enums/_ directory and the _go/enums/enums.go_ source file were chosen to
make it easy to see and navigate those package relationships.

When referencing symbols across package boundaries:

1. Add the fully qualified name of the package that provides a symbol in the
   `import` section of the _.go_ file where the symbol is to be used.

2. Qualify the symbol with just the specific package name when using it.

For example:

```go
package main

import "parasaurolophus/tutorial/enums"

var e1 enums.Enum1
```

Packages will be visible to other packages if they are either defined within the
same module, as is in this example code, or they are cited in the `require`
section of _go.mod_ file for the packages which use them. If a _go.mod_ file
requires external dependencies, use `go mod tidy` after making any changes to
make sure that the required packages are downloaded and available for use
locally and generate a corresponding `go.sum` file. Like `go.mod`, `go.sum` (if
it it exists) is a part of the definition of a module and should be checked into
version control along with all the other source code for a given project.

The `main` package is a special case. No matter where it appears in a directory
hierarchy, it exists outside of any particular package hierarchy. It exists so
that the Go tool chain can recognize the entry points of modules that contain
compilation units which produce stand-alone executables. I.e. library code
should appear in normally scoped packages while a program's entry point is
always defined as a function whose fully-qualified name is `main.main()`.

### Compilation Units

The Go tool chain operates almost exclusively at the level of specific package
directories. This allows a single module to contain the source code for multiple
libraries and executables. This can be exploited to advantage in maintaining a
large code base with many internal dependencies but needs to be handled with a
degree of caution. For example, all the packages in a given module will share a
single set of versioned dependencies such that a breaking change to a dependency
for one compilation unit could end up breaking multiple compilation units in the
same module.

The alternative is to have each compilation unit exist within its own module.
That provides the greatest flexibility for mixing and matching versions of
executable code and dependencies in production, but also vastly increases the
friction on day-to-day development. Each developer will need to manage
dependencies locally using temporary adjustments to local copies of _go.mod_
using `replace` blocks and the like as well as using _go.work_ files, none of
which should be checked into the version control system. If this is the approach
used, be prepared to spend a certain amount of time setting up your IDE at the
start of each new development task and cleaning up the messes that will
regularly be caused when broken _go.mod_ files are accidentally merged upstream.

## Unit Testing

The Go tool chain assumes that any file name ending in *_test.go* contains unit
tests. Such files are treated specially by commands like `go build`, `go doc`
and `go test`. The same is true for packages whose names end with `_test`. Such
packages are assumed to contain types and functions of use only during testing
that should be excluded from normal builds.

Within a *_test.go* file, the framework invoked by `go test` assumes that any
function with a specific signature and whose name begins with `Test` is a test
function. See [./enums/enums_test.go](./enums/enums_test.go),
[./generics/generics_test.go](./generics/generics_test.go) and the like for
examples.
