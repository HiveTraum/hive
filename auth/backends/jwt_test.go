package backends

import (
	"auth/enums"
	"auth/models"
	"context"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateSessionFromTokens(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitJWTAuthenticationBackendWithMockedInternals(ctrl)

	userID := uuid.NewV4()

	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 0,
		Value:   uuid.NewV4(),
	}

	accessToken := backend.Backend.EncodeAccessToken(ctx, userID, []string{}, secret, time.Now().Add(time.Microsecond).Unix())

	backend.
		Store.
		EXPECT().
		GetSecret(ctx, secret.Id).
		Times(1).
		Return(secret)

	status, loggedUser := backend.Backend.GetUser(ctx, accessToken)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, loggedUser.GetUserID())
	require.Equal(t, userID, loggedUser.GetUserID())
}

func TestCreateSessionFromTokensWithoutSecret(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backend := InitJWTAuthenticationBackendWithMockedInternals(ctrl)

	userID := uuid.NewV4()

	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 1,
		Value:   uuid.NewV4(),
	}

	accessToken := backend.Backend.EncodeAccessToken(ctx, userID, []string{}, secret, time.Now().Add(time.Microsecond).Unix())

	backend.
		Store.
		EXPECT().
		GetSecret(ctx, secret.Id).
		Return(nil).
		Times(1)

	status, loggedUserID := backend.Backend.GetUser(ctx, accessToken)
	require.Equal(t, enums.SecretNotFound, status)
	require.Nil(t, loggedUserID)
}
