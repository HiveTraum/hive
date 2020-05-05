package enums

const (
	// Common

	Ok                      = iota
	NotOk                   // 1
	MinimumOneFieldRequired // 2

	// Users

	UserNotFound // 3

	// Phones

	IncorrectPhoneCode            // 4
	IncorrectPhone                // 5
	PhoneNotFound                 // 6
	PhoneConfirmationCodeNotFound // 7

	// Emails

	IncorrectEmailCode            // 8
	IncorrectEmail                // 9
	EmailNotFound                 // 10
	EmailConfirmationCodeNotFound // 11
	EmailConfirmationCodeRequired // 12

	// Passwords

	PasswordRequired  // 13
	IncorrectPassword // 14
	PasswordNotFound  // 15

	// Roles

	RoleAlreadyExist // 16
	RoleNotFound     // 17

	// User Roles

	UserRoleAlreadyExist // 18
	UserRoleNotFound     // 19

	// Sessions

	SessionNotFound // 20

	// Secrets

	SecretNotFound // 21

	// Tokens

	IncorrectToken // 22
	InvalidToken   // 23

	// Credentials

	CredentialsNotProvided  // 24
	CredentialsTypeNotFound // 25

	// Auth backends

	BackendNotFound // 26
)
