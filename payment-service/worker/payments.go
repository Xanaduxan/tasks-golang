package worker

import (
	"log"
	"time"

	"github.com/Xanaduxan/tasks-golang/payment-service/internal/service/payments"
)

type PaymentWorker struct {
	service *payments.PaymentsService
}

func NewPaymentWorker(s *payments.PaymentsService) *PaymentWorker {
	return &PaymentWorker{
		service: s,
	}
}

func (w *PaymentWorker) Start() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := w.service.ProcessWaitingPayments(); err != nil {
				log.Println("payment state machine error:", err)
			}
		}
	}()
}
