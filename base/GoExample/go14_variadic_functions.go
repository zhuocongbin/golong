// [_Variadic functions_](http://en.wikipedia.org/wiki/Variadic_function)
// can be called with any number of trailing arguments.
// 可以使用任意数量的尾随参数来调用。
// For example, `fmt.Println` is a common variadic function.
// 例如,`fmt.Println` 是一个常见的变量函数。

package main

import "fmt"

// Here's a function that will take an arbitrary number of `int`s as arguments.
// 这里有一个函数可以取任意数量的int作为参数。
func sum(nums ...int) {
	fmt.Print(nums, " ")
	total := 0
	for _, num := range nums {
		total += num
	}
	fmt.Println(total)
}

func main() {

	// Variadic functions can be called in the usual way with individual arguments.
	// 变量函数可以用通常的方式和单个参数调用。
	sum(1, 2)
	sum(1, 2, 3)

	// If you already have multiple args in a slice,
	// apply them to a variadic function using
	// `func(slice...)` like this.
	nums := []int{1, 2, 3, 4}
	sum(nums...)
}
