package models

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type IAuthenticationBackendUser interface {
	GetIsAdmin() bool
	GetRoles() []string
	GetUserID() uuid.UUID
}

type AccessTokenEncoder func(_ context.Context, userID uuid.UUID, roles []string, secret *Secret, expires int64) string
