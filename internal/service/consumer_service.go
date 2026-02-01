package service

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-ai/internal/dto"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/gofiber/fiber/v2/log"
)

type IConsumerService interface {
	Consume(ctx context.Context) error
}

type consumerService struct {
	pubsub    *gochannel.GoChannel
	topicName string
}

func NewConsumerService(pubsub *gochannel.GoChannel, topicName string) IConsumerService {
	return &consumerService{
		pubsub:    pubsub,
		topicName: topicName,
	}
}

func (c *consumerService) processMessage(ctx context.Context, msg *message.Message) {
	defer msg.Nack()
	defer func() {
		if r := recover(); r != nil {
			log.Error(r)
			fmt.Println("Recovered in processMessage:", r)
		}
	}()
	var payload dto.PublishEmbedNoteMessage
	err := json.Unmarshal(msg.Payload, &payload)
	if err != nil {
		// Handle error appropriately
		fmt.Println("Error unmarshaling message:", err)
		panic(err)
	}
	// Process the payload as needed
	fmt.Println("Received message:", payload.NoteId)
	_ = payload // Example usage of payload
	msg.Ack()
}
func (c *consumerService) Consume(ctx context.Context) error {
	messages, err := c.pubsub.Subscribe(ctx, c.topicName)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			c.processMessage(ctx, msg)
		}
	}()

	return nil
}
