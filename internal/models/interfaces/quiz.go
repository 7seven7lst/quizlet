package interfaces

// QuizInterface defines the interface for Quiz models
type QuizInterface interface {
	GetID() uint
	GetQuestion() string
	GetQuizType() string
	GetCreatedByID() uint
} 