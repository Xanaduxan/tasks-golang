package worker

import (
	"log"
	"time"

	"github.com/Xanaduxan/tasks-golang/internal/queue"
	"github.com/Xanaduxan/tasks-golang/internal/service/tasks"
	"github.com/Xanaduxan/tasks-golang/internal/storage"
	"github.com/Xanaduxan/tasks-golang/metrics"
)

type TaskWorker struct {
	queue   *queue.TaskQueue
	service *tasks.Service
}

func NewTaskWorker(queue *queue.TaskQueue, s *tasks.Service) *TaskWorker {
	return &TaskWorker{
		queue:   queue,
		service: s,
	}
}

func (w *TaskWorker) Start() {
	go func() {
		for id := range w.queue.GetJobs() {
			start := time.Now()

			task, err := w.service.GetTaskForWorker(id)
			if err != nil {
				log.Println("worker error, task not found:", err)
				metrics.TaskProcessingDuration.Observe(time.Since(start).Seconds())
				continue
			}

			next, ok := nextStatusTask(task.Status)
			if !ok {
				metrics.TaskProcessingDuration.Observe(time.Since(start).Seconds())
				continue
			}

			err = w.service.UpdateTaskStatus(id, next)
			if err != nil {
				log.Println("worker error:", err)
				metrics.TaskProcessingDuration.Observe(time.Since(start).Seconds())
				continue
			}

			if next != storage.StatusDone {
				w.queue.Push(id)
			}

			metrics.TaskProcessingDuration.Observe(time.Since(start).Seconds())
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		w.enqueuePending()

		for range ticker.C {
			w.enqueuePending()
		}
	}()
}

func (w *TaskWorker) enqueuePending() {
	t, err := w.service.GetAllNotDone()
	if err != nil {
		log.Println("enqueue pending t error:", err)
		return
	}

	for _, task := range t {
		w.queue.Push(task.ID)
	}
}

func nextStatusTask(status storage.TaskStatus) (storage.TaskStatus, bool) {
	switch status {
	case storage.StatusCreated:
		return storage.StatusInProgress, true
	case storage.StatusInProgress:
		return storage.StatusDone, true
	default:
		return "", false
	}
}
