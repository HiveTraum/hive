package passwordProcessors

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/crypto/bcrypt"
)

type IPasswordProcessor interface {
	EncodePassword(context.Context, string) string
	VerifyPassword(ctx context.Context, password string, encodedPassword string) bool
}

type PasswordProcessor struct {
}

func InitPasswordProcessor() *PasswordProcessor {
	return &PasswordProcessor{}
}

func (processor *PasswordProcessor) EncodePassword(ctx context.Context, value string) string {

	span, ctx := opentracing.StartSpanFromContext(ctx, "Password encoding")
	hash, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.MinCost)
	if err != nil {
		span.LogFields(log.Error(err))
		sentry.CaptureException(err)
		return ""
	}

	span.Finish()
	return string(hash)
}

func (processor *PasswordProcessor) VerifyPassword(ctx context.Context, password string, encodedPassword string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Password verification")
	encodedPasswordBytes := []byte(encodedPassword)
	passwordBytes := []byte(password)

	err := bcrypt.CompareHashAndPassword(encodedPasswordBytes, passwordBytes)
	if err != nil {
		span.LogFields(log.Error(err))
		sentry.CaptureException(err)
		return false
	}

	return true
}
