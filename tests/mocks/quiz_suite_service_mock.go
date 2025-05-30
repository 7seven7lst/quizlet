// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/quiz_suite_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	quiz_suite "quizlet/internal/models/quiz_suite"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockQuizSuiteService is a mock of QuizSuiteService interface.
type MockQuizSuiteService struct {
	ctrl     *gomock.Controller
	recorder *MockQuizSuiteServiceMockRecorder
}

// MockQuizSuiteServiceMockRecorder is the mock recorder for MockQuizSuiteService.
type MockQuizSuiteServiceMockRecorder struct {
	mock *MockQuizSuiteService
}

// NewMockQuizSuiteService creates a new mock instance.
func NewMockQuizSuiteService(ctrl *gomock.Controller) *MockQuizSuiteService {
	mock := &MockQuizSuiteService{ctrl: ctrl}
	mock.recorder = &MockQuizSuiteServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuizSuiteService) EXPECT() *MockQuizSuiteServiceMockRecorder {
	return m.recorder
}

// AddQuizToSuite mocks base method.
func (m *MockQuizSuiteService) AddQuizToSuite(quizSuiteID, quizID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddQuizToSuite", quizSuiteID, quizID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddQuizToSuite indicates an expected call of AddQuizToSuite.
func (mr *MockQuizSuiteServiceMockRecorder) AddQuizToSuite(quizSuiteID, quizID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddQuizToSuite", reflect.TypeOf((*MockQuizSuiteService)(nil).AddQuizToSuite), quizSuiteID, quizID)
}

// CreateQuizSuite mocks base method.
func (m *MockQuizSuiteService) CreateQuizSuite(quizSuite *quiz_suite.QuizSuite) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateQuizSuite", quizSuite)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateQuizSuite indicates an expected call of CreateQuizSuite.
func (mr *MockQuizSuiteServiceMockRecorder) CreateQuizSuite(quizSuite interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateQuizSuite", reflect.TypeOf((*MockQuizSuiteService)(nil).CreateQuizSuite), quizSuite)
}

// DeleteQuizSuite mocks base method.
func (m *MockQuizSuiteService) DeleteQuizSuite(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteQuizSuite", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteQuizSuite indicates an expected call of DeleteQuizSuite.
func (mr *MockQuizSuiteServiceMockRecorder) DeleteQuizSuite(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteQuizSuite", reflect.TypeOf((*MockQuizSuiteService)(nil).DeleteQuizSuite), id)
}

// GetQuizSuite mocks base method.
func (m *MockQuizSuiteService) GetQuizSuite(id uint) (*quiz_suite.QuizSuite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuizSuite", id)
	ret0, _ := ret[0].(*quiz_suite.QuizSuite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuizSuite indicates an expected call of GetQuizSuite.
func (mr *MockQuizSuiteServiceMockRecorder) GetQuizSuite(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuizSuite", reflect.TypeOf((*MockQuizSuiteService)(nil).GetQuizSuite), id)
}

// GetUserQuizSuites mocks base method.
func (m *MockQuizSuiteService) GetUserQuizSuites(userID uint) ([]*quiz_suite.QuizSuite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserQuizSuites", userID)
	ret0, _ := ret[0].([]*quiz_suite.QuizSuite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserQuizSuites indicates an expected call of GetUserQuizSuites.
func (mr *MockQuizSuiteServiceMockRecorder) GetUserQuizSuites(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserQuizSuites", reflect.TypeOf((*MockQuizSuiteService)(nil).GetUserQuizSuites), userID)
}

// RemoveQuizFromSuite mocks base method.
func (m *MockQuizSuiteService) RemoveQuizFromSuite(quizSuiteID, quizID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveQuizFromSuite", quizSuiteID, quizID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveQuizFromSuite indicates an expected call of RemoveQuizFromSuite.
func (mr *MockQuizSuiteServiceMockRecorder) RemoveQuizFromSuite(quizSuiteID, quizID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveQuizFromSuite", reflect.TypeOf((*MockQuizSuiteService)(nil).RemoveQuizFromSuite), quizSuiteID, quizID)
}

// UpdateQuizSuite mocks base method.
func (m *MockQuizSuiteService) UpdateQuizSuite(quizSuite *quiz_suite.QuizSuite) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateQuizSuite", quizSuite)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateQuizSuite indicates an expected call of UpdateQuizSuite.
func (mr *MockQuizSuiteServiceMockRecorder) UpdateQuizSuite(quizSuite interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateQuizSuite", reflect.TypeOf((*MockQuizSuiteService)(nil).UpdateQuizSuite), quizSuite)
}
