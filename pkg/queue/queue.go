package queue

import "sync"

type Queue[T any] struct {
	c     []T
	mutex sync.Mutex
}

func (q *Queue[T]) Push(el T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.c = append(q.c, el)
}

func (q *Queue[T]) Pop() T {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	el := q.c[0]
	q.c = q.c[1:]
	return el
}

func (q *Queue[T]) Size() int {
	return len(q.c)
}
