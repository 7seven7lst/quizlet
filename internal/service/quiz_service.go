package service

import (
	"quizlet/internal/models"
	"quizlet/internal/repository"
)

type QuizService interface {
	CreateQuiz(quiz *models.Quiz) error
	GetQuizByID(id uint) (*models.Quiz, error)
	GetQuizzesByUserID(userID uint) ([]*models.Quiz, error)
	UpdateQuiz(quiz *models.Quiz) error
	DeleteQuiz(id uint) error
	AddSelection(quizID uint, selection *models.QuizSelection) error
	RemoveSelection(quizID uint, selectionID uint) error
}

type quizService struct {
	quizRepo repository.QuizRepository
}

func NewQuizService(quizRepo repository.QuizRepository) QuizService {
	return &quizService{
		quizRepo: quizRepo,
	}
}

func (s *quizService) CreateQuiz(quiz *models.Quiz) error {
	return s.quizRepo.Create(quiz)
}

func (s *quizService) GetQuizByID(id uint) (*models.Quiz, error) {
	return s.quizRepo.FindByID(id)
}

func (s *quizService) GetQuizzesByUserID(userID uint) ([]*models.Quiz, error) {
	return s.quizRepo.FindByUserID(userID)
}

func (s *quizService) UpdateQuiz(quiz *models.Quiz) error {
	// Verify the quiz exists
	existing, err := s.quizRepo.FindByID(quiz.ID)
	if err != nil {
		return err
	}

	// Only allow updating certain fields
	existing.Question = quiz.Question
	existing.QuizType = quiz.QuizType
	existing.CorrectAnswer = quiz.CorrectAnswer

	return s.quizRepo.Update(existing)
}

func (s *quizService) DeleteQuiz(id uint) error {
	return s.quizRepo.Delete(id)
}

func (s *quizService) AddSelection(quizID uint, selection *models.QuizSelection) error {
	// Verify the quiz exists
	_, err := s.quizRepo.FindByID(quizID)
	if err != nil {
		return err
	}

	return s.quizRepo.AddSelection(quizID, selection)
}

func (s *quizService) RemoveSelection(quizID uint, selectionID uint) error {
	return s.quizRepo.RemoveSelection(quizID, selectionID)
} 