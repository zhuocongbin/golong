// `for` is Go's only looping construct. Here are
// three basic types of `for` loops.
//“for”是Go的唯一循环结构。 下面是“for”循环的三种基本类型。
package main

import "fmt"

func main() {

	// The most basic type, with a single condition.
	// 最基本的类型，只有一个条件。
	i := 1
	for i <= 3 {
		fmt.Println(i)
		i = i + 1
	}

	// A classic initial/condition/after `for` loop.
	// 典型的初始/条件/后循环。
	for j := 7; j <= 9; j++ {
		fmt.Println(j)
	}

	// `for` without a condition will loop repeatedly
	// until you `break` out of the loop or `return` from
	// the enclosing function.
	//“for”没有一个条件会反复循环，直到你从封闭函数中“跳出”循环或“返回”。
	for {
		fmt.Println("loop")
		break
	}

	// You can also `continue` to the next iteration of the loop.
	// 您还可以“继续”到循环的下一个迭代。
	for n := 0; n <= 5; n++ {
		if n%2 == 0 {
			continue
		}
		fmt.Println(n)
	}
}
