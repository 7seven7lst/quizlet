package models

import (
	"time"

	"gorm.io/gorm"
)

type QuizSuite struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	CreatedByID uint           `json:"created_by_id"`
	CreatedBy   *User          `json:"created_by,omitempty"`
	Quizzes     []Quiz         `gorm:"many2many:quiz_suite_quizzes;" json:"quizzes,omitempty"`
} 