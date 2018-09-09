// _Slices_ are a key data type in Go, giving a more powerful interface to sequences than arrays.
// _Slices_是Go中的关键数据类型，它提供了比数组更强大的序列接口。

package main

import "fmt"

func main() {

	// Unlike arrays, slices are typed only by the
	// elements they contain (not the number of elements).
	// To create an empty slice with non-zero length, use
	// the builtin `make`. Here we make a slice of
	// `string`s of length `3` (initially zero-valued).
	//与数组不同，切片仅由它们所包含的元素(而不是元素的数量)来输入。
	//要创建一个非零长度的空切片，请使用builtin ' make '。
	//这里我们把长度为3的字符串(最初为零值)。
	s := make([]string, 3)
	fmt.Println("emp:", s)

	// We can set and get just like with arrays.
	// 我们可以设置和使用数组。
	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println("set:", s)
	fmt.Println("get:", s[2])

	// `len` returns the length of the slice as expected.
	//“len”将按预期返回切片的长度。
	fmt.Println("len:", len(s))

	// In addition to these basic operations, slices
	// support several more that make them richer than
	// arrays. One is the builtin `append`, which
	// returns a slice containing one or more new values.
	// Note that we need to accept a return value from
	// `append` as we may get a new slice value.
	//除了这些基本操作之外，切片还支持多一些，使它们比数组更丰富。
	//一个是builtin ' append '，它返回一个包含一个或多个新值的片段。
	//注意，我们需要接受从“append”返回的返回值，因为我们可能会得到一个新的切片值。
	s = append(s, "d")
	s = append(s, "e", "f")
	fmt.Println("apd:", s)

	// Slices can also be `copy`'d. Here we create an
	// empty slice `c` of the same length as `s` and copy
	// into `c` from `s`.
	//切片也可以是“复制”d。
	//在这里，我们创建一个与' s '相同长度的空片' c '，并从' s '复制到' c '。
	c := make([]string, len(s))
	copy(c, s)
	fmt.Println("cpy:", c)

	// Slices support a "slice" operator with the syntax
	// `slice[low:high]`. For example, this gets a slice
	// of the elements `s[2]`, `s[3]`, and `s[4]`.
	// Slices支持“切片”操作符，语法`slice[low:high]`。
	// 例如，它获取了元素的[2]'、' s[3] '和' s[4]的部分。
	l := s[2:5]
	fmt.Println("sl1:", l)

	// This slices up to (but excluding) `s[5]`.
	// 这部分是(但不包括)`s[5]`
	l = s[:5]
	fmt.Println("sl2:", l)

	// And this slices up from (and including) `s[2]`.
	//这部分来自(包括)`s[2]`.
	l = s[2:]
	fmt.Println("sl3:", l)

	// We can declare and initialize a variable for slice in a single line as well.
	// 我们也可以声明和初始化一个变量在单行中。
	t := []string{"g", "h", "i"}
	fmt.Println("dcl:", t)

	// Slices can be composed into multi-dimensional data structures.
	// The length of the inner slices can vary, unlike with multi-dimensional arrays.
	// 切片可以组成多维数据结构。 与多维数组不同，内片的长度可以变化。
	twoD := make([][]int, 3)
	for i := 0; i < 3; i++ {
		innerLen := i + 1
		twoD[i] = make([]int, innerLen)
		for j := 0; j < innerLen; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoD)
}
