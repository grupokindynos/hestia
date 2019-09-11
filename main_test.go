package main

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/utils"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	_ = godotenv.Load()
}

func TestCardsModel_Update(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := models.CardsModel{
		Db:         db,
		Collection: "cards",
	}
	err = model.Update(models.TestCard)
	assert.Nil(t, err)
}

func TestCoinsModel_UpdateCoinsData(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := models.CoinsModel{
		Db:         db,
		Collection: "coins",
	}
	err = model.UpdateCoinsData(models.TestCoinData)
	assert.Nil(t, err)
}

func TestDepositsModel_Update(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := models.DepositsModel{
		Db:         db,
		Collection: "deposits",
	}
	err = model.Update(models.TestDeposit)
	assert.Nil(t, err)
}

func TestOrdersModel_Update(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := models.OrdersModel{
		Db:         db,
		Collection: "orders",
	}
	err = model.Update(models.TestOrder)
	assert.Nil(t, err)
}


func TestShiftsModel_Update(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := models.ShiftModel{
		Db:         db,
		Collection: "shifts",
	}
	err = model.Update(models.TestShift)
	assert.Nil(t, err)
}

func TestUsersModel_UpdateUser(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := models.UsersModel{
		Db:         db,
		Collection: "users",
	}
	err = model.UpdateUser(models.TestUser)
	assert.Nil(t, err)
}

func TestVouchersModel_Update(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := models.VouchersModel{
		Db:         db,
		Collection: "vouchers",
	}
	err = model.Update(models.TestVoucher)
	assert.Nil(t, err)
}

func TestUsersModel_AddCard(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := models.UsersModel{
		Db:         db,
		Collection: "users",
	}
	err = model.AddCard(models.TestUser.ID, models.TestCard.CardCode)
	assert.Nil(t, err)
	userData, err := model.GetUserInformation(models.TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Cards, models.TestCard.CardCode)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddDeposit(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := models.UsersModel{
		Db:         db,
		Collection: "users",
	}
	err = model.AddDeposit(models.TestUser.ID, models.TestDeposit.ID)
	assert.Nil(t, err)
	userData, err := model.GetUserInformation(models.TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Deposits, models.TestDeposit.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddVoucher(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := models.UsersModel{
		Db:         db,
		Collection: "users",
	}
	err = model.AddVoucher(models.TestUser.ID, models.TestVoucher.ID)
	assert.Nil(t, err)
	userData, err := model.GetUserInformation(models.TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Vouchers, models.TestVoucher.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddShift(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := models.UsersModel{
		Db:         db,
		Collection: "users",
	}
	err = model.AddShift(models.TestUser.ID, models.TestShift.ID)
	assert.Nil(t, err)
	userData, err := model.GetUserInformation(models.TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Shifts, models.TestShift.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddOrder(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := models.UsersModel{
		Db:         db,
		Collection: "users",
	}
	err = model.AddOrder(models.TestUser.ID, models.TestOrder.ID)
	assert.Nil(t, err)
	userData, err := model.GetUserInformation(models.TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Orders, models.TestOrder.ID)
	assert.Equal(t, true, ok)
}
