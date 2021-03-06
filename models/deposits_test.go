package models

import (
	"github.com/grupokindynos/common/hestia"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDepositsModel_Update(t *testing.T) {
	err := depositsModel.Update(TestDeposit)
	assert.Nil(t, err)
}

func TestDepositsModel_Get(t *testing.T) {
	newDeposit, err := depositsModel.Get(TestDeposit.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestDeposit, newDeposit)
}

func TestDepositsModel_GetAll(t *testing.T) {
	deposits, err := depositsModel.GetAll("all")
	assert.Nil(t, err)
	assert.NotZero(t, len(deposits))
	assert.Equal(t, TestDeposit, deposits[0])
	assert.IsType(t, []hestia.Deposit{}, deposits)
}
