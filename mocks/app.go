package mocks

import (
	"auth/infrastructure"
	"github.com/golang/mock/gomock"
)

type MockApp struct {
	Store           infrastructure.StoreInterface
	ESB             infrastructure.ESBInterface
	LoginController infrastructure.LoginControllerInterface
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

func InitMockApp(ctrl *gomock.Controller) (*MockApp, *MockStoreInterface, *MockESBInterface, *MockLoginControllerInterface) {
	store := NewMockStoreInterface(ctrl)
	esb := NewMockESBInterface(ctrl)
	loginController := NewMockLoginControllerInterface(ctrl)
	return &MockApp{ESB: esb, Store: store, LoginController: loginController}, store, esb, loginController
}
