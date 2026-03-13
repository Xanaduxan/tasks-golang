package queue

import "github.com/google/uuid"

type TaskQueue struct {
	jobs chan uuid.UUID
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{jobs: make(chan uuid.UUID, 100)}
}

func (q *TaskQueue) Push(job uuid.UUID) {
	q.jobs <- job
}

func (q *TaskQueue) GetJobs() chan uuid.UUID {
	return q.jobs
}
