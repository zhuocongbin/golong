// _Switch statements_ express conditionals across many branches.
// 在许多分支中有_Switch语句表达式

package main

import "fmt"
import "time"

func main() {

	// Here's a basic `switch`.
	// 这是一个`switch`的基本
	i := 2
	fmt.Print("Write ", i, " as ")
	switch i {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 3:
		fmt.Println("three")
	}

	// You can use commas to separate multiple expressions
	// in the same `case` statement. We use the optional
	// `default` case in this example as well.
	//您可以在同一个“case”语句中使用逗号分隔多个表达式。
	// 在这个示例中，我们还使用了可选的“默认”案例。
	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("It's the weekend")
	default:
		fmt.Println("It's a weekday")
	}

	// `switch` without an expression is an alternate way
	// to express if/else logic. Here we also show how the
	// `case` expressions can be non-constants.
	//“switch”没有表达式是表示if/else逻辑的另一种方式。
	// 这里我们还展示了“case”表达式可以是非常量。
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("It's before noon")
	default:
		fmt.Println("It's after noon")
	}

	// A type `switch` compares types instead of values.  You
	// can use this to discover the type of an interface
	// value.  In this example, the variable `t` will have the
	// type corresponding to its clause.
	// A type ' switch '比较类型而不是值。
	//您可以使用它来发现接口值的类型。
	//在本例中，变量' t '将具有与它的子句相对应的类型。
	whatAmI := func(i interface{}) {
		switch t := i.(type) {
		case bool:
			fmt.Println("I'm a bool")
		case int:
			fmt.Println("I'm an int")
		default:
			fmt.Printf("Don't know type %T\n", t)
		}
	}
	whatAmI(true)
	whatAmI(1)
	whatAmI("hey")
}
