// We often want to execute Go code at some point in the
// future, or repeatedly at some interval.
// 我们通常希望在将来某个时刻执行Go代码，或者在某个时间间隔重复执行。
// Go's built-in _timer_ and _ticker_ features make both of these tasks easy.
// Go的内置_timer_和_ticker_特性使这两项任务都很容易。
//  We'll look first at timers and then at [tickers](tickers).
// 我们先看定时器，然后再看(滴答声)。

package main

import "time"
import "fmt"

func main() {

	// Timers represent a single event in the future. You
	// tell the timer how long you want to wait, and it
	// provides a channel that will be notified at that
	// time. This timer will wait 2 seconds.
	// 计时器代表了未来的单个事件。
	// 你告诉计时器你要等待多长时间，它提供了一个在那个时候会被通知的通道。
	// 这个计时器将等待2秒。
	timer1 := time.NewTimer(2 * time.Second)

	// The `<-timer1.C` blocks on the timer's channel `C`
	// until it sends a value indicating that the timer
	// expired.
	// `<-timer1.C` 在计时器的通道“C”上的块，直到它发送一个指示计时器过期的值。
	<-timer1.C
	fmt.Println("Timer 1 expired")

	// If you just wanted to wait, you could have used
	// `time.Sleep`. One reason a timer may be useful is
	// that you can cancel the timer before it expires.
	// Here's an example of that.
	//如果你只是想等，你可以用“time.Sleep”。
	//计时器可能有用的一个原因是您可以在计时器到期之前取消计时器。
	//这里有一个例子。
	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 expired")
	}()
	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stopped")
	}
}
