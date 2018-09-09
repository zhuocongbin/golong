// We can use channels to synchronize execution across goroutines.
// 我们可以使用通道在goroutines之间同步执行。
//  Here's an example of using a blocking receive to wait for a goroutine to finish.
// 这里有一个使用阻塞接收等待goroutine完成的例子。

package main

import "fmt"
import "time"

// This is the function we'll run in a goroutine.
// 这是我们将在goroutine中运行的函数。
// The`done` channel will be used to notify another goroutine that this function's work is done.
// “done”通道将用于通知另一个goroutine，该函数的工作已经完成。
func worker(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second)
	fmt.Println("done")

	// Send a value to notify that we're done.
	// 发送一个值通知我们完成了。
	done <- true
}

func main() {

	// Start a worker goroutine, giving it the channel to notify on.
	// 启动一个worker goroutine，给它一个通知的渠道。
	done := make(chan bool, 1)
	go worker(done)

	// Block until we receive a notification from the worker on the channel.
	// 阻塞，直到我们收到通道上的worker的通知。
	<-done
}
