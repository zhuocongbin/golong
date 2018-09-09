// Go supports _constants_ of character, string, boolean,
// and numeric values.
// Go支持字符、字符串、布尔值和数值的常量。

package main

import "fmt"
import "math"

// `const` declares a constant value.
//“const”声明一个常量值。
const s string = "constant"

func main() {
	fmt.Println(s)

	// A `const` statement can appear anywhere a `var`statement can.
	// const的语句可以在任何地方出现“var”语句。
	const n = 500000000

	// Constant expressions perform arithmetic with arbitrary precision.
	// 常数表达式的运算具有任意精度。
	const d = 3e20 / n
	fmt.Println(d)

	// A numeric constant has no type until it's given one, such as by an explicit cast.
	//一个数值常量在给定一个类型之前是没有类型的，比如一个显式的cast。
	fmt.Println(int64(d))

	// A number can be given a type by using it in a
	// context that requires one, such as a variable
	// assignment or function call. For example, here
	// `math.Sin` expects a `float64`.
	//一个数字可以在需要一个的上下文中使用它，比如变量赋值或函数调用。
	// 例如,这里的`math.Sin`预计"float64"。
	fmt.Println(math.Sin(n))
}
