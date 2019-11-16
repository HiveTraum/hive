package models

type User struct {
	Id      int64
	Created int64
}

type UserView struct {
	Id      int64
	Created int64
	Roles   []string
	Phones  []string
	Emails  []string
}
