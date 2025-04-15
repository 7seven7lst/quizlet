package interfaces

// QuizSuiteInterface defines the interface for QuizSuite models
type QuizSuiteInterface interface {
	GetID() uint
	GetTitle() string
	GetDescription() string
	GetCreatedByID() uint
} 