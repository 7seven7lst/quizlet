package quiz_attempt

import (
	"quizlet/internal/models/user"
	"time"
)

// CreateQuizAttemptRequest represents the request body for creating a quiz attempt
// @model CreateQuizAttemptRequest
// @Description Request body for creating a new quiz attempt
type CreateQuizAttemptRequest struct {
	// The score achieved in this attempt
	// @example 80
	// @required true
	// @minimum 0
	// @maximum 100
	Score int `json:"score" binding:"required,min=0,max=100" example:"80"`

	// Whether the attempt is completed
	// @example true
	// @required true
	Completed bool `json:"completed" binding:"required" example:"true"`
}

// UpdateQuizAttemptRequest represents the request body for updating a quiz attempt
// @model UpdateQuizAttemptRequest
// @Description Request body for updating an existing quiz attempt
type UpdateQuizAttemptRequest struct {
	// The score achieved in this attempt
	// @example 90
	// @minimum 0
	// @maximum 100
	Score *int `json:"score,omitempty" example:"90"`

	// Whether the attempt is completed
	// @example true
	Completed *bool `json:"completed,omitempty" example:"true"`
}

// QuizAttempt represents a user's attempt at a quiz suite
// @model QuizAttempt
// @Description A record of a user's attempt at completing a quiz suite
type QuizAttempt struct {
	// The unique identifier for the quiz attempt
	// @example 1
	// @readOnly true
	ID int64 `json:"id" gorm:"primaryKey" example:"1"`

	// The timestamp when the quiz attempt was created
	// @example "2024-04-17T00:00:00Z"
	// @readOnly true
	CreatedAt time.Time `json:"created_at" example:"2024-04-17T00:00:00Z"`

	// The timestamp when the quiz attempt was last updated
	// @example "2024-04-17T00:00:00Z"
	// @readOnly true
	UpdatedAt time.Time `json:"updated_at" example:"2024-04-17T00:00:00Z"`

	// The timestamp when the quiz attempt was deleted (soft delete)
	// @readOnly true
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// The ID of the user who made the attempt
	// @example 1
	// @readOnly true
	UserID int64 `json:"user_id" example:"1"`

	// The user who made the attempt
	// @readOnly true
	User *user.User `json:"user,omitempty" gorm:"foreignKey:UserID"`

	// The ID of the quiz suite being attempted
	// @example 1
	// @readOnly true
	QuizSuiteID int64 `json:"quiz_suite_id" example:"1"`

	// The score achieved in this attempt
	// @example 80
	// @minimum 0
	// @maximum 100
	Score int `json:"score" example:"80"`

	// Whether the attempt is completed
	// @example true
	Completed bool `json:"completed" example:"true"`

	// The timestamp when the attempt was started
	// @example "2024-04-17T00:00:00Z"
	// @readOnly true
	StartedAt time.Time `json:"started_at" example:"2024-04-17T00:00:00Z"`

	// The timestamp when the attempt was completed
	// @example "2024-04-17T00:00:00Z"
	// @readOnly true
	CompletedAt *time.Time `json:"completed_at,omitempty" example:"2024-04-17T00:00:00Z"`
} 