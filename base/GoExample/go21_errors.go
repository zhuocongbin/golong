// In Go it's idiomatic to communicate errors via an explicit, separate return value.
// go通过显式的、单独的返回值来传达错误是惯用的。
//  This contrasts with the exceptions used in languages like Java and Ruby and
//  the overloaded single result / error value sometimes used in C.
// 这与在Java和Ruby等语言中使用的异常和有时在C中使用的重载的单个结果/错误值形成了对比。
//  Go's approach makes it easy to see which functions return errors and to handle them using the
//  same language constructs employed for any other,non-error tasks.
//  Go的方法使我们可以很容易地看到哪些函数返回错误，并使用用于其他非错误任务的相同语言结构来处理它们。

package main

import "errors"
import "fmt"

// By convention, errors are the last return value and have type `error`, a built-in interface.
// 根据惯例，错误是最后的返回值，并且有类型“error”，一个内置的接口。
func f1(arg int) (int, error) {
	if arg == 42 {

		// `errors.New` constructs a basic `error` value with the given error message.
		// `errors.New`用给定的错误消息构造一个基本的`error`值。
		return -1, errors.New("can't work with 42")

	}

	// A `nil` value in the error position indicates that there was no error.
	// 错误位置的“nil”值表示没有错误。
	return arg + 3, nil
}

// It's possible to use custom types as `error`s by implementing the `Error()` method on them.
// 可以使用自定义类型作为“错误”，在它们上实现“error()”方法。
// Here's a variant on the example above that uses a custom type to explicitly represent an argument error.
// 在上面的例子中有一个变体，它使用自定义类型显式地表示一个参数错误。
type argError struct {
	arg  int
	prob string
}

func (e *argError) Error() string {
	return fmt.Sprintf("%d - %s", e.arg, e.prob)
}

func f2(arg int) (int, error) {
	if arg == 42 {

		// In this case we use `&argError` syntax to build a new struct, supplying values for the two fields `arg` and `prob`.
		// 在本例中，我们使用“&argError”语法构建一个新的结构，为两个字段的arg和prob提供值。
		//
		return -1, &argError{arg, "can't work with it"}
	}
	return arg + 3, nil
}

func main() {

	// The two loops below test out each of our error-returning functions.
	// 下面的两个循环测试了每个返回错误的函数。
	// Note that the use of an inline error check on the `if` line is a common idiom in Go code.
	// 注意，在“if”行中使用内联错误检查是Go代码中常见的习惯用法。
	for _, i := range []int{7, 42} {
		if r, e := f1(i); e != nil {
			fmt.Println("f1 failed:", e)
		} else {
			fmt.Println("f1 worked:", r)
		}
	}
	for _, i := range []int{7, 42} {
		if r, e := f2(i); e != nil {
			fmt.Println("f2 failed:", e)
		} else {
			fmt.Println("f2 worked:", r)
		}
	}

	// If you want to programmatically use the data in a custom error, you'll need to get the error as an
	// instance of the custom error type via type assertion.
	// 如果希望以编程方式使用自定义错误中的数据，则需要通过类型断言将错误作为定制错误类型的实例。
	_, e := f2(42)
	if ae, ok := e.(*argError); ok {
		fmt.Println(ae.arg)
		fmt.Println(ae.prob)
	}
}
