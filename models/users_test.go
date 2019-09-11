package models

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestUsersModel_GetUserInformation(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := UsersModel{
		Db:         db,
		Collection: "users",
	}
	userData, err := model.GetUserInformation(TestUser.ID)
	assert.Nil(t, err)
	assert.IsType(t, User{}, userData)
}