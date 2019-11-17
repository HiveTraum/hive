// Code generated by MockGen. DO NOT EDIT.
// Source: infrastructure.go

// Package mocks is a generated GoMock package.
package mocks

import (
	infrastructure "auth/infrastructure"
	inout "auth/inout"
	models "auth/models"
	repositories "auth/repositories"
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockStoreInterface is a mock of StoreInterface interface
type MockStoreInterface struct {
	ctrl     *gomock.Controller
	recorder *MockStoreInterfaceMockRecorder
}

// MockStoreInterfaceMockRecorder is the mock recorder for MockStoreInterface
type MockStoreInterfaceMockRecorder struct {
	mock *MockStoreInterface
}

// NewMockStoreInterface creates a new mock instance
func NewMockStoreInterface(ctrl *gomock.Controller) *MockStoreInterface {
	mock := &MockStoreInterface{ctrl: ctrl}
	mock.recorder = &MockStoreInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStoreInterface) EXPECT() *MockStoreInterfaceMockRecorder {
	return m.recorder
}

// CreateUser mocks base method
func (m *MockStoreInterface) CreateUser(ctx context.Context, query *inout.CreateUserRequestV1) (int, *models.User) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, query)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(*models.User)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockStoreInterfaceMockRecorder) CreateUser(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStoreInterface)(nil).CreateUser), ctx, query)
}

// GetUser mocks base method
func (m *MockStoreInterface) GetUser(context context.Context, id int64) *models.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", context, id)
	ret0, _ := ret[0].(*models.User)
	return ret0
}

// GetUser indicates an expected call of GetUser
func (mr *MockStoreInterfaceMockRecorder) GetUser(context, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStoreInterface)(nil).GetUser), context, id)
}

// GetUsers mocks base method
func (m *MockStoreInterface) GetUsers(context context.Context, query repositories.GetUsersQuery) []*models.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", context, query)
	ret0, _ := ret[0].([]*models.User)
	return ret0
}

// GetUsers indicates an expected call of GetUsers
func (mr *MockStoreInterfaceMockRecorder) GetUsers(context, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockStoreInterface)(nil).GetUsers), context, query)
}

// GetUsersView mocks base method
func (m *MockStoreInterface) GetUsersView(context context.Context, query repositories.GetUsersViewQuery) []*inout.GetUserViewResponseV1 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersView", context, query)
	ret0, _ := ret[0].([]*inout.GetUserViewResponseV1)
	return ret0
}

// GetUsersView indicates an expected call of GetUsersView
func (mr *MockStoreInterfaceMockRecorder) GetUsersView(context, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersView", reflect.TypeOf((*MockStoreInterface)(nil).GetUsersView), context, query)
}

// GetUserView mocks base method
func (m *MockStoreInterface) GetUserView(context context.Context, id int64) *inout.GetUserViewResponseV1 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserView", context, id)
	ret0, _ := ret[0].(*inout.GetUserViewResponseV1)
	return ret0
}

// GetUserView indicates an expected call of GetUserView
func (mr *MockStoreInterfaceMockRecorder) GetUserView(context, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserView", reflect.TypeOf((*MockStoreInterface)(nil).GetUserView), context, id)
}

// CreateOrUpdateUsersView mocks base method
func (m *MockStoreInterface) CreateOrUpdateUsersView(context context.Context, query repositories.CreateOrUpdateUsersViewQuery) []*inout.GetUserViewResponseV1 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrUpdateUsersView", context, query)
	ret0, _ := ret[0].([]*inout.GetUserViewResponseV1)
	return ret0
}

// CreateOrUpdateUsersView indicates an expected call of CreateOrUpdateUsersView
func (mr *MockStoreInterfaceMockRecorder) CreateOrUpdateUsersView(context, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrUpdateUsersView", reflect.TypeOf((*MockStoreInterface)(nil).CreateOrUpdateUsersView), context, query)
}

// CreateOrUpdateUsersViewByUsersID mocks base method
func (m *MockStoreInterface) CreateOrUpdateUsersViewByUsersID(context context.Context, id []int64) []*inout.GetUserViewResponseV1 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrUpdateUsersViewByUsersID", context, id)
	ret0, _ := ret[0].([]*inout.GetUserViewResponseV1)
	return ret0
}

// CreateOrUpdateUsersViewByUsersID indicates an expected call of CreateOrUpdateUsersViewByUsersID
func (mr *MockStoreInterfaceMockRecorder) CreateOrUpdateUsersViewByUsersID(context, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrUpdateUsersViewByUsersID", reflect.TypeOf((*MockStoreInterface)(nil).CreateOrUpdateUsersViewByUsersID), context, id)
}

// CreateOrUpdateUsersViewByRolesID mocks base method
func (m *MockStoreInterface) CreateOrUpdateUsersViewByRolesID(context context.Context, id []int64) []*inout.GetUserViewResponseV1 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrUpdateUsersViewByRolesID", context, id)
	ret0, _ := ret[0].([]*inout.GetUserViewResponseV1)
	return ret0
}

// CreateOrUpdateUsersViewByRolesID indicates an expected call of CreateOrUpdateUsersViewByRolesID
func (mr *MockStoreInterfaceMockRecorder) CreateOrUpdateUsersViewByRolesID(context, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrUpdateUsersViewByRolesID", reflect.TypeOf((*MockStoreInterface)(nil).CreateOrUpdateUsersViewByRolesID), context, id)
}

// CreateOrUpdateUsersViewByUserID mocks base method
func (m *MockStoreInterface) CreateOrUpdateUsersViewByUserID(context context.Context, id int64) []*inout.GetUserViewResponseV1 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrUpdateUsersViewByUserID", context, id)
	ret0, _ := ret[0].([]*inout.GetUserViewResponseV1)
	return ret0
}

// CreateOrUpdateUsersViewByUserID indicates an expected call of CreateOrUpdateUsersViewByUserID
func (mr *MockStoreInterfaceMockRecorder) CreateOrUpdateUsersViewByUserID(context, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrUpdateUsersViewByUserID", reflect.TypeOf((*MockStoreInterface)(nil).CreateOrUpdateUsersViewByUserID), context, id)
}

// CreateOrUpdateUsersViewByRoleID mocks base method
func (m *MockStoreInterface) CreateOrUpdateUsersViewByRoleID(context context.Context, id int64) []*inout.GetUserViewResponseV1 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrUpdateUsersViewByRoleID", context, id)
	ret0, _ := ret[0].([]*inout.GetUserViewResponseV1)
	return ret0
}

// CreateOrUpdateUsersViewByRoleID indicates an expected call of CreateOrUpdateUsersViewByRoleID
func (mr *MockStoreInterfaceMockRecorder) CreateOrUpdateUsersViewByRoleID(context, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrUpdateUsersViewByRoleID", reflect.TypeOf((*MockStoreInterface)(nil).CreateOrUpdateUsersViewByRoleID), context, id)
}

// CacheUserView mocks base method
func (m *MockStoreInterface) CacheUserView(ctx context.Context, userViews []*inout.GetUserViewResponseV1) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CacheUserView", ctx, userViews)
}

// CacheUserView indicates an expected call of CacheUserView
func (mr *MockStoreInterfaceMockRecorder) CacheUserView(ctx, userViews interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CacheUserView", reflect.TypeOf((*MockStoreInterface)(nil).CacheUserView), ctx, userViews)
}

// CreateEmail mocks base method
func (m *MockStoreInterface) CreateEmail(ctx context.Context, userId int64, value string) (int, *models.Email) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEmail", ctx, userId, value)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(*models.Email)
	return ret0, ret1
}

// CreateEmail indicates an expected call of CreateEmail
func (mr *MockStoreInterfaceMockRecorder) CreateEmail(ctx, userId, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEmail", reflect.TypeOf((*MockStoreInterface)(nil).CreateEmail), ctx, userId, value)
}

// GetEmail mocks base method
func (m *MockStoreInterface) GetEmail(ctx context.Context, email string) (int, *models.Email) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmail", ctx, email)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(*models.Email)
	return ret0, ret1
}

// GetEmail indicates an expected call of GetEmail
func (mr *MockStoreInterfaceMockRecorder) GetEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmail", reflect.TypeOf((*MockStoreInterface)(nil).GetEmail), ctx, email)
}

// CreateEmailConfirmationCode mocks base method
func (m *MockStoreInterface) CreateEmailConfirmationCode(email, code string, duration time.Duration) *models.EmailConfirmation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEmailConfirmationCode", email, code, duration)
	ret0, _ := ret[0].(*models.EmailConfirmation)
	return ret0
}

// CreateEmailConfirmationCode indicates an expected call of CreateEmailConfirmationCode
func (mr *MockStoreInterfaceMockRecorder) CreateEmailConfirmationCode(email, code, duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEmailConfirmationCode", reflect.TypeOf((*MockStoreInterface)(nil).CreateEmailConfirmationCode), email, code, duration)
}

// GetEmailConfirmationCode mocks base method
func (m *MockStoreInterface) GetEmailConfirmationCode(email string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmailConfirmationCode", email)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetEmailConfirmationCode indicates an expected call of GetEmailConfirmationCode
func (mr *MockStoreInterfaceMockRecorder) GetEmailConfirmationCode(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmailConfirmationCode", reflect.TypeOf((*MockStoreInterface)(nil).GetEmailConfirmationCode), email)
}

// CreatePassword mocks base method
func (m *MockStoreInterface) CreatePassword(ctx context.Context, userId int64, value string) (int, *models.Password) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePassword", ctx, userId, value)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(*models.Password)
	return ret0, ret1
}

// CreatePassword indicates an expected call of CreatePassword
func (mr *MockStoreInterfaceMockRecorder) CreatePassword(ctx, userId, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePassword", reflect.TypeOf((*MockStoreInterface)(nil).CreatePassword), ctx, userId, value)
}

// GetPasswords mocks base method
func (m *MockStoreInterface) GetPasswords(ctx context.Context, userId int64) []*models.Password {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPasswords", ctx, userId)
	ret0, _ := ret[0].([]*models.Password)
	return ret0
}

// GetPasswords indicates an expected call of GetPasswords
func (mr *MockStoreInterfaceMockRecorder) GetPasswords(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPasswords", reflect.TypeOf((*MockStoreInterface)(nil).GetPasswords), ctx, userId)
}

// GetLatestPassword mocks base method
func (m *MockStoreInterface) GetLatestPassword(ctx context.Context, userId int64) (int, *models.Password) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestPassword", ctx, userId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(*models.Password)
	return ret0, ret1
}

// GetLatestPassword indicates an expected call of GetLatestPassword
func (mr *MockStoreInterfaceMockRecorder) GetLatestPassword(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestPassword", reflect.TypeOf((*MockStoreInterface)(nil).GetLatestPassword), ctx, userId)
}

// CreatePhone mocks base method
func (m *MockStoreInterface) CreatePhone(ctx context.Context, userId int64, value string) (int, *models.Phone) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePhone", ctx, userId, value)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(*models.Phone)
	return ret0, ret1
}

// CreatePhone indicates an expected call of CreatePhone
func (mr *MockStoreInterfaceMockRecorder) CreatePhone(ctx, userId, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePhone", reflect.TypeOf((*MockStoreInterface)(nil).CreatePhone), ctx, userId, value)
}

// GetPhone mocks base method
func (m *MockStoreInterface) GetPhone(ctx context.Context, phone string) (int, *models.Phone) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPhone", ctx, phone)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(*models.Phone)
	return ret0, ret1
}

// GetPhone indicates an expected call of GetPhone
func (mr *MockStoreInterfaceMockRecorder) GetPhone(ctx, phone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPhone", reflect.TypeOf((*MockStoreInterface)(nil).GetPhone), ctx, phone)
}

// CreatePhoneConfirmationCode mocks base method
func (m *MockStoreInterface) CreatePhoneConfirmationCode(ctx context.Context, phone, code string, duration time.Duration) *models.PhoneConfirmation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePhoneConfirmationCode", ctx, phone, code, duration)
	ret0, _ := ret[0].(*models.PhoneConfirmation)
	return ret0
}

// CreatePhoneConfirmationCode indicates an expected call of CreatePhoneConfirmationCode
func (mr *MockStoreInterfaceMockRecorder) CreatePhoneConfirmationCode(ctx, phone, code, duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePhoneConfirmationCode", reflect.TypeOf((*MockStoreInterface)(nil).CreatePhoneConfirmationCode), ctx, phone, code, duration)
}

// GetPhoneConfirmationCode mocks base method
func (m *MockStoreInterface) GetPhoneConfirmationCode(ctx context.Context, phone string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPhoneConfirmationCode", ctx, phone)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPhoneConfirmationCode indicates an expected call of GetPhoneConfirmationCode
func (mr *MockStoreInterfaceMockRecorder) GetPhoneConfirmationCode(ctx, phone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPhoneConfirmationCode", reflect.TypeOf((*MockStoreInterface)(nil).GetPhoneConfirmationCode), ctx, phone)
}

// CreateRole mocks base method
func (m *MockStoreInterface) CreateRole(context context.Context, title string) (int, *models.Role) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRole", context, title)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(*models.Role)
	return ret0, ret1
}

// CreateRole indicates an expected call of CreateRole
func (mr *MockStoreInterfaceMockRecorder) CreateRole(context, title interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRole", reflect.TypeOf((*MockStoreInterface)(nil).CreateRole), context, title)
}

// GetRole mocks base method
func (m *MockStoreInterface) GetRole(context context.Context, id int64) (int, *models.Role) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRole", context, id)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(*models.Role)
	return ret0, ret1
}

// GetRole indicates an expected call of GetRole
func (mr *MockStoreInterfaceMockRecorder) GetRole(context, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockStoreInterface)(nil).GetRole), context, id)
}

// GetRoles mocks base method
func (m *MockStoreInterface) GetRoles(context context.Context, query repositories.GetRolesQuery) []*models.Role {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoles", context, query)
	ret0, _ := ret[0].([]*models.Role)
	return ret0
}

// GetRoles indicates an expected call of GetRoles
func (mr *MockStoreInterfaceMockRecorder) GetRoles(context, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoles", reflect.TypeOf((*MockStoreInterface)(nil).GetRoles), context, query)
}

// CreateUserRole mocks base method
func (m *MockStoreInterface) CreateUserRole(ctx context.Context, userId, roleId int64) (int, *models.UserRole) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserRole", ctx, userId, roleId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(*models.UserRole)
	return ret0, ret1
}

// CreateUserRole indicates an expected call of CreateUserRole
func (mr *MockStoreInterfaceMockRecorder) CreateUserRole(ctx, userId, roleId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserRole", reflect.TypeOf((*MockStoreInterface)(nil).CreateUserRole), ctx, userId, roleId)
}

// GetUserRoles mocks base method
func (m *MockStoreInterface) GetUserRoles(ctx context.Context, query repositories.GetUserRoleQuery) []*models.UserRole {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserRoles", ctx, query)
	ret0, _ := ret[0].([]*models.UserRole)
	return ret0
}

// GetUserRoles indicates an expected call of GetUserRoles
func (mr *MockStoreInterfaceMockRecorder) GetUserRoles(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserRoles", reflect.TypeOf((*MockStoreInterface)(nil).GetUserRoles), ctx, query)
}

// MockESBInterface is a mock of ESBInterface interface
type MockESBInterface struct {
	ctrl     *gomock.Controller
	recorder *MockESBInterfaceMockRecorder
}

// MockESBInterfaceMockRecorder is the mock recorder for MockESBInterface
type MockESBInterfaceMockRecorder struct {
	mock *MockESBInterface
}

// NewMockESBInterface creates a new mock instance
func NewMockESBInterface(ctrl *gomock.Controller) *MockESBInterface {
	mock := &MockESBInterface{ctrl: ctrl}
	mock.recorder = &MockESBInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockESBInterface) EXPECT() *MockESBInterfaceMockRecorder {
	return m.recorder
}

// OnUserChanged mocks base method
func (m *MockESBInterface) OnUserChanged(id []int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnUserChanged", id)
}

// OnUserChanged indicates an expected call of OnUserChanged
func (mr *MockESBInterfaceMockRecorder) OnUserChanged(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnUserChanged", reflect.TypeOf((*MockESBInterface)(nil).OnUserChanged), id)
}

// OnEmailCodeConfirmationCreated mocks base method
func (m *MockESBInterface) OnEmailCodeConfirmationCreated(email, code string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnEmailCodeConfirmationCreated", email, code)
}

// OnEmailCodeConfirmationCreated indicates an expected call of OnEmailCodeConfirmationCreated
func (mr *MockESBInterfaceMockRecorder) OnEmailCodeConfirmationCreated(email, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnEmailCodeConfirmationCreated", reflect.TypeOf((*MockESBInterface)(nil).OnEmailCodeConfirmationCreated), email, code)
}

// OnPhoneCodeConfirmationCreated mocks base method
func (m *MockESBInterface) OnPhoneCodeConfirmationCreated(phone, code string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnPhoneCodeConfirmationCreated", phone, code)
}

// OnPhoneCodeConfirmationCreated indicates an expected call of OnPhoneCodeConfirmationCreated
func (mr *MockESBInterfaceMockRecorder) OnPhoneCodeConfirmationCreated(phone, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnPhoneCodeConfirmationCreated", reflect.TypeOf((*MockESBInterface)(nil).OnPhoneCodeConfirmationCreated), phone, code)
}

// OnUsersViewChanged mocks base method
func (m *MockESBInterface) OnUsersViewChanged(usersView []*inout.GetUserViewResponseV1) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnUsersViewChanged", usersView)
}

// OnUsersViewChanged indicates an expected call of OnUsersViewChanged
func (mr *MockESBInterfaceMockRecorder) OnUsersViewChanged(usersView interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnUsersViewChanged", reflect.TypeOf((*MockESBInterface)(nil).OnUsersViewChanged), usersView)
}

// OnPasswordChanged mocks base method
func (m *MockESBInterface) OnPasswordChanged(userId int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnPasswordChanged", userId)
}

// OnPasswordChanged indicates an expected call of OnPasswordChanged
func (mr *MockESBInterfaceMockRecorder) OnPasswordChanged(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnPasswordChanged", reflect.TypeOf((*MockESBInterface)(nil).OnPasswordChanged), userId)
}

// OnPhoneChanged mocks base method
func (m *MockESBInterface) OnPhoneChanged(userId []int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnPhoneChanged", userId)
}

// OnPhoneChanged indicates an expected call of OnPhoneChanged
func (mr *MockESBInterfaceMockRecorder) OnPhoneChanged(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnPhoneChanged", reflect.TypeOf((*MockESBInterface)(nil).OnPhoneChanged), userId)
}

// OnEmailChanged mocks base method
func (m *MockESBInterface) OnEmailChanged(userId []int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnEmailChanged", userId)
}

// OnEmailChanged indicates an expected call of OnEmailChanged
func (mr *MockESBInterfaceMockRecorder) OnEmailChanged(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnEmailChanged", reflect.TypeOf((*MockESBInterface)(nil).OnEmailChanged), userId)
}

// OnRoleChanged mocks base method
func (m *MockESBInterface) OnRoleChanged(roleId []int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnRoleChanged", roleId)
}

// OnRoleChanged indicates an expected call of OnRoleChanged
func (mr *MockESBInterfaceMockRecorder) OnRoleChanged(roleId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnRoleChanged", reflect.TypeOf((*MockESBInterface)(nil).OnRoleChanged), roleId)
}

// MockESBDispatcherInterface is a mock of ESBDispatcherInterface interface
type MockESBDispatcherInterface struct {
	ctrl     *gomock.Controller
	recorder *MockESBDispatcherInterfaceMockRecorder
}

// MockESBDispatcherInterfaceMockRecorder is the mock recorder for MockESBDispatcherInterface
type MockESBDispatcherInterfaceMockRecorder struct {
	mock *MockESBDispatcherInterface
}

// NewMockESBDispatcherInterface creates a new mock instance
func NewMockESBDispatcherInterface(ctrl *gomock.Controller) *MockESBDispatcherInterface {
	mock := &MockESBDispatcherInterface{ctrl: ctrl}
	mock.recorder = &MockESBDispatcherInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockESBDispatcherInterface) EXPECT() *MockESBDispatcherInterfaceMockRecorder {
	return m.recorder
}

// Send mocks base method
func (m *MockESBDispatcherInterface) Send(event inout.Event) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Send", event)
}

// Send indicates an expected call of Send
func (mr *MockESBDispatcherInterfaceMockRecorder) Send(event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockESBDispatcherInterface)(nil).Send), event)
}

// MockPasswordProcessorInterface is a mock of PasswordProcessorInterface interface
type MockPasswordProcessorInterface struct {
	ctrl     *gomock.Controller
	recorder *MockPasswordProcessorInterfaceMockRecorder
}

// MockPasswordProcessorInterfaceMockRecorder is the mock recorder for MockPasswordProcessorInterface
type MockPasswordProcessorInterfaceMockRecorder struct {
	mock *MockPasswordProcessorInterface
}

// NewMockPasswordProcessorInterface creates a new mock instance
func NewMockPasswordProcessorInterface(ctrl *gomock.Controller) *MockPasswordProcessorInterface {
	mock := &MockPasswordProcessorInterface{ctrl: ctrl}
	mock.recorder = &MockPasswordProcessorInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPasswordProcessorInterface) EXPECT() *MockPasswordProcessorInterfaceMockRecorder {
	return m.recorder
}

// Encode mocks base method
func (m *MockPasswordProcessorInterface) Encode(arg0 context.Context, arg1 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Encode", arg0, arg1)
	ret0, _ := ret[0].(string)
	return ret0
}

// Encode indicates an expected call of Encode
func (mr *MockPasswordProcessorInterfaceMockRecorder) Encode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encode", reflect.TypeOf((*MockPasswordProcessorInterface)(nil).Encode), arg0, arg1)
}

// MockAppInterface is a mock of AppInterface interface
type MockAppInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAppInterfaceMockRecorder
}

// MockAppInterfaceMockRecorder is the mock recorder for MockAppInterface
type MockAppInterfaceMockRecorder struct {
	mock *MockAppInterface
}

// NewMockAppInterface creates a new mock instance
func NewMockAppInterface(ctrl *gomock.Controller) *MockAppInterface {
	mock := &MockAppInterface{ctrl: ctrl}
	mock.recorder = &MockAppInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppInterface) EXPECT() *MockAppInterfaceMockRecorder {
	return m.recorder
}

// GetStore mocks base method
func (m *MockAppInterface) GetStore() infrastructure.StoreInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStore")
	ret0, _ := ret[0].(infrastructure.StoreInterface)
	return ret0
}

// GetStore indicates an expected call of GetStore
func (mr *MockAppInterfaceMockRecorder) GetStore() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStore", reflect.TypeOf((*MockAppInterface)(nil).GetStore))
}

// GetESB mocks base method
func (m *MockAppInterface) GetESB() infrastructure.ESBInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetESB")
	ret0, _ := ret[0].(infrastructure.ESBInterface)
	return ret0
}

// GetESB indicates an expected call of GetESB
func (mr *MockAppInterfaceMockRecorder) GetESB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetESB", reflect.TypeOf((*MockAppInterface)(nil).GetESB))
}

// GetPasswordProcessor mocks base method
func (m *MockAppInterface) GetPasswordProcessor() infrastructure.PasswordProcessorInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPasswordProcessor")
	ret0, _ := ret[0].(infrastructure.PasswordProcessorInterface)
	return ret0
}

// GetPasswordProcessor indicates an expected call of GetPasswordProcessor
func (mr *MockAppInterfaceMockRecorder) GetPasswordProcessor() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPasswordProcessor", reflect.TypeOf((*MockAppInterface)(nil).GetPasswordProcessor))
}
