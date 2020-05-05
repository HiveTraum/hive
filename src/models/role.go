package models

import uuid "github.com/satori/go.uuid"

type Role struct {
	Id      uuid.UUID
	Created int64
	Title   string
}
