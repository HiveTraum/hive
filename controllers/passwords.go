package controllers

import (
	"hive/enums"
	"hive/models"
	"context"
	uuid "github.com/satori/go.uuid"
)

func (controller *Controller) CreatePassword(ctx context.Context, userId uuid.UUID, value string) (int, *models.Password) {

	value = controller.passwordProcessor.EncodePassword(ctx, value)
	if value == "" {
		return enums.IncorrectPassword, nil
	}

	status, password := controller.store.CreatePassword(ctx, userId, value)
	if status == enums.Ok {
		controller.OnPasswordChanged(userId)
	}

	return status, password
}
