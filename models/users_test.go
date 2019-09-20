package models

import (
	"github.com/grupokindynos/common/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsersModel_UpdateUser(t *testing.T) {
	err := usersModel.Update(TestUser)
	assert.Nil(t, err)
}

func TestUsersModel_AddCard(t *testing.T) {
	err := usersModel.AddCard(TestUser.ID, TestCard.CardCode)
	assert.Nil(t, err)
	userData, err := usersModel.Get(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Cards, TestCard.CardCode)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddDeposit(t *testing.T) {
	err := usersModel.AddDeposit(TestUser.ID, TestDeposit.ID)
	assert.Nil(t, err)
	userData, err := usersModel.Get(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Deposits, TestDeposit.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddVoucher(t *testing.T) {
	err := usersModel.AddVoucher(TestUser.ID, TestVoucher.ID)
	assert.Nil(t, err)
	userData, err := usersModel.Get(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Vouchers, TestVoucher.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddShift(t *testing.T) {
	err := usersModel.AddShift(TestUser.ID, TestShift.ID)
	assert.Nil(t, err)
	userData, err := usersModel.Get(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Shifts, TestShift.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddOrder(t *testing.T) {
	err := usersModel.AddOrder(TestUser.ID, TestOrder.ID)
	assert.Nil(t, err)
	userData, err := usersModel.Get(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Orders, TestOrder.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_GetUserInformation(t *testing.T) {
	userData, err := usersModel.Get(TestUser.ID)
	assert.Nil(t, err)
	assert.IsType(t, User{}, userData)
}
