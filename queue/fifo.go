package queue

import "go-snake-ai/tile"

func NewFIFO() *FIFO {
	return &FIFO{
		queue: make([]*tile.Vector, 0),
	}
}

type FIFO struct {
	queue []*tile.Vector
}

func (q *FIFO) Add(v *tile.Vector) {
	q.queue = append(q.queue, v)
}

func (q *FIFO) Pop() *tile.Vector {
	if len(q.queue) == 0 {
		return nil
	}

	v := q.queue[0]
	q.queue = q.queue[1:]

	return v
}
