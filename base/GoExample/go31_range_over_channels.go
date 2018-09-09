// In a [previous](range) example we saw how `for` and
// `range` provide iteration over basic data structures.
// We can also use this syntax to iterate over
// values received from a channel.
// 在[先前](范围)示例中，我们看到了“for”和“range”在基本数据结构上的迭代。
// 我们还可以使用此语法对从通道接收的值进行迭代。
package main

import "fmt"

func main() {

	// We'll iterate over 2 values in the `queue` channel.
	// 我们将在“队列”通道中迭代2个值。
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)

	// This `range` iterates over each element as it's
	// received from `queue`. Because we `close`d the
	// channel above, the iteration terminates after
	// receiving the 2 elements.
	// 这个“range”遍历每个元素，因为它是从“queue”接收的。
	// 因为我们关闭了上面的通道，在接收到两个元素后，迭代终止。
	for elem := range queue {
		fmt.Println(elem)
	}
}
