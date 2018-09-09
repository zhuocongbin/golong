// In this example we'll look at how to implement
// a _worker pool_ using goroutines and channels.
// 在本例中，我们将讨论如何实现_worker pool_使用goroutine和channels。
package main

import "fmt"
import "time"

// Here's the worker, of which we'll run several concurrent instances.
//这里是worker，我们将运行几个并发实例。
// These workers will receive work on the `jobs` channel and send the corresponding results on `results`.
//这些workers将在`jobs`频道上接收，并将相应的结果发送到`results`上。
//  We'll sleep a second per job to simulate an expensive task.
// 为了模拟一项高级的任务，我们每个job都要休眠一秒钟。
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func main() {

	// In order to use our pool of workers we need to send them work and collect their results.We make 2 channels for this.
	//为了使用我们的workers池，我们需要派他们工作和收集他们的结果。 我们有两个channels。
	//
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// This starts up 3 workers, initially blocked
	// because there are no jobs yet.
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Here we send 5 `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	// Finally we collect all the results of the work.
	for a := 1; a <= 5; a++ {
		<-results
	}
}
