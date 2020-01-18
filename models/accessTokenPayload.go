package models

import (
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type AccessTokenPayload struct {
	jwt.StandardClaims
	IsAdmin bool      `json:"isAdmin"`
	Roles   []string  `json:"roles"`
	UserID  uuid.UUID `json:"userID"`
}

func (payload *AccessTokenPayload) GetUserID() uuid.UUID {
	return payload.UserID
}
