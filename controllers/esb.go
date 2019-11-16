package controllers

import (
	"auth/infrastructure"
	"auth/inout"
	"context"
)

const SENDER = "auth"

type ESB struct {
	Store      infrastructure.StoreInterface
	Dispatcher infrastructure.ESBDispatcherInterface
}

// Private methods / Implementation

func (esb *ESB) onUserChanged(userId []int64) {
	ctx := context.Background()
	CreateOrUpdateUsersView(esb.Store, esb, ctx, userId)
}

func (esb *ESB) onPhoneChanged(userId []int64) {
	ctx := context.Background()
	CreateOrUpdateUsersView(esb.Store, esb, ctx, userId)
}

func (esb *ESB) onEmailChanged(userId []int64) {
	ctx := context.Background()
	CreateOrUpdateUsersView(esb.Store, esb, ctx, userId)
}

func (esb *ESB) onRoleChanged(roleId []int64) {
	ctx := context.Background()
	CreateOrUpdateUsersViewByRoles(esb.Store, esb, ctx, roleId)
}

func (esb *ESB) onUsersViewChanged(usersView []*inout.GetUserViewResponseV1) {
	identifiers := make([]int64, len(usersView))

	for i, u := range usersView {
		identifiers[i] = u.Id
	}

	event := esb.getUserViewChangedEvent(identifiers)
	esb.Dispatcher.Send(event)
}

func (esb *ESB) onPhoneCodeConfirmationCreated(phone string, code string) {
	event := esb.getPhoneConfirmationCreatedEvent(phone, code)
	esb.Dispatcher.Send(event)
}

func (esb *ESB) onEmailCodeConfirmationCreated(email string, code string) {
	event := esb.getEmailConfirmationCreatedEvent(email, code)
	esb.Dispatcher.Send(event)
}

func (esb *ESB) getUserViewChangedEvent(userId []int64) inout.Event {
	return inout.Event{
		Sender:       SENDER,
		Object:       "user",
		Action:       "changed",
		EventVersion: "1",
		DataVersion:  "1",
		Data: &inout.Event_ChangedUserViewsEvent{
			ChangedUserViewsEvent: &inout.ChangedUserViewsEventV1{
				Identifiers: userId,
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
	go esb.OnUsersViewChanged(usersView)
}

func (esb *ESB) OnPasswordChanged(userId int64) {
	// Todo tokens invalidation
}

func (esb *ESB) OnUserChanged(id []int64) {
	go esb.onUserChanged(id)
}

func (esb *ESB) OnEmailChanged(userId []int64) {
	go esb.onEmailChanged(userId)
}

func (esb *ESB) OnPhoneChanged(userId []int64) {
	go esb.onPhoneChanged(userId)
}

func (esb *ESB) OnRoleChanged(roleId []int64) {
	go esb.onRoleChanged(roleId)
}
