package grpc

import (
	"context"
	"time"

	notificationpb "github.com/Xanaduxan/tasks-golang/notification-service/pkg/pb/notification/v1"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NotificationClient struct {
	client notificationpb.NotificationServiceClient
	conn   *gogrpc.ClientConn
}

func NewNotificationClient(addr string) (*NotificationClient, error) {
	conn, err := gogrpc.NewClient(
		addr,
		gogrpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &NotificationClient{
		client: notificationpb.NewNotificationServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *NotificationClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *NotificationClient) SendNotification(userIDs []string, typ, entityID, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.client.SendNotification(ctx, &notificationpb.SendNotificationRequest{
		UserIds:  userIDs,
		Type:     typ,
		EntityId: entityID,
		Status:   status,
	})

	return err
}
