package functools

import (
	"github.com/getsentry/sentry-go"
	"golang.org/x/crypto/bcrypt"
)

type PasswordProcessor struct {
}

func (processor *PasswordProcessor) Encode(value string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.MinCost)
	if err != nil {
		sentry.CaptureException(err)
		return ""
	}

	return string(hash)
}
