package quiz_suite

import (
	"quizlet/internal/models/quiz"
	"quizlet/internal/models/user"
	"time"

	"gorm.io/gorm"
)

// CreateQuizSuiteRequest represents the request body for creating a quiz suite
// @model CreateQuizSuiteRequest
type CreateQuizSuiteRequest struct {
	// @example "My Quiz Suite"
	// @required true
	Title       string `json:"title" binding:"required" example:"My Quiz Suite"`
	
	// @example "A collection of quizzes about various topics"
	// @required true
	Description string `json:"description" binding:"required" example:"A collection of quizzes about various topics"`
}

// UpdateQuizSuiteRequest represents the request body for updating a quiz suite
type UpdateQuizSuiteRequest struct {
	Title       string `json:"title" example:"Updated Quiz Suite"`
	Description string `json:"description" example:"An updated collection of quizzes"`
}

// QuizSuite represents a collection of quizzes
// @model QuizSuite
// @Description A collection of quizzes grouped together
type QuizSuite struct {
	// The unique identifier for the quiz suite
	// @example 1
	ID          uint           `json:"id" gorm:"primaryKey" example:"1"`
	
	// The timestamp when the quiz suite was created
	// @example "2024-04-17T00:00:00Z"
	CreatedAt   time.Time      `json:"created_at" example:"2024-04-17T00:00:00Z"`
	
	// The timestamp when the quiz suite was last updated
	// @example "2024-04-17T00:00:00Z"
	UpdatedAt   time.Time      `json:"updated_at" example:"2024-04-17T00:00:00Z"`
	
	// The timestamp when the quiz suite was deleted (soft delete)
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	
	// The title of the quiz suite
	// @example "My Quiz Suite"
	// @required true
	Title       string         `json:"title" binding:"required" example:"My Quiz Suite"`
	
	// A description of the quiz suite
	// @example "A collection of quizzes about various topics"
	// @required true
	Description string         `json:"description" binding:"required" example:"A collection of quizzes about various topics"`
	
	// The ID of the user who created the quiz suite
	// @example 1
	CreatedByID uint           `json:"created_by_id" example:"1"`
	
	// The user who created the quiz suite
	CreatedBy   *user.User     `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`
	
	// The quizzes in this suite
	Quizzes     []*quiz.Quiz   `json:"quizzes,omitempty" gorm:"many2many:quiz_suite_quizzes;"`
} 