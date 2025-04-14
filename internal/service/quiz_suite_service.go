package service

import (
	"quizlet/internal/models/quiz_suite"
	"quizlet/internal/repository"
)

type QuizSuiteService interface {
	CreateQuizSuite(quizSuite *quiz_suite.QuizSuite) error
	GetQuizSuite(id uint) (*quiz_suite.QuizSuite, error)
	GetUserQuizSuites(userID uint) ([]*quiz_suite.QuizSuite, error)
	UpdateQuizSuite(quizSuite *quiz_suite.QuizSuite) error
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

func (s *quizSuiteService) CreateQuizSuite(quizSuite *quiz_suite.QuizSuite) error {
	return s.quizSuiteRepo.Create(quizSuite)
}

func (s *quizSuiteService) GetQuizSuite(id uint) (*quiz_suite.QuizSuite, error) {
	return s.quizSuiteRepo.FindByID(id)
}

func (s *quizSuiteService) GetUserQuizSuites(userID uint) ([]*quiz_suite.QuizSuite, error) {
	return s.quizSuiteRepo.FindByUserID(userID)
}

func (s *quizSuiteService) UpdateQuizSuite(quizSuite *quiz_suite.QuizSuite) error {
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
	quizSuite, err := s.quizSuiteRepo.FindByID(quizSuiteID)
	if err != nil {
		return err
	}

	quiz, err := s.quizRepo.FindByID(quizID)
	if err != nil {
		return err
	}

	// Add quiz to suite
	quizSuite.Quizzes = append(quizSuite.Quizzes, quiz)
	return s.quizSuiteRepo.Update(quizSuite)
}

func (s *quizSuiteService) RemoveQuizFromSuite(quizSuiteID uint, quizID uint) error {
	// Verify quiz suite exists
	quizSuite, err := s.quizSuiteRepo.FindByID(quizSuiteID)
	if err != nil {
		return err
	}

	// Remove quiz from suite
	for i, quiz := range quizSuite.Quizzes {
		if quiz.ID == quizID {
			quizSuite.Quizzes = append(quizSuite.Quizzes[:i], quizSuite.Quizzes[i+1:]...)
			break
		}
	}

	return s.quizSuiteRepo.Update(quizSuite)
} 