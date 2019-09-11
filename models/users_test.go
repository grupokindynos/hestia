package models

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsersModel_UpdateUser(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := UsersModel{
		Db:         db,
		Collection: "users",
	}
	err = model.UpdateUser(TestUser)
	assert.Nil(t, err)
}

func TestUsersModel_GetUserInformation(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := UsersModel{
		Db:         db,
		Collection: "users",
	}
	userData, err := model.GetUserInformation(TestUser.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestUser, userData)
}

func TestUsersModel_AddCard(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := UsersModel{
		Db:         db,
		Collection: "users",
	}
	err = model.AddCard(TestUser.ID, TestCard.CardCode)
	assert.Nil(t, err)
	userData, err := model.GetUserInformation(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Cards, TestCard.CardCode)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddDeposit(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := UsersModel{
		Db:         db,
		Collection: "users",
	}
	err = model.AddDeposit(TestUser.ID, TestDeposit.ID)
	assert.Nil(t, err)
	userData, err := model.GetUserInformation(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Deposits, TestDeposit.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddVoucher(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := UsersModel{
		Db:         db,
		Collection: "users",
	}
	err = model.AddVoucher(TestUser.ID, TestVoucher.ID)
	assert.Nil(t, err)
	userData, err := model.GetUserInformation(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Vouchers, TestVoucher.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddShift(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := UsersModel{
		Db:         db,
		Collection: "users",
	}
	err = model.AddShift(TestUser.ID, TestShift.ID)
	assert.Nil(t, err)
	userData, err := model.GetUserInformation(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Shifts, TestShift.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddOrder(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := UsersModel{
		Db:         db,
		Collection: "users",
	}
	err = model.AddOrder(TestUser.ID, TestOrder.ID)
	assert.Nil(t, err)
	userData, err := model.GetUserInformation(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Orders, TestOrder.ID)
	assert.Equal(t, true, ok)
}
