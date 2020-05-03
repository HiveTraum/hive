package repositories

import (
	"auth/config"
	"auth/enums"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateSession(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeSessions(pool, ctx)
	PurgeSecrets(pool, ctx)
	PurgeUsers(pool, ctx)
	secret := CreateSecret(pool, ctx)
	user := CreateUser(pool, ctx)
	status, session := CreateSession(pool, ctx, "123", user.Id, secret.Id, "chrome")
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, session)
	require.NotNil(t, session.RefreshToken)
	require.Equal(t, user.Id, session.UserID)
	require.Equal(t, secret.Id, session.SecretID)
	require.Equal(t, "123", session.Fingerprint)
	require.Equal(t, "chrome", session.UserAgent)
}

func TestGetSession(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeSessions(pool, ctx)
	PurgeSecrets(pool, ctx)
	PurgeUsers(pool, ctx)
	secret := CreateSecret(pool, ctx)
	user := CreateUser(pool, ctx)
	_, createdSession := CreateSession(pool, ctx, "123", user.Id, secret.Id, "chrome")
	session := GetSession(pool, ctx, "123", createdSession.RefreshToken, user.Id)
	require.NotNil(t, session)
	require.Equal(t, createdSession, session)
}
