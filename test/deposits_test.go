package test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestDepositsModel_Update(t *testing.T) {
	err := depositsCtrl.Model.Update(TestDeposit)
	assert.Nil(t, err)
}

func TestDepositsModel_Get(t *testing.T) {
	newDeposit, err := depositsCtrl.Model.Get(TestDeposit.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestDeposit, newDeposit)
}

func TestDepositsModel_GetAll(t *testing.T) {
	deposits, err := depositsCtrl.Model.GetAll()
	assert.Nil(t, err)
	assert.NotZero(t, len(deposits))
	assert.IsType(t, []models.Deposit{}, deposits)
}

func TestDepositsController_GetUserAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	deposits, err := depositsCtrl.GetUserAll(TestUser, c)
	assert.Nil(t, err)
	var depositsArray []models.Deposit
	depositsBytes, err := json.Marshal(deposits)
	assert.Nil(t, err)
	err = json.Unmarshal(depositsBytes, &depositsArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Deposit{}, deposits)
	assert.Equal(t, TestDeposit, depositsArray[0])
}

func TestDepositsController_GetUserSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "depositid", Value: TestDeposit.ID}}
	deposit, err := depositsCtrl.GetUserSingle(TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Deposit{}, deposit)
	assert.Equal(t, TestDeposit, deposit)
}

func TestDepositsController_GetAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	deposits, err := depositsCtrl.GetAll(TestUser, c)
	assert.Nil(t, err)
	var depositArray []models.Deposit
	depositBytes, err := json.Marshal(deposits)
	assert.Nil(t, err)
	err = json.Unmarshal(depositBytes, &depositArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Deposit{}, deposits)
	assert.Equal(t, TestDeposit, depositArray[0])
}

func TestDepositsController_GetSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "depositid", Value: TestDeposit.ID}}
	deposit, err := depositsCtrl.GetSingle(TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Deposit{}, deposit)
	assert.Equal(t, TestDeposit, deposit)
}
