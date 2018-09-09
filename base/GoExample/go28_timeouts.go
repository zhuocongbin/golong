// _Timeouts_ are important for programs that connect to
// external resources or that otherwise need to bound
// execution time. Implementing timeouts in Go is easy and
// elegant thanks to channels and `select`.
// _Timeouts_对于连接到外部资源的程序或需要绑定执行时间的程序非常重要。
// 通过通道和“select”，实现Go中的超时非常简单和优雅。
package main

import "time"
import "fmt"

func main() {

	// For our example, suppose we're executing an external
	// call that returns its result on a channel `c1` after 2s.
	// 对于我们的例子来说，假设我们正在执行一个外部调用，它在2s之后返回其在通道c1上的结果。
	c1 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "result 1"
	}()

	// Here's the `select` implementing a timeout.
	// `res := <-c1` awaits the result and `<-Time.After`
	// awaits a value to be sent after the timeout of
	// 1s. Since `select` proceeds with the first
	// receive that's ready, we'll take the timeout case
	// if the operation takes more than the allowed 1s.
	//这里是“选择”实现超时。 “res:= <-c1”等待结果和<时间。 “等待在1的超时后发送的值”。
	//由于“select”是第一个接收到的，所以如果操作超过了允许的1s，我们将接受超时。
	select {
	case res := <-c1:
		fmt.Println(res)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout 1")
	}

	// If we allow a longer timeout of 3s, then the receive
	// from `c2` will succeed and we'll print the result.
	//如果我们允许3秒的超时，那么c2的接收将会成功，我们将打印结果。
	c2 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "result 2"
	}()
	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(3 * time.Second):
		fmt.Println("timeout 2")
	}
}
