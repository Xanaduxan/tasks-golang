package worker

import (
	"log"
	"time"

	"github.com/Xanaduxan/tasks-golang/internal/queue"
	"github.com/Xanaduxan/tasks-golang/internal/service/deliveries"
	"github.com/Xanaduxan/tasks-golang/internal/storage"
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
			err := w.service.UpdateDeliveryStatus(id, storage.StatusAccepted)
			if err != nil {
				log.Println("worker error:", err)
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
