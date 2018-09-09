// Go's _structs_ are typed collections of fields.
// Go的_structs_是字段的类型化集合。
// They're useful for grouping data together to form records.
// 它们对于将数据分组到一起形成记录非常有用。

package main

import "fmt"

// This `person` struct type has `name` and `age` fields.
// 这个“person”结构类型有“name”和“age”字段。
type person struct {
	name string
	age  int
}

func main() {

	// This syntax creates a new struct.
	// 这种语法创建了一个新的结构。
	fmt.Println(person{"Bob", 20})

	// You can name the fields when initializing a struct.
	// 您可以在初始化结构时命名字段。
	fmt.Println(person{name: "Alice", age: 30})

	// Omitted fields will be zero-valued.
	// 省略字段将为零值。
	fmt.Println(person{name: "Fred"})

	// An `&` prefix yields a pointer to the struct.
	// 一个' & '前缀产生一个指向结构的指针。
	fmt.Println(&person{name: "Ann", age: 40})

	// Access struct fields with a dot.
	// 使用一个点访问struct字段。
	s := person{name: "Sean", age: 50}
	fmt.Println(s.name)

	// You can also use dots with struct pointers - the pointers are automatically dereferenced.
	// 您还可以使用带有struct指针的点——指针会自动取消引用。
	sp := &s
	fmt.Println(sp.age)

	// Structs are mutable.
	sp.age = 51
	fmt.Println(sp.age)
}
