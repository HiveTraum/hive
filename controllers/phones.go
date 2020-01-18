package controllers

import (
	"auth/enums"
	"auth/infrastructure"
	"auth/models"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/nyaruka/phonenumbers"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"strconv"
	"time"
)

func getPhone(phone string) string {
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

func checkPhoneConfirmationCode(store infrastructure.StoreInterface, ctx context.Context, phone string, code string) int {
	cachedCode := store.GetPhoneConfirmationCode(ctx, phone)
	if cachedCode == "" {
		return enums.PhoneNotFound
	} else if cachedCode != code {
		return enums.IncorrectPhoneCode
	}

	return enums.Ok
}

func validatePhone(store infrastructure.StoreInterface, ctx context.Context, phone string, code string) (int, string) {
	phone = getPhone(phone)

	if phone == "" {
		return enums.IncorrectPhone, ""
	}

	status := checkPhoneConfirmationCode(store, ctx, phone, code)

	return status, phone
}

func CreatePhone(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, phone string, code string, userId uuid.UUID) (int, *models.Phone) {

	status, phone := validatePhone(store, ctx, phone, code)

	if status != enums.Ok {
		return status, nil
	}

	_, oldPhone := store.GetPhone(ctx, phone)

	identifiers := []uuid.UUID{userId}

	if oldPhone != nil {
		identifiers = append(identifiers, oldPhone.UserId)
	}

	status, phoneObject := store.CreatePhone(ctx, userId, phone)
	esb.OnPhoneChanged(identifiers)
	return status, phoneObject
}

func getRandomCode() string {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func CreatePhoneConfirmation(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, phone string) (int, *models.PhoneConfirmation) {

	phone = getPhone(phone)

	if phone == "" {
		return enums.IncorrectPhone, nil
	}

	code := getRandomCode()
	phoneConfirmation := store.CreatePhoneConfirmationCode(ctx, phone, code, time.Minute*15)
	esb.OnPhoneCodeConfirmationCreated(phone, code)
	return enums.Ok, phoneConfirmation
}
