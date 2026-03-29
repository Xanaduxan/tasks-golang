package tasks

type Notifier interface {
	SendNotification(userIDs []string, typ, entityID, status string) error
}
