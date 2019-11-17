package controllers

import (
	"auth/enums"
	"auth/infrastructure"
	"auth/inout"
	"auth/models"
	"context"
)

func CreateUser(
	store infrastructure.StoreInterface,
	esb infrastructure.ESBInterface,
	passwordProcessor infrastructure.PasswordProcessorInterface,
	ctx context.Context,
	body *inout.CreateUserRequestV1) (int, *models.User) {

	var identifiers []models.UserID

	if body.Password == "" {
		return enums.PasswordRequired, nil
	}

	if body.Email == "" && body.Phone == "" {
		return enums.MinimumOneFieldRequired, nil
	}

	body.Password = passwordProcessor.Encode(ctx, body.Password)
	if body.Password == "" {
		return enums.IncorrectPassword, nil
	}

	if len(body.Phone) > 0 {
		phoneStatus, phone := validatePhone(store, ctx, body.Phone, body.PhoneCode)
		if phoneStatus != enums.Ok {
			return phoneStatus, nil
		}

		body.Phone = phone
		_, oldPhone := store.GetPhone(ctx, phone)
		if oldPhone != nil {
			identifiers = append(identifiers, oldPhone.UserId)
		}
	}

	if len(body.Email) > 0 {
		emailStatus, email := validateEmail(store, body.Email, body.EmailCode)
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
