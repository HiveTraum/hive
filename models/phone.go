package models

import uuid "github.com/satori/go.uuid"

type Phone struct {
	Id      uuid.UUID
	Created int64
	UserId  uuid.UUID
	Value   string
}

type PhoneConfirmation struct {
	Created int64
	Expire  int64
	Phone   string
	Code    string
}
