package repository

import (
	"quizlet/internal/models"

	"gorm.io/gorm"
)

type QuizSuiteRepository interface {
	Create(quizSuite *models.QuizSuite) error
	FindByID(id uint) (*models.QuizSuite, error)
	FindByUserID(userID uint) ([]*models.QuizSuite, error)
	Update(quizSuite *models.QuizSuite) error
	Delete(id uint) error
	AddQuiz(quizSuiteID uint, quizID uint) error
	RemoveQuiz(quizSuiteID uint, quizID uint) error
}

type quizSuiteRepository struct {
	db *gorm.DB
}

func NewQuizSuiteRepository(db *gorm.DB) QuizSuiteRepository {
	return &quizSuiteRepository{db: db}
}

func (r *quizSuiteRepository) Create(quizSuite *models.QuizSuite) error {
	return r.db.Create(quizSuite).Error
}

func (r *quizSuiteRepository) FindByID(id uint) (*models.QuizSuite, error) {
	var quizSuite models.QuizSuite
	err := r.db.Preload("Quizzes").Preload("CreatedBy").First(&quizSuite, id).Error
	if err != nil {
		return nil, err
	}
	return &quizSuite, nil
}

func (r *quizSuiteRepository) FindByUserID(userID uint) ([]*models.QuizSuite, error) {
	var quizSuites []*models.QuizSuite
	err := r.db.Where("created_by_id = ?", userID).Preload("Quizzes").Find(&quizSuites).Error
	if err != nil {
		return nil, err
	}
	return quizSuites, nil
}

func (r *quizSuiteRepository) Update(quizSuite *models.QuizSuite) error {
	return r.db.Save(quizSuite).Error
}

func (r *quizSuiteRepository) Delete(id uint) error {
	return r.db.Delete(&models.QuizSuite{}, id).Error
}

func (r *quizSuiteRepository) AddQuiz(quizSuiteID uint, quizID uint) error {
	return r.db.Model(&models.QuizSuite{ID: quizSuiteID}).Association("Quizzes").Append(&models.Quiz{ID: quizID})
}

func (r *quizSuiteRepository) RemoveQuiz(quizSuiteID uint, quizID uint) error {
	return r.db.Model(&models.QuizSuite{ID: quizSuiteID}).Association("Quizzes").Delete(&models.Quiz{ID: quizID})
} 