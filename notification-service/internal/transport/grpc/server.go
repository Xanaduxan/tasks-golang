package grpc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Xanaduxan/tasks-golang/notification-service/internal/transport/websocket"
	notificationpb "github.com/Xanaduxan/tasks-golang/notification-service/pkg/pb/notification/v1"
	"github.com/google/uuid"
)

type Server struct {
	notificationpb.UnimplementedNotificationServiceServer
	manager *websocket.Manager
}

func NewServer(manager *websocket.Manager) *Server {
	return &Server{manager: manager}
}

func (s *Server) SendNotification(
	ctx context.Context,
	req *notificationpb.SendNotificationRequest,
) (*notificationpb.SendNotificationResponse, error) {
	var data interface{}

	switch req.GetType() {
	case "task.status_updated":
		taskID, err := uuid.Parse(req.GetEntityId())
		if err != nil {
			return nil, err
		}
		data = websocket.TaskStatusUpdatedData{
			TaskID: taskID,
			Status: req.GetStatus(),
		}

	case "delivery.status_updated":
		deliveryID, err := uuid.Parse(req.GetEntityId())
		if err != nil {
			return nil, err
		}
		data = websocket.DeliveryStatusUpdatedData{
			DeliveryID: deliveryID,
			Status:     req.GetStatus(),
		}

	default:
		data = map[string]string{
			"entity_id": req.GetEntityId(),
			"status":    req.GetStatus(),
		}
	}

	msg := websocket.Message{
		Type:      req.GetType(),
		Data:      data,
		Timestamp: time.Now(),
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	for _, userIDStr := range req.GetUserIds() {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			continue
		}
		s.manager.SendToUser(userID, payload)
	}

	return &notificationpb.SendNotificationResponse{Ok: true}, nil
}
