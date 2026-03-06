package queue

import "github.com/google/uuid"

type DeliveryQueue struct {
	jobs chan uuid.UUID
}

func NewDeliveryQueue() *DeliveryQueue {
	return &DeliveryQueue{jobs: make(chan uuid.UUID, 100)}
}

func (q *DeliveryQueue) Push(job uuid.UUID) {
	q.jobs <- job
}

func (q *DeliveryQueue) GetJobs() chan uuid.UUID {
	return q.jobs
}
