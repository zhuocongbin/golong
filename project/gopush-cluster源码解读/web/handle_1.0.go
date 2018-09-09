

package main

import (
	log "github.com/alecthomas/log4go"
	myrpc "github.com/Terry-Mao/gopush-cluster/rpc"
	"net/http"
	"strconv"
	"time"
)

// GetServer handle for server get
func GetServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	params := r.URL.Query()
	key := params.Get("k")
	callback := params.Get("cb")
	protoStr := params.Get("p")
	res := map[string]interface{}{"ret": OK}
	defer retWrite(w, r, res, callback, time.Now())
	if key == "" {
		res["ret"] = ParamErr
		return
	}
	// Match a push-server with the value computed through ketama algorithm
	node := myrpc.GetComet(key)
	if node == nil {
		res["ret"] = NotFoundServer
		return
	}
	addrs, ret := getProtoAddr(node, protoStr)
	if ret != OK {
		res["ret"] = ret
		return
	}
	res["data"] = map[string]interface{}{"server": addrs[0]}
	return
}

// GetServer2 handle for server get.
func GetServer2(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	params := r.URL.Query()
	key := params.Get("k")
	callback := params.Get("cb")
	protoStr := params.Get("p")
	res := map[string]interface{}{"ret": OK}
	defer retWrite(w, r, res, callback, time.Now())
	if key == "" {
		res["ret"] = ParamErr
		return
	}
	// Match a push-server with the value computed through ketama algorithm
	node := myrpc.GetComet(key)
	if node == nil {
		res["ret"] = NotFoundServer
		return
	}
	addrs, ret := getProtoAddr(node, protoStr)
	if ret != OK {
		res["ret"] = ret
		return
	}
	// give client a addr list, client do the best choice
	res["data"] = map[string]interface{}{"server": addrs}
	return
}

// GetOfflineMsg get offline mesage http handler.
func GetOfflineMsg(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	params := r.URL.Query()
	key := params.Get("k")
	midStr := params.Get("m")
	callback := params.Get("cb")
	res := map[string]interface{}{"ret": OK}
	defer retWrite(w, r, res, callback, time.Now())
	if key == "" || midStr == "" {
		res["ret"] = ParamErr
		return
	}
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		res["ret"] = ParamErr
		log.Error("strconv.ParseInt(\"%s\", 10, 64) error(%v)", midStr, err)
		return
	}
	// RPC get offline messages
	reply := &myrpc.MessageGetResp{}
	args := &myrpc.MessageGetPrivateArgs{MsgId: mid, Key: key}
	client := myrpc.MessageRPC.Get()
	if client == nil {
		log.Error("no message node found")
		res["ret"] = InternalErr
		return
	}
	if err := client.Call(myrpc.MessageServiceGetPrivate, args, reply); err != nil {
		log.Error("myrpc.MessageRPC.Call(\"%s\", \"%v\", reply) error(%v)", myrpc.MessageServiceGetPrivate, args, err)
		res["ret"] = InternalErr
		return
	}
	if len(reply.Msgs) == 0 {
		return
	}
	res["data"] = map[string]interface{}{"msgs": reply.Msgs}
	return
}

// GetTime get server time http handler.
func GetTime(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	params := r.URL.Query()
	callback := params.Get("cb")
	res := map[string]interface{}{"ret": OK}
	now := time.Now()
	defer retWrite(w, r, res, callback, now)
	res["data"] = map[string]interface{}{"timeid": now.UnixNano() / 100}
	return
}
