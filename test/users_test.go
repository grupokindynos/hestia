package test

import (
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsersModel_UpdateUser(t *testing.T) {
	err := userCtrl.Model.UpdateUser(TestUser)
	assert.Nil(t, err)
}

func TestUsersModel_AddCard(t *testing.T) {
	err := userCtrl.Model.AddCard(TestUser.ID, TestCard.CardCode)
	assert.Nil(t, err)
	userData, err := userCtrl.Model.GetUserInformation(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Cards, TestCard.CardCode)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddDeposit(t *testing.T) {
	err := userCtrl.Model.AddDeposit(TestUser.ID, TestDeposit.ID)
	assert.Nil(t, err)
	userData, err := userCtrl.Model.GetUserInformation(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Deposits, TestDeposit.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddVoucher(t *testing.T) {
	err := userCtrl.Model.AddVoucher(TestUser.ID, TestVoucher.ID)
	assert.Nil(t, err)
	userData, err := userCtrl.Model.GetUserInformation(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Vouchers, TestVoucher.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddShift(t *testing.T) {
	err := userCtrl.Model.AddShift(TestUser.ID, TestShift.ID)
	assert.Nil(t, err)
	userData, err := userCtrl.Model.GetUserInformation(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Shifts, TestShift.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_AddOrder(t *testing.T) {
	err := userCtrl.Model.AddOrder(TestUser.ID, TestOrder.ID)
	assert.Nil(t, err)
	userData, err := userCtrl.Model.GetUserInformation(TestUser.ID)
	assert.Nil(t, err)
	ok := utils.Contains(userData.Orders, TestOrder.ID)
	assert.Equal(t, true, ok)
}

func TestUsersModel_GetUserInformation(t *testing.T) {
	userData, err := userCtrl.Model.GetUserInformation(TestUser.ID)
	assert.Nil(t, err)
	assert.IsType(t, models.User{}, userData)
}
