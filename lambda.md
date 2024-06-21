# What Your CS Professor Might Have Failed to Tell You

## All Programming Languages Are Modeled on the Lambda Calculus

To understand what all programming languages have in common, first learn to
think this way<sup>[1](#church)</sup>:

> _let y = &lambda;x.+ x 1 (2_
>
> _y - 3 = 0_

It will then not take long for it to feel entirely natural to think this
way<sup>[2](#lisp)</sup>:

```
> (let ((y ((lambda (x) (+ x 1)) 2)))
    (equal? (- y 3) 0))
#t
```

Which then makes it easy to think any of these ways:

```
irb(main):002:0> y = ->(x) { x + 1 }.(2)
=> 3
irb(main):003:0> y - 3 == 0
=> true
```

```
let y = ((x) => { return x + 1 }) (2)
console.log(y - 3 == 0)
```

```
y := func(x int) int { return x + 1 }(2)
fmt.Println(y-3 == 0)
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

```
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

```
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
            return v;
        }
    }
    return 0;
}
```

That program defines and invokes the function `f`, which takes an array of
integers and adds 1 to each element, using C's standard idiom for arrays
involving "pointers" and "pointer arithmetic". The programm then checks the
contents of the array after `f` returns, exiting with a status code indicating
failure corresponding to any element whose value is incorrect or with a status
code indicating success if all elements pass validation.

Now consider the following Go program:

```
package main

import "os"

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
			os.Exit(v)
		}
	}
}
```

The preceding two programs implement as close to identical logic, using as
close to identical data types as the syntax and semantics of the two languages
allow. Both are barely a notch above "hello world" in features and
functionality. But note that while the executable emitted by the GNU C compiler
targeting Linux on ARM64 is 69K bytes, the corresponding Go executable is 1.6M
bytes. This ratio of a Go executable being several orders of magnitude larger
than a functionally equivalent C program is typical, not just for Go but for
any "higher level" compiled language. (Calculations of the relative size and
efficiency of programs executed by interpreted languages like JavaScript or
Ruby are far more complicated and nuanced affairs, but suffice it to say that
such languages cannot by any reasonable measure be regarded as significantly
more efficient than compiled high-level languages like C# or Go. Such
calculations for virtual-machine based compiled languages like Java -- which
are first compiled, then interpreted at run time -- are even more fraught, but
all the same realities are revelaed when assessing their memory utilization and
run time performace.)

Here is the output of `objdump -d` for the preceding program written in C:

```

example:     file format elf64-littleaarch64


Disassembly of section .init:

0000000000000580 <_init>:
 580:	d503201f 	nop
 584:	a9bf7bfd 	stp	x29, x30, [sp, #-16]!
 588:	910003fd 	mov	x29, sp
 58c:	9400002a 	bl	634 <call_weak_fn>
 590:	a8c17bfd 	ldp	x29, x30, [sp], #16
 594:	d65f03c0 	ret

Disassembly of section .plt:

00000000000005a0 <.plt>:
 5a0:	a9bf7bf0 	stp	x16, x30, [sp, #-16]!
 5a4:	f00000f0 	adrp	x16, 1f000 <__FRAME_END__+0x1e6d8>
 5a8:	f947fe11 	ldr	x17, [x16, #4088]
 5ac:	913fe210 	add	x16, x16, #0xff8
 5b0:	d61f0220 	br	x17
 5b4:	d503201f 	nop
 5b8:	d503201f 	nop
 5bc:	d503201f 	nop

00000000000005c0 <__libc_start_main@plt>:
 5c0:	90000110 	adrp	x16, 20000 <__libc_start_main@GLIBC_2.34>
 5c4:	f9400211 	ldr	x17, [x16]
 5c8:	91000210 	add	x16, x16, #0x0
 5cc:	d61f0220 	br	x17

00000000000005d0 <__cxa_finalize@plt>:
 5d0:	90000110 	adrp	x16, 20000 <__libc_start_main@GLIBC_2.34>
 5d4:	f9400611 	ldr	x17, [x16, #8]
 5d8:	91002210 	add	x16, x16, #0x8
 5dc:	d61f0220 	br	x17

00000000000005e0 <__gmon_start__@plt>:
 5e0:	90000110 	adrp	x16, 20000 <__libc_start_main@GLIBC_2.34>
 5e4:	f9400a11 	ldr	x17, [x16, #16]
 5e8:	91004210 	add	x16, x16, #0x10
 5ec:	d61f0220 	br	x17

00000000000005f0 <abort@plt>:
 5f0:	90000110 	adrp	x16, 20000 <__libc_start_main@GLIBC_2.34>
 5f4:	f9400e11 	ldr	x17, [x16, #24]
 5f8:	91006210 	add	x16, x16, #0x18
 5fc:	d61f0220 	br	x17

Disassembly of section .text:

0000000000000600 <_start>:
 600:	d503201f 	nop
 604:	d280001d 	mov	x29, #0x0                   	// #0
 608:	d280001e 	mov	x30, #0x0                   	// #0
 60c:	aa0003e5 	mov	x5, x0
 610:	f94003e1 	ldr	x1, [sp]
 614:	910023e2 	add	x2, sp, #0x8
 618:	910003e6 	mov	x6, sp
 61c:	f00000e0 	adrp	x0, 1f000 <__FRAME_END__+0x1e6d8>
 620:	f947ec00 	ldr	x0, [x0, #4056]
 624:	d2800003 	mov	x3, #0x0                   	// #0
 628:	d2800004 	mov	x4, #0x0                   	// #0
 62c:	97ffffe5 	bl	5c0 <__libc_start_main@plt>
 630:	97fffff0 	bl	5f0 <abort@plt>

0000000000000634 <call_weak_fn>:
 634:	f00000e0 	adrp	x0, 1f000 <__FRAME_END__+0x1e6d8>
 638:	f947e800 	ldr	x0, [x0, #4048]
 63c:	b4000040 	cbz	x0, 644 <call_weak_fn+0x10>
 640:	17ffffe8 	b	5e0 <__gmon_start__@plt>
 644:	d65f03c0 	ret
 648:	d503201f 	nop
 64c:	d503201f 	nop

0000000000000650 <deregister_tm_clones>:
 650:	90000100 	adrp	x0, 20000 <__libc_start_main@GLIBC_2.34>
 654:	9100c000 	add	x0, x0, #0x30
 658:	90000101 	adrp	x1, 20000 <__libc_start_main@GLIBC_2.34>
 65c:	9100c021 	add	x1, x1, #0x30
 660:	eb00003f 	cmp	x1, x0
 664:	540000c0 	b.eq	67c <deregister_tm_clones+0x2c>  // b.none
 668:	f00000e1 	adrp	x1, 1f000 <__FRAME_END__+0x1e6d8>
 66c:	f947e021 	ldr	x1, [x1, #4032]
 670:	b4000061 	cbz	x1, 67c <deregister_tm_clones+0x2c>
 674:	aa0103f0 	mov	x16, x1
 678:	d61f0200 	br	x16
 67c:	d65f03c0 	ret

0000000000000680 <register_tm_clones>:
 680:	90000100 	adrp	x0, 20000 <__libc_start_main@GLIBC_2.34>
 684:	9100c000 	add	x0, x0, #0x30
 688:	90000101 	adrp	x1, 20000 <__libc_start_main@GLIBC_2.34>
 68c:	9100c021 	add	x1, x1, #0x30
 690:	cb000021 	sub	x1, x1, x0
 694:	d37ffc22 	lsr	x2, x1, #63
 698:	8b810c41 	add	x1, x2, x1, asr #3
 69c:	9341fc21 	asr	x1, x1, #1
 6a0:	b40000c1 	cbz	x1, 6b8 <register_tm_clones+0x38>
 6a4:	f00000e2 	adrp	x2, 1f000 <__FRAME_END__+0x1e6d8>
 6a8:	f947f042 	ldr	x2, [x2, #4064]
 6ac:	b4000062 	cbz	x2, 6b8 <register_tm_clones+0x38>
 6b0:	aa0203f0 	mov	x16, x2
 6b4:	d61f0200 	br	x16
 6b8:	d65f03c0 	ret
 6bc:	d503201f 	nop

00000000000006c0 <__do_global_dtors_aux>:
 6c0:	a9be7bfd 	stp	x29, x30, [sp, #-32]!
 6c4:	910003fd 	mov	x29, sp
 6c8:	f9000bf3 	str	x19, [sp, #16]
 6cc:	90000113 	adrp	x19, 20000 <__libc_start_main@GLIBC_2.34>
 6d0:	3940c260 	ldrb	w0, [x19, #48]
 6d4:	35000140 	cbnz	w0, 6fc <__do_global_dtors_aux+0x3c>
 6d8:	f00000e0 	adrp	x0, 1f000 <__FRAME_END__+0x1e6d8>
 6dc:	f947e400 	ldr	x0, [x0, #4040]
 6e0:	b4000080 	cbz	x0, 6f0 <__do_global_dtors_aux+0x30>
 6e4:	90000100 	adrp	x0, 20000 <__libc_start_main@GLIBC_2.34>
 6e8:	f9401400 	ldr	x0, [x0, #40]
 6ec:	97ffffb9 	bl	5d0 <__cxa_finalize@plt>
 6f0:	97ffffd8 	bl	650 <deregister_tm_clones>
 6f4:	52800020 	mov	w0, #0x1                   	// #1
 6f8:	3900c260 	strb	w0, [x19, #48]
 6fc:	f9400bf3 	ldr	x19, [sp, #16]
 700:	a8c27bfd 	ldp	x29, x30, [sp], #32
 704:	d65f03c0 	ret
 708:	d503201f 	nop
 70c:	d503201f 	nop

0000000000000710 <frame_dummy>:
 710:	17ffffdc 	b	680 <register_tm_clones>

0000000000000714 <f>:
 714:	d10083ff 	sub	sp, sp, #0x20
 718:	f90007e0 	str	x0, [sp, #8]
 71c:	b90007e1 	str	w1, [sp, #4]
 720:	b9001fff 	str	wzr, [sp, #28]
 724:	1400000a 	b	74c <f+0x38>
 728:	f94007e0 	ldr	x0, [sp, #8]
 72c:	91001001 	add	x1, x0, #0x4
 730:	f90007e1 	str	x1, [sp, #8]
 734:	b9400001 	ldr	w1, [x0]
 738:	11000421 	add	w1, w1, #0x1
 73c:	b9000001 	str	w1, [x0]
 740:	b9401fe0 	ldr	w0, [sp, #28]
 744:	11000400 	add	w0, w0, #0x1
 748:	b9001fe0 	str	w0, [sp, #28]
 74c:	b9401fe1 	ldr	w1, [sp, #28]
 750:	b94007e0 	ldr	w0, [sp, #4]
 754:	6b00003f 	cmp	w1, w0
 758:	54fffe83 	b.cc	728 <f+0x14>  // b.lo, b.ul, b.last
 75c:	d503201f 	nop
 760:	d503201f 	nop
 764:	910083ff 	add	sp, sp, #0x20
 768:	d65f03c0 	ret

000000000000076c <main>:
 76c:	a9bc7bfd 	stp	x29, x30, [sp, #-64]!
 770:	910003fd 	mov	x29, sp
 774:	b9001fe0 	str	w0, [sp, #28]
 778:	f9000be1 	str	x1, [sp, #16]
 77c:	90000000 	adrp	x0, 0 <__abi_tag-0x278>
 780:	91206001 	add	x1, x0, #0x818
 784:	9100a3e0 	add	x0, sp, #0x28
 788:	f9400022 	ldr	x2, [x1]
 78c:	f9000002 	str	x2, [x0]
 790:	b9400821 	ldr	w1, [x1, #8]
 794:	b9000801 	str	w1, [x0, #8]
 798:	9100a3e0 	add	x0, sp, #0x28
 79c:	52800061 	mov	w1, #0x3                   	// #3
 7a0:	97ffffdd 	bl	714 <f>
 7a4:	14000010 	b	7e4 <main+0x78>
 7a8:	b9403fe0 	ldr	w0, [sp, #60]
 7ac:	11000400 	add	w0, w0, #0x1
 7b0:	b9003be0 	str	w0, [sp, #56]
 7b4:	b9403fe0 	ldr	w0, [sp, #60]
 7b8:	d37ef400 	lsl	x0, x0, #2
 7bc:	9100a3e1 	add	x1, sp, #0x28
 7c0:	b8606820 	ldr	w0, [x1, x0]
 7c4:	b9403be1 	ldr	w1, [sp, #56]
 7c8:	6b00003f 	cmp	w1, w0
 7cc:	54000060 	b.eq	7d8 <main+0x6c>  // b.none
 7d0:	b9403be0 	ldr	w0, [sp, #56]
 7d4:	14000008 	b	7f4 <main+0x88>
 7d8:	b9403fe0 	ldr	w0, [sp, #60]
 7dc:	11000400 	add	w0, w0, #0x1
 7e0:	b9003fe0 	str	w0, [sp, #60]
 7e4:	b9403fe0 	ldr	w0, [sp, #60]
 7e8:	7100081f 	cmp	w0, #0x2
 7ec:	54fffde9 	b.ls	7a8 <main+0x3c>  // b.plast
 7f0:	52800000 	mov	w0, #0x0                   	// #0
 7f4:	a8c47bfd 	ldp	x29, x30, [sp], #64
 7f8:	d65f03c0 	ret

Disassembly of section .fini:

00000000000007fc <_fini>:
 7fc:	d503201f 	nop
 800:	a9bf7bfd 	stp	x29, x30, [sp, #-16]!
 804:	910003fd 	mov	x29, sp
 808:	a8c17bfd 	ldp	x29, x30, [sp], #16
 80c:	d65f03c0 	ret
```

It is 203 lines long, including whitespace and annotations added by the
`objdump` program when formatting the disassembled program. Because of the
particular `objdump` command line options used, not shown is a certain amount
of overhead inserted by the GNU compiler and required by the Linux ELF file
format for every program, e.g. zero-padding for alignment of various sections
in memory.

The corresponding disassembled Go executable (not shown) is 104,775 lines long,
an increase of three (3) orders of magnitude over that for the C
executable. It is so very much larger because even the most trivial Go program
must drag in library dependencies and include large amounts of
application-level support for garbage collection, non-triveal data management
for built-in constructs like _slices_ (elastic sized views of fixed size
arrays), and so on. In short, the conveniences provided by languages like Go
come at a very high price in terms of memory utilization and run time
performance.

A careful reader might claim that the comparison between the two programs is
not fair, because (unlike any real world example), the given C program does not
depend on any standard library functions (note the lack of any
`#include` statements) while the Go program depends on the `os` package (in
order to invoke `os.Exit`). Similarly, those familiar with the semantics of
both C and Go will point out how much richer the functionality of Go's slices
are in comparison to arrays in C. But that is exactly the point of these
examples.

C's semantics have very low intrinsic overhead over that of writing the
equivalent functionality directly in assembly language as can be seen easily
from the output of `objdump`[<sup>3</sup>](#risc). A _pointer_ in C is
literally just the address of a memory cell, which can be manipulated using
arithmetic operations just as one would when writing a program in assembly
language that needed to access data stored in successive RAM locations.

By contrast, Go's semantics result in even the simplest program incurring very
substantial overhead, whether or not that overhead is particularly useful for
the given program. Many real-world C programs have no need to depend on library
code that is much bigger than a few scanning, formatting and memory management
functions from `stdio.h` and `stdlib.h`. Even such programs would be orders of
magnitude smaller and more efficient at run time compared to their corresponding Go
programs. Programmers should be conscious of the trade-offs of the
"convenience" of not having to keep track of the length of arrays and not
having to remember to call `free` (or, in C++, `delete`) from time to time.

## Conclusion

The bottom line is that higher-level features come at the cost of less
efficient programs when looked at from the point of view of the average amount
of memory or speed of execution per line of code. Sometimes that higher cost is
worth it, but only for programs which actually need those features and would
otherwise simply have to re-invent some number of wheels provided directly by a
higher-level language. But often a given program's requirements do not actually
benefit from those higher level language features, which then simply become
unnecessary bloat. This is why every programmer should know multiple languages
and be willing and able to use an appropriate language for a given task, as in
the old adage, "if all you have is a hammer, every problem looks like a nail."

As a corollary, being able to at least read the assembly language output of
tools like `objdump` can be a very useful skill for debugging and for learning,
not just reverse engineering (which is rarely worth the effort in the real
world). Compilers like `cc` allow you to mix and match C/C++ and assembly
source code in a single build easily and naturally, opening up opportunities
for optimizations and access to hardware level features that are inaccessible
in higher level languages and interpreted run time platforms[<sup>4</sup>](#no-really).

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
represents a smaller "unit" of logic, compilers have more opportunities to
apply optimizations when compiling for RISC. This is why, after decades of
dominance by CISC, RISC has become the standard style for microprocessor design
and why, for example, mobile devices with slower CPU clocks and smaller caches
can perform as well or better, while consuming less power and generating less
heat, compared to the desktop CPU's of yore (i.e. not so many years ago).

<a id="no-really"><sup>4</sup></a> Back in the day, the author of this document
wrote entire mission-critical applications in 680x0 and 80x86 assembly
languages in scientific, aerospace and defense domains where bugs could cost
lives or cause billions of dollars' worth of damage. This was during the same
time in which he was a member of teams doing R&D in the first AI boom of the
1980's, developing machine learning algorithms and expert systems in various
dialects of Lisp, Smalltalk and Prolog. He has continued to practice what
is preached here throughout the intervening decades.
