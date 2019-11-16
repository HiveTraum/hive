package models

type Phone struct {
	Id      int64
	Created int64
	UserId  int64
	Value   string
}

type PhoneConfirmation struct {
	Created int64
	Expire  int64
	Phone   string
	Code    string
}
