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

func (q *Queue) Push(v interface{}) *list.Element {
	e := q.v.PushBack(v)
	q.Count++
	return e
}

func (q *Queue) Pop() interface{} {
	front := q.v.Front()
	if front == nil {
		return nil
	}

	q.Count--

	return q.v.Remove(front)
}

func (q *Queue) Peek() interface{} {
	return q.v.Front().Value
}

func (q *Queue) Remove(item *list.Element) interface{} {
	return q.v.Remove(item)
}
