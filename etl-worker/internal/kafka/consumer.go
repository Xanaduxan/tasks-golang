package kafka

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Xanaduxan/tasks-golang/etl-worker/internal/events"
	kafkago "github.com/segmentio/kafka-go"
)

type Processor interface {
	Process(event events.TaskEvent) error
}

type Consumer struct {
	reader    *kafkago.Reader
	processor Processor
}

func NewConsumer(brokers []string, groupID string, processor Processor) (*Consumer, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers are required")
	}
	if groupID == "" {
		return nil, errors.New("kafka group id is required")
	}
	if processor == nil {
		return nil, errors.New("processor is required")
	}

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   events.TaskEventsTopic,
	})

	return &Consumer{
		reader:    reader,
		processor: processor,
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}

		var event events.TaskEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			continue
		}

		if err := c.processor.Process(event); err != nil {
			return err
		}
	}
}

func (c *Consumer) Close() error {
	if c.reader == nil {
		return nil
	}

	return c.reader.Close()
}
