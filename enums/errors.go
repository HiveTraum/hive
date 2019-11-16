package enums

const (
	// Common

	Ok = iota
	NotOk
	MinimumOneFieldRequired

	// Users

	UserNotFound

	// Phones

	IncorrectPhoneCode
	IncorrectPhone
	PhoneNotFound

	// Emails

	IncorrectEmailCode
	IncorrectEmail
	EmailNotFound

	// Passwords

	PasswordRequired
	IncorrectPassword
)
