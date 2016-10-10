package queue

import (
	"container/list"
	"sync"
)

type Queue struct {
	data    *list.List
	rwMutex sync.RWMutex
}

func New() *Queue {
	return &Queue{
		data: list.New(),
	}
}

func (q *Queue) Add(elem interface{}) {
	q.rwMutex.Lock()
	q.data.PushBack(elem)
	q.rwMutex.Unlock()
}

func (q *Queue) Remove() interface{} {
	q.rwMutex.Lock()
	element := q.data.Front()
	defer q.rwMutex.Unlock()
	if element != nil {
		return element.Value
	}
	return nil
}
