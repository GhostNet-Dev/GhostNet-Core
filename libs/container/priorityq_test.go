package container

import (
	"container/heap"
	"fmt"
	"math/rand"
	"testing"
)

func TestPriorityQ(t *testing.T) {
	pq := &PriorityQueue{}
	heap.Init(pq)

	for i := 0; i < 10; i++ {
		item := &Item{priority: rand.Uint64()}
		heap.Push(pq, item)
		pq.Update(item, fmt.Sprint(item.index), item.priority)
	}

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*Item)
		fmt.Println(item.priority, "---", item.value)
	}
}
