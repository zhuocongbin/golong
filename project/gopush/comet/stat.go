
package main

import (
	log "github.com/alecthomas/log4go"
	"encoding/json"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	// server
	startTime int64 // process start unixnano
	// channel
	ChStat = &ChannelStat{}
	// message
	MsgStat = &MessageStat{}
	// connection
	ConnStat = &ConnectionStat{}
)

// Channel stat info //通道统计信息
type ChannelStat struct {
	Access uint64 // total access count
	Create uint64 // total create count
	Delete uint64 // total delete count
}

func (s *ChannelStat) IncrAccess() {
	atomic.AddUint64(&s.Access, 1)
}

func (s *ChannelStat) IncrCreate() {
	atomic.AddUint64(&s.Create, 1)
}

func (s *ChannelStat) IncrDelete() {
	atomic.AddUint64(&s.Delete, 1)
}

// Stat get the channle stat info //获取通道统计信息
func (s *ChannelStat) Stat() []byte {
	res := map[string]interface{}{}
	res["access"] = s.Access
	res["create"] = s.Create
	res["delete"] = s.Delete
	res["current"] = UserChannel.Count()
	return jsonRes(res)
}

// Message stat info // 消息统计信息
type MessageStat struct {
	Succeed uint64 // total push message succeed count
	Failed  uint64 // total push message failed count
}

func (s *MessageStat) IncrSucceed(delta uint64) {
	atomic.AddUint64(&s.Succeed, delta)
}

func (s *MessageStat) IncrFailed(delta uint64) {
	atomic.AddUint64(&s.Failed, delta)
}

// Stat get the message stat info //获取消息统计信息
func (s *MessageStat) Stat() []byte {
	res := map[string]interface{}{}
	res["succeed"] = s.Succeed
	res["failed"] = s.Failed
	res["total"] = s.Succeed + s.Failed
	return jsonRes(res)
}

// Connection stat info // 连接统计信息
type ConnectionStat struct {
	Add    uint64 // total add connection count
	Remove uint64 // total remove connection count
}

func (s *ConnectionStat) IncrAdd() {
	atomic.AddUint64(&s.Add, 1)
}

func (s *ConnectionStat) IncrRemove() {
	atomic.AddUint64(&s.Remove, 1)
}

// Stat get the connection stat info // 获取连接统计信息
func (s *ConnectionStat) Stat() []byte {
	res := map[string]interface{}{}
	res["add"] = s.Add
	res["remove"] = s.Remove
	res["current"] = s.Add - s.Remove
	return jsonRes(res)
}

func statListen(bind string) {
	httpServeMux := http.NewServeMux()
	httpServeMux.HandleFunc("/stat", StatHandle)
	if err := http.ListenAndServe(bind, httpServeMux); err != nil {
		log.Error("http.ListenAdServe(\"%s\") error(%v)", bind, err)
		panic(err)
	}
}

// start stats, called at process start
func StartStats() {
	startTime = time.Now().UnixNano()
	for _, bind := range Conf.StatBind {
		log.Info("start stat listen addr:\"%s\"", bind)
		go statListen(bind)
	}
}

// memory stats //内存统计
func MemStats() []byte {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	// general
	res := map[string]interface{}{}
	res["alloc"] = m.Alloc
	res["total_alloc"] = m.TotalAlloc
	res["sys"] = m.Sys
	res["lookups"] = m.Lookups
	res["mallocs"] = m.Mallocs
	res["frees"] = m.Frees
	// heap
	res["heap_alloc"] = m.HeapAlloc
	res["heap_sys"] = m.HeapSys
	res["heap_idle"] = m.HeapIdle
	res["heap_inuse"] = m.HeapInuse
	res["heap_released"] = m.HeapReleased
	res["heap_objects"] = m.HeapObjects
	// low-level fixed-size struct alloctor
	res["stack_inuse"] = m.StackInuse
	res["stack_sys"] = m.StackSys
	res["mspan_inuse"] = m.MSpanInuse
	res["mspan_sys"] = m.MSpanSys
	res["mcache_inuse"] = m.MCacheInuse
	res["mcache_sys"] = m.MCacheSys
	res["buckhash_sys"] = m.BuckHashSys
	// GC
	res["next_gc"] = m.NextGC
	res["last_gc"] = m.LastGC
	res["pause_total_ns"] = m.PauseTotalNs
	res["pause_ns"] = m.PauseNs
	res["num_gc"] = m.NumGC
	res["enable_gc"] = m.EnableGC
	res["debug_gc"] = m.DebugGC
	res["by_size"] = m.BySize
	return jsonRes(res)
}

// golang stats  // golang统计
func GoStats() []byte {
	res := map[string]interface{}{}
	res["compiler"] = runtime.Compiler
	res["arch"] = runtime.GOARCH
	res["os"] = runtime.GOOS
	res["max_procs"] = runtime.GOMAXPROCS(-1)
	res["root"] = runtime.GOROOT()
	res["cgo_call"] = runtime.NumCgoCall()
	res["goroutine_num"] = runtime.NumGoroutine()
	res["version"] = runtime.Version()
	return jsonRes(res)
}

// server stats  //服务器统计
func ServerStats() []byte {
	res := map[string]interface{}{}
	res["uptime"] = time.Now().UnixNano() - startTime
	hostname, _ := os.Hostname()
	res["hostname"] = hostname
	wd, _ := os.Getwd()
	res["wd"] = wd
	res["ppid"] = os.Getppid()
	res["pid"] = os.Getpid()
	res["pagesize"] = os.Getpagesize()
	if usr, err := user.Current(); err != nil {
		log.Error("user.Current() error(%v)", err)
		res["group"] = ""
		res["user"] = ""
	} else {
		res["group"] = usr.Gid
		res["user"] = usr.Uid
	}
	return jsonRes(res)
}

// configuration info  // 配置文件信息
func ConfigInfo() []byte {
	byteJson, err := json.MarshalIndent(Conf, "", "    ")
	if err != nil {
		log.Error("json.MarshalIndent(\"%v\", \"\", \"    \") error(%v)", Conf, err)
		return nil
	}
	return byteJson
}

// jsonRes format the output  // 格式化输出
func jsonRes(res map[string]interface{}) []byte {
	byteJson, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Error("json.MarshalIndent(\"%v\", \"\", \"    \") error(%v)", res, err)
		return nil
	}
	return byteJson
}

func ChInfoStat(key string) []byte {
	res := map[string]interface{}{}
	if ch, err := UserChannel.Get(key, false); err == nil {
		if sch, ok := ch.(*SeqChannel); ok {
			res["channel"] = map[string]interface{}{"conn": sch.conn.Len()}
		} else {
			return nil
		}
	} else {
		return nil
	}
	return jsonRes(res)
}

// StatHandle get stat info by http  //通过http得到统计信息
func StatHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	params := r.URL.Query()
	types := params.Get("type")
	res := []byte{}
	switch types {
	case "memory":
		res = MemStats()
	case "server":
		res = ServerStats()
	case "golang":
		res = GoStats()
	case "config":
		res = ConfigInfo()
	case "channel":
		key := params.Get("key")
		if key == "" {
			res = ChStat.Stat()
		} else {
			res = ChInfoStat(key)
		}
	case "message":
		res = MsgStat.Stat()
	case "connection":
		res = ConnStat.Stat()
	default:
		http.Error(w, "Not Found", 404)
	}
	if res != nil {
		if _, err := w.Write(res); err != nil {
			log.Error("w.Write(\"%s\") error(%v)", string(res), err)
		}
	}
}
