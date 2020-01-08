package models

import "github.com/dgrijalva/jwt-go"

type AccessTokenPayload struct {
	jwt.StandardClaims
	IsAdmin bool     `json:"isAdmin"`
	Roles   []string `json:"roles"`
	UserID  UserID   `json:"userID"`
}

func (payload *AccessTokenPayload) GetUserID() UserID {
	return payload.UserID
}
