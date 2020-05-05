package controllers

import (
	"hive/enums"
	"hive/functools"
	"hive/models"
	"hive/repositories"
	"context"
	uuid "github.com/satori/go.uuid"
)

func (controller *Controller) CreateUser(ctx context.Context, password, email, emailCode, phone, phoneCode string) (int, *models.User) {

	var identifiers []uuid.UUID

	if password == "" {
		return enums.PasswordRequired, nil
	}

	if email == "" && phone == "" {
		return enums.MinimumOneFieldRequired, nil
	}

	password = controller.passwordProcessor.EncodePassword(ctx, password)
	if password == "" {
		return enums.IncorrectPassword, nil
	}

	if len(phone) > 0 {
		phone = functools.NormalizePhone(phone)
		if phone == "" {
			return enums.IncorrectPhone, nil
		}

		cachedCode := controller.store.GetPhoneConfirmationCode(ctx, phone)
		if cachedCode == "" {
			return enums.PhoneNotFound, nil
		} else if cachedCode != phoneCode {
			return enums.IncorrectPhoneCode, nil
		}

		_, oldPhone := controller.store.GetPhone(ctx, phone)
		if oldPhone != nil {
			identifiers = append(identifiers, oldPhone.UserId)
		}
	}

	if len(email) > 0 {
		emailStatus, email := controller.validateEmail(ctx, email, emailCode)
		if emailStatus != enums.Ok {
			return emailStatus, nil
		}

		_, oldEmail := controller.store.GetEmail(ctx, email)
		if oldEmail != nil {
			identifiers = append(identifiers, oldEmail.UserId)
		}
	}

	status, user := controller.store.CreateUser(ctx, password, email, phone)
	identifiers = append(identifiers, user.Id)
	controller.OnUserChanged(identifiers)
	return status, user
}

func (controller *Controller) DeleteUser(ctx context.Context, id uuid.UUID) (int, *models.User) {
	status, deletedUser := controller.store.DeleteUser(ctx, id)
	if status == enums.Ok {
		controller.OnUserChanged([]uuid.UUID{deletedUser.Id})
	}

	return status, deletedUser
}

func (controller *Controller) GetUsers(ctx context.Context, query repositories.GetUsersQuery) []*models.User {
	return controller.store.GetUsers(ctx, query)
}

func (controller *Controller) GetUser(ctx context.Context, id uuid.UUID) *models.User {
	return controller.store.GetUser(ctx, id)
}
