// Go有各种各样的值类型，包括字符串、整数、浮点数、布尔值等等。
//这里有几个基本的例子。

package main

import "fmt"

func main() {

	// Strings, which can be added together with `+`.
	//多个字符串，可以用“+”添加在一起。
	fmt.Println("go" + "lang")

	// Integers and floats.整数和浮点数。
	fmt.Println("1+1 =", 1+1)
	fmt.Println("7.0/3.0 =", 7.0/3.0)

	// Booleans, with boolean operators as you'd expect.
	//使用布尔运算符。
	fmt.Println(true && false)
	fmt.Println(true || false)
	fmt.Println(!true)
}
