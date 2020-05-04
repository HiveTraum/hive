package models

import uuid "github.com/satori/go.uuid"

type Session struct {
	Id           uuid.UUID
	RefreshToken uuid.UUID
	Fingerprint  string
	UserID       uuid.UUID
	SecretID     uuid.UUID
	Created      int64
	Expires      int64
	UserAgent    string
	AccessToken  string
}
