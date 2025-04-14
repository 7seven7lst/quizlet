package repository

import (
	"quizlet/internal/models/quiz_suite"

	"gorm.io/gorm"
)

type QuizSuiteRepository interface {
	Create(quizSuite *quiz_suite.QuizSuite) error
	FindByID(id uint) (*quiz_suite.QuizSuite, error)
	FindByUserID(userID uint) ([]*quiz_suite.QuizSuite, error)
	Update(quizSuite *quiz_suite.QuizSuite) error
	Delete(id uint) error
}

type quizSuiteRepository struct {
	db *gorm.DB
}

func NewQuizSuiteRepository(db *gorm.DB) QuizSuiteRepository {
	return &quizSuiteRepository{db: db}
}

func (r *quizSuiteRepository) Create(quizSuite *quiz_suite.QuizSuite) error {
	return r.db.Create(quizSuite).Error
}

func (r *quizSuiteRepository) FindByID(id uint) (*quiz_suite.QuizSuite, error) {
	var quizSuite quiz_suite.QuizSuite
	err := r.db.Preload("Quizzes").Preload("CreatedBy").First(&quizSuite, id).Error
	if err != nil {
		return nil, err
	}
	return &quizSuite, nil
}

func (r *quizSuiteRepository) FindByUserID(userID uint) ([]*quiz_suite.QuizSuite, error) {
	var quizSuites []*quiz_suite.QuizSuite
	err := r.db.Where("created_by_id = ?", userID).Preload("Quizzes").Find(&quizSuites).Error
	if err != nil {
		return nil, err
	}
	return quizSuites, nil
}

func (r *quizSuiteRepository) Update(quizSuite *quiz_suite.QuizSuite) error {
	return r.db.Save(quizSuite).Error
}

func (r *quizSuiteRepository) Delete(id uint) error {
	return r.db.Delete(&quiz_suite.QuizSuite{}, id).Error
} 