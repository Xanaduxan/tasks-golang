package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Xanaduxan/tasks-golang/task-service/internal/events"
	kafkago "github.com/segmentio/kafka-go"
)

type Publisher struct {
	writer *kafkago.Writer
}

func NewPublisher(brokers []string) (*Publisher, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers are required")
	}

	writer := &kafkago.Writer{
		Addr:         kafkago.TCP(brokers...),
		Topic:        events.TaskEventsTopic,
		Balancer:     &kafkago.Hash{},
		RequiredAcks: kafkago.RequireOne,
		Async:        false,
	}

	return &Publisher{
		writer: writer,
	}, nil
}

func (p *Publisher) PublishTaskEvent(ctx context.Context, event events.TaskEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(event.UserID.String()),
		Value: payload,
		Time:  time.Now(),
	})
}

func (p *Publisher) Close() error {
	if p.writer == nil {
		return nil
	}

	return p.writer.Close()
}
