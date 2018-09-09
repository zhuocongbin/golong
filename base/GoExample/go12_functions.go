// _Functions_ are central in Go. We'll learn about functions with a few different examples.
// _Functions_是Go的中心。 我们将用几个不同的例子来学习函数。

package main

import "fmt"

// Here's a function that takes two `int`s and returns their sum as an `int`.
// 这里有一个函数，它取两个' int '并返回它们的和作为' int '。
func plus(a int, b int) int {

	// Go requires explicit returns, i.e. it won't automatically return the value of the last expression.
	// Go需要显式的返回，即它不会自动返回最后一个表达式的值。
	//
	return a + b
}

// When you have multiple consecutive parameters of the same type,
// you may omit the type name for the like-typed parameters up to the final parameter that declares the type.
// 当你有相同类型的多个连续参数时，您可以将like类型参数的类型名称省略到声明类型的最终参数中。
func plusPlus(a, b, c int) int {
	return a + b + c
}

func main() {

	// Call a function just as you'd expect, with `name(args)`.
	// 调用一个函数，正如您所期望的那样，使用`name(args)`。
	res := plus(1, 2)
	fmt.Println("1+2 =", res)

	res = plusPlus(1, 2, 3)
	fmt.Println("1+2+3 =", res)
}
