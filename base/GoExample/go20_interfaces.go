// _Interfaces_ are named collections of method signatures.
// _Interfaces_被命名为方法签名集合。

package main

import "fmt"
import "math"

// Here's a basic interface for geometric shapes.
// 这是几何图形的基本界面。
type geometry interface {
	area() float64
	perim() float64
}

// For our example we'll implement this interface on `rect` and `circle` types.
// 对于我们的示例，我们将在“rect”和“circle”类型上实现这个接口。
type rect struct {
	width, height float64
}
type circle struct {
	radius float64
}

// To implement an interface in Go, we just need to implement all the methods in the interface.
// 为了实现一个接口，我们只需要实现接口中的所有方法。
// Here we implement `geometry` on `rect`s.

func (r rect) area() float64 {
	return r.width * r.height
}
func (r rect) perim() float64 {
	return 2*r.width + 2*r.height
}

// The implementation for `circle`s.
// circle的的实现。
func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}

// If a variable has an interface type, then we can call methods that are in the named interface.
// 如果一个变量有一个接口类型，那么我们可以调用命名接口中的方法。
// Here's a generic `measure` function taking advantage of this to work on any `geometry`.
// 这里有一个通用的“度量”函数，利用它来处理任何“几何”。
func measure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perim())
}

func main() {
	r := rect{width: 3, height: 4}
	c := circle{radius: 5}

	// The `circle` and `rect` struct types both implement the `geometry` interface
	// “circle”和“rect”结构类型都实现了“几何”界面。
	// so we can use instances of these structs as arguments to `measure`.
	// 所以我们可以用这些结构的实例作为参数来衡量
	measure(r)
	measure(c)
}
