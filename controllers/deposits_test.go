package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestDepositsController_GetUserAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	deposits, err := depositsCtrl.GetUserAll(models.TestUser, c)
	assert.Nil(t, err)
	var depositsArray []models.Deposit
	depositsBytes, err := json.Marshal(deposits)
	assert.Nil(t, err)
	err = json.Unmarshal(depositsBytes, &depositsArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Deposit{}, deposits)
	assert.Equal(t, models.TestDeposit, depositsArray[0])
}

func TestDepositsController_GetUserSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "depositid", Value: models.TestDeposit.ID}}
	deposit, err := depositsCtrl.GetUserSingle(models.TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Deposit{}, deposit)
	assert.Equal(t, models.TestDeposit, deposit)
}

func TestDepositsController_GetAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	deposits, err := depositsCtrl.GetAll(models.TestUser, c)
	assert.Nil(t, err)
	var depositArray []models.Deposit
	depositBytes, err := json.Marshal(deposits)
	assert.Nil(t, err)
	err = json.Unmarshal(depositBytes, &depositArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Deposit{}, deposits)
	assert.Equal(t, models.TestDeposit, depositArray[0])
}

func TestDepositsController_GetSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "depositid", Value: models.TestDeposit.ID}}
	deposit, err := depositsCtrl.GetSingle(models.TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Deposit{}, deposit)
	assert.Equal(t, models.TestDeposit, deposit)
}
