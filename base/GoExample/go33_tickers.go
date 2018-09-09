// [Timers](timers) are for when you want to do
// something once in the future - _tickers_ are for when
// you want to do something repeatedly at regular
// intervals. Here's an example of a ticker that ticks
// periodically until we stop it.
// (定时器)是当你想在将来做某事的时候。
// - _tickers_是当你想要定时重复做某事的时候。
// 这里有个例子，它会周期性地滴答作响，直到我们停止它。

package main

import "time"
import "fmt"

func main() {

	// Tickers use a similar mechanism to timers: a
	// channel that is sent values. Here we'll use the
	// `range` builtin on the channel to iterate over
	// the values as they arrive every 500ms.
	// Tickers使用与计时器类似的机制:发送值的通道。
	// 在这里，我们将在通道上使用“range”内置函数来迭代每500毫秒到达的值。
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()

	// Tickers can be stopped like timers. Once a ticker
	// is stopped it won't receive any more values on its
	// channel. We'll stop ours after 1600ms.
	// 滴答声可以像计时器一样停止。
	// 一旦代码停止，它就不会在其通道上接收更多的值。
	// 我们会在1600ms后停下来。
	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}
