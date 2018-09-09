// _Maps_ are Go's built-in [associative data type](http://en.wikipedia.org/wiki/Associative_array)
// (sometimes called _hashes_ or _dicts_ in other languages).

package main

import "fmt"

func main() {

	// To create an empty map, use the builtin `make`:
	// `make(map[key-type]val-type)`.
	m := make(map[string]int)

	// Set key/value pairs using typical `name[key] = val`syntax.
	// 用`name[key] = val`语法来设置键值对
	m["k1"] = 7
	m["k2"] = 13

	// Printing a map with e.g. `fmt.Println` will show all of its key/value pairs.
	// 用`fmt.Println` 打印一个map将显示其所有键/值对。
	fmt.Println("map:", m)

	// Get a value for a key with `name[key]`.
	// 用`name[key]`获取一个键的值
	v1 := m["k1"]
	fmt.Println("v1: ", v1)

	// The builtin `len` returns the number of key/value pairs when called on a map.
	// 当调用map时，builtin ' len '返回键/值对的数目。
	fmt.Println("len:", len(m))

	// The builtin `delete` removes key/value pairs from a map.
	// builtin ' delete '从map中删除键/值对。
	delete(m, "k2")
	fmt.Println("map:", m)

	// The optional second return value when getting a value from a map indicates if the key was present in the map.
	// 当从映射中获取值时，可选的第二个返回值指示该键是否存在于映射中。
	//  This can be used to disambiguate between missing keys and keys with zero values like `0` or `""`.
	// 这可以用于在丢失的键和键之间消除歧义，比如`0` or `""`。
	//  Here we didn't need the value itself, so we ignored it with the _blank identifier_`_`.
	// 这里我们不需要值本身，所以我们用_blank标识符“_”来忽略它。
	//
	_, prs := m["k2"]
	fmt.Println("prs:", prs)

	// You can also declare and initialize a new map in the same line with this syntax.
	// 您还可以使用该语法声明和初始化一个新映射。
	n := map[string]int{"foo": 1, "bar": 2}
	fmt.Println("map:", n)
}
