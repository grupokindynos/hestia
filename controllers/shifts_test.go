package controllers

import (
	"encoding/json"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShiftsController_GetUserAll(t *testing.T) {
	shifts, err := shiftCtrl.GetAll(models.TestUser, TestParams)
	assert.Nil(t, err)
	var shiftsArray []hestia.Shift
	shiftBytes, err := json.Marshal(shifts)
	assert.Nil(t, err)
	err = json.Unmarshal(shiftBytes, &shiftsArray)
	assert.Nil(t, err)
	assert.IsType(t, []hestia.Shift{}, shifts)
	assert.Equal(t, models.TestShift, shiftsArray[0])
}

func TestShiftsController_GetUserSingle(t *testing.T) {
	shift, err := shiftCtrl.GetSingle(models.TestUser, TestParams)
	assert.Nil(t, err)
	assert.IsType(t, hestia.Shift{}, shift)
	assert.Equal(t, models.TestShift, shift)
}

/*func TestShiftsController_GetAll(t *testing.T) {
	shifts, err := shiftCtrl.GetAll(models.TestUser, TestParamsAdmin)
	assert.Nil(t, err)
	var shiftsArray []hestia.Shift
	shiftBytes, err := json.Marshal(shifts)
	assert.Nil(t, err)
	err = json.Unmarshal(shiftBytes, &shiftsArray)
	assert.Nil(t, err)
	assert.IsType(t, []hestia.Shift{}, shifts)
	assert.Equal(t, models.TestShift, shiftsArray[0])
}*/

func TestShiftsController_GetSingle(t *testing.T) {
	shift, err := shiftCtrl.GetSingle(models.TestUser, TestParamsAdmin)
	assert.Nil(t, err)
	assert.IsType(t, hestia.Shift{}, shift)
	assert.Equal(t, models.TestShift, shift)
}
