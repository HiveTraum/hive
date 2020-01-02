package models

type RoleID int64

type Role struct {
	Id      RoleID
	Created int64
	Title   string
}
