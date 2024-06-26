# What Your CS Professor Might Have Failed to Tell You

## All Programming Languages Are Modeled on the Lambda Calculus

To understand what all programming languages have in common, first learn to
think this way<sup>[1](#church)</sup>:

**Alonzo Church's Lambda Calculus**

> _let y = &lambda;x.+ x 1 (2_
>
> _y - 3 = 0_

It will then not take long for it to feel entirely natural to think this
way<sup>[2](#lisp)</sup>:

**Scheme**

```scheme
(let ((y ((lambda (x) (+ x 1)) 2)))
    (equal? (- y 3) 0))
#t
```

Which then makes it easy to think any of these ways:

**Ruby**

```ruby
irb(main):002:0> y = ->(x) { x + 1 }.(2)
=> 3
irb(main):003:0> y - 3 == 0
=> true
```

**JavaScript**

```javascript
let y = ((x) => { return x + 1 }) (2)
console.log(y - 3 == 0)
```

**Go**

```go
y := func(x int) int { return x + 1 }(2)
fmt.Println(y-3 == 0)
```

**C++**

```c++
auto f = [](int x) -> int
{
    return x + 1;
};
int y = f(2);
std::cout << (y - 3 == 0) << std::endl;
```

And, yes, that is why anonymous functions in any language are often referred to
as "lambdas" and why AWS chose _Lambda_ as the brand name for its supposedly
"serverless" SOA platform.

---

### Notes

<a id="church"><sup>1</sup></a> That is how Alonzo Church would have written a
simple example on a chalk board (yes, chalk) when teaching about his Lambda
Calculus, as witnessed by the author of this document _circa_ 1980 while Church
was still teaching at UCLA. It can be read in English as:

> _Let y be the result of passing 2 to a function which adds 1 to its argument.
> As a consequence, y minus 3 is equal to 0._

And, no, there are no typos in that example of the Lambda Calculus. Church was
an enthusistic proponent of "prefix notation" (sometimes referred to as "Polish
notation" after the school of mathematicians from Poland who first popularized
it). Thus he used the form _+ x 1_ where traditional mathematicians and most
programming languages would use _x + 1_.

Further, among the advantages of prefix notation is that it greatly reduces the
need for grouping symbols like parentheses, and it entirely eliminates the need
for pairs of such symbols. That is why when separating the definition of a
function using his &lambda; operator from the arguments to which the function
was applied Church would only write a left parenthesis as in _&lambda;x.+ x 1
(2_

Note that the Lambda Calculus was developed to address problems in a field of
mathematics known as Computability Theory well before modern digital computing
devices existed or languages were devised in which to write programs for them.
The very term, _computer_, derives from the fact that the designs of the first
generation of such hardware (e.g. ENIAC and its successor EDVAC) were derived
consciously and explicitly by reference to Church's and Turing's work in
Computability Theory. I.e. Computer Science was historically considered the
applied math complement to Computability Theory, and its continuation into the
realm of practical engineering (hence the term, _software engineering_ despite
how little resemblance the latter has to other kinds of "engineering").

<a id="lisp"><sup>2</sup></a> Despite claims about FORTRAN (which, in fairness,
did see commercial application earlier), work on creating the world's first
high-level programming languages began with LISP, which was originally
conceived of as directly embodying Church's Lambda Calculus to the extent
possible given the limitations of a digital computing device with finite
storage capacity.

---

## All Programs Are Really Implemented in Machine Code

To really understand how to deal with corner cases and realistically complex
resource management and performance issues, first learn to think this way:

```
f:
    sub     sp,     sp,         #0x20
    str     x0,     [sp,        #8]
    str     w1,     [sp,        #4]
    str     wzr,    [sp,        #28]
    b       74c <f+0x38>
    ldr     x0,     [sp,        #8]
    add     x1,     x0,         #0x4
    str     x1,     [sp,        #8]
    ldr     w1,     [x0]
    add     w1,     w1,         #0x1
    str     w1,     [x0]
    ldr     w0,     [sp,        #28]
    add     w0,     w0,         #0x1
    str     w0,     [sp,        #28]
    ldr     w1,     [sp,        #28]
    ldr     w0,     [sp,        #4]
    cmp     w1,     w0
    b.cc    728 <f+0x14>  // b.lo, b.ul, b.last
    nop
    nop
    add     sp,     sp,     #0x20
    ret
```

Then understand why that is the literal implementation emitted by a C compiler
for:

```c
void f(int *x, unsigned n)
{
    for (unsigned i = 0; i < n; ++i)
    {
        *(x++) += 1;
    }
}
```

This lets one be aware of how altogether more, and more complex, machine code
is emitted for "higher level" languages with "sophisticated" (which is the
polite way of saying "bloated") run time features like automatic memory
allocation and deallocation, built-in types with "rich" (another euphism)
features such as collections with elastic capacity and the like.

To illustrate this point, consider the following complete C program, from which
the preceding snippet was taken:

```c
#include <stdio.h>

void f(int *x, unsigned n)
{
    for (unsigned i = 0; i < n; ++i)
    {
        *(x++) += 1;
    }
}

int main(int argc, char **argv)
{
    int y[3] = {0, 1, 2};
    f(y, 3);
    for (unsigned i; i < 3; ++i)
    {
        int v = i + 1;
        if (y[i] != v)
        {
            printf("expected %d, got %d\n", v, y[i]);
            return v;
        }
    }
    printf("success\n");
    return 0;
}
```

That program defines and invokes the function `f`, which takes an array of
integers and adds 1 to each element, using C's standard idiom for arrys
involving "pointers" and "pointer arithmetic". The programm then checks the
contents of the array after `f` returns, exiting with a status code indicating
failure corresponding to any element whose value is incorrect or with a status
code indicating success if all elements pass validation.

Now consider the following Go program:

```go
package main

import (
	"fmt"
	"os"
)

func f(x []int) {
	for i := 0; i < len(x); i += 1 {
		x[i] += 1
	}
}

func main() {
	y := []int{0, 1, 2}
	f(y)
	for i := 0; i < len(y); i += 1 {
		v := i + 1
		if y[i] != v {
			fmt.Printf("expected %d, got %d\n", v, y[i])
			os.Exit(v)
		}
	}
	fmt.Printf("success\n")
}
```

The preceding two programs implement as close to identical logic, using as
close to identical data types as the syntax and semantics of the two languages
allow. Both are barely a notch above "hello world" in features and
functionality. But note that while the executable emitted by the GNU C compiler
targeting Linux on ARM64 is 69K bytes, the corresponding Go executable is 1.9M
bytes. This ratio of a Go executable being several orders of magnitude larger
than a functionally equivalent C program is typical, not just for Go but for
any "higher level" compiled language. (Calculations of the relative size and
efficiency of programs executed by interpreted languages like JavaScript or
Ruby are far more complicated and nuanced affairs, but suffice it to say that
such languages cannot by any reasonable measure be regarded as significantly
more efficient than compiled high-level languages like Go. Such calculations
for virtual-machine based compiled languages like Java -- which are first
compiled, then interpreted at run time -- are even more fraught, but all the
same realities are revelaed when assessing their memory utilization and run
time performace.)

Here is the output of `objdump -d` for the C program:

```

example:     file format elf64-littleaarch64


Disassembly of section .init:

00000000000005f0 <_init>:
 5f0:	d503201f 	nop
 5f4:	a9bf7bfd 	stp	x29, x30, [sp, #-16]!
 5f8:	910003fd 	mov	x29, sp
 5fc:	9400003e 	bl	6f4 <call_weak_fn>
 600:	a8c17bfd 	ldp	x29, x30, [sp], #16
 604:	d65f03c0 	ret

Disassembly of section .plt:

0000000000000610 <.plt>:
 610:	a9bf7bf0 	stp	x16, x30, [sp, #-16]!
 614:	f00000f0 	adrp	x16, 1f000 <__FRAME_END__+0x1e5c8>
 618:	f947fe11 	ldr	x17, [x16, #4088]
 61c:	913fe210 	add	x16, x16, #0xff8
 620:	d61f0220 	br	x17
 624:	d503201f 	nop
 628:	d503201f 	nop
 62c:	d503201f 	nop

0000000000000630 <__libc_start_main@plt>:
 630:	90000110 	adrp	x16, 20000 <__libc_start_main@GLIBC_2.34>
 634:	f9400211 	ldr	x17, [x16]
 638:	91000210 	add	x16, x16, #0x0
 63c:	d61f0220 	br	x17

0000000000000640 <__cxa_finalize@plt>:
 640:	90000110 	adrp	x16, 20000 <__libc_start_main@GLIBC_2.34>
 644:	f9400611 	ldr	x17, [x16, #8]
 648:	91002210 	add	x16, x16, #0x8
 64c:	d61f0220 	br	x17

0000000000000650 <__gmon_start__@plt>:
 650:	90000110 	adrp	x16, 20000 <__libc_start_main@GLIBC_2.34>
 654:	f9400a11 	ldr	x17, [x16, #16]
 658:	91004210 	add	x16, x16, #0x10
 65c:	d61f0220 	br	x17

0000000000000660 <abort@plt>:
 660:	90000110 	adrp	x16, 20000 <__libc_start_main@GLIBC_2.34>
 664:	f9400e11 	ldr	x17, [x16, #24]
 668:	91006210 	add	x16, x16, #0x18
 66c:	d61f0220 	br	x17

0000000000000670 <puts@plt>:
 670:	90000110 	adrp	x16, 20000 <__libc_start_main@GLIBC_2.34>
 674:	f9401211 	ldr	x17, [x16, #32]
 678:	91008210 	add	x16, x16, #0x20
 67c:	d61f0220 	br	x17

0000000000000680 <printf@plt>:
 680:	90000110 	adrp	x16, 20000 <__libc_start_main@GLIBC_2.34>
 684:	f9401611 	ldr	x17, [x16, #40]
 688:	9100a210 	add	x16, x16, #0x28
 68c:	d61f0220 	br	x17

Disassembly of section .text:

00000000000006c0 <_start>:
 6c0:	d503201f 	nop
 6c4:	d280001d 	mov	x29, #0x0                   	// #0
 6c8:	d280001e 	mov	x30, #0x0                   	// #0
 6cc:	aa0003e5 	mov	x5, x0
 6d0:	f94003e1 	ldr	x1, [sp]
 6d4:	910023e2 	add	x2, sp, #0x8
 6d8:	910003e6 	mov	x6, sp
 6dc:	f00000e0 	adrp	x0, 1f000 <__FRAME_END__+0x1e5c8>
 6e0:	f947ec00 	ldr	x0, [x0, #4056]
 6e4:	d2800003 	mov	x3, #0x0                   	// #0
 6e8:	d2800004 	mov	x4, #0x0                   	// #0
 6ec:	97ffffd1 	bl	630 <__libc_start_main@plt>
 6f0:	97ffffdc 	bl	660 <abort@plt>

00000000000006f4 <call_weak_fn>:
 6f4:	f00000e0 	adrp	x0, 1f000 <__FRAME_END__+0x1e5c8>
 6f8:	f947e800 	ldr	x0, [x0, #4048]
 6fc:	b4000040 	cbz	x0, 704 <call_weak_fn+0x10>
 700:	17ffffd4 	b	650 <__gmon_start__@plt>
 704:	d65f03c0 	ret
 708:	d503201f 	nop
 70c:	d503201f 	nop

0000000000000710 <deregister_tm_clones>:
 710:	90000100 	adrp	x0, 20000 <__libc_start_main@GLIBC_2.34>
 714:	91010000 	add	x0, x0, #0x40
 718:	90000101 	adrp	x1, 20000 <__libc_start_main@GLIBC_2.34>
 71c:	91010021 	add	x1, x1, #0x40
 720:	eb00003f 	cmp	x1, x0
 724:	540000c0 	b.eq	73c <deregister_tm_clones+0x2c>  // b.none
 728:	f00000e1 	adrp	x1, 1f000 <__FRAME_END__+0x1e5c8>
 72c:	f947e021 	ldr	x1, [x1, #4032]
 730:	b4000061 	cbz	x1, 73c <deregister_tm_clones+0x2c>
 734:	aa0103f0 	mov	x16, x1
 738:	d61f0200 	br	x16
 73c:	d65f03c0 	ret

0000000000000740 <register_tm_clones>:
 740:	90000100 	adrp	x0, 20000 <__libc_start_main@GLIBC_2.34>
 744:	91010000 	add	x0, x0, #0x40
 748:	90000101 	adrp	x1, 20000 <__libc_start_main@GLIBC_2.34>
 74c:	91010021 	add	x1, x1, #0x40
 750:	cb000021 	sub	x1, x1, x0
 754:	d37ffc22 	lsr	x2, x1, #63
 758:	8b810c41 	add	x1, x2, x1, asr #3
 75c:	9341fc21 	asr	x1, x1, #1
 760:	b40000c1 	cbz	x1, 778 <register_tm_clones+0x38>
 764:	f00000e2 	adrp	x2, 1f000 <__FRAME_END__+0x1e5c8>
 768:	f947f042 	ldr	x2, [x2, #4064]
 76c:	b4000062 	cbz	x2, 778 <register_tm_clones+0x38>
 770:	aa0203f0 	mov	x16, x2
 774:	d61f0200 	br	x16
 778:	d65f03c0 	ret
 77c:	d503201f 	nop

0000000000000780 <__do_global_dtors_aux>:
 780:	a9be7bfd 	stp	x29, x30, [sp, #-32]!
 784:	910003fd 	mov	x29, sp
 788:	f9000bf3 	str	x19, [sp, #16]
 78c:	90000113 	adrp	x19, 20000 <__libc_start_main@GLIBC_2.34>
 790:	39410260 	ldrb	w0, [x19, #64]
 794:	35000140 	cbnz	w0, 7bc <__do_global_dtors_aux+0x3c>
 798:	f00000e0 	adrp	x0, 1f000 <__FRAME_END__+0x1e5c8>
 79c:	f947e400 	ldr	x0, [x0, #4040]
 7a0:	b4000080 	cbz	x0, 7b0 <__do_global_dtors_aux+0x30>
 7a4:	90000100 	adrp	x0, 20000 <__libc_start_main@GLIBC_2.34>
 7a8:	f9401c00 	ldr	x0, [x0, #56]
 7ac:	97ffffa5 	bl	640 <__cxa_finalize@plt>
 7b0:	97ffffd8 	bl	710 <deregister_tm_clones>
 7b4:	52800020 	mov	w0, #0x1                   	// #1
 7b8:	39010260 	strb	w0, [x19, #64]
 7bc:	f9400bf3 	ldr	x19, [sp, #16]
 7c0:	a8c27bfd 	ldp	x29, x30, [sp], #32
 7c4:	d65f03c0 	ret
 7c8:	d503201f 	nop
 7cc:	d503201f 	nop

00000000000007d0 <frame_dummy>:
 7d0:	17ffffdc 	b	740 <register_tm_clones>

00000000000007d4 <f>:
 7d4:	d10083ff 	sub	sp, sp, #0x20
 7d8:	f90007e0 	str	x0, [sp, #8]
 7dc:	b90007e1 	str	w1, [sp, #4]
 7e0:	b9001fff 	str	wzr, [sp, #28]
 7e4:	1400000a 	b	80c <f+0x38>
 7e8:	f94007e0 	ldr	x0, [sp, #8]
 7ec:	91001001 	add	x1, x0, #0x4
 7f0:	f90007e1 	str	x1, [sp, #8]
 7f4:	b9400001 	ldr	w1, [x0]
 7f8:	11000421 	add	w1, w1, #0x1
 7fc:	b9000001 	str	w1, [x0]
 800:	b9401fe0 	ldr	w0, [sp, #28]
 804:	11000400 	add	w0, w0, #0x1
 808:	b9001fe0 	str	w0, [sp, #28]
 80c:	b9401fe1 	ldr	w1, [sp, #28]
 810:	b94007e0 	ldr	w0, [sp, #4]
 814:	6b00003f 	cmp	w1, w0
 818:	54fffe83 	b.cc	7e8 <f+0x14>  // b.lo, b.ul, b.last
 81c:	d503201f 	nop
 820:	d503201f 	nop
 824:	910083ff 	add	sp, sp, #0x20
 828:	d65f03c0 	ret

000000000000082c <main>:
 82c:	a9bc7bfd 	stp	x29, x30, [sp, #-64]!
 830:	910003fd 	mov	x29, sp
 834:	b9001fe0 	str	w0, [sp, #28]
 838:	f9000be1 	str	x1, [sp, #16]
 83c:	90000000 	adrp	x0, 0 <__abi_tag-0x278>
 840:	9124a001 	add	x1, x0, #0x928
 844:	9100a3e0 	add	x0, sp, #0x28
 848:	f9400022 	ldr	x2, [x1]
 84c:	f9000002 	str	x2, [x0]
 850:	b9400821 	ldr	w1, [x1, #8]
 854:	b9000801 	str	w1, [x0, #8]
 858:	9100a3e0 	add	x0, sp, #0x28
 85c:	52800061 	mov	w1, #0x3                   	// #3
 860:	97ffffdd 	bl	7d4 <f>
 864:	14000019 	b	8c8 <main+0x9c>
 868:	b9403fe0 	ldr	w0, [sp, #60]
 86c:	11000400 	add	w0, w0, #0x1
 870:	b9003be0 	str	w0, [sp, #56]
 874:	b9403fe0 	ldr	w0, [sp, #60]
 878:	d37ef400 	lsl	x0, x0, #2
 87c:	9100a3e1 	add	x1, sp, #0x28
 880:	b8606820 	ldr	w0, [x1, x0]
 884:	b9403be1 	ldr	w1, [sp, #56]
 888:	6b00003f 	cmp	w1, w0
 88c:	54000180 	b.eq	8bc <main+0x90>  // b.none
 890:	b9403fe0 	ldr	w0, [sp, #60]
 894:	d37ef400 	lsl	x0, x0, #2
 898:	9100a3e1 	add	x1, sp, #0x28
 89c:	b8606820 	ldr	w0, [x1, x0]
 8a0:	2a0003e2 	mov	w2, w0
 8a4:	b9403be1 	ldr	w1, [sp, #56]
 8a8:	90000000 	adrp	x0, 0 <__abi_tag-0x278>
 8ac:	91242000 	add	x0, x0, #0x908
 8b0:	97ffff74 	bl	680 <printf@plt>
 8b4:	b9403be0 	ldr	w0, [sp, #56]
 8b8:	1400000b 	b	8e4 <main+0xb8>
 8bc:	b9403fe0 	ldr	w0, [sp, #60]
 8c0:	11000400 	add	w0, w0, #0x1
 8c4:	b9003fe0 	str	w0, [sp, #60]
 8c8:	b9403fe0 	ldr	w0, [sp, #60]
 8cc:	7100081f 	cmp	w0, #0x2
 8d0:	54fffcc9 	b.ls	868 <main+0x3c>  // b.plast
 8d4:	90000000 	adrp	x0, 0 <__abi_tag-0x278>
 8d8:	91248000 	add	x0, x0, #0x920
 8dc:	97ffff65 	bl	670 <puts@plt>
 8e0:	52800000 	mov	w0, #0x0                   	// #0
 8e4:	a8c47bfd 	ldp	x29, x30, [sp], #64
 8e8:	d65f03c0 	ret

Disassembly of section .fini:

00000000000008ec <_fini>:
 8ec:	d503201f 	nop
 8f0:	a9bf7bfd 	stp	x29, x30, [sp, #-16]!
 8f4:	910003fd 	mov	x29, sp
 8f8:	a8c17bfd 	ldp	x29, x30, [sp], #16
 8fc:	d65f03c0 	ret
```

It is 227 lines long, including whitespace and annotations added by the
`objdump` program when formatting the disassembled program. Because of the
particular `objdump` command line option used, not shown is a certain amount of
overhead inserted by the GNU compiler and required by the Linux ELF file format
for every program, e.g. zero-padding for alignment of various sections in
memory.

The corresponding disassembled Go executable (not shown) is 134,668 lines long,
an increase of three (3) orders of magnitude over that for the C executable. It
is so very much larger because even the most trivial Go program must drag in
library dependencies and include large amounts of application-level support for
garbage collection, non-triveal data management for built-in constructs like
_slices_ (elastic sized views of fixed size arrays), and so on. In short, the
conveniences provided by languages like Go come at a very high price in terms
of memory utilization and run time performance.

C's semantics have very low intrinsic overhead over that of writing the
equivalent functionality directly in assembly language as can be seen easily
from the output of `objdump`[<sup>3</sup>](#risc). A _pointer_ in C is
literally just the address of a memory cell, which can be manipulated using
arithmetic operations just as one would when writing a program in assembly
language that needed to access successive values stored in RAM.

By contrast, Go's semantics result in even the simplest program incurring very
substantial overhead, whether or not that overhead is particularly useful for
the given program. Many real-world C programs have no need to depend on library
code that is much bigger than a few scanning, formatting and memory management
functions from `stdio.h` and `stdlib.h`. Even such programs would be orders of
magnitude smaller and more efficient at run time than their corresponding Go
programs. Programmers should be conscious of the trade-offs of the
"convenience" of not having to keep track of the length of arrays and not
having to remember to call `free` (or, in C++, `delete`) from time to time.

As a final example, consider this C++ version of the same logic:

```c++
#include <iostream>
#include <array>

void f(auto &x)
{
    for (unsigned i = 0; i < x.size(); ++i)
    {
        x[i] += 1;
    }
}

int main(int, char **)
{
    auto y = std::array<int, 3>{0, 1, 2};
    f(y);
    for (unsigned i; i < 3; ++i)
    {
        int v = i + 1;
        if (y[i] != v)
        {
            std::cout << "expected " << v << ", got " << y[i];
            return v;
        }
    }
    std::cout << "success" << std::endl;
    return 0;
}
```

While the C++ standard library class, `std::array`, is not nearly as flexible
as Go's slices, a great deal of the time all you really care about is a method
like `size()` (the equivalent of Go's `len()` function). Note that the C++
version's executable is only 1K larger than the C version (70K vs. 69K) and its
disassembly is only a little longer (333 lines vs. 227). In other words, C++
provides many of the advantages of high-level languages while introducing only
very marginal "bloat" compared to plain C.

## Conclusion

The bottom line is that higher-level features come at the cost of less
efficient programs when looked at from the point of view of the average amount
of memory consumed or speed of execution per line of code. Sometimes that
higher cost is worth it, but only for programs which actually need those
features and would otherwise simply have to re-invent some number of wheels
provided directly by a higher-level language. But often a given program's
requirements do not actually benefit from those higher level language features,
which then simply become unnecessary bloat. This is why every programmer should
know multiple languages and be willing and able to use an appropriate language
for a given task, as in the old adage, "if all you have is a hammer, every
problem looks like a nail."

As a corollary, being able to at least read the assembly language output of
tools like `objdump` can be a very useful skill for debugging and for learning,
not just reverse engineering (which is rarely worth the effort in the real
world). Compilers like `cc` / `gcc` / `g++` allow you to mix and match C/C++
and assembly source code in a single build easily and naturally, opening up
opportunities for optimizations and access to hardware level features that are
inaccessible in higher level languages and interpreted run time
platforms[<sup>4</sup>](#no-really).

---

### Notes

<a id="risc"><sup>3</sup></a> Ok, "easily" might be an exaggeration. ARM is a
RISC (Reduced Instruction Set Computer) architecture. The disassembled
executable for a CISC (Complex Instruction Set Computer) architecture CPU would
typically be somewhat smaller, i.e. there would be fewer individual machine
language instructions in the executable file for a C program compiled from the
same source code. The individual instructions for a CISC architecture CPU
provide more complex behavior compared to RISC, resulting in shorter, more
readable assembly language programs. But since a single CISC operation does
more on average compared to a single RISC operation, the average CISC
instruction also takes more time to execute. And since each RISC instruction
represents a smaller unit of logic, compilers have more opportunities to apply
optimizations when compiling for RISC. This is why, after decades of dominance
by CISC, RISC has become the standard style for microprocessor design and why,
for example, mobile devices with slower CPU clocks and smaller caches can
perform as well or better, while consuming less power and generating less heat,
compared to the desktop CPU's of yore (i.e. not so many years ago).

<a id="no-really"><sup>4</sup></a> Back in the day, the author of this document
was responsible for entire mission-critical applications written entirely in
680x0 and 80x86 assembly languages in scientific, aerospace and defense domains
where bugs could cost lives or cause billions of dollars' worth of damage. This
was during the same time in which he was a member of teams doing R&D in the
first AI boom of the 1980's, developing machine learning algorithms and expert
systems in various dialects of Lisp, Smalltalk and Prolog. No, really.
