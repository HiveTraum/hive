package models

import uuid "github.com/satori/go.uuid"

type Secret struct {
	Id      uuid.UUID
	Created int64
	Value   uuid.UUID
}
