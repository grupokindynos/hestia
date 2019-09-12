package main

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/controllers"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/utils"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

var userCtrl controllers.UsersController

func init() {
	_ = godotenv.Load("../.env")
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	// Init DB models
	usersModel := &models.UsersModel{Db: db, Collection: "users"}

	// Init Controllers
	userCtrl = controllers.UsersController{Model: usersModel}
}

func TestUsersModel_UpdateUser(t *testing.T) {
	err := userCtrl.Model.UpdateUser(models.TestUser)
	assert.Nil(t, err)
}

func TestUsersModel_AddCard(t *testing.T) {
	err := userCtrl.Model.AddCard(models.TestUser.ID, models.TestCard.CardCode)
	assert.Nil(t, err)
	userData, err := userCtrl.Model.GetUserInformation(models.TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Cards, models.TestCard.CardCode)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddDeposit(t *testing.T) {
	err := userCtrl.Model.AddDeposit(models.TestUser.ID, models.TestDeposit.ID)
	assert.Nil(t, err)
	userData, err := userCtrl.Model.GetUserInformation(models.TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Deposits, models.TestDeposit.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddVoucher(t *testing.T) {
	err := userCtrl.Model.AddVoucher(models.TestUser.ID, models.TestVoucher.ID)
	assert.Nil(t, err)
	userData, err := userCtrl.Model.GetUserInformation(models.TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Vouchers, models.TestVoucher.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddShift(t *testing.T) {
	err := userCtrl.Model.AddShift(models.TestUser.ID, models.TestShift.ID)
	assert.Nil(t, err)
	userData, err := userCtrl.Model.GetUserInformation(models.TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Shifts, models.TestShift.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddOrder(t *testing.T) {
	err := userCtrl.Model.AddOrder(models.TestUser.ID, models.TestOrder.ID)
	assert.Nil(t, err)
	userData, err := userCtrl.Model.GetUserInformation(models.TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Orders, models.TestOrder.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_GetUserInformation(t *testing.T) {
	userData, err := userCtrl.Model.GetUserInformation(models.TestUser.ID)
	assert.Nil(t, err)
	assert.IsType(t, models.User{}, userData)
}

