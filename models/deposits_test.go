package models

import (
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
	deposits, err := depositsModel.GetAll()
	assert.Nil(t, err)
	assert.NotZero(t, len(deposits))
	assert.IsType(t, []Deposit{}, deposits)
}

