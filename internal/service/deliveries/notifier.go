package deliveries

import "github.com/Xanaduxan/tasks-golang/internal/events"

type Notifier interface {
	NotifyDeliveryStatusUpdated(event events.DeliveryStatusUpdated) error
}
