package backends

import (
	"auth/enums"
	"auth/models"
	"context"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateSessionFromTokens(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitJWTAuthenticationBackendWithMockedInternals(ctrl)

	userID := uuid.NewV4()
	fingerprint := "123"
	refreshToken := "321"
	accessToken := "123321"

	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 0,
		Value:   uuid.NewV4(),
	}

	backend.
		Store.
		EXPECT().
		GetSession(ctx, fingerprint, refreshToken, userID).
		Times(1).
		Return(&models.Session{
			RefreshToken: refreshToken,
			Fingerprint:  fingerprint,
			UserID:       userID,
			SecretID:     secret.Id,
			Created:      1,
			UserAgent:    "chrome",
			AccessToken:  accessToken,
		})

	status, loggedUser := backend.Backend.GetUser(ctx, accessToken, "")
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, loggedUser.GetUserID())
	require.Equal(t, userID, loggedUser.GetUserID())
}

func TestCreateSessionFromTokensWithIncorrectAuth(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitJWTAuthenticationBackendWithMockedInternals(ctrl)

	fingerprint := "123"
	refreshToken := "321"
	accessToken := "123321"

	backend.loginController.
		EXPECT().
		Login(ctx, accessToken).
		Times(1).
		Return(enums.IncorrectToken, nil, ctx)

	status, loggedUserID := controller.getUserFromTokens(ctx, accessToken, fingerprint, refreshToken)
	require.Equal(t, enums.IncorrectToken, status)
	require.Nil(t, loggedUserID)
}

func TestCreateSessionFromTokensWithoutSecret(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	controller, _, _, loginController, _ := InitControllerWithMockedInternals(ctrl)

	fingerprint := "123"
	refreshToken := "321"
	accessToken := "123321"

	loginController.
		EXPECT().
		Login(ctx, accessToken).
		Times(1).
		Return(enums.SecretNotFound, nil, ctx)

	status, loggedUserID := controller.getUserFromTokens(ctx, accessToken, fingerprint, refreshToken)
	require.Equal(t, enums.SecretNotFound, status)
	require.Nil(t, loggedUserID)
}
