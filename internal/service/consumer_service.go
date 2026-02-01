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
	notebookRepository      repository.INotebookRepository
	pubsub                  *gochannel.GoChannel
	noteEmbeddingRepository repository.INoteEmbeddingRepository
	topicName               string
	noteRepository          repository.INoteRepository
}

func NewConsumerService(pubsub *gochannel.GoChannel, notebookRepository repository.INotebookRepository, topicName string, noteRepository repository.INoteRepository, noteEmbeddingRepository repository.INoteEmbeddingRepository) IConsumerService {
	return &consumerService{
		notebookRepository:      notebookRepository,
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
	notebook, err := c.notebookRepository.GetById(ctx, note.NotebookId)
	if err != nil {
		panic(err)
	}
	noteUpdatedAt := "-"
	if note.UpdatedAt != nil {
		noteUpdatedAt = note.UpdatedAt.Format(time.RFC3339)
	}
	content := fmt.Sprintf(`
	Note Title: %s
	Notebook Title: %s

	%s

	Created At: %s
	Updated At: %s
	`, note.Title,
		notebook.Name,
		note.Content,
		note.CreatedAt.Format(time.RFC3339),
		noteUpdatedAt,
	)

	res, err := embedding.GetGeminiEmbedding(
		os.Getenv("GOOGLE_GEMINI_API_KEY"),
		content,
	)
	if err != nil {
		panic(err)
	}

	noteEmbedding := entity.NoteEmbedding{
		Id:             uuid.New(),
		Document:       content,
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
