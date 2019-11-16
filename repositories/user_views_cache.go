package repositories

import (
	"auth/enums"
	"auth/inout"
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v7"
	"github.com/golang/protobuf/proto"
	"github.com/opentracing/opentracing-go"
	"time"
)

func getUserKey(id int64) string {
	return fmt.Sprintf("%s:%d", enums.UserView, id)
}

func GetUserViewFromCache(cache *redis.Client, ctx context.Context, id int64) *inout.GetUserViewResponseV1 {

	span, ctx := opentracing.StartSpanFromContext(ctx, "Get user view from cache")

	key := getUserKey(id)

	value, err := cache.Get(key).Bytes()
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	var userView inout.GetUserViewResponseV1

	err = proto.Unmarshal(value, &userView)

	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	span.Finish()

	return &userView
}

func CacheUserView(cache *redis.Client, ctx context.Context, userViews []*inout.GetUserViewResponseV1) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "Cache user views")

	if userViews == nil {
		return
	}

	pipeline := cache.TxPipeline()
	for _, uv := range userViews {
		data, err := proto.Marshal(uv)
		if err != nil {
			sentry.CaptureException(err)
			continue
		}

		pipeline.Set(getUserKey(uv.Id), data, time.Hour*48)
	}

	_, err := pipeline.Exec()

	if err != nil {
		sentry.CaptureException(err)
	}

	span.Finish()
}
