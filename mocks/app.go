package mocks

import (
	"auth/infrastructure"
	"github.com/golang/mock/gomock"
)

type MockApp struct {
	Store             infrastructure.StoreInterface
	ESB               infrastructure.ESBInterface
	PasswordProcessor infrastructure.PasswordProcessorInterface
}

func (app *MockApp) GetStore() infrastructure.StoreInterface {
	return app.Store
}

func (app *MockApp) GetESB() infrastructure.ESBInterface {
	return app.ESB
}

func (app *MockApp) GetPasswordProcessor() infrastructure.PasswordProcessorInterface {
	return app.PasswordProcessor
}

func InitMockApp(ctrl *gomock.Controller) (*MockApp, *MockStoreInterface, *MockESBInterface, *MockPasswordProcessorInterface) {
	store := NewMockStoreInterface(ctrl)
	esb := NewMockESBInterface(ctrl)
	passwordProcessor := NewMockPasswordProcessorInterface(ctrl)
	return &MockApp{ESB: esb, Store: store, PasswordProcessor: passwordProcessor}, store, esb, passwordProcessor
}
