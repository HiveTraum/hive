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
	PhoneConfirmationCodeNotFound

	// Emails

	IncorrectEmailCode
	IncorrectEmail
	EmailNotFound
	EmailConfirmationCodeNotFound

	// Passwords

	PasswordRequired
	IncorrectPassword

	// Roles

	RoleAlreadyExist
	RoleNotFound

	// User Roles

	UserRoleAlreadyExist
	UserRoleNotFound

	// Sessions

	SessionNotFound

	// Secrets

	SecretNotFound

	// Tokens

	IncorrectToken

	// Credentials

	CredentialsNotProvided
)
