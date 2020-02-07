package controllers

import (
	"auth/config"
	"auth/enums"
	"auth/functools"
	"auth/infrastructure"
	"auth/models"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type LoginController struct {
	Store infrastructure.StoreInterface
}

// Password

func (controller *LoginController) EncodePassword(ctx context.Context, value string) string {

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

func (controller *LoginController) VerifyPassword(ctx context.Context, password string, encodedPassword string) bool {
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

// Access Token

func (controller *LoginController) EncodeAccessToken(_ context.Context, userID uuid.UUID, roles []string, secret *models.Secret, expires time.Time) string {

	claims := models.AccessTokenPayload{
		UserID:   userID,
		Roles:    roles,
		IsAdmin:  functools.Contains(config.GetEnvironment().AdminRole, roles),
		SecretID: secret.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
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

func (controller *LoginController) DecodeAccessTokenWithoutValidation(_ context.Context, tokenValue string) (int, *models.AccessTokenPayload) {
	parser := jwt.Parser{
		SkipClaimsValidation: false,
	}

	token, _, err := parser.ParseUnverified(tokenValue, &models.AccessTokenPayload{})

	if err != nil {
		sentry.CaptureException(err)
		return enums.IncorrectToken, nil
	}

	if claims, ok := token.Claims.(*models.AccessTokenPayload); ok {
		return enums.Ok, claims
	} else {
		return enums.IncorrectToken, nil
	}
}

func (controller *LoginController) DecodeAccessToken(_ context.Context, tokenValue string, secret uuid.UUID) (int, *models.AccessTokenPayload) {

	token, err := jwt.ParseWithClaims(tokenValue, &models.AccessTokenPayload{}, func(token *jwt.Token) (interface{}, error) {
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
	} else if claims, ok := token.Claims.(*models.AccessTokenPayload); ok {
		return enums.Ok, claims
	} else {
		return enums.IncorrectToken, nil
	}
}

// Login

func (controller *LoginController) Login(ctx context.Context, accessToken string) (int, *models.AccessTokenPayload) {
	status, unverifiedPayload := controller.DecodeAccessTokenWithoutValidation(ctx, accessToken)
	if status != enums.Ok {
		return status, nil
	}

	secret := controller.Store.GetSecret(ctx, unverifiedPayload.SecretID)
	if secret == nil {
		return enums.SecretNotFound, nil
	}

	status, payload := controller.DecodeAccessToken(ctx, accessToken, secret.Value)
	if status != enums.Ok {
		return status, nil
	}

	return enums.Ok, payload
}
