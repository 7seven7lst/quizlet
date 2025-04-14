package repository

import (
	"quizlet/internal/models/quiz"

	"gorm.io/gorm"
)

type QuizRepository interface {
	Create(quiz *quiz.Quiz) error
	FindByID(id uint) (*quiz.Quiz, error)
	FindByUserID(userID uint) ([]*quiz.Quiz, error)
	Update(quiz *quiz.Quiz) error
	Delete(id uint) error
	AddSelection(quizID uint, selection quiz.QuizSelection) error
	RemoveSelection(quizID uint, selectionID uint) error
}

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) QuizRepository {
	return &quizRepository{db: db}
}

func (r *quizRepository) Create(quiz *quiz.Quiz) error {
	return r.db.Create(quiz).Error
}

func (r *quizRepository) FindByID(id uint) (*quiz.Quiz, error) {
	var quiz quiz.Quiz
	err := r.db.Preload("Selections").Preload("CreatedBy").First(&quiz, id).Error
	if err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *quizRepository) FindByUserID(userID uint) ([]*quiz.Quiz, error) {
	var quizzes []*quiz.Quiz
	err := r.db.Where("created_by_id = ?", userID).Preload("Selections").Find(&quizzes).Error
	if err != nil {
		return nil, err
	}
	return quizzes, nil
}

func (r *quizRepository) Update(quiz *quiz.Quiz) error {
	return r.db.Save(quiz).Error
}

func (r *quizRepository) Delete(id uint) error {
	return r.db.Delete(&quiz.Quiz{}, id).Error
}

func (r *quizRepository) AddSelection(quizID uint, selection quiz.QuizSelection) error {
	selection.QuizID = quizID
	return r.db.Create(&selection).Error
}

func (r *quizRepository) RemoveSelection(quizID uint, selectionID uint) error {
	return r.db.Where("quiz_id = ? AND id = ?", quizID, selectionID).Delete(&quiz.QuizSelection{}).Error
} 