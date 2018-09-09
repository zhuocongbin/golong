// Go supports _methods_ defined on struct types.
// Go支持_methods_上定义结构体类型。
package main

import "fmt"

type rect struct {
	width, height int
}

// This `area` method has a _receiver type_ of `*rect`.
// 这个“区域”方法有一个“*rect”的_receiver类型。
func (r *rect) area() int {
	return r.width * r.height
}

// Methods can be defined for either pointer or value receiver types. Here's an example of a value receiver.
// Methods可以为指针或值接收类型定义。 这是一个值接收器的例子。
func (r rect) perim() int {
	return 2*r.width + 2*r.height
}

func main() {
	r := rect{width: 10, height: 5}

	// Here we call the 2 methods defined for our struct.
	// 这里我们调用为结构定义的2个方法。
	fmt.Println("area: ", r.area())
	fmt.Println("perim:", r.perim())

	// Go automatically handles conversion between values and pointers for method calls.
	// Go自动处理方法调用的值和指针之间的转换。
	// You may want to use a pointer receiver type to avoid copying on method calls
	// or to allow the method to mutate the receiving struct.
	// 您可能想要使用一个指针接收器类型，以避免在方法调用上复制或允许方法对接收结构进行突变。
	rp := &r
	fmt.Println("area: ", rp.area())
	fmt.Println("perim:", rp.perim())
}
