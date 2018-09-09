// Go supports [_anonymous functions_](http://en.wikipedia.org/wiki/Anonymous_function),
// which can form <a href="http://en.wikipedia.org/wiki/Closure_(computer_science)">
// <em>closures 闭包</em></a>.
// Anonymous functions are useful when you want to define a function inline without having to name it.
// 当您想要在不需要命名函数的情况下定义一个函数时，匿名函数是很有用的。

package main

import "fmt"

// This function `intSeq` returns another function, which we define anonymously in the body of `intSeq`.
// 这个函数“intSeq”返回另一个函数，我们在“intSeq”的主体中匿名定义。
// The returned function _closes over_ the variable `i` to form a closure.
// 返回的函数_关闭了变量' i '以形成一个闭包。
//
func intSeq() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

func main() {

	// We call `intSeq`, assigning the result (a function) to `nextInt`.
	// 我们调用“intSeq”，将结果(一个函数)赋给“nextInt”。
	// This function value captures its own `i` value, which will be updated each timewe call `nextInt`.
	// 这个函数值捕捉它自己的“i”值，每次我们调用“nextInt”时，它都会被更新。
	nextInt := intSeq()

	// See the effect of the closure by calling `nextInt`a few times.
	// 通过调用“nextInt”来查看关闭的效果。
	fmt.Println(nextInt())
	fmt.Println(nextInt())
	fmt.Println(nextInt())

	// To confirm that the state is unique to that particular function, create and test a new one.
	// 要确认状态对于特定的函数是唯一的，创建并测试一个新的函数。
	newInts := intSeq()
	fmt.Println(newInts())
}
