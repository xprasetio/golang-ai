package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateNotebookRequest struct {
	Name     string     `json:"name" validate:"required,min=3"`
	ParentId *uuid.UUID `json:"parent_id" validate:"omitempty,uuid4"`
}

type CreateNotebookResponse struct {
	Id uuid.UUID `json:"id"`
}

type ShowNotebookResponse struct {
	Id        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	ParentId  *uuid.UUID `json:"parent_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type UpdateNotebookRequest struct {
	Id   uuid.UUID `json:"id" validate:"required,uuid4"`
	Name string    `json:"name" validate:"required,min=3"`
}
type UpdateNotebookResponse struct {
	Id uuid.UUID `json:"id"`
}
type MoveNotebookRequest struct {
	Id       uuid.UUID  `json:"id" validate:"required,uuid4"`
	ParentId *uuid.UUID `json:"parent_id"`
}

type MoveNotebookResponse struct {
	Id uuid.UUID `json:"id"`
}
type GetAllNotebooksResponse struct {
	Id        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	ParentId  *uuid.UUID `json:"parent_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
