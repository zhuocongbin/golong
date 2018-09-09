// In Go, an _array_ is a numbered sequence of elements of a specific length.
// 在Go中，_array_是特定长度的元素的编号序列。

package main

import "fmt"

func main() {

	// Here we create an array `a` that will hold exactly
	// 5 `int`s. The type of elements and length are both
	// part of the array's type. By default an array is
	// zero-valued, which for `int`s means `0`s.
	//在这里，我们创建一个数组' a '，它将恰好保持5 ' int '。
	//元素的类型和长度都是数组类型的一部分。
	//默认情况下，数组是零值的，这对于“int”的意思是0。
	var a [5]int
	fmt.Println("emp:", a)

	// We can set a value at an index using the
	// `array[index] = value` syntax, and get a value with `array[index]`.
	// 我们可以使用' array[index] = value '语法在索引中设置一个值，并使用' array[index] '获取一个值。
	a[4] = 100
	fmt.Println("set:", a)
	fmt.Println("get:", a[4])

	// The builtin `len` returns the length of an array.
	fmt.Println("len:", len(a))

	// Use this syntax to declare and initialize an array in one line.
	// 使用此语法在一行中声明和初始化一个数组。
	b := [5]int{1, 2, 3, 4, 5}
	fmt.Println("dcl:", b)

	// Array types are one-dimensional,
	// but you can compose types to build multi-dimensional data structures.
	// 数组类型是一维的，但是您可以组合类型来构建多维数据结构。
	var twoD [2][3]int
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoD)
}
