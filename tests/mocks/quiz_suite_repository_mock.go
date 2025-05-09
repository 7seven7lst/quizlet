// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/quiz_suite_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	quiz_suite "quizlet/internal/models/quiz_suite"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockQuizSuiteRepository is a mock of QuizSuiteRepository interface.
type MockQuizSuiteRepository struct {
	ctrl     *gomock.Controller
	recorder *MockQuizSuiteRepositoryMockRecorder
}

// MockQuizSuiteRepositoryMockRecorder is the mock recorder for MockQuizSuiteRepository.
type MockQuizSuiteRepositoryMockRecorder struct {
	mock *MockQuizSuiteRepository
}

// NewMockQuizSuiteRepository creates a new mock instance.
func NewMockQuizSuiteRepository(ctrl *gomock.Controller) *MockQuizSuiteRepository {
	mock := &MockQuizSuiteRepository{ctrl: ctrl}
	mock.recorder = &MockQuizSuiteRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuizSuiteRepository) EXPECT() *MockQuizSuiteRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockQuizSuiteRepository) Create(quizSuite *quiz_suite.QuizSuite) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", quizSuite)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockQuizSuiteRepositoryMockRecorder) Create(quizSuite interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockQuizSuiteRepository)(nil).Create), quizSuite)
}

// Delete mocks base method.
func (m *MockQuizSuiteRepository) Delete(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockQuizSuiteRepositoryMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockQuizSuiteRepository)(nil).Delete), id)
}

// FindByID mocks base method.
func (m *MockQuizSuiteRepository) FindByID(id uint) (*quiz_suite.QuizSuite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(*quiz_suite.QuizSuite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockQuizSuiteRepositoryMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockQuizSuiteRepository)(nil).FindByID), id)
}

// FindByUserID mocks base method.
func (m *MockQuizSuiteRepository) FindByUserID(userID uint) ([]*quiz_suite.QuizSuite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserID", userID)
	ret0, _ := ret[0].([]*quiz_suite.QuizSuite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserID indicates an expected call of FindByUserID.
func (mr *MockQuizSuiteRepositoryMockRecorder) FindByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserID", reflect.TypeOf((*MockQuizSuiteRepository)(nil).FindByUserID), userID)
}

// Update mocks base method.
func (m *MockQuizSuiteRepository) Update(quizSuite *quiz_suite.QuizSuite) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", quizSuite)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockQuizSuiteRepositoryMockRecorder) Update(quizSuite interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockQuizSuiteRepository)(nil).Update), quizSuite)
}
