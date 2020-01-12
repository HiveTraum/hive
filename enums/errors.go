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

	// Roles

	RoleAlreadyExist // 14
	RoleNotFound // 15

	// User Roles

	UserRoleAlreadyExist // 16
	UserRoleNotFound // 17

	// Sessions

	SessionNotFound // 18

	// Secrets

	SecretNotFound // 19

	// Tokens

	IncorrectToken // 20
	InvalidToken   // 21

	// Credentials

	CredentialsNotProvided // 20
)
