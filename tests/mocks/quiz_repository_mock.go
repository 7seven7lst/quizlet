// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/quiz_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	quiz "quizlet/internal/models/quiz"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockQuizRepository is a mock of QuizRepository interface.
type MockQuizRepository struct {
	ctrl     *gomock.Controller
	recorder *MockQuizRepositoryMockRecorder
}

// MockQuizRepositoryMockRecorder is the mock recorder for MockQuizRepository.
type MockQuizRepositoryMockRecorder struct {
	mock *MockQuizRepository
}

// NewMockQuizRepository creates a new mock instance.
func NewMockQuizRepository(ctrl *gomock.Controller) *MockQuizRepository {
	mock := &MockQuizRepository{ctrl: ctrl}
	mock.recorder = &MockQuizRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuizRepository) EXPECT() *MockQuizRepositoryMockRecorder {
	return m.recorder
}

// AddSelection mocks base method.
func (m *MockQuizRepository) AddSelection(quizID uint, selection quiz.QuizSelection) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSelection", quizID, selection)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSelection indicates an expected call of AddSelection.
func (mr *MockQuizRepositoryMockRecorder) AddSelection(quizID, selection interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSelection", reflect.TypeOf((*MockQuizRepository)(nil).AddSelection), quizID, selection)
}

// Create mocks base method.
func (m *MockQuizRepository) Create(quiz *quiz.Quiz) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", quiz)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockQuizRepositoryMockRecorder) Create(quiz interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockQuizRepository)(nil).Create), quiz)
}

// Delete mocks base method.
func (m *MockQuizRepository) Delete(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockQuizRepositoryMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockQuizRepository)(nil).Delete), id)
}

// FindByID mocks base method.
func (m *MockQuizRepository) FindByID(id uint) (*quiz.Quiz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(*quiz.Quiz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockQuizRepositoryMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockQuizRepository)(nil).FindByID), id)
}

// FindByUserID mocks base method.
func (m *MockQuizRepository) FindByUserID(userID uint) ([]*quiz.Quiz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserID", userID)
	ret0, _ := ret[0].([]*quiz.Quiz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserID indicates an expected call of FindByUserID.
func (mr *MockQuizRepositoryMockRecorder) FindByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserID", reflect.TypeOf((*MockQuizRepository)(nil).FindByUserID), userID)
}

// RemoveSelection mocks base method.
func (m *MockQuizRepository) RemoveSelection(quizID, selectionID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSelection", quizID, selectionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSelection indicates an expected call of RemoveSelection.
func (mr *MockQuizRepositoryMockRecorder) RemoveSelection(quizID, selectionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSelection", reflect.TypeOf((*MockQuizRepository)(nil).RemoveSelection), quizID, selectionID)
}

// Update mocks base method.
func (m *MockQuizRepository) Update(quiz *quiz.Quiz) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", quiz)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockQuizRepositoryMockRecorder) Update(quiz interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockQuizRepository)(nil).Update), quiz)
}
