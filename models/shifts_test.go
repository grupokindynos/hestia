package models

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/stretchr/testify/assert"
	"testing"
)



func TestShiftsModel_Update(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := ShiftModel{
		Db:         db,
		Collection: "shifts",
	}
	err = model.Update(TestShift)
	assert.Nil(t, err)
}

func TestShiftsModel_Get(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := ShiftModel{
		Db:         db,
		Collection: "shifts",
	}
	newShift, err := model.Get(TestShift.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestShift, newShift)
}

func TestShiftsModel_GetAll(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := ShiftModel{
		Db:         db,
		Collection: "shifts",
	}
	shifts, err := model.GetAll()
	assert.Nil(t, err)
	assert.NotZero(t, len(shifts))
	assert.IsType(t, []Shift{}, shifts)
}
