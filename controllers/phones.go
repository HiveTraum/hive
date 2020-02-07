package controllers

import (
	"auth/enums"
	"auth/functools"
	"auth/infrastructure"
	"auth/models"
	"context"
	uuid "github.com/satori/go.uuid"
	"time"
)

func CreatePhone(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, phone string, code string, userId uuid.UUID, phoneCountryCode string) (int, *models.Phone) {

	phone = functools.NormalizePhone(phone, phoneCountryCode)
	if phone == "" {
		return enums.IncorrectPhone, nil
	}

	cachedCode := store.GetPhoneConfirmationCode(ctx, phone)
	if cachedCode == "" {
		return enums.PhoneNotFound, nil
	} else if cachedCode != code {
		return enums.IncorrectPhoneCode, nil
	}

	_, oldPhone := store.GetPhone(ctx, phone)

	identifiers := []uuid.UUID{userId}

	if oldPhone != nil {
		identifiers = append(identifiers, oldPhone.UserId)
	}

	status, phoneObject := store.CreatePhone(ctx, userId, phone, phoneCountryCode)
	esb.OnPhoneChanged(identifiers)
	return status, phoneObject
}

func CreatePhoneConfirmation(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, phone string, countryCode string) (int, *models.PhoneConfirmation) {

	phone = functools.NormalizePhone(phone, countryCode)
	if phone == "" {
		return enums.IncorrectPhone, nil
	}

	code := store.GetRandomCodeForPhoneConfirmation()
	phoneConfirmation := store.CreatePhoneConfirmationCode(ctx, phone, code, time.Minute*15)
	esb.OnPhoneCodeConfirmationCreated(phone, code)
	return enums.Ok, phoneConfirmation
}
