// When using channels as function parameters, you can
// specify if a channel is meant to only send or receive
// values. This specificity increases the type-safety of
// the program.
// 当使用通道作为函数参数时，可以指定通道是否仅用于发送或接收值。
// 这种特性增加了程序的类型安全性。
package main

import "fmt"

// This `ping` function only accepts a channel for sending values.
//  It would be a compile-time error to try to receive on this channel.
// 这个“ping”函数只接受一个发送值的通道。试图在这个通道上接收它将是一个编译时错误。
func ping(pings chan<- string, msg string) {
	pings <- msg
}

// The `pong` function accepts one channel for receives (`pings`) and a second for sends (`pongs`).
// “pong”功能接受一个接收通道(“ping”)，第二个用于发送(“pongs”)。
func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

func main() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}
