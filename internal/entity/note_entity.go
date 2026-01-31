package entity

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	Id         uuid.UUID  `json:"id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	NotebookId uuid.UUID  `json:"notebook_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
	IsDeleted  bool       `json:"is_deleted"`
}
