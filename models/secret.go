package models

type SecretID int64

type Secret struct {
	Id      SecretID
	Created int64
	Value   string
}
