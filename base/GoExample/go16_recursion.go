// Go supports
// <a href="http://en.wikipedia.org/wiki/Recursion_(computer_science)">
// <em>recursive functions 递归函数</em></a>.
// Here's a classic factorial example.
// 这是一个经典的阶乘例子。

package main

import "fmt"

// This `fact` function calls itself until it reaches the base case of `fact(0)`.
// 这个“事实”函数调用它自己，直到它到达“fact(0)”的基本情况。
func fact(n int) int {
	if n == 0 {
		return 1
	}
	return n * fact(n-1)
}

func main() {
	fmt.Println(fact(7))
}
