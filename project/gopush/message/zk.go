
package main

import (
	log "github.com/alecthomas/log4go"
	"encoding/json"
	"github.com/Terry-Mao/gopush-cluster/rpc"
	myzk "github.com/Terry-Mao/gopush-cluster/zk"
	"github.com/samuel/go-zookeeper/zk"
)

// InitZK create zookeeper root path, and register a temp node.
func InitZK() (*zk.Conn, error) {
	conn, err := myzk.Connect(Conf.ZookeeperAddr, Conf.ZookeeperTimeout)
	if err != nil {
		log.Error("zk.Connect() error(%v)", err)
		return nil, err
	}
	if err = myzk.Create(conn, Conf.ZookeeperPath); err != nil {
		log.Error("zk.Create() error(%v)", err)
		return conn, err
	}
	nodeInfo := rpc.MessageNodeInfo{}
	nodeInfo.Rpc = Conf.RPCBind
	nodeInfo.Weight = Conf.NodeWeight
	data, err := json.Marshal(nodeInfo)
	if err != nil {
		log.Error("json.Marshal(() error(%v)", err)
		return conn, err
	}
	log.Debug("zk data: \"%s\"", string(data))
	// tcp, websocket and rpc bind address store in the zk
	if err = myzk.RegisterTemp(conn, Conf.ZookeeperPath, data); err != nil {
		log.Error("zk.RegisterTemp() error(%v)", err)
		return conn, err
	}
	return conn, nil
}
