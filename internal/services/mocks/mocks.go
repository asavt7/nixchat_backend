// Code generated by MockGen. DO NOT EDIT.
// Source: ./services.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"
	time "time"

	model "github.com/asavt7/nixchat_backend/internal/model"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserService) CreateUser(user model.User, password string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user, password)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserServiceMockRecorder) CreateUser(user, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserService)(nil).CreateUser), user, password)
}

// FindByUsernameOrEmail mocks base method.
func (m *MockUserService) FindByUsernameOrEmail(username, email string) ([]model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUsernameOrEmail", username, email)
	ret0, _ := ret[0].([]model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUsernameOrEmail indicates an expected call of FindByUsernameOrEmail.
func (mr *MockUserServiceMockRecorder) FindByUsernameOrEmail(username, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUsernameOrEmail", reflect.TypeOf((*MockUserService)(nil).FindByUsernameOrEmail), username, email)
}

// GetAll mocks base method.
func (m *MockUserService) GetAll(arg0 model.PagedQuery) ([]model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0)
	ret0, _ := ret[0].([]model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockUserServiceMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockUserService)(nil).GetAll), arg0)
}

// GetByID mocks base method.
func (m *MockUserService) GetByID(userID uuid.UUID) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", userID)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockUserServiceMockRecorder) GetByID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUserService)(nil).GetByID), userID)
}

// GetByUsername mocks base method.
func (m *MockUserService) GetByUsername(username string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUsername", username)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUsername indicates an expected call of GetByUsername.
func (mr *MockUserServiceMockRecorder) GetByUsername(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUsername", reflect.TypeOf((*MockUserService)(nil).GetByUsername), username)
}

// Update mocks base method.
func (m *MockUserService) Update(userID uuid.UUID, updateInput model.UpdateUserInfo) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", userID, updateInput)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockUserServiceMockRecorder) Update(userID, updateInput interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserService)(nil).Update), userID, updateInput)
}

// MockAuthorizationService is a mock of AuthorizationService interface.
type MockAuthorizationService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationServiceMockRecorder
}

// MockAuthorizationServiceMockRecorder is the mock recorder for MockAuthorizationService.
type MockAuthorizationServiceMockRecorder struct {
	mock *MockAuthorizationService
}

// NewMockAuthorizationService creates a new mock instance.
func NewMockAuthorizationService(ctrl *gomock.Controller) *MockAuthorizationService {
	mock := &MockAuthorizationService{ctrl: ctrl}
	mock.recorder = &MockAuthorizationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorizationService) EXPECT() *MockAuthorizationServiceMockRecorder {
	return m.recorder
}

// CheckUserCredentials mocks base method.
func (m *MockAuthorizationService) CheckUserCredentials(username, password string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserCredentials", username, password)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUserCredentials indicates an expected call of CheckUserCredentials.
func (mr *MockAuthorizationServiceMockRecorder) CheckUserCredentials(username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserCredentials", reflect.TypeOf((*MockAuthorizationService)(nil).CheckUserCredentials), username, password)
}

// GenerateTokens mocks base method.
func (m *MockAuthorizationService) GenerateTokens(userID uuid.UUID) (string, string, time.Time, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateTokens", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(time.Time)
	ret3, _ := ret[3].(time.Time)
	ret4, _ := ret[4].(error)
	return ret0, ret1, ret2, ret3, ret4
}

// GenerateTokens indicates an expected call of GenerateTokens.
func (mr *MockAuthorizationServiceMockRecorder) GenerateTokens(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateTokens", reflect.TypeOf((*MockAuthorizationService)(nil).GenerateTokens), userID)
}

// GetAccessSigningKey mocks base method.
func (m *MockAuthorizationService) GetAccessSigningKey() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccessSigningKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAccessSigningKey indicates an expected call of GetAccessSigningKey.
func (mr *MockAuthorizationServiceMockRecorder) GetAccessSigningKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccessSigningKey", reflect.TypeOf((*MockAuthorizationService)(nil).GetAccessSigningKey))
}

// IsNeedToRefresh mocks base method.
func (m *MockAuthorizationService) IsNeedToRefresh(claims *model.Claims) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNeedToRefresh", claims)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNeedToRefresh indicates an expected call of IsNeedToRefresh.
func (mr *MockAuthorizationServiceMockRecorder) IsNeedToRefresh(claims interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNeedToRefresh", reflect.TypeOf((*MockAuthorizationService)(nil).IsNeedToRefresh), claims)
}

// Logout mocks base method.
func (m *MockAuthorizationService) Logout(accessTokenClaims *model.Claims) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", accessTokenClaims)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockAuthorizationServiceMockRecorder) Logout(accessTokenClaims interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAuthorizationService)(nil).Logout), accessTokenClaims)
}

// ParseAccessTokenToClaims mocks base method.
func (m *MockAuthorizationService) ParseAccessTokenToClaims(token string) (*model.Claims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseAccessTokenToClaims", token)
	ret0, _ := ret[0].(*model.Claims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseAccessTokenToClaims indicates an expected call of ParseAccessTokenToClaims.
func (mr *MockAuthorizationServiceMockRecorder) ParseAccessTokenToClaims(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseAccessTokenToClaims", reflect.TypeOf((*MockAuthorizationService)(nil).ParseAccessTokenToClaims), token)
}

// ParseRefreshTokenToClaims mocks base method.
func (m *MockAuthorizationService) ParseRefreshTokenToClaims(token string) (*model.Claims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseRefreshTokenToClaims", token)
	ret0, _ := ret[0].(*model.Claims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseRefreshTokenToClaims indicates an expected call of ParseRefreshTokenToClaims.
func (mr *MockAuthorizationServiceMockRecorder) ParseRefreshTokenToClaims(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseRefreshTokenToClaims", reflect.TypeOf((*MockAuthorizationService)(nil).ParseRefreshTokenToClaims), token)
}

// ValidateAccessToken mocks base method.
func (m *MockAuthorizationService) ValidateAccessToken(accessTokenClaims *model.Claims) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAccessToken", accessTokenClaims)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateAccessToken indicates an expected call of ValidateAccessToken.
func (mr *MockAuthorizationServiceMockRecorder) ValidateAccessToken(accessTokenClaims interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAccessToken", reflect.TypeOf((*MockAuthorizationService)(nil).ValidateAccessToken), accessTokenClaims)
}

// ValidateRefreshToken mocks base method.
func (m *MockAuthorizationService) ValidateRefreshToken(accessTokenClaims *model.Claims) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateRefreshToken", accessTokenClaims)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateRefreshToken indicates an expected call of ValidateRefreshToken.
func (mr *MockAuthorizationServiceMockRecorder) ValidateRefreshToken(accessTokenClaims interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateRefreshToken", reflect.TypeOf((*MockAuthorizationService)(nil).ValidateRefreshToken), accessTokenClaims)
}
