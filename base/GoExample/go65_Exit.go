// Use `os.Exit` to immediately exit with a given status.
// 使用`os.Exit`一个给定的状态立即退出。

package main

import "fmt"
import "os"

func main() {

	// `defer`s will not be run when using `os.Exit`, so this `fmt.Println` will never be called.
	// 当使用`os.Exit`的时候`defer`s 不会被运行，因此这个`fmt.Println`永远不会被调用。
	defer fmt.Println("!")

	// Exit with status 3.
	// 以3的状态退出
	os.Exit(3)
}

// Note that unlike e.g. C, Go does not use an integer return value from `main` to indicate exit status.
// 注意不像例如C，GO不用整型返回值
// If you'd like to exit with a non-zero status you should use `os.Exit`.
//
