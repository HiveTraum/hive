package controllers

import (
	"auth/infrastructure"
	"auth/inout"
	"context"
)

type ESB struct {
	Store infrastructure.StoreInterface
}

func (esb *ESB) OnUserChanged(id []int64) {

}

func (esb *ESB) OnEmailCodeConfirmationCreated(email string, code string) {

}

func (esb *ESB) OnPhoneCodeConfirmationCreated(phone string, code string) {

}

func (esb *ESB) OnUsersViewChanged(usersView []*inout.GetUserViewResponseV1) {

}

func (esb *ESB) OnEmailChanged(userId []int64) {
	ctx := context.Background()
	CreateOrUpdateUsersView(esb.Store, esb, ctx, userId)
}

func (esb *ESB) OnPhoneChanged(userId []int64) {
	ctx := context.Background()
	CreateOrUpdateUsersView(esb.Store, esb, ctx, userId)
}

func (esb *ESB) OnPasswordChanged(userId int64) {

}

func (esb *ESB) OnRoleChanged(roleId []int64) {
	ctx := context.Background()
	CreateOrUpdateUsersViewByRoles(esb.Store, esb, ctx, roleId)
}
