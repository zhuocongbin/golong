// _Closing_ a channel indicates that no more values
// will be sent on it. This can be useful to communicate
// completion to the channel's receivers.
// _Closing_一个通道表示将不再发送任何值。 这对于将完成信息传递给信道的接收者是很有用的。
package main

import "fmt"

// In this example we'll use a `jobs` channel to
// communicate work to be done from the `main()` goroutine
// to a worker goroutine. When we have no more jobs for
// the worker we'll `close` the `jobs` channel.
// 在本例中，我们将使用一个“jobs”通道来将工作从“main()”“goroutine”传递给一个工人goroutine。
// 当我们不再为工人提供工作时，我们将关闭“就业”频道。
func main() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	// Here's the worker goroutine. It repeatedly receives
	// from `jobs` with `j, more := <-jobs`. In this
	// special 2-value form of receive, the `more` value
	// will be `false` if `jobs` has been `close`d and all
	// values in the channel have already been received.
	// We use this to notify on `done` when we've worked
	// all our jobs.
	//这是工人goroutine。 它不断从“工作”中获得“j”，更多:= <-jobs。
	//在这种特殊的2值接收形式中，如果“jobs”已经“关闭”，并且通道中的所有值都已经收到，那么“more”值将是“false”。
	//当我们工作了所有的工作时，我们用这个来通知“完成”。
	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				done <- true
				return
			}
		}
	}()

	// This sends 3 jobs to the worker over the `jobs`channel, then closes it.
	// 这将给工作人员发送3个工作机会，然后关闭它。
	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("sent job", j)
	}
	close(jobs)
	fmt.Println("sent all jobs")

	// We await the worker using the [synchronization](channel-synchronization) approach we saw earlier.
	// 我们使用前面看到的(同步) (通道同步)方法等待worker。
	<-done
}
