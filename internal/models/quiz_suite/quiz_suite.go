package quiz_suite

import (
	"quizlet/internal/models/quiz"
	"quizlet/internal/models/user"
	"time"

	"gorm.io/gorm"
)

type QuizSuite struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	CreatedByID uint           `json:"created_by_id"`
	CreatedBy   *user.User     `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`
	Quizzes     []*quiz.Quiz   `json:"quizzes,omitempty" gorm:"many2many:quiz_suite_quizzes;"`
} 