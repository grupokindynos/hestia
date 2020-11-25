package controllers

import (
	"encoding/json"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDepositsController_GetUserAll(t *testing.T) {
	deposits, err := depositsCtrl.GetAll(models.TestUser, TestParams)
	assert.Nil(t, err)
	var depositsArray []hestia.Deposit
	depositsBytes, err := json.Marshal(deposits)
	assert.Nil(t, err)
	err = json.Unmarshal(depositsBytes, &depositsArray)
	assert.Nil(t, err)
	assert.IsType(t, []hestia.Deposit{}, deposits)
	// assert.Equal(t, models.TestDeposit, depositsArray[0])
}

func TestDepositsController_GetUserSingle(t *testing.T) {
	deposit, err := depositsCtrl.GetSingle(models.TestUser, TestParams)
	assert.Nil(t, err)
	assert.IsType(t, hestia.Deposit{}, deposit)
	assert.Equal(t, models.TestDeposit, deposit)
}

func TestDepositsController_GetAll(t *testing.T) {
	deposits, err := depositsCtrl.GetAll(models.TestUser, TestParamsAdmin)
	assert.Nil(t, err)
	var depositArray []hestia.Deposit
	depositBytes, err := json.Marshal(deposits)
	assert.Nil(t, err)
	err = json.Unmarshal(depositBytes, &depositArray)
	assert.Nil(t, err)
	assert.IsType(t, []hestia.Deposit{}, deposits)
	assert.Equal(t, models.TestDeposit, depositArray[0])
}

func TestDepositsController_GetSingle(t *testing.T) {
	deposit, err := depositsCtrl.GetSingle(models.TestUser, TestParamsAdmin)
	assert.Nil(t, err)
	assert.IsType(t, hestia.Deposit{}, deposit)
	assert.Equal(t, models.TestDeposit, deposit)
}
