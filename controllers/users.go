package controllers

import (
	"auth/enums"
	"auth/infrastructure"
	"auth/inout"
	"auth/models"
	"context"
)

func CreateUser(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, body *inout.CreateUserRequestV1) (int, *models.User) {

	var identifiers []int64

	if body.Password == "" {
		return enums.PasswordRequired, nil
	}

	if body.Email == "" && body.Phone == "" {
		return enums.MinimumOneFieldRequired, nil
	}

	if len(body.Phone) > 0 {
		phoneStatus, phone := validatePhone(store, body.Phone, body.PhoneCode)
		if phoneStatus != enums.Ok {
			return phoneStatus, nil
		}

		body.Phone = phone
		_, oldPhone := store.GetPhone(ctx, phone)
		if oldPhone != nil {
			identifiers = append(identifiers, oldPhone.Id)
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
			identifiers = append(identifiers, oldEmail.Id)
		}
	}

	status, user := store.CreateUser(ctx, body)
	identifiers = append(identifiers, user.Id)
	esb.OnUserChanged(identifiers)
	return status, user
}
