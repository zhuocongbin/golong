
package main

import (
	log "github.com/alecthomas/log4go"
	"errors"
	"github.com/Terry-Mao/gopush-cluster/hlist"
	"github.com/Terry-Mao/gopush-cluster/id"
	myrpc "github.com/Terry-Mao/gopush-cluster/rpc"
	"sync"
)

var (
	ErrMessageSave   = errors.New("Message set failed")
	ErrMessageGet    = errors.New("Message get failed")
	ErrMessageRPC    = errors.New("Message RPC not init")
	ErrAssectionConn = errors.New("Assection type Connection failed")
)

// Sequence Channel struct.
type SeqChannel struct {
	// Mutex
	mutex *sync.Mutex
	// client conn double linked-list
	conn *hlist.Hlist
	// Remove time id or lazy New
	// timeID *id.TimeID
	// token
	token *Token
}

// New a user seq stored message channel. //用户seq存储的消息通道。
func NewSeqChannel() *SeqChannel {
	ch := &SeqChannel{
		mutex: &sync.Mutex{},
		conn:  hlist.New(),
		//timeID: id.NewTimeID(),
		token: nil,
	}
	// save memory
	if Conf.Auth {
		ch.token = NewToken()
	}
	return ch
}

// AddToken implements the Channel AddToken method. //实现通道AddToken方法。
func (c *SeqChannel) AddToken(key, token string) error {
	if !Conf.Auth {
		return nil
	}
	c.mutex.Lock()
	if err := c.token.Add(token); err != nil {
		c.mutex.Unlock()
		log.Error("user_key:\"%s\" c.token.Add(\"%s\") error(%v)", key, token, err)
		return err
	}
	c.mutex.Unlock()
	return nil
}

// AuthToken implements the Channel AuthToken method. // 实现Channel AuthToken方法。
func (c *SeqChannel) AuthToken(key, token string) bool {
	if !Conf.Auth {
		return true
	}
	c.mutex.Lock()
	if err := c.token.Auth(token); err != nil {
		c.mutex.Unlock()
		log.Error("user_key:\"%s\" c.token.Auth(\"%s\") error(%v)", key, token, err)
		return false
	}
	c.mutex.Unlock()
	return true
}

// WriteMsg implements the Channel WriteMsg method. //实现频道WriteMsg方法。
func (c *SeqChannel) WriteMsg(key string, m *myrpc.Message) (err error) {
	c.mutex.Lock()
	err = c.writeMsg(key, m)
	c.mutex.Unlock()
	return
}

// writeMsg write msg to conn.
func (c *SeqChannel) writeMsg(key string, m *myrpc.Message) (err error) {
	var (
		oldMsg, msg, sendMsg []byte
	)
	// push message
	for e := c.conn.Front(); e != nil; e = e.Next() {
		conn, _ := e.Value.(*Connection)
		// if version empty then use old protocol
		// 如果版本为空，则使用旧协议
		if conn.Version == "" {
			if oldMsg == nil {
				if oldMsg, err = m.OldBytes(); err != nil {
					return
				}
			}
			sendMsg = oldMsg
		} else {
			if msg == nil {
				if msg, err = m.Bytes(); err != nil {
					return
				}
			}
			sendMsg = msg
		}
		// TODO use goroutine
		conn.Write(key, sendMsg)
	}
	return
}

// PushMsg implements the Channel PushMsg method. // 实现Channel PushMsg方法。
func (c *SeqChannel) PushMsg(key string, m *myrpc.Message, expire uint) (err error) {
	client := myrpc.MessageRPC.Get()
	if client == nil {
		return ErrMessageRPC
	}
	c.mutex.Lock()
	// private message need persistence // 私有消息需要持久性
	// if message expired no need persistence, only send online message rewrite message id
	// 如果消息过期不需要持久性，只发送在线消息重写消息id
	//m.MsgId = c.timeID.ID()
	m.MsgId = id.Get()
	if m.GroupId != myrpc.PublicGroupId && expire > 0 {
		args := &myrpc.MessageSavePrivateArgs{Key: key, Msg: m.Msg, MsgId: m.MsgId, Expire: expire}
		ret := 0
		if err = client.Call(myrpc.MessageServiceSavePrivate, args, &ret); err != nil {
			c.mutex.Unlock()
			log.Error("%s(\"%s\", \"%v\", &ret) error(%v)", myrpc.MessageServiceSavePrivate, key, args, err)
			return
		}
	}
	// push message
	if err = c.writeMsg(key, m); err != nil {
		c.mutex.Unlock()
		log.Error("c.WriteMsg(\"%s\", m) error(%v)", key, err)
		return
	}
	c.mutex.Unlock()
	return
}

// AddConn implements the Channel AddConn method. // 实现通道AddConn方法。
func (c *SeqChannel) AddConn(key string, conn *Connection) (*hlist.Element, error) {
	c.mutex.Lock()
	if c.conn.Len()+1 > Conf.MaxSubscriberPerChannel {
		c.mutex.Unlock()
		log.Error("user_key:\"%s\" exceed conn", key)
		return nil, ErrMaxConn
	}
	// send first heartbeat to tell client service is ready for accept heartbeat
	// 发送第一个心跳告诉客户端服务准备接受心跳
	if _, err := conn.Conn.Write(HeartbeatReply); err != nil {
		c.mutex.Unlock()
		log.Error("user_key:\"%s\" write first heartbeat to client error(%v)", key, err)
		return nil, err
	}
	// add conn
	conn.Buf = make(chan []byte, Conf.MsgBufNum)
	conn.HandleWrite(key)
	e := c.conn.PushFront(conn)
	c.mutex.Unlock()
	ConnStat.IncrAdd()
	log.Info("user_key:\"%s\" add conn = %d", key, c.conn.Len())
	return e, nil
}

// RemoveConn implements the Channel RemoveConn method.
// 实现通道RemoveConn方法。
func (c *SeqChannel) RemoveConn(key string, e *hlist.Element) error {
	c.mutex.Lock()
	tmp := c.conn.Remove(e)
	c.mutex.Unlock()
	conn, ok := tmp.(*Connection)
	if !ok {
		return ErrAssectionConn
	}
	close(conn.Buf)
	ConnStat.IncrRemove()
	log.Info("user_key:\"%s\" remove conn = %d", key, c.conn.Len())
	return nil
}

// Close implements the Channel Close method. //实现通道关闭方法。
func (c *SeqChannel) Close() error {
	c.mutex.Lock()
	for e := c.conn.Front(); e != nil; e = e.Next() {
		if conn, ok := e.Value.(*Connection); !ok {
			c.mutex.Unlock()
			return ErrAssectionConn
		} else {
			if err := conn.Conn.Close(); err != nil {
				// ignore close error //忽略关闭错误
				log.Warn("conn.Close() error(%v)", err)
			}
		}
	}
	c.mutex.Unlock()
	return nil
}
