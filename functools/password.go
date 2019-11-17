package functools

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/crypto/bcrypt"
)

type PasswordProcessor struct {
}

func (processor *PasswordProcessor) Encode(ctx context.Context, value string) string {

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
