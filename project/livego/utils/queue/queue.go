package queue

import (
	"github.com/livego/av"
	"sync"
)

// Queue is a basic FIFO queue for Messages.
// Queue是消息的基本FIFO队列。
type Queue struct {
	maxSize int

	list  []*av.Packet
	mutex sync.Mutex
}

// NewQueue returns a new Queue. If maxSize is greater than zero the queue will not grow more than the defined size.
// NewQueue 返回一个新的队列。 如果maxSize大于0，队列将不会超过定义的大小。
func NewQueue(maxSize int) *Queue {
	return &Queue{
		maxSize: maxSize,
	}
}

// Push adds a message to the queue.
// Push 添加一个笑消息到队列中
func (q *Queue) Push(msg *av.Packet) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.list) == q.maxSize {
		q.pop()
	}

	q.list = append(q.list, msg)
}

// Pop removes and returns a message from the queue in first to last order.
// Pop 从队列中移除并返回一条消息，在第一个到最后一个顺序。
func (q *Queue) Pop() *av.Packet {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.list) == 0 {
		return nil
	}

	return q.pop()
}

func (q *Queue) pop() *av.Packet {
	x := len(q.list) - 1
	msg := q.list[x]
	q.list = q.list[:x]
	return msg
}

// Len returns the length of the queue.
// Len 返回队列的长度
func (q *Queue) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.list)
}

// All returns and removes all messages from the queue.
// All 返回并删除队列中的所有消息。
func (q *Queue) All() []*av.Packet {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	cache := q.list
	q.list = nil
	return cache
}
