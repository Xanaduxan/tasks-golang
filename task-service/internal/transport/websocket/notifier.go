package websocket

import (
	"encoding/json"
	"time"

	events2 "github.com/Xanaduxan/tasks-golang/task-service/internal/events"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/storage"
	"github.com/google/uuid"
)

type Notifier struct {
	manager      *Manager
	groupMembers *storage.GroupMemberStorage
}

func NewNotifier(manager *Manager, groupMembers *storage.GroupMemberStorage) *Notifier {
	return &Notifier{
		manager:      manager,
		groupMembers: groupMembers,
	}
}

func (n *Notifier) NotifyDeliveryStatusUpdated(event events2.DeliveryStatusUpdated) error {
	msg := Message{
		Type: "delivery.status_updated",
		Data: DeliveryStatusUpdatedData{
			DeliveryID: event.DeliveryID,
			Status:     string(event.Status),
		},
		Timestamp: time.Now(),
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	n.manager.SendToUser(event.UserID, payload)
	return nil
}
func (n *Notifier) NotifyTaskStatusUpdated(event events2.TaskStatusUpdated) error {
	msg := Message{
		Type: "task.status_updated",
		Data: TaskStatusUpdatedData{
			TaskID: event.TaskID,
			Status: string(event.Status),
		},
		Timestamp: time.Now(),
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	recipients := map[uuid.UUID]struct{}{
		event.UserID: {},
	}

	if event.GroupID != nil {
		members, err := n.groupMembers.GetByGroupID(*event.GroupID)
		if err != nil {
			return err
		}

		for _, member := range members {
			recipients[member.UserID] = struct{}{}
		}
	}

	for userID := range recipients {
		n.manager.SendToUser(userID, payload)
	}

	return nil
}
