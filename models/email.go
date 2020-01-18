package models

import uuid "github.com/satori/go.uuid"

type EmailID int64

type Email struct {
	Id      uuid.UUID
	Created int64
	UserId  uuid.UUID
	Value   string
}

type EmailConfirmation struct {
	Created int64
	Expire  int64
	Email   string
	Code    string
}
