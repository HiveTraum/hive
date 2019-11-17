package models

type EmailID = int64

type Email struct {
	Id      EmailID
	Created int64
	UserId  int64
	Value   string
}

type EmailConfirmation struct {
	Created int64
	Expire  int64
	Email   string
	Code    string
}
