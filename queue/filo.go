package queue

import "go-snake-ai/tile"

func NewFILO() *FILO {
	return &FILO{
		queue: make([]*tile.Vector, 0),
	}
}

type FILO struct {
	queue []*tile.Vector
}

func (q *FILO) Add(v *tile.Vector) {
	q.queue = append(q.queue, v)
}

func (q *FILO) Pop() *tile.Vector {
	qLen := len(q.queue)
	if qLen == 0 {
		return nil
	}

	v := q.queue[qLen-1]
	q.queue = q.queue[:qLen-1]

	return v
}
