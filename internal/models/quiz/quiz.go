package quiz

import (
	"time"

	"gorm.io/gorm"
	"quizlet/internal/models/user"
)

type QuizType string

const (
	QuizTypeSingleChoice QuizType = "single_choice"
	QuizTypeMultiChoice  QuizType = "multi_choice"
	QuizTypeTrueFalse    QuizType = "true_false"
)

type Quiz struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Question      string         `gorm:"not null" json:"question"`
	QuizType      QuizType       `gorm:"not null" json:"quiz_type"`
	CreatedByID   uint           `gorm:"not null" json:"created_by_id"`
	CreatedBy     *user.User     `json:"created_by,omitempty"`
	Selections    []QuizSelection `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"selections,omitempty"`
}

type QuizSelection struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	QuizID        uint           `json:"quiz_id"`
	Quiz          *Quiz          `json:"quiz,omitempty"`
	SelectionText string         `gorm:"not null" json:"selection_text"`
	IsCorrect     bool           `gorm:"not null" json:"is_correct"`
} 