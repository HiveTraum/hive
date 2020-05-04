package backends

import (
	"auth/models"
	"context"
)

type IAuthenticationBackend interface {
	GetUser(ctx context.Context, token string) (int, models.IAuthenticationBackendUser)
}
