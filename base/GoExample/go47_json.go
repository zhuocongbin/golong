// Go offers built-in support for JSON encoding and decoding, including to and from built-in and custom data types.
// Go为JSON编码和解码提供内置支持，包括内置和自定义数据类型。

package main

import "encoding/json"
import "fmt"
import "os"

// We'll use these two structs to demonstrate encoding and decoding of custom types below.
// 我们将使用这两个结构来演示自定义类型的编码和解码。
type response1 struct {
	Page   int
	Fruits []string
}
type response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

func main() {

	// First we'll look at encoding basic data types to JSON strings. Here are some examples for atomic values.
	// 首先，我们将研究如何将基本数据类型编码为JSON字符串。 这里有一些原子值的例子。

	bolB, _ := json.Marshal(true)
	fmt.Println(string(bolB))

	intB, _ := json.Marshal(1)
	fmt.Println(string(intB))

	fltB, _ := json.Marshal(2.34)
	fmt.Println(string(fltB))

	strB, _ := json.Marshal("gopher")
	fmt.Println(string(strB))

	// And here are some for slices and maps, which encode to JSON arrays and objects as you'd expect.
	// 下面是一些用于切片和映射的代码，它们可以像您期望的那样对JSON数组和对象进行编码。
	slcD := []string{"apple", "peach", "pear"}
	slcB, _ := json.Marshal(slcD)
	fmt.Println(string(slcB))

	mapD := map[string]int{"apple": 5, "lettuce": 7}
	mapB, _ := json.Marshal(mapD)
	fmt.Println(string(mapB))

	// The JSON package can automatically encode your custom data types.
	// It will only include exported fields in the encoded output and will by default use those names as the JSON keys.
	// JSON包可以自动编码自定义数据类型。 它将只包括已编码输出的导出字段，默认情况下将使用这些名称作为JSON键。

	res1D := &response1{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res1B, _ := json.Marshal(res1D)
	fmt.Println(string(res1B))

	// You can use tags on struct field declarations to customize the encoded JSON key names.
	//  Check the definition of `response2` above to see an example of such tags.
	// 您可以在struct字段声明上使用标记来定制编码的JSON密钥名。
	//检查上面的“response2”的定义，看看这些标签的例子。

	res2D := &response2{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res2B, _ := json.Marshal(res2D)
	fmt.Println(string(res2B))

	// Now let's look at decoding JSON data into Go values.
	//  Here's an example for a generic data structure.
	// 现在让我们看看将JSON数据解码为Go值。
	// 这里有一个通用数据结构的例子。
	byt := []byte(`{"num":6.13,"strs":["a","b"]}`)

	// We need to provide a variable where the JSON package can put the decoded data.
	// This `map[string]interface{}` will hold a map of strings to arbitrary data types.
	// 我们需要提供一个变量，其中JSON包可以放置解码数据。
	//	这个“map[string]接口{}”将为任意数据类型保留一个字符串映射。

	var dat map[string]interface{}

	// Here's the actual decoding, and a check for associated errors.
	// 这是实际的解码，并检查相关的错误。
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)

	// In order to use the values in the decoded map,we'll need to cast them to their appropriate type.
	// 为了在解码映射中使用值，我们需要将它们转换为合适的类型。
	// For example here we cast the value in `num` to the expected `float64` type.
	// 例如，我们将“num”的值转换为预期的“float64”类型。
	num := dat["num"].(float64)
	fmt.Println(num)

	// Accessing nested data requires a series of casts.
	// 访问嵌套数据需要一系列的类型转换。
	strs := dat["strs"].([]interface{})
	str1 := strs[0].(string)
	fmt.Println(str1)

	// We can also decode JSON into custom data types.
	// This has the advantages of adding additional type-safety to our programs and eliminating the
	// need for type assertions when accessing the decoded data.
	// 我们还可以将JSON解码为自定义数据类型。
	//	这具有在我们的程序中添加额外的类型安全的优点，并且在访问解码数据时消除了类型断言的需要。

	str := `{"page": 1, "fruits": ["apple", "peach"]}`
	res := response2{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res)
	fmt.Println(res.Fruits[0])

	// In the examples above we always used bytes and strings as intermediates between the data and JSON representation on standard out.
	// 在上面的示例中，我们总是使用字节和字符串作为数据和标准输出的JSON表示之间的中介。
	//  We can also stream JSON encodings directly to `os.Writer`s like `os.Stdout` or even HTTP response bodies.
	// 我们还可以将JSON编码直接流到“os”。作家就像操作系统。Stdout，甚至HTTP响应体。
	//
	enc := json.NewEncoder(os.Stdout)
	d := map[string]int{"apple": 5, "lettuce": 7}
	enc.Encode(d)
}
