package controllers

import (
	"auth/config"
	"auth/enums"
	"auth/functools"
	"auth/infrastructure"
	"auth/inout"
	"auth/models"
	"context"
	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	"github.com/nyaruka/phonenumbers"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/crypto/bcrypt"
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

// Input Normalization

func (controller *LoginController) NormalizeEmail(_ context.Context, email string) string {
	err := checkmail.ValidateFormat(email)

	if err != nil {
		sentry.CaptureException(err)
		return ""
	}

	return email
}

func (controller *LoginController) NormalizePhone(_ context.Context, phone string) string {
	num, err := phonenumbers.Parse(phone, "RU")

	if err != nil {
		sentry.CaptureException(err)
		return ""
	}

	if num == nil || !phonenumbers.IsPossibleNumber(num) {
		return ""
	}

	return phonenumbers.Format(num, phonenumbers.E164)
}

// Access Token

func (controller *LoginController) DecodeAccessToken(_ context.Context, tokenValue string, secret string) (int, *models.AccessTokenPayload) {

	var keyFunc jwt.Keyfunc

	if secret != "" {
		keyFunc = func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		}
	}

	token, err := jwt.ParseWithClaims(tokenValue, &models.AccessTokenPayload{}, keyFunc)

	if err != nil {
		sentry.CaptureException(err)
		return enums.IncorrectToken, nil
	}

	if claims, ok := token.Claims.(*models.AccessTokenPayload); ok && token.Valid {
		return enums.Ok, claims
	} else {
		return enums.IncorrectToken, nil
	}
}

func (controller *LoginController) EncodeAccessToken(_ context.Context, userID models.UserID, roles []string, secret string) string {

	claims := models.AccessTokenPayload{
		UserID:  userID,
		Roles:   roles,
		IsAdmin: functools.Contains(config.AdminRole, roles),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			NotBefore: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		sentry.CaptureException(err)
		return ""
	}

	return ss
}

// Login

func (controller *LoginController) LoginByTokens(ctx context.Context, refreshToken string, accessToken string, fingerprint string) (int, *models.User) {

	status, unverifiedPayload := controller.DecodeAccessToken(ctx, accessToken, "")
	if status != enums.Ok {
		return status, nil
	}

	session := controller.Store.GetSession(ctx, fingerprint, refreshToken, unverifiedPayload.UserID)
	if session == nil {
		return enums.SessionNotFound, nil
	}

	secret := controller.Store.GetSecret(ctx, session.SecretID)
	if secret == nil {
		return enums.SecretNotFound, nil
	}

	status, payload := controller.DecodeAccessToken(ctx, accessToken, secret.Value)
	if status != enums.Ok {
		return status, nil
	}

	return enums.Ok, controller.Store.GetUser(ctx, payload.GetUserID())
}

func (controller *LoginController) LoginByEmail(ctx context.Context, emailValue string, emailCode string, passwordValue string) (int, *models.User) {

	emailValue = controller.NormalizeEmail(ctx, emailValue)
	if emailValue == "" {
		return enums.IncorrectEmail, nil
	}

	code := controller.Store.GetEmailConfirmationCode(ctx, emailValue)
	if code == "" {
		return enums.EmailConfirmationCodeNotFound, nil
	} else if code != emailCode {
		return enums.IncorrectEmailCode, nil
	}

	status, email := controller.Store.GetEmail(ctx, emailValue)
	if status != enums.Ok {
		return status, nil
	}

	status, password := controller.Store.GetLatestPassword(ctx, email.UserId)
	if status != enums.Ok {
		return status, nil
	}

	passwordVerified := controller.VerifyPassword(ctx, passwordValue, password.Value)
	if !passwordVerified {
		return enums.IncorrectPassword, nil
	}

	return enums.Ok, controller.Store.GetUser(ctx, email.UserId)
}

func (controller *LoginController) LoginByPhone(ctx context.Context, phoneValue string, phoneCode string, passwordValue string) (int, *models.User) {
	phoneValue = controller.NormalizePhone(ctx, phoneValue)
	if phoneValue == "" {
		return enums.IncorrectPhone, nil
	}

	code := controller.Store.GetPhoneConfirmationCode(ctx, phoneValue)
	if code == "" {
		return enums.PhoneConfirmationCodeNotFound, nil
	} else if code != phoneCode {
		return enums.IncorrectPhoneCode, nil
	}

	status, phone := controller.Store.GetPhone(ctx, phoneValue)
	if status != enums.Ok {
		return status, nil
	}

	status, password := controller.Store.GetLatestPassword(ctx, phone.UserId)
	if status != enums.Ok {
		return status, nil
	}

	passwordVerified := controller.VerifyPassword(ctx, passwordValue, password.Value)
	if !passwordVerified {
		return enums.IncorrectPassword, nil
	}

	return enums.Ok, controller.Store.GetUser(ctx, phone.UserId)
}

func (controller *LoginController) Login(ctx context.Context, credentials inout.CreateSessionRequestV1) (int, *models.User) {
	var status int
	var user *models.User

	switch credentials.Type {
	case inout.CreateSessionRequestV1_Email:
		status, user = controller.LoginByEmail(ctx, credentials.Email, credentials.EmailCode, credentials.Password)
		break
	case inout.CreateSessionRequestV1_Phone:
		status, user = controller.LoginByPhone(ctx, credentials.Phone, credentials.PhoneCode, credentials.Password)
		break
	case inout.CreateSessionRequestV1_Token:
		status, user = controller.LoginByTokens(ctx, credentials.RefreshToken, credentials.AccessToken, credentials.Fingerprint)
		break
	default:
		return enums.CredentialsNotProvided, nil
	}

	if user == nil {
		return enums.UserNotFound, nil
	}

	return status, user
}
