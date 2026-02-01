package entity

import (
	"time"

	"github.com/google/uuid"
)

type NoteEmbedding struct {
	Id             uuid.UUID  `json:"id"`
	Document       string     `json:"document"`
	EmbeddingValue []float32    `json:"embedding_value"`
	NotekId        uuid.UUID  `json:"note_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
	IsDeleted      bool       `json:"is_deleted"`
}
