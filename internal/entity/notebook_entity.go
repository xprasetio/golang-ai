package entity

import (
	"time"

	"github.com/google/uuid"
)

type Notebook struct {
	Id        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	ParentId  *uuid.UUID `json:"parent_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	IsDeleted bool       `json:"is_deleted"`
}
