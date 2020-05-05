// Code generated by MockGen. DO NOT EDIT.
// Source: main.go

// Package api is a generated GoMock package.
package api

import (
	gomock "github.com/golang/mock/gomock"
	auth "hive/auth"
	controllers "hive/controllers"
	http "net/http"
	reflect "reflect"
)

// MockIAPI is a mock of IAPI interface
type MockIAPI struct {
	ctrl     *gomock.Controller
	recorder *MockIAPIMockRecorder
}

// MockIAPIMockRecorder is the mock recorder for MockIAPI
type MockIAPIMockRecorder struct {
	mock *MockIAPI
}

// NewMockIAPI creates a new mock instance
func NewMockIAPI(ctrl *gomock.Controller) *MockIAPI {
	mock := &MockIAPI{ctrl: ctrl}
	mock.recorder = &MockIAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIAPI) EXPECT() *MockIAPIMockRecorder {
	return m.recorder
}

// GetAuthenticationController mocks base method
func (m *MockIAPI) GetAuthenticationController() auth.IAuthenticationController {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthenticationController")
	ret0, _ := ret[0].(auth.IAuthenticationController)
	return ret0
}

// GetAuthenticationController indicates an expected call of GetAuthenticationController
func (mr *MockIAPIMockRecorder) GetAuthenticationController() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthenticationController", reflect.TypeOf((*MockIAPI)(nil).GetAuthenticationController))
}

// GetController mocks base method
func (m *MockIAPI) GetController() controllers.IController {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetController")
	ret0, _ := ret[0].(controllers.IController)
	return ret0
}

// GetController indicates an expected call of GetController
func (mr *MockIAPIMockRecorder) GetController() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetController", reflect.TypeOf((*MockIAPI)(nil).GetController))
}

// CreateEmailV1 mocks base method
func (m *MockIAPI) CreateEmailV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateEmailV1", w, r)
}

// CreateEmailV1 indicates an expected call of CreateEmailV1
func (mr *MockIAPIMockRecorder) CreateEmailV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEmailV1", reflect.TypeOf((*MockIAPI)(nil).CreateEmailV1), w, r)
}

// CreateEmailConfirmationV1 mocks base method
func (m *MockIAPI) CreateEmailConfirmationV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateEmailConfirmationV1", w, r)
}

// CreateEmailConfirmationV1 indicates an expected call of CreateEmailConfirmationV1
func (mr *MockIAPIMockRecorder) CreateEmailConfirmationV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEmailConfirmationV1", reflect.TypeOf((*MockIAPI)(nil).CreateEmailConfirmationV1), w, r)
}

// CreatePasswordV1 mocks base method
func (m *MockIAPI) CreatePasswordV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreatePasswordV1", w, r)
}

// CreatePasswordV1 indicates an expected call of CreatePasswordV1
func (mr *MockIAPIMockRecorder) CreatePasswordV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePasswordV1", reflect.TypeOf((*MockIAPI)(nil).CreatePasswordV1), w, r)
}

// CreatePhoneV1 mocks base method
func (m *MockIAPI) CreatePhoneV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreatePhoneV1", w, r)
}

// CreatePhoneV1 indicates an expected call of CreatePhoneV1
func (mr *MockIAPIMockRecorder) CreatePhoneV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePhoneV1", reflect.TypeOf((*MockIAPI)(nil).CreatePhoneV1), w, r)
}

// CreatePhoneConfirmationV1 mocks base method
func (m *MockIAPI) CreatePhoneConfirmationV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreatePhoneConfirmationV1", w, r)
}

// CreatePhoneConfirmationV1 indicates an expected call of CreatePhoneConfirmationV1
func (mr *MockIAPIMockRecorder) CreatePhoneConfirmationV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePhoneConfirmationV1", reflect.TypeOf((*MockIAPI)(nil).CreatePhoneConfirmationV1), w, r)
}

// CreateRoleV1 mocks base method
func (m *MockIAPI) CreateRoleV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateRoleV1", w, r)
}

// CreateRoleV1 indicates an expected call of CreateRoleV1
func (mr *MockIAPIMockRecorder) CreateRoleV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRoleV1", reflect.TypeOf((*MockIAPI)(nil).CreateRoleV1), w, r)
}

// GetRolesV1 mocks base method
func (m *MockIAPI) GetRolesV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetRolesV1", w, r)
}

// GetRolesV1 indicates an expected call of GetRolesV1
func (mr *MockIAPIMockRecorder) GetRolesV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRolesV1", reflect.TypeOf((*MockIAPI)(nil).GetRolesV1), w, r)
}

// GetRoleV1 mocks base method
func (m *MockIAPI) GetRoleV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetRoleV1", w, r)
}

// GetRoleV1 indicates an expected call of GetRoleV1
func (mr *MockIAPIMockRecorder) GetRoleV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoleV1", reflect.TypeOf((*MockIAPI)(nil).GetRoleV1), w, r)
}

// GetSecretV1 mocks base method
func (m *MockIAPI) GetSecretV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetSecretV1", w, r)
}

// GetSecretV1 indicates an expected call of GetSecretV1
func (mr *MockIAPIMockRecorder) GetSecretV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretV1", reflect.TypeOf((*MockIAPI)(nil).GetSecretV1), w, r)
}

// CreateSessionV1 mocks base method
func (m *MockIAPI) CreateSessionV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateSessionV1", w, r)
}

// CreateSessionV1 indicates an expected call of CreateSessionV1
func (mr *MockIAPIMockRecorder) CreateSessionV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSessionV1", reflect.TypeOf((*MockIAPI)(nil).CreateSessionV1), w, r)
}

// GetUserV1 mocks base method
func (m *MockIAPI) GetUserV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetUserV1", w, r)
}

// GetUserV1 indicates an expected call of GetUserV1
func (mr *MockIAPIMockRecorder) GetUserV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserV1", reflect.TypeOf((*MockIAPI)(nil).GetUserV1), w, r)
}

// DeleteUserV1 mocks base method
func (m *MockIAPI) DeleteUserV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteUserV1", w, r)
}

// DeleteUserV1 indicates an expected call of DeleteUserV1
func (mr *MockIAPIMockRecorder) DeleteUserV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserV1", reflect.TypeOf((*MockIAPI)(nil).DeleteUserV1), w, r)
}

// GetUsersV1 mocks base method
func (m *MockIAPI) GetUsersV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetUsersV1", w, r)
}

// GetUsersV1 indicates an expected call of GetUsersV1
func (mr *MockIAPIMockRecorder) GetUsersV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersV1", reflect.TypeOf((*MockIAPI)(nil).GetUsersV1), w, r)
}

// CreateUserV1 mocks base method
func (m *MockIAPI) CreateUserV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateUserV1", w, r)
}

// CreateUserV1 indicates an expected call of CreateUserV1
func (mr *MockIAPIMockRecorder) CreateUserV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserV1", reflect.TypeOf((*MockIAPI)(nil).CreateUserV1), w, r)
}

// GetUserViewV1 mocks base method
func (m *MockIAPI) GetUserViewV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetUserViewV1", w, r)
}

// GetUserViewV1 indicates an expected call of GetUserViewV1
func (mr *MockIAPIMockRecorder) GetUserViewV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserViewV1", reflect.TypeOf((*MockIAPI)(nil).GetUserViewV1), w, r)
}

// GetUsersViewV1 mocks base method
func (m *MockIAPI) GetUsersViewV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetUsersViewV1", w, r)
}

// GetUsersViewV1 indicates an expected call of GetUsersViewV1
func (mr *MockIAPIMockRecorder) GetUsersViewV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersViewV1", reflect.TypeOf((*MockIAPI)(nil).GetUsersViewV1), w, r)
}

// DeleteUserRoleV1 mocks base method
func (m *MockIAPI) DeleteUserRoleV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteUserRoleV1", w, r)
}

// DeleteUserRoleV1 indicates an expected call of DeleteUserRoleV1
func (mr *MockIAPIMockRecorder) DeleteUserRoleV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserRoleV1", reflect.TypeOf((*MockIAPI)(nil).DeleteUserRoleV1), w, r)
}

// CreateUserRoleV1 mocks base method
func (m *MockIAPI) CreateUserRoleV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateUserRoleV1", w, r)
}

// CreateUserRoleV1 indicates an expected call of CreateUserRoleV1
func (mr *MockIAPIMockRecorder) CreateUserRoleV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserRoleV1", reflect.TypeOf((*MockIAPI)(nil).CreateUserRoleV1), w, r)
}

// GetUserRolesV1 mocks base method
func (m *MockIAPI) GetUserRolesV1(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetUserRolesV1", w, r)
}

// GetUserRolesV1 indicates an expected call of GetUserRolesV1
func (mr *MockIAPIMockRecorder) GetUserRolesV1(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserRolesV1", reflect.TypeOf((*MockIAPI)(nil).GetUserRolesV1), w, r)
}
