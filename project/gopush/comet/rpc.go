
package main

import (
	log "github.com/alecthomas/log4go"
	"errors"
	"github.com/Terry-Mao/gopush-cluster/id"
	myrpc "github.com/Terry-Mao/gopush-cluster/rpc"
	"net"
	"net/rpc"
	"sync"
)

var (
	ErrMigrate = errors.New("migrate nodes don't include self")
)

// StartRPC start rpc listen.
// 开始rpc监听
func StartRPC() error {
	c := &CometRPC{}
	rpc.Register(c)
	for _, bind := range Conf.RPCBind {
		log.Info("start listen rpc addr: \"%s\"", bind)
		go rpcListen(bind)
	}

	return nil
}

func rpcListen(bind string) {
	l, err := net.Listen("tcp", bind)
	if err != nil {
		log.Error("net.Listen(\"tcp\", \"%s\") error(%v)", bind, err)
		panic(err)
	}
	// if process exit, then close the rpc bind
	// 如果进程退出，则关闭rpc绑定
	defer func() {
		log.Info("rpc addr: \"%s\" close", bind)
		if err := l.Close(); err != nil {
			log.Error("listener.Close() error(%v)", err)
		}
	}()
	rpc.Accept(l)
}

// Channel RPC
type CometRPC struct {
}

// New expored a method for creating new channel.
// 介绍了一种创建新通道的方法。
func (c *CometRPC) New(args *myrpc.CometNewArgs, ret *int) error {
	if args == nil || args.Key == "" {
		return myrpc.ErrParam
	}
	// create a new channel for the user
	// 为用户创建一个新通道
	ch, _, err := UserChannel.New(args.Key)
	if err != nil {
		log.Error("UserChannel.New(\"%s\") error(%v)", args.Key, err)
		return err
	}
	if err = ch.AddToken(args.Key, args.Token); err != nil {
		log.Error("ch.AddToken(\"%s\", \"%s\") error(%v)", args.Key, args.Token)
		return err
	}
	return nil
}

// Close expored a method for closing new channel.
// 关闭出口关闭新通道的方法。
func (c *CometRPC) Close(key string, ret *int) error {
	if key == "" {
		return myrpc.ErrParam
	}
	// close the channle for the user // 为用户关闭通道
	ch, err := UserChannel.Delete(key)
	if err != nil {
		log.Error("UserChannel.Delete(\"%s\") error(%v)", key, err)
		return err
	}
	// ignore channel close error, only log a warnning
	// 忽略通道关闭错误，只记录警告
	if err := ch.Close(); err != nil {
		log.Error("ch.Close() error(%v)", err)
		return err
	}
	return nil
}

// PushPrivate expored a method for publishing a user private message for the channel.
// PushPrivate导出用于为通道发布用户私有消息的方法。
// if it`s going failed then it`ll return an error
// 如果失败了，它将返回一个错误
func (c *CometRPC) PushPrivate(args *myrpc.CometPushPrivateArgs, ret *int) error {
	if args == nil || args.Key == "" {
		return myrpc.ErrParam
	}
	// get a user channel //获得用户通道
	ch, _, err := UserChannel.New(args.Key)
	if err != nil {
		log.Error("UserChannel.New(\"%s\") error(%v)", args.Key, err)
		return err
	}
	// use the channel push message
	// 使用通道推送消息
	m := &myrpc.Message{Msg: args.Msg}
	if err = ch.PushMsg(args.Key, m, args.Expire); err != nil {
		log.Error("ch.PushMsg(\"%s\", \"%v\") error(%v)", args.Key, m, err)
		return err
	}
	return nil
}

// batchChannel is use for PushPrivates.
// batchChannel是为PushPrivates服务的。
type batchChannel struct {
	Keys []string
	Chs  map[string]Channel
}

// PushPrivates expored a method for publishing a user multiple private message for the channel.
// PushPrivates导出了为通道发布用户多个私有消息的方法。
// because of it`s going asynchronously in this method, so it won`t return an error to caller.
// 因为在这个方法中它是异步进行的，所以它不会向调用者返回错误。
func (c *CometRPC) PushPrivates(args *myrpc.CometPushPrivatesArgs, rw *myrpc.CometPushPrivatesResp) error {
	if args == nil {
		return myrpc.ErrParam
	}
	bucketMap := make(map[*ChannelBucket]*batchChannel, Conf.ChannelBucket)
	for _, key := range args.Keys {
		// get channel
		ch, bp, err := UserChannel.New(key)
		if err != nil {
			log.Error("UserChannel.New(\"%s\") error(%v)", key, err)
			// log failed keys.
			rw.FKeys = append(rw.FKeys, key)
			continue
		}
		if bucket, ok := bucketMap[bp]; !ok {
			bucketMap[bp] = &batchChannel{
				Keys: []string{key},
				Chs:  map[string]Channel{key: ch},
			}
		} else {
			// ignore duplicate key
			if _, ok := bucket.Chs[key]; !ok {
				bucket.Chs[key] = ch
				bucket.Keys = append(bucket.Keys, key)
			}
		}
	}
	// every bucket start a goroutine, return till all bucket gorouint finish
	// 每个桶开始一个goroutine，返回直到所有桶gorouint完成
	wg := &sync.WaitGroup{}
	wg.Add(len(bucketMap))
	// stored every gorouint failed keys
	// 存储每个gorouint失败键
	fKeysList := make([][]string, len(bucketMap))
	ti := 0
	for tb, tm := range bucketMap {
		go func(b *ChannelBucket, m *batchChannel, i int) {
			defer wg.Done()
			c := myrpc.MessageRPC.Get()
			if c == nil {
				// static slice is thread-safe
				// log all keys
				fKeysList[i] = m.Keys
				log.Debug("fkeys len:%d", len(m.Keys))
				return
			}
			b.Lock()
			defer b.Unlock()
			timeId := id.Get()
			msg := &myrpc.Message{Msg: args.Msg, MsgId: timeId}
			// private message need persistence
			// 私有消息需要持久性
			// if message expired no need persistence, only send online message rewrite message id
			// 如果消息过期不需要持久性，只发送在线消息重写消息id
			resp := &myrpc.MessageSavePrivatesResp{}
			if args.Expire > 0 {
				args := &myrpc.MessageSavePrivatesArgs{Keys: m.Keys, Msg: args.Msg, MsgId: timeId, Expire: args.Expire}
				if err := c.Call(myrpc.MessageServiceSavePrivates, args, resp); err != nil {
					log.Error("%s(\"%v\", \"%v\", &ret) error(%v)", myrpc.MessageServiceSavePrivates, m.Keys, args, err)
					// static slice is thread-safe
					// 静态片是线程安全的
					fKeysList[i] = m.Keys
					log.Debug("fkeys len:%d", len(m.Keys))
					return
				}
				fKeysList[i] = resp.FKeys
				log.Debug("fkeys len:%d", len(resp.FKeys))
			}
			// delete the failed keys
			// 删除失败的键
			for _, fk := range resp.FKeys {
				delete(m.Chs, fk)
			}
			// get all channels from batchChannel chs.
			// 从批量通道chs获取所有通道。
			for key, ch := range m.Chs {
				if err := ch.WriteMsg(key, msg); err != nil {
					// ignore online push error, cause offline msg succeed
					// 忽略在线推送错误，导致离线msg成功
					log.Error("ch.WriteMsg(\"%s\", \"%s\") error(%v)", key, string(msg.Msg), err)
					continue
				}
			}
		}(tb, tm, ti)
		ti++
	}
	wg.Wait()
	// merge all failed keys
	// 合并所有失败键
	for _, k := range fKeysList {
		rw.FKeys = append(rw.FKeys, k...)
	}
	return nil
}

// Migrate update the inner hashring and node info.
// 更新内部hashring和节点信息。
func (c *CometRPC) Migrate(args *myrpc.CometMigrateArgs, ret *int) error {
	return UserChannel.Migrate(args.Nodes)
}

// Ping check health.
func (c *CometRPC) Ping(args int, ret *int) error {
	log.Debug("ping ok")
	return nil
}
