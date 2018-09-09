// A _goroutine_ is a lightweight thread of execution.
// _goroutine_是一个轻量级执行线程。

package main

import "fmt"

func f(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}

func main() {

	// Suppose we have a function call `f(s)`.
	// Here's how we'd call that in the usual way, running it synchronously.
	// 假设我们有一个函数叫f(s) 我们用通常的方式调用它，同步运行它。
	f("direct")

	// To invoke this function in a goroutine, use `go f(s)`.
	//  This new goroutine will execute concurrently with the calling one.
	// 要在goroutine中调用此函数，请使用“go f(s)”。 这个新的goroutine将与调用one同时执行。
	go f("goroutine")

	// You can also start a goroutine for an anonymous function call.
	// 您还可以为一个匿名函数调用启动goroutine。
	go func(msg string) {
		fmt.Println(msg)
	}("going")

	// Our two function calls are running asynchronously in separate goroutines now, so execution falls through to here.
	// This `Scanln` requires we press a key before the program exits.
	// 我们的两个函数调用现在在单独的goroutines中异步运行，所以执行会在这里执行。
	// 这个“Scanln”要求我们在程序退出前按下一个键。

	fmt.Scanln()
	fmt.Println("done")
}
