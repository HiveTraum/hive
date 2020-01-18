package models

import uuid "github.com/satori/go.uuid"

type UserRole struct {
	Id      uuid.UUID
	Created int64
	UserId  uuid.UUID
	RoleId  uuid.UUID
}
