package functools

import (
	"github.com/badoux/checkmail"
	"github.com/getsentry/sentry-go"
	"github.com/ttacon/libphonenumber"
)

func NormalizeEmail(email string) string {
	err := checkmail.ValidateFormat(email)

	if err != nil {
		sentry.CaptureException(err)
		return ""
	}

	return email
}

func NormalizePhone(phone string) string {

	num, err := libphonenumber.Parse(phone, "")

	if err != nil {
		sentry.CaptureException(err)
		return ""
	}

	if num == nil || !libphonenumber.IsPossibleNumber(num) {
		return ""
	}

	return libphonenumber.Format(num, libphonenumber.INTERNATIONAL)
}
