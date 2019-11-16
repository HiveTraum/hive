package controllers

import (
	"auth/enums"
	"auth/infrastructure"
	"auth/models"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/nyaruka/phonenumbers"
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

func checkPhoneConfirmationCode(store infrastructure.StoreInterface, phone string, code string) int {
	cachedCode := store.GetPhoneConfirmationCode(phone)
	if cachedCode == "" {
		return enums.PhoneNotFound
	} else if cachedCode != code {
		return enums.IncorrectPhoneCode
	}

	return enums.Ok
}

func validatePhone(store infrastructure.StoreInterface, phone string, code string) (int, string) {
	phone = getPhone(phone)

	if phone == "" {
		return enums.IncorrectPhone, ""
	}

	status := checkPhoneConfirmationCode(store, phone, code)

	return status, phone
}

func CreatePhone(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, phone string, code string, userId int64) (int, *models.Phone) {

	status, phone := validatePhone(store, phone, code)

	if status != enums.Ok {
		return status, nil
	}

	_, oldPhone := store.GetPhone(ctx, phone)

	identifiers := []int64{userId}

	if oldPhone != nil {
		identifiers = append(identifiers, oldPhone.UserId)
	}

	status, phoneObject := store.CreatePhone(ctx, userId, phone)
	go esb.OnPhoneChanged(identifiers)
	return status, phoneObject
}

func getRandomCode() string {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func CreatePhoneConfirmation(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, phone string) (int, *models.PhoneConfirmation) {

	phone = getPhone(phone)

	if phone == "" {
		return enums.IncorrectPhone, nil
	}

	code := getRandomCode()
	phoneConfirmation := store.CreatePhoneConfirmationCode(phone, code, time.Minute*15)
	go esb.OnPhoneCodeConfirmationCreated(phone, code)
	return enums.Ok, phoneConfirmation
}
