package join

// QuizSuiteQuiz represents the many-to-many relationship between QuizSuite and Quiz
type QuizSuiteQuiz struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	QuizSuiteID  uint `json:"quiz_suite_id" gorm:"not null"`
	QuizID       uint `json:"quiz_id" gorm:"not null"`
	QuizSuite    interface{} `json:"quiz_suite,omitempty" gorm:"foreignKey:QuizSuiteID"`
	Quiz         interface{} `json:"quiz,omitempty" gorm:"foreignKey:QuizID"`
} 