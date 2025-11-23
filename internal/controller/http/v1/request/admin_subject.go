package request

import "github.com/evrone/go-clean-template/internal/entity"

// AdminSubjectCreateRequest is payload for creating subjects.
type AdminSubjectCreateRequest struct {
	Exam entity.ExamCategory `json:"exam" validate:"required"`
	Name string              `json:"name" validate:"required"`
}
