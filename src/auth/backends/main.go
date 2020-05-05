package backends

import (
	"hive/models"
	"context"
)

type IAuthenticationBackend interface {
	GetUser(ctx context.Context, token string) (int, models.IAuthenticationBackendUser)
}
