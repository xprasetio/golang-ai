package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateNoteRequest struct {
	Title      string    `json:"title" validate:"required,min=3"`
	Content    string    `json:"content"`
	NotebookId uuid.UUID `json:"notebook_id" validate:"required,uuid4"`
}

type CreateNoteResponse struct {
	Id uuid.UUID `json:"id"`
}
type ShowNoteResponse struct {
	Id         uuid.UUID  `json:"id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	NotebookId uuid.UUID  `json:"notebook_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type UpdateNoteRequest struct {
	Id      uuid.UUID
	Title   string `json:"title" validate:"required,min=3"`
	Content string `json:"content"`
}
type UpdateNoteResponse struct {
	Id         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	NotebookId uuid.UUID `json:"notebook_id"`
}
