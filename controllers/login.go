package controllers

import (
	"auth/config"
	"auth/enums"
	"auth/functools"
	"auth/infrastructure"
	"auth/inout"
	"auth/models"
	"context"
	"errors"
	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	"github.com/nyaruka/phonenumbers"
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

func (controller *LoginController) EncodeAccessToken(_ context.Context, userID uuid.UUID, roles []string, secret *models.Secret, expires time.Time) string {

	claims := models.AccessTokenPayload{
		UserID:   userID,
		Roles:    roles,
		IsAdmin:  functools.Contains(config.AdminRole, roles),
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

func (controller *LoginController) LoginByTokens(ctx context.Context, accessToken string) (int, uuid.UUID) {

	status, unverifiedPayload := controller.DecodeAccessTokenWithoutValidation(ctx, accessToken)
	if status != enums.Ok {
		return status, uuid.Nil
	}

	secret := controller.Store.GetSecret(ctx, unverifiedPayload.SecretID)
	if secret == nil {
		return enums.SecretNotFound, uuid.Nil
	}

	status, payload := controller.DecodeAccessToken(ctx, accessToken, secret.Value)
	if status != enums.Ok {
		return status, uuid.Nil
	}

	return enums.Ok, payload.UserID
}

func (controller *LoginController) LoginByEmailAndCode(ctx context.Context, emailValue string, emailCode string) (int, uuid.UUID) {

	emailValue = controller.NormalizeEmail(ctx, emailValue)
	if emailValue == "" {
		return enums.IncorrectEmail, uuid.Nil
	}

	code := controller.Store.GetEmailConfirmationCode(ctx, emailValue)
	if code == "" {
		return enums.EmailConfirmationCodeNotFound, uuid.Nil
	} else if code != emailCode {
		return enums.IncorrectEmailCode, uuid.Nil
	}

	status, email := controller.Store.GetEmail(ctx, emailValue)
	if status != enums.Ok {
		return status, uuid.Nil
	}
	if email == nil {
		return enums.EmailNotFound, uuid.Nil
	}

	return enums.Ok, email.UserId
}

func (controller *LoginController) LoginByEmailAndPassword(ctx context.Context, emailValue string, passwordValue string) (int, uuid.UUID) {

	emailValue = controller.NormalizeEmail(ctx, emailValue)
	if emailValue == "" {
		return enums.IncorrectEmail, uuid.Nil
	}

	status, email := controller.Store.GetEmail(ctx, emailValue)
	if status != enums.Ok {
		return status, uuid.Nil
	}
	if email == nil {
		return enums.EmailNotFound, uuid.Nil
	}

	status, password := controller.Store.GetLatestPassword(ctx, email.UserId)
	if status != enums.Ok {
		return status, uuid.Nil
	}
	if password == nil {
		return enums.PasswordNotFound, uuid.Nil
	}

	passwordVerified := controller.VerifyPassword(ctx, passwordValue, password.Value)
	if !passwordVerified {
		return enums.IncorrectPassword, uuid.Nil
	}

	return enums.Ok, password.UserId
}

func (controller *LoginController) LoginByPhoneAndPassword(ctx context.Context, phoneValue string, passwordValue string) (int, uuid.UUID) {
	phoneValue = controller.NormalizePhone(ctx, phoneValue)
	if phoneValue == "" {
		return enums.IncorrectPhone, uuid.Nil
	}

	status, phone := controller.Store.GetPhone(ctx, phoneValue)
	if status != enums.Ok {
		return status, uuid.Nil
	}
	if phone == nil {
		return enums.PhoneNotFound, uuid.Nil
	}

	status, password := controller.Store.GetLatestPassword(ctx, phone.UserId)
	if status != enums.Ok {
		return status, uuid.Nil
	}
	if password == nil {
		return enums.PasswordNotFound, uuid.Nil
	}

	passwordVerified := controller.VerifyPassword(ctx, passwordValue, password.Value)
	if !passwordVerified {
		return enums.IncorrectPassword, uuid.Nil
	}

	return enums.Ok, password.UserId
}

func (controller *LoginController) LoginByPhoneAndCode(ctx context.Context, phoneValue string, phoneCode string) (int, uuid.UUID) {
	phoneValue = controller.NormalizePhone(ctx, phoneValue)
	if phoneValue == "" {
		return enums.IncorrectPhone, uuid.Nil
	}

	code := controller.Store.GetPhoneConfirmationCode(ctx, phoneValue)
	if code == "" {
		return enums.PhoneConfirmationCodeNotFound, uuid.Nil
	} else if code != phoneCode {
		return enums.IncorrectPhoneCode, uuid.Nil
	}

	status, phone := controller.Store.GetPhone(ctx, phoneValue)
	if status != enums.Ok {
		return status, uuid.Nil
	}
	if phone == nil {
		return enums.PhoneNotFound, uuid.Nil
	}

	return enums.Ok, phone.UserId
}

func (controller *LoginController) Login(ctx context.Context, credentials inout.CreateSessionRequestV1) (int, uuid.UUID) {
	var status int
	var user uuid.UUID

	switch credentials.Data.(type) {
	case *inout.CreateSessionRequestV1_Tokens_:
		tokens := credentials.GetTokens()
		status, user = controller.LoginByTokens(ctx, tokens.AccessToken)
	case *inout.CreateSessionRequestV1_EmailAndPassword_:
		emailAndPassword := credentials.GetEmailAndPassword()
		status, user = controller.LoginByEmailAndPassword(ctx, emailAndPassword.Email, emailAndPassword.Password)
	case *inout.CreateSessionRequestV1_EmailAndCode_:
		emailAndCode := credentials.GetEmailAndCode()
		status, user = controller.LoginByEmailAndCode(ctx, emailAndCode.Email, emailAndCode.Code)
	case *inout.CreateSessionRequestV1_PhoneAndPassword_:
		phoneAndPassword := credentials.GetPhoneAndPassword()
		status, user = controller.LoginByPhoneAndPassword(ctx, phoneAndPassword.Phone, phoneAndPassword.Password)
	case *inout.CreateSessionRequestV1_PhoneAndCode_:
		phoneAndCode := credentials.GetPhoneAndCode()
		status, user = controller.LoginByPhoneAndPassword(ctx, phoneAndCode.Phone, phoneAndCode.Code)
	default:
		return enums.CredentialsNotProvided, uuid.Nil
	}

	return status, user
}
