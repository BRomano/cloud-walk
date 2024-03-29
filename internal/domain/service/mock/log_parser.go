// Code generated by MockGen. DO NOT EDIT.
// Source: log_parser.go
//
// Generated by this command:
//
//	mockgen -source log_parser.go -destination mock/log_parser.go -package=mock
//
// Package mock is a generated GoMock package.
package mock

import (
	domain "cloud-walk/internal/domain"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockLogParserService is a mock of LogParserService interface.
type MockLogParserService struct {
	ctrl     *gomock.Controller
	recorder *MockLogParserServiceMockRecorder
}

// MockLogParserServiceMockRecorder is the mock recorder for MockLogParserService.
type MockLogParserServiceMockRecorder struct {
	mock *MockLogParserService
}

// NewMockLogParserService creates a new mock instance.
func NewMockLogParserService(ctrl *gomock.Controller) *MockLogParserService {
	mock := &MockLogParserService{ctrl: ctrl}
	mock.recorder = &MockLogParserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogParserService) EXPECT() *MockLogParserServiceMockRecorder {
	return m.recorder
}

// GetKillsByMeans mocks base method.
func (m *MockLogParserService) GetKillsByMeans(ctx context.Context, gameID int, logger []byte) (map[string]domain.MatchDeathStatistics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKillsByMeans", ctx, gameID, logger)
	ret0, _ := ret[0].(map[string]domain.MatchDeathStatistics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKillsByMeans indicates an expected call of GetKillsByMeans.
func (mr *MockLogParserServiceMockRecorder) GetKillsByMeans(ctx, gameID, logger any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKillsByMeans", reflect.TypeOf((*MockLogParserService)(nil).GetKillsByMeans), ctx, gameID, logger)
}

// GetMatchesStatistics mocks base method.
func (m *MockLogParserService) GetMatchesStatistics(ctx context.Context, gameID int, logger []byte) (map[string]domain.MatchStatistics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatchesStatistics", ctx, gameID, logger)
	ret0, _ := ret[0].(map[string]domain.MatchStatistics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatchesStatistics indicates an expected call of GetMatchesStatistics.
func (mr *MockLogParserServiceMockRecorder) GetMatchesStatistics(ctx, gameID, logger any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchesStatistics", reflect.TypeOf((*MockLogParserService)(nil).GetMatchesStatistics), ctx, gameID, logger)
}
