package models

type Password struct {
	Id      int64
	Created int64
	UserId  int64
	Value   string
}
