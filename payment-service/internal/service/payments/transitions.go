package payments

import "github.com/Xanaduxan/tasks-golang/payment-service/internal/storage"

var transitions = map[storage.PaymentStatus][]storage.PaymentStatus{
	storage.PaymentStatusNew: {
		storage.PaymentStatusWaitingValidation2,
	},
	storage.PaymentStatusWaitingValidation2: {
		storage.PaymentStatusFailed,
		storage.PaymentStatusReadyForClosure,
	},
	storage.PaymentStatusReadyForClosure: {
		storage.PaymentStatusClosed,
	},
}

func isValidStatusTransition(from, to storage.PaymentStatus) bool {
	allowed, ok := transitions[from]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}

	return false
}
