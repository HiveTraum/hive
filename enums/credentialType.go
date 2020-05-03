package enums

type CredentialType int

const (
	EmailAndPassword CredentialType = 0
	EmailAndCode     CredentialType = 1
	PhoneAndPassword CredentialType = 2
	PhoneAndCode     CredentialType = 3
)
