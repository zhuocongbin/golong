// Go provides built-in support for [base64 encoding/decoding]
// Go为[base64编码/解码]提供内置支持 (http://en.wikipedia.org/wiki/Base64).

package main

// This syntax imports the `encoding/base64` package with the `b64` name instead of the default `base64`.
//  It'll save us some space below.
// 该语法将“编码/base64”包导入“b64”名称，而不是默认的“base64”，它将为我们节省一些空间。
import b64 "encoding/base64"
import "fmt"

func main() {

	// Here's the `string` we'll encode/decode. 这是我们要编码/解码的“字符串”。
	data := "abc123!?$*&()'-=@~"

	// Go supports both standard and URL-compatible base64.
	//  Here's how to encode using the standard encoder.
	//  The encoder requires a `[]byte` so we cast our `string` to that type.
	// Go支持标准兼容和url兼容的base64。下面是如何使用标准编码器进行编码的方法。
	//编码器需要一个“[]字节”，因此我们将“字符串”转换为该类型。

	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(sEnc)

	// Decoding may return an error, which you can check if you don't already know the input to be well-formed.
	// 解码可能会返回一个错误，如果您不知道输入的格式良好，您可以检查这个错误。

	sDec, _ := b64.StdEncoding.DecodeString(sEnc)
	fmt.Println(string(sDec))
	fmt.Println()

	// This encodes/decodes using a URL-compatible base64 format.
	// 这个编码/解码使用一个url兼容的base64格式。
	uEnc := b64.URLEncoding.EncodeToString([]byte(data))
	fmt.Println(uEnc)
	uDec, _ := b64.URLEncoding.DecodeString(uEnc)
	fmt.Println(string(uDec))
}
