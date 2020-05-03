package backends

import (
	"auth/models"
	"context"
)

type IAuthenticationBackend interface {
	GetUser(ctx context.Context, headerToken string, cookieToken string) (int, models.IAuthenticationBackendUser)
}
