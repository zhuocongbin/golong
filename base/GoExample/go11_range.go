// _range_ iterates over elements in a variety of data structures.
//  Let's see how to use `range` with some of the data structures we've already learned.
// _range_遍历各种数据结构中的元素。 让我们看看如何使用我们已经学过的一些数据结构来使用“range”。

package main

import "fmt"

func main() {

	// Here we use `range` to sum the numbers in a slice.Arrays work like this too.
	// 在这里，我们使用“range”来将数字加起来。 数组也是这样工作的。
	nums := []int{2, 3, 4}
	sum := 0
	for _, num := range nums {
		sum += num
	}
	fmt.Println("sum:", sum)

	// `range` on arrays and slices provides both the index and value for each entry.
	//  Above we didn't need the index, so we ignored it with the blank identifier `_`.
	// Sometimes we actually want the indexes though.
	// 在数组和片上的范围提供了每个条目的索引和值。
	// 上面我们不需要索引，所以我们用空白标识符“_”来忽略它。
	// 有时我们实际上想要索引。

	for i, num := range nums {
		if num == 3 {
			fmt.Println("index:", i)
		}
	}

	// `range` on map iterates over key/value pairs.
	// `range`在map上迭代键/值对。
	kvs := map[string]string{"a": "apple", "b": "banana"}
	for k, v := range kvs {
		fmt.Printf("%s -> %s\n", k, v)
	}

	// `range` can also iterate over just the keys of a map.
	for k := range kvs {
		fmt.Println("key:", k)
	}

	// `range` on strings iterates over Unicode code points.
	//  The first value is the starting byte index of the `rune` and the second the `rune` itself.
	//
	for i, c := range "go" {
		fmt.Println(i, c)
	}
}
