package main

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

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