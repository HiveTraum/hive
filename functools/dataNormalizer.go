package functools

import (
	"github.com/badoux/checkmail"
	"github.com/getsentry/sentry-go"
	"github.com/nyaruka/phonenumbers"
)

func NormalizeEmail(email string) string {
	err := checkmail.ValidateFormat(email)

	if err != nil {
		sentry.CaptureException(err)
		return ""
	}

	return email
}

func NormalizePhone(phone string, countryCode string) string {
	num, err := phonenumbers.Parse(phone, countryCode)

	if err != nil {
		sentry.CaptureException(err)
		return ""
	}

	if num == nil || !phonenumbers.IsPossibleNumber(num) {
		return ""
	}

	return phonenumbers.Format(num, phonenumbers.E164)
}
