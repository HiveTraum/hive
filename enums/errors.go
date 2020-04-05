package enums

const (
	// Common

	Ok = iota
	NotOk // 1
	MinimumOneFieldRequired // 2

	// Users

	UserNotFound // 3

	// Phones

	IncorrectPhoneCode // 4
	IncorrectPhone // 5
	PhoneNotFound // 6
	PhoneConfirmationCodeNotFound // 7

	// Emails

	IncorrectEmailCode // 8
	IncorrectEmail // 9
	EmailNotFound // 10
	EmailConfirmationCodeNotFound // 11

	// Passwords

	PasswordRequired // 12
	IncorrectPassword // 13
	PasswordNotFound // 14

	// Roles

	RoleAlreadyExist // 15
	RoleNotFound // 16

	// User Roles

	UserRoleAlreadyExist // 17
	UserRoleNotFound // 18

	// Sessions

	SessionNotFound // 19

	// Secrets

	SecretNotFound // 20

	// Tokens

	IncorrectToken // 21
	InvalidToken   // 22

	// Credentials

	CredentialsNotProvided // 23

	// Auth Backends

	BackendNotFound // 24
)
