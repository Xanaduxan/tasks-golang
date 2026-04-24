package payments

import (
	"log"
	"math/rand"
	"time"

	"github.com/Xanaduxan/tasks-golang/payment-service/internal/storage"
)

func (s *PaymentsService) ProcessWaitingPayments() error {

	payments, err := s.payments.GetPaymentsWaitingForValidation2()
	if err != nil {
		return err
	}
	for _, payment := range payments {
		err := s.processSinglePayment(payment)
		if err != nil {
			log.Printf("error processing payment %s: %v", payment.ID, err)
			continue
		}
	}

	return nil
}
func (s *PaymentsService) processSinglePayment(payment storage.Payment) error {
	pTime := payment.WaitingForValidation2At
	if pTime.Valid && time.Since(pTime.Time) > time.Hour {
		_, err := s.payments.DeleteByID(payment.ID)
		if err != nil {
			return err
		}
		return nil
	}

	if s.validation2(payment) {
		payment.Status = storage.PaymentStatusReadyForClosure

		_, err := s.payments.Update(payment)
		if err != nil {
			log.Printf("error updating payment %s: %v", payment.ID, err)
			return err
		}
		return nil
	}

	payment.Attempts++
	if payment.Attempts >= 3 {
		payment.Status = storage.PaymentStatusFailed
	}

	_, err := s.payments.Update(payment)
	if err != nil {
		log.Printf("error updating payment %s: %v", payment.ID, err)
		return err
	}

	return nil
}

func (s *PaymentsService) validation2(payment storage.Payment) bool {

	return rand.Intn(100) < 30
}
