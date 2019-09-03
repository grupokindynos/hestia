package controllers

import (
	"github.com/grupokindynos/hestia/models"
)

/*

	UsersController is a safe-access query for user information
	Database Structure:

	users/
		UID/
          	userData

*/

type User struct{}

type UsersController struct {
	Model *models.UsersModel
}
