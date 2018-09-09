// Basic sends and receives on channels are blocking.
// However, we can use `select` with a `default` clause to
// implement _non-blocking_ sends, receives, and even
// non-blocking multi-way `select`s.
// 基本的发送和接收在通道上是阻塞。
// 但是，我们可以使用“select”来实现_non-blocking_发送、接收，甚至非阻塞多路选择。
package main

import "fmt"

func main() {
	messages := make(chan string)
	signals := make(chan bool)

	// Here's a non-blocking receive. If a value is
	// available on `messages` then `select` will take
	// the `<-messages` `case` with that value. If not
	// it will immediately take the `default` case.
	//  这是一个非阻塞接收。
	//  如果“消息”上有一个值，那么“select”将使用“<-message”的“case”值。
	//	如果不是，它将立即采取“默认”的情况。
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no message received")
	}

	// A non-blocking send works similarly. Here `msg`
	// cannot be sent to the `messages` channel, because
	// the channel has no buffer and there is no receiver.
	// Therefore the `default` case is selected.
	// 非阻塞发送也类似。 这里的“msg”不能发送到“消息”通道，因为通道没有缓冲区，也没有接收器。
	// 因此，选择“默认”情况。
	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("sent message", msg)
	default:
		fmt.Println("no message sent")
	}

	// We can use multiple `case`s above the `default`
	// clause to implement a multi-way non-blocking
	// select. Here we attempt non-blocking receives
	// on both `messages` and `signals`.
	// 我们可以使用多个' case '在' default '子句中实现多路非阻塞选择。
	// 在这里，我们尝试非阻塞接收在“消息”和“信号”上。
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}
