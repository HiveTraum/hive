package controllers

import (
	"hive/enums"
	"hive/functools"
	"hive/models"
	"context"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (controller *Controller) CreatePhone(ctx context.Context, phone string, code string, userId uuid.UUID) (int, *models.Phone) {

	phone = functools.NormalizePhone(phone)
	if phone == "" {
		return enums.IncorrectPhone, nil
	}

	cachedCode := controller.store.GetPhoneConfirmationCode(ctx, phone)
	if cachedCode == "" {
		return enums.PhoneNotFound, nil
	} else if cachedCode != code {
		return enums.IncorrectPhoneCode, nil
	}

	_, oldPhone := controller.store.GetPhone(ctx, phone)

	identifiers := []uuid.UUID{userId}

	if oldPhone != nil {
		identifiers = append(identifiers, oldPhone.UserId)
	}

	status, phoneObject := controller.store.CreatePhone(ctx, userId, phone)
	controller.OnPhoneChanged(identifiers)
	return status, phoneObject
}

func (controller *Controller) CreatePhoneConfirmation(ctx context.Context, phone string) (int, *models.PhoneConfirmation) {

	phone = functools.NormalizePhone(phone)
	if phone == "" {
		return enums.IncorrectPhone, nil
	}

	code := controller.store.GetRandomCodeForPhoneConfirmation()
	phoneConfirmation := controller.store.CreatePhoneConfirmationCode(ctx, phone, code, time.Minute*15)
	controller.OnPhoneCodeConfirmationCreated(phone, code)
	return enums.Ok, phoneConfirmation
}
