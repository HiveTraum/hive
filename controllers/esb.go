package controllers

import (
	"auth/functools"
	"auth/infrastructure"
	"auth/inout"
	"auth/models"
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	uuid "github.com/satori/go.uuid"
	"strings"
)

const SENDER = "auth"

type ESB struct {
	Store      infrastructure.StoreInterface
	Dispatcher infrastructure.ESBDispatcherInterface
}

// Private methods / Implementation

func (esb *ESB) onUserChanged(userId []uuid.UUID) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnUserChanged")
	CreateOrUpdateUsersView(esb.Store, esb, ctx, userId)
	span.LogFields(log.String("user_id", strings.Join(functools.UUIDListToStringList(userId), ", ")))
	span.Finish()
	ctx.Done()
}

func (esb *ESB) onPhoneChanged(userId []uuid.UUID) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnPhoneChanged")
	CreateOrUpdateUsersView(esb.Store, esb, ctx, userId)
	span.LogFields(log.String("user_id", strings.Join(functools.UUIDListToStringList(userId), ", ")))
	span.Finish()
	ctx.Done()
}

func (esb *ESB) onEmailChanged(userId []uuid.UUID) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnEmailChanged")
	CreateOrUpdateUsersView(esb.Store, esb, ctx, userId)
	span.LogFields(log.String("user_id", strings.Join(functools.UUIDListToStringList(userId), ", ")))
	span.Finish()
	ctx.Done()
}

func (esb *ESB) onRoleChanged(roleId []uuid.UUID) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnRoleChanged")
	CreateOrUpdateUsersViewByRoles(esb.Store, esb, ctx, roleId)
	span.LogFields(log.String("role_id", strings.Join(functools.UUIDListToStringList(roleId), "")))
	span.Finish()
	ctx.Done()
}

func (esb *ESB) onUsersViewChanged(usersView []*models.UserView) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnUserViewChanged")
	identifiers := make([]uuid.UUID, len(usersView))

	for i, u := range usersView {
		identifiers[i] = u.Id
	}

	esb.Store.CacheUserView(ctx, usersView)
	event := esb.getUserViewChangedEvent(identifiers)
	esb.Dispatcher.Send(event)
	span.LogFields(log.String("user_id", functools.UUIDListToString(identifiers, ", ")))
	span.Finish()
	ctx.Done()
}

func (esb *ESB) onPhoneCodeConfirmationCreated(phone string, code string) {
	event := esb.getPhoneConfirmationCreatedEvent(phone, code)
	esb.Dispatcher.Send(event)
}

func (esb *ESB) onEmailCodeConfirmationCreated(email string, code string) {
	event := esb.getEmailConfirmationCreatedEvent(email, code)
	esb.Dispatcher.Send(event)
}

func (esb *ESB) getUserViewChangedEvent(userId []uuid.UUID) inout.Event {

	return inout.Event{
		Sender:       SENDER,
		Object:       "user",
		Action:       "changed",
		EventVersion: "1",
		DataVersion:  "1",
		Data: &inout.Event_ChangedUserViewsEvent{
			ChangedUserViewsEvent: &inout.ChangedUserViewsEventV1{
				Identifiers: functools.UUIDListToStringList(userId),
			}},
	}
}

func (esb *ESB) getPhoneConfirmationCreatedEvent(phone string, code string) inout.Event {
	return inout.Event{
		Sender:       SENDER,
		Object:       "phoneConfirmation",
		Action:       "created",
		EventVersion: "1",
		DataVersion:  "1",
		Data: &inout.Event_PhoneConfirmationEvent{
			PhoneConfirmationEvent: &inout.CreatePhoneConfirmationEventV1{
				Phone: phone,
				Code:  code,
			},
		},
	}
}

func (esb *ESB) getEmailConfirmationCreatedEvent(email string, code string) inout.Event {
	return inout.Event{
		Sender:       SENDER,
		Object:       "emailConfirmation",
		Action:       "created",
		EventVersion: "1",
		DataVersion:  "1",
		Data: &inout.Event_EmailConfirmationEvent{
			EmailConfirmationEvent: &inout.CreateEmailConfirmationEventV1{
				Email: email,
				Code:  code,
			},
		},
	}
}

// Public methods / Header

func (esb *ESB) OnEmailCodeConfirmationCreated(email string, code string) {
	go esb.onEmailCodeConfirmationCreated(email, code)
}

func (esb *ESB) OnPhoneCodeConfirmationCreated(phone string, code string) {
	go esb.onPhoneCodeConfirmationCreated(phone, code)
}

func (esb *ESB) OnUsersViewChanged(usersView []*models.UserView) {
	go esb.onUsersViewChanged(usersView)
}

func (esb *ESB) OnPasswordChanged(userId uuid.UUID) {
	// Todo tokens invalidation
}

func (esb *ESB) OnUserChanged(id []uuid.UUID) {
	go esb.onUserChanged(id)
}

func (esb *ESB) OnEmailChanged(userId []uuid.UUID) {
	go esb.onEmailChanged(userId)
}

func (esb *ESB) OnPhoneChanged(userId []uuid.UUID) {
	go esb.onPhoneChanged(userId)
}

func (esb *ESB) OnRoleChanged(roleId []uuid.UUID) {
	go esb.onRoleChanged(roleId)
}
