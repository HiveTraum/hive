package mocks

import (
	"auth/infrastructure"
	"github.com/golang/mock/gomock"
)

type MockApp struct {
	Store             infrastructure.StoreInterface
	ESB               infrastructure.ESBInterface
	PasswordProcessor infrastructure.LoginControllerInterface
}

func (app *MockApp) GetStore() infrastructure.StoreInterface {
	return app.Store
}

func (app *MockApp) GetESB() infrastructure.ESBInterface {
	return app.ESB
}

func (app *MockApp) GetLoginController() infrastructure.LoginControllerInterface {
	return app.PasswordProcessor
}

func InitMockApp(ctrl *gomock.Controller) (*MockApp, *MockStoreInterface, *MockESBInterface, *MockLoginControllerInterface) {
	store := NewMockStoreInterface(ctrl)
	esb := NewMockESBInterface(ctrl)
	passwordProcessor := NewMockLoginControllerInterface(ctrl)
	return &MockApp{ESB: esb, Store: store, PasswordProcessor: passwordProcessor}, store, esb, passwordProcessor
}
