package config

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	zerolog "github.com/rs/zerolog/log"
	"strconv"
	"sync"
	"time"
)

var client *redis.Client
var onceRedis sync.Once

const (
	RedisOpenTracingHookStartSpanKey = "RedisOpenTracingHookStartSpanKey"
)

type OpenTracingHook struct {
}

func (o OpenTracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	startTime := opentracing.StartTime(time.Now())
	var arg string
	if len(cmd.Args()) > 1 {
		argStr, argOk := cmd.Args()[1].(string)
		if argOk {
			arg = argStr
		}
	}
	span, ctx := opentracing.StartSpanFromContext(ctx, cmd.Name()+" "+arg, startTime)
	for i, a := range cmd.Args() {
		switch value := a.(type) {
		case string:
			span.LogFields(log.String("arg "+strconv.Itoa(i), value))
		}
	}
	ctx = context.WithValue(ctx, RedisOpenTracingHookStartSpanKey, span)
	return ctx, nil
}

func (o OpenTracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span, spanOk := ctx.Value(RedisOpenTracingHookStartSpanKey).(opentracing.Span)
	if spanOk {
		if err := cmd.Err(); err != nil {
			span.LogFields(log.Error(err))
		}
		span.Finish()
	}

	return nil
}

func (o OpenTracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	startTime := opentracing.StartTime(time.Now())
	span, ctx := opentracing.StartSpanFromContext(ctx, "Redis Pipeline", startTime)
	for i, c := range cmds {
		span.LogFields(log.String("cmd "+strconv.Itoa(i), c.Name()))
	}
	ctx = context.WithValue(ctx, RedisOpenTracingHookStartSpanKey, span)
	return ctx, nil
}

func (o OpenTracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	span, spanOk := ctx.Value(RedisOpenTracingHookStartSpanKey).(opentracing.Span)
	if spanOk {
		for _, c := range cmds {
			span.LogFields(log.Error(c.Err()))
		}
		span.Finish()
	}

	return nil
}

type SentryRedisHook struct {
}

func (s SentryRedisHook) BeforeProcess(ctx context.Context, _ redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (s SentryRedisHook) AfterProcess(_ context.Context, cmd redis.Cmder) error {
	if err := cmd.Err(); err != nil && err != redis.Nil {
		sentry.CaptureException(err)
	}
	return nil
}

func (s SentryRedisHook) BeforeProcessPipeline(ctx context.Context, _ []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (s SentryRedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	for _, c := range cmds {
		_ = s.AfterProcess(ctx, c)
	}
	return nil
}

func InitRedis(environment *Environment) *redis.Client {

	onceRedis.Do(func() {

		config, err := redis.ParseURL(environment.RedisURI)
		if err != nil {
			panic(err)
		}

		client = redis.NewClient(config)
		client.AddHook(OpenTracingHook{})
		client.AddHook(SentryRedisHook{})
		_, err = client.Ping().Result()
		if err != nil {
			panic(err)
		}
		zerolog.Log().Msg("Redis connection successfully initiated")
	})

	return client
}
