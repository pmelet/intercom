package main

import "container/list"

type Queue struct {
	items list.List
}

func (q *Queue) Push(v interface{}) {
	q.items.PushBack(v)
}

func (q *Queue) Pop() interface{} {
	e := q.items.Front()
	if e == nil {
		return nil
	}
	q.items.Remove(e)
	return e.Value
}
