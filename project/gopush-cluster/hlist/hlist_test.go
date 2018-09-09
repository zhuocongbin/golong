

package hlist

import (
	"fmt"
	"testing"
)

func TestHlist(t *testing.T) {
	l := New()
	first := l.Front()
	if first != nil {
		t.Error("first != nil")
	}
	l.PushFront(1)
	l.PushFront(2)
	first = l.Front()
	if first == nil {
		t.Error("first == nil")
	}
	if i, ok := first.Value.(int); !ok {
		t.Error("first.Value assection failed")
	} else {
		if i != 2 {
			t.Errorf("i value error: %d", i)
		}
	}
	if next := first.Next(); next == nil {
		t.Error("next == nil")
	} else {
		if i, ok := next.Value.(int); !ok {
			t.Error("next.Value assection failed")
		} else {
			if i != 1 {
				t.Errorf("i value error: %d", i)
			}
		}
	}
	if l.Len() != 2 {
		t.Errorf("length error")
	}
	l.PushFront(3)
	l.PushFront(4)
	l.PushFront(5)
	l.PushFront(6)
	first = l.Front()
	if l.Len() != 6 {
		t.Errorf("length error")
	}
	for e := l.Front(); e != nil; e = e.Next() {
		if i, ok := e.Value.(int); !ok {
			t.Error("e.Value assection failed")
		} else {
			fmt.Println(i)
		}
	}
	fmt.Println("------")
	if i, ok := l.Remove(first).(int); !ok {
		t.Error("first.Value assection failed")
	} else {
		if i != 6 {
			t.Errorf("i value error: %d", i)
		}
	}
	if l.Len() != 5 {
		t.Errorf("length error")
	}
	for e := l.Front(); e != nil; e = e.Next() {
		if i, ok := e.Value.(int); !ok {
			t.Error("e.Value assection failed")
		} else {
			fmt.Println(i)
		}
	}
	second := l.Front()    // 5
	thrid := second.Next() // 4
	fourth := thrid.Next() // 3
	fifth := fourth.Next() // 2
	sixth := fifth.Next()  // 1
	l.Remove(second)
	l.Remove(thrid)
	l.Remove(fourth)
	l.Remove(fifth)
	l.Remove(sixth)
	if l.Len() != 0 {
		t.Errorf("length error")
	}
	first = l.Front()
	if first != nil {
		t.Error("first != nil")
	}
	for e := l.Front(); e != nil; e = e.Next() {
		if i, ok := e.Value.(int); !ok {
			t.Error("e.Value assection failed")
		} else {
			fmt.Println(i)
		}
	}

	e := l.PushFront(7)
	if i, ok := e.Value.(int); !ok {
		t.Error("e.Value assection failed")
	} else {
		if i != 7 {
			t.Error("i value error: %d", i)
		}
	}
}
