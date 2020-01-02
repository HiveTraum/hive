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
	PhoneAlreadyExist

	// Emails

	IncorrectEmailCode
	IncorrectEmail
	EmailNotFound

	// Passwords

	PasswordRequired
	IncorrectPassword

	// Roles

	RoleAlreadyExist
	RoleNotFound

	// User Roles

	UserRoleAlreadyExist
	UserRoleNotFound
)
