package app

import (
	"auth/config"
	"auth/enums"
	"auth/infrastructure"
	"auth/inout"
	"context"
	"strings"
	"time"
)

func InitialAdmin(store infrastructure.StoreInterface) {
	env := config.GetEnvironment()
	if env.InitialAdmin == "" {
		return
	}

	ctx := context.Background()

	emailAndPassword := strings.Split(env.InitialAdmin, ":")
	emailValue := emailAndPassword[0]
	passwordValue := emailAndPassword[1]
	status, email := store.GetEmail(ctx, emailValue)

	if status != enums.Ok || email != nil {
		return
	}

	store.CreateEmailConfirmationCode(ctx, emailValue, env.TestConfirmationCode, time.Minute)
	_, user := store.CreateUser(ctx, &inout.CreateUserRequestV1{
		Password:  passwordValue,
		Email:     emailValue,
		EmailCode: env.TestConfirmationCode,
	})
	_, role := store.GetAdminRole(ctx)
	store.CreateUserRole(ctx, user.Id, role.Id)
}
