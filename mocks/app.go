package mocks

import (
	"auth/infrastructure"
	"github.com/golang/mock/gomock"
)

type MockApp struct {
	Store             infrastructure.StoreInterface
	ESB               infrastructure.ESBInterface
	LoginController   infrastructure.LoginControllerInterface
	PasswordProcessor infrastructure.PasswordProcessorInterface
}

func (app *MockApp) GetPasswordProcessor() infrastructure.PasswordProcessorInterface {
	return app.PasswordProcessor
}

func (app *MockApp) GetStore() infrastructure.StoreInterface {
	return app.Store
}

func (app *MockApp) GetESB() infrastructure.ESBInterface {
	return app.ESB
}

func (app *MockApp) GetLoginController() infrastructure.LoginControllerInterface {
	return app.LoginController
}

func InitMockApp(ctrl *gomock.Controller) (*MockApp, *MockStoreInterface, *MockESBInterface, *MockLoginControllerInterface, *MockPasswordProcessorInterface) {
	store := NewMockStoreInterface(ctrl)
	esb := NewMockESBInterface(ctrl)
	loginController := NewMockLoginControllerInterface(ctrl)
	passwordProcessor := NewMockPasswordProcessorInterface(ctrl)
	return &MockApp{ESB: esb, Store: store, LoginController: loginController, PasswordProcessor: passwordProcessor}, store, esb, loginController, passwordProcessor
}
