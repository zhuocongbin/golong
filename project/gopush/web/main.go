
package main

import (
	log "github.com/alecthomas/log4go"
	"flag"
	"github.com/Terry-Mao/gopush-cluster/perf"
	"github.com/Terry-Mao/gopush-cluster/process"
	"github.com/Terry-Mao/gopush-cluster/ver"
	"runtime"
)

func main() {
	var err error
	// Parse cmd-line arguments   解析命令行参数
	flag.Parse()
	log.Info("web ver: \"%s\" start", ver.Version)
	if err = InitConfig(); err != nil {
		panic(err)
	}
	// Set max routine  设置最大例行
	runtime.GOMAXPROCS(Conf.MaxProc)
	// init log 初始化日志
	log.LoadConfiguration(Conf.Log)
	defer log.Close()
	// init zookeeper
	zkConn, err := InitZK()
	if err != nil {
		if zkConn != nil {
			zkConn.Close()
		}
		panic(err)
	}
	// start pprof http
	perf.Init(Conf.PprofBind)
	// start http listen.
	StartHTTP()
	// process init
	if err = process.Init(Conf.User, Conf.Dir, Conf.PidFile); err != nil {
		panic(err)
	}
	// init signals, block wait signals
	signalCH := InitSignal()
	HandleSignal(signalCH)
	log.Info("web stop")
}
