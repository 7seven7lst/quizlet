// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/quiz_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	quiz "quizlet/internal/models/quiz"
	service "quizlet/internal/service"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockQuizService is a mock of QuizService interface.
type MockQuizService struct {
	ctrl     *gomock.Controller
	recorder *MockQuizServiceMockRecorder
}

// MockQuizServiceMockRecorder is the mock recorder for MockQuizService.
type MockQuizServiceMockRecorder struct {
	mock *MockQuizService
}

// NewMockQuizService creates a new mock instance.
func NewMockQuizService(ctrl *gomock.Controller) *MockQuizService {
	mock := &MockQuizService{ctrl: ctrl}
	mock.recorder = &MockQuizServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuizService) EXPECT() *MockQuizServiceMockRecorder {
	return m.recorder
}

// AddSelection mocks base method.
func (m *MockQuizService) AddSelection(quizID uint, selection service.QuizSelection) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSelection", quizID, selection)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSelection indicates an expected call of AddSelection.
func (mr *MockQuizServiceMockRecorder) AddSelection(quizID, selection interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSelection", reflect.TypeOf((*MockQuizService)(nil).AddSelection), quizID, selection)
}

// CreateQuiz mocks base method.
func (m *MockQuizService) CreateQuiz(quiz *quiz.Quiz) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateQuiz", quiz)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateQuiz indicates an expected call of CreateQuiz.
func (mr *MockQuizServiceMockRecorder) CreateQuiz(quiz interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateQuiz", reflect.TypeOf((*MockQuizService)(nil).CreateQuiz), quiz)
}

// DeleteQuiz mocks base method.
func (m *MockQuizService) DeleteQuiz(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteQuiz", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteQuiz indicates an expected call of DeleteQuiz.
func (mr *MockQuizServiceMockRecorder) DeleteQuiz(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteQuiz", reflect.TypeOf((*MockQuizService)(nil).DeleteQuiz), id)
}

// GetQuizByID mocks base method.
func (m *MockQuizService) GetQuizByID(id uint) (*quiz.Quiz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuizByID", id)
	ret0, _ := ret[0].(*quiz.Quiz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuizByID indicates an expected call of GetQuizByID.
func (mr *MockQuizServiceMockRecorder) GetQuizByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuizByID", reflect.TypeOf((*MockQuizService)(nil).GetQuizByID), id)
}

// GetQuizzesByUserID mocks base method.
func (m *MockQuizService) GetQuizzesByUserID(userID uint) ([]*quiz.Quiz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuizzesByUserID", userID)
	ret0, _ := ret[0].([]*quiz.Quiz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuizzesByUserID indicates an expected call of GetQuizzesByUserID.
func (mr *MockQuizServiceMockRecorder) GetQuizzesByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuizzesByUserID", reflect.TypeOf((*MockQuizService)(nil).GetQuizzesByUserID), userID)
}

// RemoveSelection mocks base method.
func (m *MockQuizService) RemoveSelection(quizID, selectionID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSelection", quizID, selectionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSelection indicates an expected call of RemoveSelection.
func (mr *MockQuizServiceMockRecorder) RemoveSelection(quizID, selectionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSelection", reflect.TypeOf((*MockQuizService)(nil).RemoveSelection), quizID, selectionID)
}

// UpdateQuiz mocks base method.
func (m *MockQuizService) UpdateQuiz(quiz *quiz.Quiz) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateQuiz", quiz)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateQuiz indicates an expected call of UpdateQuiz.
func (mr *MockQuizServiceMockRecorder) UpdateQuiz(quiz interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateQuiz", reflect.TypeOf((*MockQuizService)(nil).UpdateQuiz), quiz)
}
