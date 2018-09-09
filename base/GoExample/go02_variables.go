//在Go中，_variables_被编译器显式地声明和使用，例如检查函数调用的类型正确性。

package main

import "fmt"

func main() {

	// “var”声明一个或多个变量。
	var a = "initial"
	fmt.Println(a)

	// 你可以同时声明多个变量。
	var b, c int = 1, 2
	fmt.Println(b, c)

	// Go 将推断初始化变量的类型。
	var d = true
	fmt.Println(d)

	// 没有相应初始化声明的变量是_zero-valued_。
	//例如，“int”的零值为“0”。
	var e int
	fmt.Println(e)

	//“:=”语法是对变量进行声明和初始化的缩写，
	// 例如，在本例中为“var f string = "short"。
	f := "short"
	fmt.Println(f)
}
