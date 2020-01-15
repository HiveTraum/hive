package models

type UserID string

type User struct {
	Id      UserID
	Created int64
}

type UserView struct {
	Id      UserID
	Created int64
	Roles   []string
	Phones  []string
	Emails  []string
}
