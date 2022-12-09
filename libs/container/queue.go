package container

import (
	"container/list"
)

type Queue struct {
	v     *list.List
	Count uint32
}

func NewQueue() *Queue {
	return &Queue{
		v:     list.New(),
		Count: 0,
	}
}

func (q *Queue) Push(v interface{}) {
	q.v.PushBack(v)
	q.Count++
}

func (q *Queue) Pop() interface{} {
	front := q.v.Front()
	if front == nil {
		return nil
	}

	q.Count--

	return q.v.Remove(front)
}
