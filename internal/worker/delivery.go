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
			time.Sleep(1 * time.Second)

			err := w.service.UpdateDeliveryStatus(id, storage.StatusAccepted)
			if err != nil {
				log.Println("worker error:", err)
			}
		}
	}()
}
