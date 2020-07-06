package base

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

// Model is base model for all model
type Model struct {
	ID        uuid.UUID  `json:"id"`
	CreatedBy uuid.UUID  `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedBy *uuid.UUID `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}
