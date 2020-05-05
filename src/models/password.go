package models

import uuid "github.com/satori/go.uuid"

type Password struct {
	Id      uuid.UUID
	Created int64
	UserId  uuid.UUID
	Value   string
}
