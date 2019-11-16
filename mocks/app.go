package mocks

import (
	"auth/infrastructure"
	"github.com/golang/mock/gomock"
)

type MockApp struct {
	Store infrastructure.StoreInterface
	ESB   infrastructure.ESBInterface
}

func (app *MockApp) GetStore() infrastructure.StoreInterface {
	return app.Store
}

func (app *MockApp) GetESB() infrastructure.ESBInterface {
	return app.ESB
}

func InitMockApp(ctrl *gomock.Controller) (*MockApp, *MockStoreInterface, *MockESBInterface) {
	store := NewMockStoreInterface(ctrl)
	esb := NewMockESBInterface(ctrl)
	return &MockApp{ESB: esb, Store: store}, store, esb
}
