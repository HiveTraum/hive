package models

import uuid "github.com/satori/go.uuid"

type User struct {
	Id      uuid.UUID
	Created int64
}

type UserView struct {
	Id      uuid.UUID
	Created int64
	Roles   []string
	Phones  []string
	Emails  []string
	RolesID []uuid.UUID
}
