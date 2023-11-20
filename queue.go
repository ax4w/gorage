package Gorage

import "sync"

type queue struct {
	sync.Mutex
	d *data
	n *queue
}

func newQueue(d *data) *queue {
	return &queue{
		d: d,
	}
}

func (q *queue) append(d *data) {
	q.Lock()
	t := q
	for t.n != nil {
		t = t.n
	}
	t.n = &queue{
		d: d,
	}
	q.Unlock()
}

func (q *queue) Tail() *data {
	q.Lock()
	t := q
	for t.n != nil {
		t = t.n
	}
	q.Unlock()
	return t.d
}

func (q *queue) Head() *data {
	if q == nil {
		return nil
	}
	q.Lock()
	d := q.d
	q.Unlock()
	return d
}

func (q *queue) Shift() *queue {
	q.Lock()
	n := q.n
	q.Unlock()
	return n
}
