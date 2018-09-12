

package main

import (
	log "github.com/alecthomas/log4go"
	"errors"
	"github.com/Terry-Mao/gopush-cluster/hash"
	"github.com/Terry-Mao/gopush-cluster/hlist"
	"github.com/Terry-Mao/gopush-cluster/ketama"
	myrpc "github.com/Terry-Mao/gopush-cluster/rpc"
	"sync"
)

var (
	//通道不存在错误
	ErrChannelNotExist = errors.New("Channle not exist")
	//未知连接协议错误
	ErrConnProto       = errors.New("Unknown connection protocol")
	ErrChannelKey      = errors.New("Key not belong this comet")
	UserChannel        *ChannelList
	CometRing          *ketama.HashRing
	nodeWeightMap      = map[string]int{}
)

// The subscriber interface.//订阅者接口
type Channel interface {
	// WriteMsg push a message to the subscriber.//向订阅服务器推送消息。
	WriteMsg(key string, m *myrpc.Message) error
	// PushMsg push a message to the subscriber.
	PushMsg(key string, m *myrpc.Message, expire uint) error
	// Add a token for one subscriber //为一个订阅服务器添加token
	// The request token not equal the subscriber token will return errors.
	// 请求令牌符不等于订阅者令牌符将返回错误。
	AddToken(key, token string) error
	// Auth auth the access token.
	// The request token not match the subscriber token will return errors.
	// 请求令牌与订阅方令牌不匹配将返回错误。
	AuthToken(key, token string) bool
	// AddConn add a connection for the subscriber.
	// Exceed the max number of subscribers per key will return errors.
	// 超过每个键的最大订户数将返回错误。
	AddConn(key string, conn *Connection) (*hlist.Element, error)
	// RemoveConn remove a connection for the  subscriber.
	RemoveConn(key string, e *hlist.Element) error
	// Expire expire the channle and clean data. //终止通道和干净数据。
	Close() error
}

// Channel bucket.
type ChannelBucket struct {
	Data  map[string]Channel
	mutex *sync.Mutex
}

// Channel list.
type ChannelList struct {
	Channels []*ChannelBucket
}

// Lock lock the bucket mutex.
func (c *ChannelBucket) Lock() {
	c.mutex.Lock()
}

// Unlock unlock the bucket mutex.
func (c *ChannelBucket) Unlock() {
	c.mutex.Unlock()
}

// NewChannelList create a new channel bucket set.
func NewChannelList() *ChannelList {
	l := &ChannelList{Channels: []*ChannelBucket{}}
	// split hashmap to many bucket
	log.Debug("create %d ChannelBucket", Conf.ChannelBucket)
	for i := 0; i < Conf.ChannelBucket; i++ {
		c := &ChannelBucket{
			Data:  map[string]Channel{},
			mutex: &sync.Mutex{},
		}
		l.Channels = append(l.Channels, c)
	}
	return l
}

// Count get the bucket total channel count. //获取桶总通道数。
func (l *ChannelList) Count() int {
	c := 0
	for i := 0; i < Conf.ChannelBucket; i++ {
		c += len(l.Channels[i].Data)
	}
	return c
}

// bucket return a channelBucket use murmurhash3.
func (l *ChannelList) Bucket(key string) *ChannelBucket {
	h := hash.NewMurmur3C()
	h.Write([]byte(key))
	idx := uint(h.Sum32()) & uint(Conf.ChannelBucket-1)
	log.Debug("user_key:\"%s\" hit channel bucket index:%d", key, idx)
	return l.Channels[idx]
}

// validate check the key is belong to this comet.
func (l *ChannelList) validate(key string) error {
	if len(nodeWeightMap) == 0 {
		log.Debug("no node found")
		return ErrChannelKey
	}
	node := CometRing.Hash(key)
	log.Debug("match node:%s hash node:%s", Conf.ZookeeperCometNode, node)
	if Conf.ZookeeperCometNode != node {
		log.Warn("user_key:\"%s\" node:%s not match this node:%s", key, node, Conf.ZookeeperCometNode)
		return ErrChannelKey
	}
	return nil
}

// New create a user channle. //创建一个用户通道。
func (l *ChannelList) New(key string) (Channel, *ChannelBucket, error) {
	// validate
	if err := l.validate(key); err != nil {
		return nil, nil, err
	}
	// get a channel bucket
	b := l.Bucket(key)
	b.Lock()
	if c, ok := b.Data[key]; ok {
		b.Unlock()
		ChStat.IncrAccess()
		log.Info("user_key:\"%s\" refresh channel bucket expire time", key)
		return c, b, nil
	} else {
		c = NewSeqChannel()
		b.Data[key] = c
		b.Unlock()
		ChStat.IncrCreate()
		log.Info("user_key:\"%s\" create a new channel", key)
		return c, b, nil
	}
}

// Get a user channel from ChannleList.
// 从通道列表中获得一个用户通道。
func (l *ChannelList) Get(key string, newOne bool) (Channel, error) {
	// validate
	if err := l.validate(key); err != nil {
		return nil, err
	}
	// get a channel bucket
	b := l.Bucket(key)
	b.Lock()
	if c, ok := b.Data[key]; !ok {
		if !Conf.Auth && newOne {
			c = NewSeqChannel()
			b.Data[key] = c
			b.Unlock()
			ChStat.IncrCreate()
			log.Info("user_key:\"%s\" create a new channel", key)
			return c, nil
		} else {
			b.Unlock()
			log.Warn("user_key:\"%s\" channle not exists", key)
			return nil, ErrChannelNotExist
		}
	} else {
		b.Unlock()
		ChStat.IncrAccess()
		log.Info("user_key:\"%s\" refresh channel bucket expire time", key)
		return c, nil
	}
}

// Delete a user channel from ChannleList.
// 从通道列表中删除一个用户通道。
func (l *ChannelList) Delete(key string) (Channel, error) {
	// get a channel bucket
	b := l.Bucket(key)
	b.Lock()
	if c, ok := b.Data[key]; !ok {
		b.Unlock()
		log.Warn("user_key:\"%s\" delete channle not exists", key)
		return nil, ErrChannelNotExist
	} else {
		delete(b.Data, key)
		b.Unlock()
		ChStat.IncrDelete()
		log.Info("user_key:\"%s\" delete channel", key)
		return c, nil
	}
}

// Close close all channel.//关闭所有通道。
func (l *ChannelList) Close() {
	log.Info("channel close")
	chs := make([]Channel, 0, l.Count())
	for _, c := range l.Channels {
		c.Lock()
		for _, c := range c.Data {
			chs = append(chs, c)
		}
		c.Unlock()
	}
	// close all channels
	for _, c := range chs {
		if err := c.Close(); err != nil {
			log.Error("c.Close() error(%v)", err)
		}
	}
}

// Migrate migrate portion of connections which don't belong to this comet.
// 迁移不属于此comet的连接部分。
func (l *ChannelList) Migrate(nw map[string]int) (err error) {
	migrate := false
	// check new/update node
	for k, v := range nw {
		weight, ok := nodeWeightMap[k]
		// not found or weight change
		if !ok || weight != v {
			migrate = true
			break
		}
	}
	// check del node
	if !migrate {
		for k, _ := range nodeWeightMap {
			// node deleted
			if _, ok := nw[k]; !ok {
				migrate = true
				break
			}
		}
	}
	if !migrate {
		return
	}
	// init ketama
	ring := ketama.NewRing(ketama.Base)
	for node, weight := range nw {
		ring.AddNode(node, weight)
	}
	ring.Bake()
	// atomic update
	nodeWeightMap = nw
	CometRing = ring
	// get all the channel lock
	channels := []Channel{}
	for i, c := range l.Channels {
		c.Lock()
		for k, v := range c.Data {
			hn := ring.Hash(k)
			if hn != Conf.ZookeeperCometNode {
				channels = append(channels, v)
				delete(c.Data, k)
				log.Debug("migrate delete channel key \"%s\"", k)
			}
		}
		c.Unlock()
		log.Debug("migrate channel bucket:%d finished", i)
	}
	// close all the migrate channels //关闭所有迁移通道
	log.Info("close all the migrate channels")
	for _, channel := range channels {
		if err := channel.Close(); err != nil {
			log.Error("channel.Close() error(%v)", err)
			continue
		}
	}
	log.Info("close all the migrate channels finished")
	return
}
