package functools

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNormalizeEmail(t *testing.T) {
	t.Parallel()
	email := "mail@mail.com"
	normalizedEmail := NormalizeEmail(email)
	require.NotEmpty(t, normalizedEmail)
	require.Equal(t, email, normalizedEmail)
}

func TestNormalizeIncorrectEmail(t *testing.T) {
	t.Parallel()
	incorrectEmail := "mail"
	normalizedEmail := NormalizeEmail(incorrectEmail)
	require.Empty(t, normalizedEmail)
	require.NotEqual(t, incorrectEmail, normalizedEmail)
	require.Equal(t, "", normalizedEmail)
}

func TestNormalizePhone(t *testing.T) {
	t.Parallel()
	phone := "+79234567890"
	normalizedPhone := NormalizePhone(phone)
	require.NotEmpty(t, normalizedPhone)
	require.Equal(t, "+7 923 456-78-90", normalizedPhone)
}

func TestNormalizeIncorrectPhone(t *testing.T) {
	t.Parallel()
	phone := "+123"
	normalizedPhone := NormalizePhone(phone)
	require.Empty(t, normalizedPhone)
	require.NotEqual(t, phone, normalizedPhone)
	require.Equal(t, "", normalizedPhone)
}
