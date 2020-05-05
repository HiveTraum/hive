package repositories

import (
	"hive/enums"
	"hive/functools"
	"hive/inout"
	"hive/models"
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v7"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/proto"
	"time"
)

func getUserKey(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", enums.UserView, id.String())
}

func GetUserViewFromCache(cache *redis.Client, ctx context.Context, id uuid.UUID) *models.UserView {

	key := getUserKey(id)

	value, err := cache.WithContext(ctx).Get(key).Bytes()
	if err != nil {
		return nil
	}

	var userView inout.UserViewCache

	err = proto.Unmarshal(value, &userView)

	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

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

	if userViews == nil {
		return
	}

	userViewsCache := make([]*inout.UserViewCache, len(userViews))

	for i, uv := range userViews {
		userViewsCache[i] = &inout.UserViewCache{
			Id:      uv.Id.Bytes(),
			Created: uv.Created,
			Roles:   uv.Roles,
			Phones:  uv.Phones,
			Emails:  uv.Emails,
			RolesID: functools.UUIDSliceToByteArraySlice(uv.RolesID),
		}
	}

	identifiers := make([]uuid.UUID, len(userViews))

	pipeline := cache.WithContext(ctx).TxPipeline()
	for i, uv := range userViewsCache {
		data, err := proto.Marshal(uv)
		if err != nil {
			continue
		}

		userID := uuid.FromBytesOrNil(uv.Id)

		identifiers[i] = userID
		pipeline.Set(getUserKey(userID), data, time.Hour*48)
	}

	_, err := pipeline.Exec()

	if err != nil {
		sentry.CaptureException(err)
	}
}
