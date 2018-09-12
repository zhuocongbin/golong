

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
	// parse cmd-line arguments  //解析cmd-line参数
	flag.Parse()
	log.Info("comet ver: \"%s\" start", ver.Version)
	// init config //初始化配置
	if err := InitConfig(); err != nil {
		panic(err)
	}
	// set max routine //设置最大线程
	runtime.GOMAXPROCS(Conf.MaxProc)
	// init log //初始化日志
	log.LoadConfiguration(Conf.Log)
	defer log.Close()
	// start pprof
	perf.Init(Conf.PprofBind)
	// create channel //创建通道
	// if process exit, close channel //如果进程存在，关闭通道
	UserChannel = NewChannelList()
	defer UserChannel.Close()
	// start stats
	StartStats()
	// start rpc
	if err := StartRPC(); err != nil {
		panic(err)
	}
	// start comet
	if err := StartComet(); err != nil {
		panic(err)
	}
	// init zookeeper
	zkConn, err := InitZK()
	if err != nil {
		if zkConn != nil {
			zkConn.Close()
		}
		panic(err)
	}
	// process init
	if err = process.Init(Conf.User, Conf.Dir, Conf.PidFile); err != nil {
		panic(err)
	}
	// init signals, block wait signals
	signalCH := InitSignal()
	HandleSignal(signalCH)
	// exit
	log.Info("comet stop")
}
