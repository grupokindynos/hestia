package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestShiftsController_GetUserAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	shifts, err := shiftCtrl.GetUserAll(models.TestUser, c)
	assert.Nil(t, err)
	var shiftsArray []models.Shift
	shiftBytes, err := json.Marshal(shifts)
	assert.Nil(t, err)
	err = json.Unmarshal(shiftBytes, &shiftsArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Shift{}, shifts)
	assert.Equal(t, models.TestShift, shiftsArray[0])
}

func TestShiftsController_GetUserSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "shiftid", Value: models.TestShift.ID}}
	shift, err := shiftCtrl.GetUserSingle(models.TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Shift{}, shift)
	assert.Equal(t, models.TestShift, shift)
}

func TestShiftsController_GetAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	shifts, err := shiftCtrl.GetAll(models.TestUser, c)
	assert.Nil(t, err)
	var shiftsArray []models.Shift
	shiftBytes, err := json.Marshal(shifts)
	assert.Nil(t, err)
	err = json.Unmarshal(shiftBytes, &shiftsArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Shift{}, shifts)
	assert.Equal(t, models.TestShift, shiftsArray[0])
}

func TestShiftsController_GetSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "shiftid", Value: models.TestShift.ID}}
	shift, err := shiftCtrl.GetSingle(models.TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Shift{}, shift)
	assert.Equal(t, models.TestShift, shift)
}