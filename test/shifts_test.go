package test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestShiftsModel_Update(t *testing.T) {
	err := shiftCtrl.Model.Update(TestShift)
	assert.Nil(t, err)
}

func TestShiftsModel_Get(t *testing.T) {
	newShift, err := shiftCtrl.Model.Get(TestShift.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestShift, newShift)
}

func TestShiftsModel_GetAll(t *testing.T) {
	shifts, err := shiftCtrl.Model.GetAll()
	assert.Nil(t, err)
	assert.NotZero(t, len(shifts))
	assert.IsType(t, []models.Shift{}, shifts)
}

func TestShiftsController_GetUserAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	shifts, err := shiftCtrl.GetUserAll(TestUser, c)
	assert.Nil(t, err)
	var shiftsArray []models.Shift
	shiftBytes, err := json.Marshal(shifts)
	assert.Nil(t, err)
	err = json.Unmarshal(shiftBytes, &shiftsArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Shift{}, shifts)
	assert.Equal(t, TestShift, shiftsArray[0])
}

func TestShiftsController_GetUserSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "shiftid", Value: TestShift.ID}}
	shift, err := shiftCtrl.GetUserSingle(TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Shift{}, shift)
	assert.Equal(t, TestShift, shift)
}

func TestShiftsController_GetAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	shifts, err := shiftCtrl.GetAll(TestUser, c)
	assert.Nil(t, err)
	var shiftsArray []models.Shift
	shiftBytes, err := json.Marshal(shifts)
	assert.Nil(t, err)
	err = json.Unmarshal(shiftBytes, &shiftsArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Shift{}, shifts)
	assert.Equal(t, TestShift, shiftsArray[0])
}

func TestShiftsController_GetSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "shiftid", Value: TestShift.ID}}
	shift, err := shiftCtrl.GetSingle(TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Shift{}, shift)
	assert.Equal(t, TestShift, shift)
}
