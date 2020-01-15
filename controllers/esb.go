package controllers

import (
	"auth/functools"
	"auth/infrastructure"
	"auth/inout"
	"auth/models"
	"auth/modelsFunctools"
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

const SENDER = "auth"

type ESB struct {
	Store      infrastructure.StoreInterface
	Dispatcher infrastructure.ESBDispatcherInterface
}

// Private methods / Implementation

func (esb *ESB) onUserChanged(userId []models.UserID) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnUserChanged")
	CreateOrUpdateUsersView(esb.Store, esb, ctx, userId)
	span.LogFields(log.String("user_id", strings.Join(modelsFunctools.UserIDListToStringList(userId), ", ")))
	span.Finish()
	ctx.Done()
}

func (esb *ESB) onPhoneChanged(userId []models.UserID) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnPhoneChanged")
	CreateOrUpdateUsersView(esb.Store, esb, ctx, userId)
	span.LogFields(log.String("user_id", strings.Join(modelsFunctools.UserIDListToStringList(userId), ", ")))
	span.Finish()
	ctx.Done()
}

func (esb *ESB) onEmailChanged(userId []models.UserID) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnEmailChanged")
	CreateOrUpdateUsersView(esb.Store, esb, ctx, userId)
	span.LogFields(log.String("user_id", strings.Join(modelsFunctools.UserIDListToStringList(userId), ", ")))
	span.Finish()
	ctx.Done()
}

func (esb *ESB) onRoleChanged(roleId []models.RoleID) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnRoleChanged")
	CreateOrUpdateUsersViewByRoles(esb.Store, esb, ctx, roleId)
	span.LogFields(log.String("role_id", functools.Int64SliceToString(modelsFunctools.RoleIDListToInt64List(roleId), ", ")))
	span.Finish()
	ctx.Done()
}

func (esb *ESB) onUsersViewChanged(usersView []*inout.GetUserViewResponseV1) {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "OnUserViewChanged")
	identifiers := make([]models.UserID, len(usersView))

	for i, u := range usersView {
		identifiers[i] = models.UserID(u.Id)
	}

	esb.Store.CacheUserView(ctx, usersView)
	event := esb.getUserViewChangedEvent(identifiers)
	esb.Dispatcher.Send(event)
	span.LogFields(log.String("user_id", modelsFunctools.UserIDListToString(identifiers, ", ")))
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

func (esb *ESB) getUserViewChangedEvent(userId []models.UserID) inout.Event {

	int64list := modelsFunctools.UserIDListToStringList(userId)

	return inout.Event{
		Sender:       SENDER,
		Object:       "user",
		Action:       "changed",
		EventVersion: "1",
		DataVersion:  "1",
		Data: &inout.Event_ChangedUserViewsEvent{
			ChangedUserViewsEvent: &inout.ChangedUserViewsEventV1{
				Identifiers: int64list,
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

func (esb *ESB) OnUsersViewChanged(usersView []*inout.GetUserViewResponseV1) {
	go esb.onUsersViewChanged(usersView)
}

func (esb *ESB) OnPasswordChanged(userId models.UserID) {
	// Todo tokens invalidation
}

func (esb *ESB) OnUserChanged(id []models.UserID) {
	go esb.onUserChanged(id)
}

func (esb *ESB) OnEmailChanged(userId []models.UserID) {
	go esb.onEmailChanged(userId)
}

func (esb *ESB) OnPhoneChanged(userId []models.UserID) {
	go esb.onPhoneChanged(userId)
}

func (esb *ESB) OnRoleChanged(roleId []models.RoleID) {
	go esb.onRoleChanged(roleId)
}
