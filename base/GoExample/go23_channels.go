// _Channels_ are the pipes that connect concurrent goroutines.
//_Channels_是连接并行goroutines的管道。
// You can send values into channels from one goroutine and receive those values into another goroutine.
// 您可以将值发送到一个goroutine的通道，并将这些值接收到另一个goroutine。
//

package main

import "fmt"

func main() {

	// Create a new channel with `make(chan val-type)`.Channels are typed by the values they convey.
	// 用“make(chan vall -type)”创建一个新的通道。 通道由它们传递的值输入。
	messages := make(chan string)

	// _Send_ a value into a channel using the `channel <-`syntax.
	// Here we send `"ping"`  to the `messages`channel we made above, from a new goroutine.
	// 使用' channel <- '语法将值发送到通道。
	// 在这里，我们将“ping”发送到我们上面的“消息”通道，从一个新的goroutine。
	go func() { messages <- "ping" }()

	// The `<-channel` syntax _receives_ a value from the channel.
	//  Here we'll receive the `"ping"` message we sent above and print it out.
	// “<-通道”语法_receives_一个来自通道的值。
	// 在这里，我们将收到我们上面发送的“ping”信息并打印出来。
	msg := <-messages
	fmt.Println(msg)
}
