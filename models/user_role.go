package models

type UserRoleID int64

type UserRole struct {
	Id      UserRoleID
	Created int64
	UserId  UserID
	RoleId  RoleID
}
