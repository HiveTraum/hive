package controllers

import (
	"hive/eventDispatchers"
	"hive/functools"
	"hive/inout"
	"hive/models"
	"hive/stores"
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type EventController struct {
	Store      stores.IStore
	Dispatcher eventDispatchers.IEventDispatcher
	controller IController
}

// Private methods / Implementation

func (controller *Controller) onUserChanged(userId []uuid.UUID) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnUserChanged")
	controller.CreateOrUpdateUsersView(ctx, userId)
	span.LogFields(log.String("user_id", strings.Join(functools.UUIDListToStringList(userId), ", ")))
	span.Finish()
	ctx.Done()
}

func (controller *Controller) onPhoneChanged(userId []uuid.UUID) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnPhoneChanged")
	controller.CreateOrUpdateUsersView(ctx, userId)
	span.LogFields(log.String("user_id", strings.Join(functools.UUIDListToStringList(userId), ", ")))
	span.Finish()
	ctx.Done()
}

func (controller *Controller) onEmailChanged(userId []uuid.UUID) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnEmailChanged")
	controller.CreateOrUpdateUsersView(ctx, userId)
	span.LogFields(log.String("user_id", strings.Join(functools.UUIDListToStringList(userId), ", ")))
	span.Finish()
	ctx.Done()
}

func (controller *Controller) onRoleChanged(roleId []uuid.UUID) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnRoleChanged")
	controller.CreateOrUpdateUsersViewByRoles(ctx, roleId)
	span.LogFields(log.String("role_id", strings.Join(functools.UUIDListToStringList(roleId), "")))
	span.Finish()
	ctx.Done()
}

func (controller *Controller) onUsersViewChanged(usersView []*models.UserView) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnUserViewChanged")
	identifiers := make([]uuid.UUID, len(usersView))

	for i, u := range usersView {
		identifiers[i] = u.Id
	}

	controller.store.CacheUserView(ctx, usersView)
	controller.dispatcher.Send("userView", 1, &inout.ChangedUserViewsEventV1{
		Identifiers: functools.UUIDListToStringList(identifiers),
	})
	span.LogFields(log.String("user_id", functools.UUIDListToString(identifiers, ", ")))
	span.Finish()
	ctx.Done()
}

func (controller *Controller) onPhoneCodeConfirmationCreated(phone string, code string) {
	controller.dispatcher.Send("phoneConfirmation", 1, &inout.CreatePhoneConfirmationEventV1{
		Phone: phone,
		Code:  code,
	})
}

func (controller *Controller) onEmailCodeConfirmationCreated(email string, code string) {
	controller.dispatcher.Send("emailConfirmation", 1, &inout.CreateEmailConfirmationEventV1{
		Email: email,
		Code:  code,
	})
}

func (controller *Controller) onSecretCreatedV1(secret *models.Secret) {
	controller.dispatcher.Send("secret", 1, &inout.SecretCreatedV1{
		Id:      secret.Id.Bytes(),
		Created: secret.Created,
		Value:   secret.Value.Bytes(),
	})
}

// Public methods / Header

func (controller *Controller) OnEmailCodeConfirmationCreated(email string, code string) {
	controller.onEmailCodeConfirmationCreated(email, code)
}

func (controller *Controller) OnPhoneCodeConfirmationCreated(phone string, code string) {
	controller.onPhoneCodeConfirmationCreated(phone, code)
}

func (controller *Controller) OnUsersViewChanged(usersView []*models.UserView) {
	controller.onUsersViewChanged(usersView)
}

func (controller *Controller) OnPasswordChanged(userId uuid.UUID) {
	// Todo tokens invalidation
}

func (controller *Controller) OnUserChanged(id []uuid.UUID) {
	controller.onUserChanged(id)
}

func (controller *Controller) OnEmailChanged(userId []uuid.UUID) {
	controller.onEmailChanged(userId)
}

func (controller *Controller) OnPhoneChanged(userId []uuid.UUID) {
	controller.onPhoneChanged(userId)
}

func (controller *Controller) OnRoleChanged(roleId []uuid.UUID) {
	controller.onRoleChanged(roleId)
}

func (controller *Controller) OnSecretCreatedV1(secret *models.Secret) {
	controller.onSecretCreatedV1(secret)
}
