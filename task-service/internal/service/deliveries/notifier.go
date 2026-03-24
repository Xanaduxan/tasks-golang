package deliveries

type Notifier interface {
	SendNotification(userIDs []string, typ, entityID, status string) error
}
