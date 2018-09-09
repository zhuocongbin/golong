// Go's `sort` package implements sorting for builtins
// and user-defined types. We'll look at sorting for builtins first.
 // Go的“sort”包实现了对builtins和用户定义类型的排序。
// 我们先来看看构建的排序。

package main

import "fmt"
import "sort"

func main() {

	// Sort methods are specific to the builtin type;
	// here's an example for strings. Note that sorting is
	// in-place, so it changes the given slice and doesn't return a new one.
	// 排序方法是特定于构建类型的;这里有一个字符串示例。
	//注意，排序是正确的，所以它会改变给定的切片，不会返回一个新的。
	strs := []string{"c", "a", "b"}
	sort.Strings(strs)
	fmt.Println("Strings:", strs)

	// An example of sorting `int`s.
	ints := []int{7, 2, 4}
	sort.Ints(ints)
	fmt.Println("Ints:   ", ints)

	// We can also use `sort` to check if a slice is already in sorted order.
	// 我们还可以使用“sort”来检查是否已经排序了一个切片。
	s := sort.IntsAreSorted(ints)
	fmt.Println("Sorted: ", s)
}
