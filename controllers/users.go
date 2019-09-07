package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
)

/*

	UsersController is a safe-access query for user information
	Database Structure:

	users/
		UID/
          	userData

*/

type UsersController struct {
	Model *models.UsersModel
}

func (uc *UsersController) GetSelfInfo(uid string, params gin.Params) (interface{}, error) {
	return nil, nil
}

func (uc *UsersController) UpdateSelfInfo(uid string, params gin.Params) (interface{}, error) {
	return nil, nil
}

func (uc *UsersController) GetUserInfo(params gin.Params) (interface{}, error) {
	return nil, nil
}

func (uc *UsersController) UpdateUserInfo(params gin.Params) (interface{}, error) {
	return nil, nil
}
