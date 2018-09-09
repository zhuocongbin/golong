// Go has built-in support for _multiple return values_.
// This feature is used often in idiomatic Go, for example
// to return both result and error values from a function.
//GO已经内置支持多个返回值。 这个特性通常在惯用的Go中使用，例如从函数返回结果和错误值。
package main

import "fmt"

// The `(int, int)` in this function signature shows that the function returns 2 `int`s.
// 在这个函数签名中，“(int, int)”显示函数返回2 个int。
func vals() (int, int) {
	return 3, 7
}

func main() {

	// Here we use the 2 different return values from the call with _multiple assignment_.
	// 这里我们使用来自多个赋值的调用的两个不同的返回值。
	a, b := vals()
	fmt.Println(a)
	fmt.Println(b)

	// If you only want a subset of the returned values,use the blank identifier `_`.
	// 如果您只需要返回值的一个子集，则使用空白标识符“_”。
	_, c := vals()
	fmt.Println(c)
}
