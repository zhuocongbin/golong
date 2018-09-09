// By default channels are _unbuffered_, meaning that they will only accept sends (`chan <-`)
//  if there is a corresponding receive (`<- chan`) ready to receive the sent value.
// 默认通道是_unbuffered_，这意味着它们只接受发送(' chan <- ')，如果有相应的receive (' <- chan ')准备接收发送的值。
//  _Buffered channels_ accept a limited number of  values without a corresponding receiver for those values.
// _Buffered channels_接受有限数量的值，而没有相应的接收器。
//

package main

import "fmt"

func main() {

	// Here we `make` a channel of strings buffering up to 2 values.
	// 在这里，我们使一个字符串的通道缓冲到两个值。
	messages := make(chan string, 2)

	// Because this channel is buffered, we can send these values into the channel without a corresponding concurrent receive.
	// 由于这个通道是缓冲的，我们可以将这些值发送到通道中，而不需要相应的并发接收。
	messages <- "buffered"
	messages <- "channel"

	// Later we can receive these two values as usual.
	// 稍后我们可以像往常一样接收到这两个值。
	fmt.Println(<-messages)
	fmt.Println(<-messages)
}
