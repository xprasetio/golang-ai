package dto

import "github.com/google/uuid"

type CreateNotebookRequest struct {
	Name     string     `json:"name" validate:"required,min=3"`
	ParentId *uuid.UUID `json:"parent_id" validate:"omitempty,uuid4"`
}

type CreateNotebookResponse struct {
	Id uuid.UUID `json:"id"`
}
