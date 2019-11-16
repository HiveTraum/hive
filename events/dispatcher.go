package events

import (
	"auth/inout"
	"encoding/json"
	"fmt"
)

const (
	PhoneConfirmationEvent = "phoneConfirmationCreated"
	EmailConfirmationEvent = "emailConfirmationCreated"
	UserViewsChangedEvent  = "userViewsChangedEvent"
)

type Dispatcher interface {
	Send(message inout.EventV1)
}

type PrintDispatcher struct {
}

func (dispatcher *PrintDispatcher) Send(message inout.EventV1) {
	s, _ := json.Marshal(message)
	fmt.Println(string(s))
}

type EventDispatcher struct {
	Dispatcher
}

func (eventDispatcher *EventDispatcher) SendCreatePhoneConfirmation(phone string, code string) {
	eventDispatcher.Dispatcher.Send(inout.EventV1{
		Title: PhoneConfirmationEvent,
		TestOneof: &inout.EventV1_PhoneConfirmationEvent{
			PhoneConfirmationEvent: &inout.CreatePhoneConfirmationEventV1{
				Phone: phone,
				Code:  code,
			},
		},
	})
}

func (eventDispatcher *EventDispatcher) SendCreateEmailConfirmation(email string, code string) {
	eventDispatcher.Dispatcher.Send(inout.EventV1{
		Title: EmailConfirmationEvent,
		TestOneof: &inout.EventV1_EmailConfirmationEvent{
			EmailConfirmationEvent: &inout.CreateEmailConfirmationEventV1{
				Email: email,
				Code:  code,
			},
		},
	})
}

func (eventDispatcher *EventDispatcher) SendUserViewsChanged(userViews []*inout.GetUserViewResponseV1) {
	identifiers := make([]int64, len(userViews))
	for i, v := range userViews {
		identifiers[i] = v.Id
	}

	changedUserViewsEvent := &inout.ChangedUserViewsEventV1{Identifiers: identifiers,}
	testOnOf := &inout.EventV1_ChangedUserViewsEvent{ChangedUserViewsEvent: changedUserViewsEvent}
	event := inout.EventV1{Title: UserViewsChangedEvent, TestOneof: testOnOf}

	eventDispatcher.Dispatcher.Send(event)
}

func (eventDispatcher *EventDispatcher) SendUserViewChanged(userView *inout.GetUserViewResponseV1) {
	eventDispatcher.SendUserViewsChanged([]*inout.GetUserViewResponseV1{userView})
}
