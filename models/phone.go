package models

type PhoneID int64

type Phone struct {
	Id      PhoneID
	Created int64
	UserId  UserID
	Value   string
}

type PhoneConfirmation struct {
	Created int64
	Expire  int64
	Phone   string
	Code    string
}
