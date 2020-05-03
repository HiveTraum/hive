package models

import uuid "github.com/satori/go.uuid"

type IAuthenticationBackendUser interface {
	GetIsAdmin() bool
	GetRoles() []string
	GetUserID() uuid.UUID
}
