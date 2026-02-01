package service

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type IPublisherService interface {
	Publish(ctx context.Context, payload []byte) error
}

type publisherService struct{
	pubSub *gochannel.GoChannel

	topicName string
}

func NewPublisherService(pubSub *gochannel.GoChannel, topicName string) IPublisherService {
	return &publisherService{
		pubSub:    pubSub,
		topicName: topicName,
	}
}

func (p *publisherService) Publish(ctx context.Context, payload []byte) error {
	
	err := p.pubSub.Publish(p.topicName, message.NewMessage(watermill.NewUUID(), payload))
	if err != nil {
		return err
	}
	return nil
}
