package postgresRepository

import (
	"hive/config"
	"hive/repositories"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateSession(t *testing.T) {
	env := config.InitEnvironment()
	pool := config.InitPool(nil, env)
	repo := InitPostgresRepository(pool, env)
	ctx := context.Background()
	PurgeSessions(pool, ctx)
	PurgeSecrets(pool, ctx)
	PurgeUsers(pool, ctx)
	secret := repo.CreateSecret(ctx)
	user := repositories.CreateUser(pool, ctx)
	session := repo.CreateSession(ctx, user.Id, secret.Id, "123", "chrome")
	require.NotNil(t, session)
	require.NotNil(t, session.RefreshToken)
	require.Equal(t, user.Id, session.UserID)
	require.Equal(t, secret.Id, session.SecretID)
	require.Equal(t, "123", session.Fingerprint)
	require.Equal(t, "chrome", session.UserAgent)
}

func TestGetSession(t *testing.T) {
	env := config.InitEnvironment()
	pool := config.InitPool(nil, env)
	repo := InitPostgresRepository(pool, env)
	ctx := context.Background()
	PurgeSessions(pool, ctx)
	PurgeSecrets(pool, ctx)
	PurgeUsers(pool, ctx)
	secret := repo.CreateSecret(ctx)
	user := repositories.CreateUser(pool, ctx)
	createdSession := repo.CreateSession(ctx, user.Id, secret.Id, "123", "chrome")
	session := repo.GetSession(ctx, createdSession.Id)
	require.NotNil(t, session)
	require.Equal(t, createdSession, session)
}

func TestPostgresRepository_DeleteSession(t *testing.T) {
	env := config.InitEnvironment()
	pool := config.InitPool(nil, env)
	repo := InitPostgresRepository(pool, env)
	ctx := context.Background()
	PurgeSessions(pool, ctx)
	PurgeSecrets(pool, ctx)
	PurgeUsers(pool, ctx)
	secret := repo.CreateSecret(ctx)
	user := repositories.CreateUser(pool, ctx)
	createdSession := repo.CreateSession(ctx, user.Id, secret.Id, "123", "chrome")
	deletedSession := repo.DeleteSession(ctx, createdSession.Id)
	session := repo.GetSession(ctx, createdSession.Id)
	require.NotNil(t, deletedSession)
	require.Equal(t, createdSession, deletedSession)
	require.Nil(t, session)
}