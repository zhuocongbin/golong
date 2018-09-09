// In the previous example we looked at
// [spawning external processes](spawning-processes). We
// do this when we need an external process accessible to
// a running Go process. Sometimes we just want to
// completely replace the current Go process with another
// (perhaps non-Go) one. To do this we'll use Go's
// implementation of the classic
// <a href="http://en.wikipedia.org/wiki/Exec_(operating_system)"><code>exec</code></a>
// function.

package main

import "syscall"
import "os"
import "os/exec"

func main() {

	// For our example we'll exec `ls`. Go requires an
	// absolute path to the binary we want to execute, so
	// we'll use `exec.LookPath` to find it (probably `/bin/ls`).
	// 对于我们的例子，我们将执行“ls”。Go要求绝对路径到我们想要执行的二进制文件。
	// 我们将使用的执行。 寻找它(可能是“/bin/ls”)。

	binary, lookErr := exec.LookPath("ls")
	if lookErr != nil {
		panic(lookErr)
	}

	// `Exec` requires arguments in slice form (as apposed to one big string).
	//  We'll give `ls` a few common arguments.
	//  Note that the first argument should be the program name.
	//“Exec”要求以切片形式进行参数(就像对一个大字符串一样)。
	// 我们会给出一些常见的论点。
	// 注意第一个参数应该是程序名。
	args := []string{"ls", "-a", "-l", "-h"}

	// `Exec` also needs a set of [environment variables](environment-variables)
	// to use. Here we just provide our current environment.
	//Exec还需要一组[环境变量](环境变量)来使用。 这里我们只提供当前的环境。
	env := os.Environ()

	// Here's the actual `syscall.Exec` call. If this call is
	// successful, the execution of our process will end
	// here and be replaced by the `/bin/ls -a -l -h`
	// process. If there is an error we'll get a return value.
	// 这是实际的系统调用。高管的电话。
	//	如果这个调用成功，我们的进程的执行将在这里结束，并被“/bin/ls -a -l -h”进程所替代。
	//  如果出现错误，我们将获得返回值。
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}
