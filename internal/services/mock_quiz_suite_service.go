package services

import (
	"quizlet/internal/models/quiz_suite"
	"github.com/stretchr/testify/mock"
)

type MockQuizSuiteService struct {
	mock.Mock
}

func NewMockQuizSuiteService() *MockQuizSuiteService {
	return &MockQuizSuiteService{}
}

func (m *MockQuizSuiteService) CreateQuizSuite(quizSuite *quiz_suite.QuizSuite) error {
	args := m.Called(quizSuite)
	return args.Error(0)
}

func (m *MockQuizSuiteService) GetQuizSuite(id uint) (*quiz_suite.QuizSuite, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*quiz_suite.QuizSuite), args.Error(1)
}

func (m *MockQuizSuiteService) GetUserQuizSuites(userID uint) ([]*quiz_suite.QuizSuite, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*quiz_suite.QuizSuite), args.Error(1)
}

func (m *MockQuizSuiteService) UpdateQuizSuite(quizSuite *quiz_suite.QuizSuite) error {
	args := m.Called(quizSuite)
	return args.Error(0)
}

func (m *MockQuizSuiteService) DeleteQuizSuite(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockQuizSuiteService) AddQuizToSuite(quizSuiteID uint, quizID uint) error {
	args := m.Called(quizSuiteID, quizID)
	return args.Error(0)
}

func (m *MockQuizSuiteService) RemoveQuizFromSuite(quizSuiteID uint, quizID uint) error {
	args := m.Called(quizSuiteID, quizID)
	return args.Error(0)
}

func (m *MockQuizSuiteService) GetQuizSuites(userID uint) ([]*quiz_suite.QuizSuite, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*quiz_suite.QuizSuite), args.Error(1)
} 