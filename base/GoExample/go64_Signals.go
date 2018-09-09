// Sometimes we'd like our Go programs to intelligently handle
// 有时候，我们希望我们的Go程序能智能地处理。
// [Unix signals](http://en.wikipedia.org/wiki/Unix_signal).
// For example, we might want a server to gracefully shutdown when it receives a `SIGTERM`, or a command-line
// tool to stop processing input if it receives a `SIGINT`.
// 例如，如果服务器接收到“SIGTERM”，或者命令行工具在接收到“SIGINT”时停止处理输入，我们可能希望它能优雅地关闭。
// Here's how to handle signals in Go with channels.
// 下面是如何处理与通道的信号。

package main

import "fmt"
import "os"
import "os/signal"
import "syscall"

func main() {

	// Go signal notification works by sending `os.Signal` values on a channel.
	//go的信号通知的工作原理是在通道上发送`os.Signal`的值
	// We'll create a channel to receive these notifications (we'll also make one to notify us when the program can exit).
	// 我们将创建一个通道来接收这些通知(我们还将在程序可以退出时通知我们)。
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// `signal.Notify` registers the given channel to receive notifications of the specified signals.
	// `signal.Notify` 注册在给定的通道上接收指定信号的通知。
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// This goroutine executes a blocking receive for signals.
	// 这个goroutine执行一个阻塞接收信号。
	//  When it gets one it'll print it out and then notify the program that it can finish.
	// 当它得到一个，它会打印出来然后通知程序它可以完成。
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	// The program will wait here until it gets the
	// expected signal (as indicated by the goroutine
	// above sending a value on `done`) and then exit.
	//程序将在这里等待，直到它得到预期的信号(正如上面的goroutine所示，在“done”中发送一个值)然后退出。
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
