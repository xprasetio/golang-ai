package service

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-ai/internal/dto"
	"golang-ai/internal/entity"
	"golang-ai/internal/repository"
	"golang-ai/pkg/embedding"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type IConsumerService interface {
	Consume(ctx context.Context) error
}

type consumerService struct {
	pubsub                  *gochannel.GoChannel
	noteEmbeddingRepository repository.INoteEmbeddingRepository
	topicName               string
	noteRepository          repository.INoteRepository
}

func NewConsumerService(pubsub *gochannel.GoChannel, topicName string, noteRepository repository.INoteRepository, noteEmbeddingRepository repository.INoteEmbeddingRepository) IConsumerService {
	return &consumerService{
		pubsub:                  pubsub,
		topicName:               topicName,
		noteRepository:          noteRepository,
		noteEmbeddingRepository: noteEmbeddingRepository,
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
	note, err := c.noteRepository.GetById(ctx, payload.NoteId)
	if err != nil {
		fmt.Println("Error fetching note:", err)
		panic(err)
	}

	res, err := embedding.GetGeminiEmbedding(
		os.Getenv("GOOGLE_GEMINI_API_KEY"),
		note.Content,
	)
	if err != nil {
		panic(err)
	}
	// embeddingFloat32 := make([]float32, len(res.Embedding.Values))
	// for i, v := range res.Embedding.Values {
	// 	embeddingFloat32[i] = float32(v)
	// }
	noteEmbedding := entity.NoteEmbedding{
		Id:             uuid.New(),
		Document:       note.Content,
		EmbeddingValue: res.Embedding.Values,
		NotekId:        note.Id,
		CreatedAt:      time.Now(),
	}
	err = c.noteEmbeddingRepository.Create(ctx, &noteEmbedding)
	if err != nil {
		panic(err)
	}
	// Process the payload as needed
	fmt.Println("Received message:", res.Embedding.Values)
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
