package repositories

import (
	"auth/enums"
	"auth/functools"
	"auth/inout"
	"auth/models"
	"auth/modelsFunctools"
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v7"
	"github.com/golang/protobuf/proto"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	uuid "github.com/satori/go.uuid"
	"time"
)

func getUserKey(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", enums.UserView, id.String())
}

func GetUserViewFromCache(cache *redis.Client, ctx context.Context, id uuid.UUID) *models.UserView {

	span, ctx := opentracing.StartSpanFromContext(ctx, "Get user view from cache")

	key := getUserKey(id)

	value, err := cache.Get(key).Bytes()
	if err != nil {
		span.LogFields(log.Error(err))
		sentry.CaptureException(err)
		return nil
	}

	var userView inout.UserViewCache

	err = proto.Unmarshal(value, &userView)

	if err != nil {
		span.LogFields(log.Error(err))
		sentry.CaptureException(err)
		return nil
	}

	span.Finish()

	rolesID := make([]uuid.UUID, len(userView.RolesID))
	for i, id := range userView.RolesID {
		rolesID[i] = uuid.FromBytesOrNil(id)
	}

	return &models.UserView{
		Id:      uuid.FromBytesOrNil(userView.Id),
		Created: userView.Created,
		Roles:   userView.Roles,
		Phones:  userView.Phones,
		Emails:  userView.Emails,
		RolesID: rolesID,
	}
}

func CacheUserView(cache *redis.Client, ctx context.Context, userViews []*models.UserView) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "Cache user views")

	if userViews == nil {
		return
	}

	userViewsCache := make([]*inout.UserViewCache, len(userViews))

	for i, uv := range userViews {
		userViewsCache[i] = &inout.UserViewCache{
			Id:                   uv.Id.Bytes(),
			Created:              uv.Created,
			Roles:                uv.Roles,
			Phones:               uv.Phones,
			Emails:               uv.Emails,
			RolesID:              functools.UUIDSliceToByteArraySlice(uv.RolesID),
		}
	}

	identifiers := make([]uuid.UUID, len(userViews))

	pipeline := cache.TxPipeline()
	for i, uv := range userViewsCache {
		data, err := proto.Marshal(uv)
		if err != nil {
			span.LogFields(log.Error(err))
			sentry.CaptureException(err)
			continue
		}

		userID := uuid.FromBytesOrNil(uv.Id)

		identifiers[i] = userID
		pipeline.Set(getUserKey(userID), data, time.Hour*48)
	}

	_, err := pipeline.Exec()

	if err != nil {
		span.LogFields(log.Error(err))
		sentry.CaptureException(err)
	}

	span.LogFields(log.String("user_id", modelsFunctools.UserIDListToString(identifiers, ", ")))
	span.Finish()
}
