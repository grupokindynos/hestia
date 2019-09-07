package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
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

func (uc *UsersController) GetSelfInfo(userInfo models.User, c *gin.Context) (interface{}, error) {
	user, err := uc.Model.GetUserInformation(userInfo.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UsersController) UpdateSelfInfo(userInfo models.User, c *gin.Context) (interface{}, error) {
	err := uc.Model.UpdateUser(userInfo)
	if err != nil {
		return nil, err
	}
	return true, nil
}

func (uc *UsersController) GetUserInfo(c *gin.Context) (interface{}, error) {
	uid, ok := c.Params.Get("uid")
	if !ok {
		return nil, config.ErrorMissingUID
	}
	user, err := uc.Model.GetUserInformation(uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UsersController) UpdateUserInfo(c *gin.Context) (interface{}, error) {
	var userDataRaw []byte
	var userData models.User
	_, err := c.Request.Body.Read(userDataRaw)
	if err != nil {
		return nil, config.ErrorReadBody
	}
	err = json.Unmarshal(userDataRaw, &userData)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	err = uc.Model.UpdateUser(userData)
	if err != nil {
		return nil, err
	}
	return true, nil
}
