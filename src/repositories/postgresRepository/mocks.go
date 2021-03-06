// Code generated by MockGen. DO NOT EDIT.
// Source: main.go

// Package postgresRepository is a generated GoMock package.
package postgresRepository

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	go_uuid "github.com/satori/go.uuid"
	models "hive/models"
	reflect "reflect"
)

// MockIPostgresRepository is a mock of IPostgresRepository interface
type MockIPostgresRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIPostgresRepositoryMockRecorder
}

// MockIPostgresRepositoryMockRecorder is the mock recorder for MockIPostgresRepository
type MockIPostgresRepositoryMockRecorder struct {
	mock *MockIPostgresRepository
}

// NewMockIPostgresRepository creates a new mock instance
func NewMockIPostgresRepository(ctrl *gomock.Controller) *MockIPostgresRepository {
	mock := &MockIPostgresRepository{ctrl: ctrl}
	mock.recorder = &MockIPostgresRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIPostgresRepository) EXPECT() *MockIPostgresRepositoryMockRecorder {
	return m.recorder
}

// CreateSecret mocks base method
func (m *MockIPostgresRepository) CreateSecret(ctx context.Context) *models.Secret {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSecret", ctx)
	ret0, _ := ret[0].(*models.Secret)
	return ret0
}

// CreateSecret indicates an expected call of CreateSecret
func (mr *MockIPostgresRepositoryMockRecorder) CreateSecret(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSecret", reflect.TypeOf((*MockIPostgresRepository)(nil).CreateSecret), ctx)
}

// GetSecret mocks base method
func (m *MockIPostgresRepository) GetSecret(ctx context.Context, id go_uuid.UUID) *models.Secret {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", ctx, id)
	ret0, _ := ret[0].(*models.Secret)
	return ret0
}

// GetSecret indicates an expected call of GetSecret
func (mr *MockIPostgresRepositoryMockRecorder) GetSecret(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockIPostgresRepository)(nil).GetSecret), ctx, id)
}

// CreateSession mocks base method
func (m *MockIPostgresRepository) CreateSession(ctx context.Context, userID, secretID go_uuid.UUID, fingerprint, userAgent string) *models.Session {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, userID, secretID, fingerprint, userAgent)
	ret0, _ := ret[0].(*models.Session)
	return ret0
}

// CreateSession indicates an expected call of CreateSession
func (mr *MockIPostgresRepositoryMockRecorder) CreateSession(ctx, userID, secretID, fingerprint, userAgent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockIPostgresRepository)(nil).CreateSession), ctx, userID, secretID, fingerprint, userAgent)
}

// DeleteSession mocks base method
func (m *MockIPostgresRepository) DeleteSession(ctx context.Context, id go_uuid.UUID) *models.Session {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", ctx, id)
	ret0, _ := ret[0].(*models.Session)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession
func (mr *MockIPostgresRepositoryMockRecorder) DeleteSession(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockIPostgresRepository)(nil).DeleteSession), ctx, id)
}
