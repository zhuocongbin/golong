// Go's _select_ lets you wait on multiple channel operations.
// Go的_select_让您等待多个通道操作。
//  Combining goroutines and channels with select is a powerful feature of Go.
// 结合goroutines和channel与select是一个强大的功能。

package main

import "time"
import "fmt"

func main() {

	// For our example we'll select across two channels.
	// 对于我们的示例，我们将在两个通道中进行选择。
	c1 := make(chan string)
	c2 := make(chan string)

	// Each channel will receive a value after some amount
	// of time, to simulate e.g. blocking RPC operations
	// executing in concurrent goroutines.
	//每个通道在一定时间后都会收到一个值，以模拟在并行的goroutines中执行的阻塞RPC操作。
	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	// We'll use `select` to await both of these values simultaneously, printing each one as it arrives.
	// 我们将使用“select”来同时等待这两个值，并在它到达时打印每个值。
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}
}
