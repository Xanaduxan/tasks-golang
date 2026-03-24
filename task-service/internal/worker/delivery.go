package worker

import (
	"log"
	"time"

	"github.com/Xanaduxan/tasks-golang/task-service/internal/queue"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/service/deliveries"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/storage"
)

type DeliveryWorker struct {
	queue   *queue.DeliveryQueue
	service *deliveries.Service
}

func NewDeliveryWorker(queue *queue.DeliveryQueue, s *deliveries.Service) *DeliveryWorker {
	return &DeliveryWorker{
		queue:   queue,
		service: s,
	}
}

func (w *DeliveryWorker) Start() {

	go func() {
		for id := range w.queue.GetJobs() {
			delivery, err := w.service.GetDeliveryByID(id)
			if err != nil {
				log.Println("worker error, delivery not found:", err)
				continue
			}

			next, ok := nextStatus(delivery.Status)
			if !ok {
				continue
			}

			err = w.service.UpdateDeliveryStatus(id, next)
			if err != nil {
				log.Println("worker error:", err)
				continue
			}

			if next != storage.StatusAccepted {

				w.queue.Push(id)
			}
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

func (w *DeliveryWorker) enqueuePending() {
	d, err := w.service.GetDeliveriesNotAccepted()
	if err != nil {
		log.Println("enqueue pending deliveries error:", err)
		return
	}

	for _, delivery := range d {
		w.queue.Push(delivery.ID)
	}
}
func nextStatus(status storage.DeliveryStatus) (storage.DeliveryStatus, bool) {
	switch status {
	case storage.StatusAwaiting:
		return storage.StatusProcessing, true
	case storage.StatusProcessing:
		return storage.StatusChecked, true
	case storage.StatusChecked:
		return storage.StatusOnPath, true
	case storage.StatusOnPath:
		return storage.StatusAccepted, true
	default:
		return "", false
	}
}
