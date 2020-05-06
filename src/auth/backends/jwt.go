package backends

import (
	"hive/config"
	"hive/enums"
	"hive/functools"
	"hive/models"
	"hive/stores"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

type JWTAuthenticationBackend struct {
	store       stores.IStore
	environment *config.Environment
}

type JWTAuthenticationBackendUser struct {
	jwt.StandardClaims
	IsAdmin  bool      `json:"isAdmin"`
	Roles    []string  `json:"roles"`
	UserID   uuid.UUID `json:"userID"`
	SecretID uuid.UUID `json:"secretID"`
}

func (user JWTAuthenticationBackendUser) GetIsAdmin() bool {
	return user.IsAdmin
}

func (user JWTAuthenticationBackendUser) GetRoles() []string {
	return user.Roles
}

func (user JWTAuthenticationBackendUser) GetUserID() uuid.UUID {
	return user.UserID
}

func (backend JWTAuthenticationBackend) EncodeAccessToken(_ context.Context, userID uuid.UUID, roles []string, secret *models.Secret, expires int64) string {

	claims := JWTAuthenticationBackendUser{
		UserID:   userID,
		Roles:    roles,
		IsAdmin:  functools.Contains(config.AdminRole, roles),
		SecretID: secret.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires,
			NotBefore: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret.Value.Bytes())
	if err != nil {
		sentry.CaptureException(err)
		return ""
	}

	return ss
}

func (backend JWTAuthenticationBackend) DecodeAccessTokenWithoutValidation(_ context.Context, tokenValue string) (int, *JWTAuthenticationBackendUser) {
	parser := jwt.Parser{
		SkipClaimsValidation: false,
	}

	token, _, err := parser.ParseUnverified(tokenValue, &JWTAuthenticationBackendUser{})

	if err != nil {
		sentry.CaptureException(err)
		return enums.IncorrectToken, nil
	}

	if claims, ok := token.Claims.(*JWTAuthenticationBackendUser); ok {
		return enums.Ok, claims
	} else {
		return enums.IncorrectToken, nil
	}
}

func (backend JWTAuthenticationBackend) DecodeAccessToken(_ context.Context, tokenValue string, secret uuid.UUID) (int, *JWTAuthenticationBackendUser) {

	token, err := jwt.ParseWithClaims(tokenValue, &JWTAuthenticationBackendUser{}, func(token *jwt.Token) (interface{}, error) {
		return secret.Bytes(), nil
	})

	var e *jwt.ValidationError

	if err != nil {
		if errors.As(err, &e) && strings.Contains(e.Error(), "expired") {
			return enums.InvalidToken, nil
		} else {
			sentry.CaptureException(err)
			return enums.IncorrectToken, nil
		}
	}

	if !token.Valid {
		return enums.InvalidToken, nil
	} else if claims, ok := token.Claims.(*JWTAuthenticationBackendUser); ok {
		return enums.Ok, claims
	} else {
		return enums.IncorrectToken, nil
	}
}

func (backend JWTAuthenticationBackend) GetUser(ctx context.Context, token string) (int, models.IAuthenticationBackendUser) {

	status, unverifiedPayload := backend.DecodeAccessTokenWithoutValidation(ctx, token)
	if status != enums.Ok {
		return status, nil
	}

	secret := backend.store.GetSecret(ctx, unverifiedPayload.SecretID)
	if secret == nil {
		return enums.SecretNotFound, nil
	}

	status, payload := backend.DecodeAccessToken(ctx, token, secret.Value)
	if status != enums.Ok {
		return status, nil
	}

	return enums.Ok, payload
}

func InitJWTAuthenticationBackend(store stores.IStore, environment *config.Environment) *JWTAuthenticationBackend {
	return &JWTAuthenticationBackend{store: store, environment: environment}
}

type JWTAuthenticationBackendWithMockedInternals struct {
	Backend *JWTAuthenticationBackend
	Store   *stores.MockIStore
}

func InitJWTAuthenticationBackendWithMockedInternals(ctrl *gomock.Controller) *JWTAuthenticationBackendWithMockedInternals {
	store := stores.NewMockIStore(ctrl)
	return &JWTAuthenticationBackendWithMockedInternals{
		Backend: InitJWTAuthenticationBackend(store, config.InitEnvironment()),
		Store:   store,
	}
}
