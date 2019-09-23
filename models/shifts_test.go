package models

import (
	"github.com/grupokindynos/common/hestia"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShiftsModel_Update(t *testing.T) {
	err := shiftsModel.Update(TestShift)
	assert.Nil(t, err)
}

func TestShiftsModel_Get(t *testing.T) {
	newShift, err := shiftsModel.Get(TestShift.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestShift, newShift)
}

func TestShiftsModel_GetAll(t *testing.T) {
	shifts, err := shiftsModel.GetAll("all")
	assert.Nil(t, err)
	assert.NotZero(t, len(shifts))
	assert.IsType(t, []hestia.Shift{}, shifts)
}
