package repository

import (
	"quizlet/internal/models"

	"gorm.io/gorm"
)

type QuizRepository interface {
	Create(quiz *models.Quiz) error
	FindByID(id uint) (*models.Quiz, error)
	FindByUserID(userID uint) ([]*models.Quiz, error)
	Update(quiz *models.Quiz) error
	Delete(id uint) error
	AddSelection(quizID uint, selection *models.QuizSelection) error
	RemoveSelection(quizID uint, selectionID uint) error
}

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) QuizRepository {
	return &quizRepository{db: db}
}

func (r *quizRepository) Create(quiz *models.Quiz) error {
	return r.db.Create(quiz).Error
}

func (r *quizRepository) FindByID(id uint) (*models.Quiz, error) {
	var quiz models.Quiz
	err := r.db.Preload("Selections").Preload("CreatedBy").First(&quiz, id).Error
	if err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *quizRepository) FindByUserID(userID uint) ([]*models.Quiz, error) {
	var quizzes []*models.Quiz
	err := r.db.Where("created_by_id = ?", userID).Preload("Selections").Find(&quizzes).Error
	if err != nil {
		return nil, err
	}
	return quizzes, nil
}

func (r *quizRepository) Update(quiz *models.Quiz) error {
	return r.db.Save(quiz).Error
}

func (r *quizRepository) Delete(id uint) error {
	return r.db.Delete(&models.Quiz{}, id).Error
}

func (r *quizRepository) AddSelection(quizID uint, selection *models.QuizSelection) error {
	selection.QuizID = quizID
	return r.db.Create(selection).Error
}

func (r *quizRepository) RemoveSelection(quizID uint, selectionID uint) error {
	return r.db.Where("quiz_id = ? AND id = ?", quizID, selectionID).Delete(&models.QuizSelection{}).Error
} 