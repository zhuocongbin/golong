

package heap

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	h := NewMinheap(2)
	// add
	fmt.Println("------------- add ----------------")
	h.Add(&Element{Key: 20, Value: 1})
	fmt.Printf("Size: %d, Max: %d\n", h.Size(), h.Max())
	for i := 0; i < h.Size(); i++ {
		fmt.Printf("Key: %d, Value: %d\n", h.items[i].Key, h.items[i].Value)
	}
	h.Add(&Element{Key: 15, Value: 1})
	fmt.Printf("Size: %d, Max: %d\n", h.Size(), h.Max())
	for i := 0; i < h.Size(); i++ {
		fmt.Printf("Key: %d, Value: %d\n", h.items[i].Key, h.items[i].Value)
	}
	h.Add(&Element{Key: 2, Value: 1})
	fmt.Printf("Size: %d, Max: %d\n", h.Size(), h.Max())
	for i := 0; i < h.Size(); i++ {
		fmt.Printf("Key: %d, Value: %d\n", h.items[i].Key, h.items[i].Value)
	}
	h.Add(&Element{Key: 14, Value: 1})
	fmt.Printf("Size: %d, Max: %d\n", h.Size(), h.Max())
	for i := 0; i < h.Size(); i++ {
		fmt.Printf("Key: %d, Value: %d\n", h.items[i].Key, h.items[i].Value)
	}
	h.Add(&Element{Key: 10, Value: 1})
	fmt.Printf("Size: %d, Max: %d\n", h.Size(), h.Max())
	for i := 0; i < h.Size(); i++ {
		fmt.Printf("Key: %d, Value: %d\n", h.items[i].Key, h.items[i].Value)
	}
	// poll
	fmt.Println("------------- poll ----------------")
	e := h.Poll()
	fmt.Printf("FETCH Key: %d, Value: %d\n", e.Key, e.Value)
	fmt.Printf("Size: %d, Max: %d\n", h.Size(), h.Max())
	for i := 0; i < h.Size(); i++ {
		fmt.Printf("Key: %d, Value: %d\n", h.items[i].Key, h.items[i].Value)
	}
	e = h.Poll()
	fmt.Printf("FETCH Key: %d, Value: %d\n", e.Key, e.Value)
	fmt.Printf("Size: %d, Max: %d\n", h.Size(), h.Max())
	for i := 0; i < h.Size(); i++ {
		fmt.Printf("Key: %d, Value: %d\n", h.items[i].Key, h.items[i].Value)
	}
	e = h.Poll()
	fmt.Printf("FETCH Key: %d, Value: %d\n", e.Key, e.Value)
	fmt.Printf("Size: %d, Max: %d\n", h.Size(), h.Max())
	for i := 0; i < h.Size(); i++ {
		fmt.Printf("Key: %d, Value: %d\n", h.items[i].Key, h.items[i].Value)
	}
	e = h.Poll()
	fmt.Printf("FETCH Key: %d, Value: %d\n", e.Key, e.Value)
	fmt.Printf("Size: %d, Max: %d\n", h.Size(), h.Max())
	for i := 0; i < h.Size(); i++ {
		fmt.Printf("Key: %d, Value: %d\n", h.items[i].Key, h.items[i].Value)
	}
	e = h.Poll()
	fmt.Printf("FETCH Key: %d, Value: %d\n", e.Key, e.Value)
	fmt.Printf("Size: %d, Max: %d\n", h.Size(), h.Max())
	for i := 0; i < h.Size(); i++ {
		fmt.Printf("Key: %d, Value: %d\n", h.items[i].Key, h.items[i].Value)
	}
}
