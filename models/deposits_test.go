package models

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

var TestDeposit = Deposit{
	ID:  "TEST-DEPOSIT",
	UID: "XYZ12345678910",
	Payment: Payment{
		Address:       "FAKE-ADDR",
		Amount:        "123123123",
		Coin:          "POLIS",
		Txid:          "FAKE-TXID",
		Confirmations: "0",
	},
	AmountInPeso: "100",
	CardCode:     "TEST-CARDCODE",
	Status:       "COMPLETED",
	Timestamp:    "000000000000",
}

func TestDepositsModel_Update(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := DepositsModel{
		Db:         db,
		Collection: "deposits",
	}
	err = model.Update(TestDeposit)
	assert.Nil(t, err)
}

func TestDepositsModel_Get(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := DepositsModel{
		Db:         db,
		Collection: "deposits",
	}
	newDeposit, err := model.Get(TestDeposit.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestDeposit, newDeposit)
}

func TestDepositsModel_GetAll(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := DepositsModel{
		Db:         db,
		Collection: "deposits",
	}
	deposits, err := model.GetAll()
	assert.Nil(t, err)
	assert.NotZero(t, len(deposits))
	assert.IsType(t, []Deposit{}, deposits)
}
