package app

import (
	"auth/config"
	"auth/infrastructure"
	"context"
)

func InitialAdminRole(store infrastructure.StoreInterface) {
	env := config.GetEnvironment()
	ctx := context.Background()
	_, role := store.GetAdminRole(ctx)
	if role != nil {
		_, _ = store.CreateRole(ctx, env.AdminRole)
	}
}
