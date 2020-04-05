package controllers

import (
	"auth/enums"
	"auth/functools"
	"auth/infrastructure"
	"auth/inout"
	"auth/models"
	"context"
	uuid "github.com/satori/go.uuid"
)

func CreateUser(
	store infrastructure.StoreInterface,
	esb infrastructure.ESBInterface,
	passwordProcessor infrastructure.PasswordProcessorInterface,
	ctx context.Context,
	body *inout.CreateUserResponseV1_Request) (int, *models.User) {

	var identifiers []uuid.UUID

	if body.Password == "" {
		return enums.PasswordRequired, nil
	}

	if body.Email == "" && body.Phone == "" {
		return enums.MinimumOneFieldRequired, nil
	}

	body.Password = passwordProcessor.EncodePassword(ctx, body.Password)
	if body.Password == "" {
		return enums.IncorrectPassword, nil
	}

	if len(body.Phone) > 0 {
		phone := functools.NormalizePhone(body.Phone, body.PhoneCountryCode)
		if phone == "" {
			return enums.IncorrectPhone, nil
		}

		cachedCode := store.GetPhoneConfirmationCode(ctx, phone)
		if cachedCode == "" {
			return enums.PhoneNotFound, nil
		} else if cachedCode != body.PhoneCode {
			return enums.IncorrectPhoneCode, nil
		}

		body.Phone = phone
		_, oldPhone := store.GetPhone(ctx, phone)
		if oldPhone != nil {
			identifiers = append(identifiers, oldPhone.UserId)
		}
	}

	if len(body.Email) > 0 {
		emailStatus, email := validateEmail(ctx, store, body.Email, body.EmailCode)
		if emailStatus != enums.Ok {
			return emailStatus, nil
		}

		body.Email = email
		_, oldEmail := store.GetEmail(ctx, email)
		if oldEmail != nil {
			identifiers = append(identifiers, oldEmail.UserId)
		}
	}

	status, user := store.CreateUser(ctx, body)
	identifiers = append(identifiers, user.Id)
	esb.OnUserChanged(identifiers)
	return status, user
}

func DeleteUser(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, id uuid.UUID) (int, *models.User) {
	status, deletedUser := store.DeleteUser(ctx, id)
	if status == enums.Ok {
		esb.OnUserChanged([]uuid.UUID{deletedUser.Id})
	}

	return status, deletedUser
}
