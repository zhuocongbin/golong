
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
	flag.Parse()
	log.Info("message ver: \"%s\" start", ver.Version)
	if err := InitConfig(); err != nil {
		panic(err)
	}
	// Set max routine
	runtime.GOMAXPROCS(Conf.MaxProc)
	// init log
	log.LoadConfiguration(Conf.Log)
	defer log.Close()
	// start pprof http
	perf.Init(Conf.PprofBind)
	// Initialize redis
	if err := InitStorage(); err != nil {
		panic(err)
	}
	// init rpc service
	if err := InitRPC(); err != nil {
		panic(err)
	}
	// init zookeeper
	zk, err := InitZK()
	if err != nil {
		if zk != nil {
			zk.Close()
		}
		panic(err)
	}
	// process init
	if err = process.Init(Conf.User, Conf.Dir, Conf.PidFile); err != nil {
		panic(err)
	}
	// init signals, block wait signals
	sig := InitSignal()
	HandleSignal(sig)
	// exit
	log.Info("message stop")
}
