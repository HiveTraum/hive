package models

type Email struct {
	Id      int64
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
