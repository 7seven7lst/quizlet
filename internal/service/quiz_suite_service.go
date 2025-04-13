package service

import (
	"quizlet/internal/models"
	"quizlet/internal/repository"
)

type QuizSuiteService interface {
	CreateQuizSuite(quizSuite *models.QuizSuite) error
	GetQuizSuiteByID(id uint) (*models.QuizSuite, error)
	GetQuizSuitesByUserID(userID uint) ([]*models.QuizSuite, error)
	UpdateQuizSuite(quizSuite *models.QuizSuite) error
	DeleteQuizSuite(id uint) error
	AddQuizToSuite(quizSuiteID uint, quizID uint) error
	RemoveQuizFromSuite(quizSuiteID uint, quizID uint) error
}

type quizSuiteService struct {
	quizSuiteRepo repository.QuizSuiteRepository
	quizRepo      repository.QuizRepository
}

func NewQuizSuiteService(quizSuiteRepo repository.QuizSuiteRepository, quizRepo repository.QuizRepository) QuizSuiteService {
	return &quizSuiteService{
		quizSuiteRepo: quizSuiteRepo,
		quizRepo:      quizRepo,
	}
}

func (s *quizSuiteService) CreateQuizSuite(quizSuite *models.QuizSuite) error {
	return s.quizSuiteRepo.Create(quizSuite)
}

func (s *quizSuiteService) GetQuizSuiteByID(id uint) (*models.QuizSuite, error) {
	return s.quizSuiteRepo.FindByID(id)
}

func (s *quizSuiteService) GetQuizSuitesByUserID(userID uint) ([]*models.QuizSuite, error) {
	return s.quizSuiteRepo.FindByUserID(userID)
}

func (s *quizSuiteService) UpdateQuizSuite(quizSuite *models.QuizSuite) error {
	// Verify the quiz suite exists
	existing, err := s.quizSuiteRepo.FindByID(quizSuite.ID)
	if err != nil {
		return err
	}

	// Only allow updating certain fields
	existing.Title = quizSuite.Title
	existing.Description = quizSuite.Description

	return s.quizSuiteRepo.Update(existing)
}

func (s *quizSuiteService) DeleteQuizSuite(id uint) error {
	return s.quizSuiteRepo.Delete(id)
}

func (s *quizSuiteService) AddQuizToSuite(quizSuiteID uint, quizID uint) error {
	// Verify both quiz suite and quiz exist
	_, err := s.quizSuiteRepo.FindByID(quizSuiteID)
	if err != nil {
		return err
	}

	_, err = s.quizRepo.FindByID(quizID)
	if err != nil {
		return err
	}

	return s.quizSuiteRepo.AddQuiz(quizSuiteID, quizID)
}

func (s *quizSuiteService) RemoveQuizFromSuite(quizSuiteID uint, quizID uint) error {
	return s.quizSuiteRepo.RemoveQuiz(quizSuiteID, quizID)
} 