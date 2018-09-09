// Branching with `if` and `else` in Go is straight-forward.
// 在go中“if”和“else”的分支是直接的。

package main

import "fmt"

func main() {

	// Here's a basic example.
	// 这是基本例子
	if 7%2 == 0 {
		fmt.Println("7 is even")
	} else {
		fmt.Println("7 is odd")
	}

	// You can have an `if` statement without an else.
	// 你可以在没有其他条件的情况下使用if语句。
	if 8%4 == 0 {
		fmt.Println("8 is divisible by 4")
	}

	// A statement can precede conditionals; any variables
	// declared in this statement are available in all branches.
	//声明可以先于条件; 此语句中声明的任何变量都可以在所有分支中使用。
	if num := 9; num < 0 {
		fmt.Println(num, "is negative")
	} else if num < 10 {
		fmt.Println(num, "has 1 digit")
	} else {
		fmt.Println(num, "has multiple digits")
	}
}

// Note that you don't need parentheses around conditions in Go, but that the braces are required.
//注意，你不需要括号周围的条件， 但这需要大括号。
