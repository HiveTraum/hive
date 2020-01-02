package models

type PasswordID int64

type Password struct {
	Id      PasswordID
	Created int64
	UserId  UserID
	Value   string
}
