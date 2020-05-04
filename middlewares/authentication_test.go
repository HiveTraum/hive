package middlewares

import (
	"auth/auth"
	"auth/auth/backends"
	"auth/enums"
	"auth/models"
	"auth/repositories"
	"bytes"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationMiddleware(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authController := auth.NewMockIAuthenticationController(ctrl)
	middleware := AuthenticationMiddleware(authController)

	request := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte{}))
	recorder := httptest.NewRecorder()
	ctx := request.Context()

	userID := uuid.NewV4()

	authController.
		EXPECT().
		Login(ctx, request).
		Times(1).
		Return(enums.Ok, &backends.BasicAuthenticationBackendUser{
			IsAdmin: true,
			Roles:   nil,
			UserID:  userID,
		})

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(1111)
	}), true)
	handler.ServeHTTP(recorder, request)
	result := recorder.Result()
	require.Equal(t, 1111, result.StatusCode)
}

func TestAuthenticationMiddlewareUserInContext(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authController := auth.NewMockIAuthenticationController(ctrl)
	middleware := AuthenticationMiddleware(authController)

	request := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte{}))
	recorder := httptest.NewRecorder()
	ctx := request.Context()

	userID := uuid.NewV4()

	authController.
		EXPECT().
		Login(ctx, request).
		Times(1).
		Return(enums.Ok, &backends.BasicAuthenticationBackendUser{
			IsAdmin: true,
			Roles:   nil,
			UserID:  userID,
		})

	var user models.IAuthenticationBackendUser

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user = repositories.GetUserFromContext(r.Context())
	}), true)
	handler.ServeHTTP(recorder, request)
	require.NotNil(t, user)
	require.Equal(t, userID, user.GetUserID())
}

func TestAuthenticationMiddlewareUnauthenticated(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authController := auth.NewMockIAuthenticationController(ctrl)
	middleware := AuthenticationMiddleware(authController)

	request := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte{}))
	recorder := httptest.NewRecorder()
	ctx := request.Context()

	authController.
		EXPECT().
		Login(ctx, request).
		Times(1).
		Return(enums.Ok, nil)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(1111)
	}), true)
	handler.ServeHTTP(recorder, request)
	result := recorder.Result()
	require.Equal(t, http.StatusUnauthorized, result.StatusCode)
}

func TestAuthenticationMiddlewareUnauthenticatedAndUserNotInContext(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authController := auth.NewMockIAuthenticationController(ctrl)
	middleware := AuthenticationMiddleware(authController)

	request := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte{}))
	recorder := httptest.NewRecorder()
	ctx := request.Context()

	authController.
		EXPECT().
		Login(ctx, request).
		Times(1).
		Return(enums.Ok, nil)

	var user models.IAuthenticationBackendUser
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user = repositories.GetUserFromContext(r.Context())
	}), true)
	handler.ServeHTTP(recorder, request)
	require.Nil(t, user)
}

func TestAuthenticationMiddlewareAuthNotRequiredAndNotAuthenticated(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authController := auth.NewMockIAuthenticationController(ctrl)
	middleware := AuthenticationMiddleware(authController)

	request := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte{}))
	recorder := httptest.NewRecorder()
	ctx := request.Context()

	authController.
		EXPECT().
		Login(ctx, request).
		Times(1).
		Return(enums.Ok, nil)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(1111)
	}), false)
	handler.ServeHTTP(recorder, request)
	result := recorder.Result()
	require.Equal(t, 1111, result.StatusCode)
}
