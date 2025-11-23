package request

import (
	"github.com/google/uuid"
)

// AdminTopicCreateRequest is payload for creating topics.
type AdminTopicCreateRequest struct {
	SubjectID uuid.UUID `json:"subjectId" validate:"required"`
	Name      string    `json:"name" validate:"required"`
}
