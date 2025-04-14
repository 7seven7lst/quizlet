package services

import (
	"quizlet/internal/models/quiz_suite"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockQuizSuiteService struct {
	mock.Mock
	quizSuites []*quiz_suite.QuizSuite
}

func NewMockQuizSuiteService() *MockQuizSuiteService {
	return &MockQuizSuiteService{
		quizSuites: make([]*quiz_suite.QuizSuite, 0),
	}
}

func (m *MockQuizSuiteService) CreateQuizSuite(quizSuite *quiz_suite.QuizSuite) error {
	if quizSuite == nil {
		return gorm.ErrInvalidDB
	}
	m.quizSuites = append(m.quizSuites, quizSuite)
	return nil
}

func (m *MockQuizSuiteService) GetQuizSuite(id uint) (*quiz_suite.QuizSuite, error) {
	for _, qs := range m.quizSuites {
		if qs.ID == id {
			return qs, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockQuizSuiteService) GetUserQuizSuites(userID uint) ([]*quiz_suite.QuizSuite, error) {
	var userQuizSuites []*quiz_suite.QuizSuite
	for _, qs := range m.quizSuites {
		if qs.CreatedByID == userID {
			userQuizSuites = append(userQuizSuites, qs)
		}
	}
	return userQuizSuites, nil
}

func (m *MockQuizSuiteService) UpdateQuizSuite(quizSuite *quiz_suite.QuizSuite) error {
	if quizSuite == nil {
		return gorm.ErrInvalidDB
	}
	for i, qs := range m.quizSuites {
		if qs.ID == quizSuite.ID {
			m.quizSuites[i] = quizSuite
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *MockQuizSuiteService) DeleteQuizSuite(id uint) error {
	for i, qs := range m.quizSuites {
		if qs.ID == id {
			m.quizSuites = append(m.quizSuites[:i], m.quizSuites[i+1:]...)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *MockQuizSuiteService) AddQuizToSuite(quizSuiteID uint, quizID uint) error {
	quizSuite, err := m.GetQuizSuite(quizSuiteID)
	if err != nil {
		return err
	}
	quizSuite.Quizzes = append(quizSuite.Quizzes, &quiz.Quiz{ID: quizID})
	return nil
}

func (m *MockQuizSuiteService) RemoveQuizFromSuite(quizSuiteID uint, quizID uint) error {
	quizSuite, err := m.GetQuizSuite(quizSuiteID)
	if err != nil {
		return err
	}
	for i, quiz := range quizSuite.Quizzes {
		if quiz.ID == quizID {
			quizSuite.Quizzes = append(quizSuite.Quizzes[:i], quizSuite.Quizzes[i+1:]...)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *MockQuizSuiteService) GetQuizSuites(userID uint) ([]models.QuizSuite, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.QuizSuite), args.Error(1)
}

func (m *MockQuizSuiteService) UpdateQuizSuite(quizSuite *models.QuizSuite) (*models.QuizSuite, error) {
	args := m.Called(quizSuite)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.QuizSuite), args.Error(1)
} 