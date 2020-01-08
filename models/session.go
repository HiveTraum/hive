package models

type Session struct {
	RefreshToken string
	Fingerprint  string
	UserID       UserID
	SecretID     SecretID
	Created      int64
	Expires      int64
	UserAgent    string
}
