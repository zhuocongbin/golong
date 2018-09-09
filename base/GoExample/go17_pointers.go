// Go supports <em><a href="http://en.wikipedia.org/wiki/Pointer_(computer_programming)"
// >pointers</a></em>,
// allowing you to pass references to values and records within your program.
// 允许您在程序中传递对值和记录的引用。

package main

import "fmt"

// We'll show how pointers work in contrast to values with 2 functions: `zeroval` and `zeroptr`.
// 我们将展示指针如何与具有2个功能的值形成对比:“zeroval”和“zeroptr”。
// `zeroval` has an `int` parameter, so arguments will be passed to it by value.
// “zeroval”有一个“int”参数，因此参数将按值传递给它。
// `zeroval` will get a copy of `ival` distinct from the one in the calling function.
// “zeroval”将获得与调用函数中的“ival”不同的副本。
func zeroval(ival int) {
	ival = 0
}

// `zeroptr` in contrast has an `*int` parameter, meaning that it takes an `int` pointer.
// “zeroptr”有一个“*int”参数，这意味着它需要一个“int”指针。
// The `*iptr` code in the function body then _dereferences_ the pointer from its memory address to the current value at that address.
// 函数体中的“*iptr”代码然后将指针从其内存地址转换为该地址的当前值。
// Assigning a value to a dereferenced pointer changes the value at the referenced address.
// 将一个值赋值给一个取消引用的指针会改变引用地址的值。
func zeroptr(iptr *int) {
	*iptr = 0
}

func main() {
	i := 1
	fmt.Println("initial:", i)

	zeroval(i)
	fmt.Println("zeroval:", i)

	// The `&i` syntax gives the memory address of `i`,i.e. a pointer to `i`.
	// “&i”语法给出了“i”的内存地址，即“i”。 `i`的指针。
	zeroptr(&i)
	fmt.Println("zeroptr:", i)

	// Pointers can be printed too.
	fmt.Println("pointer:", &i)
}
