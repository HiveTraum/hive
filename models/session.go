package models

import uuid "github.com/satori/go.uuid"

type Session struct {
	RefreshToken string
	Fingerprint  string
	UserID       uuid.UUID
	SecretID     uuid.UUID
	Created      int64
	Expires      int64
	UserAgent    string
	AccessToken  string
}
