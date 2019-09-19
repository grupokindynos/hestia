package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/jws"
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

func (uc *UsersController) GetAll(userInfo models.User, c *gin.Context, admin bool) (interface{}, error) {
	if admin {
		return uc.Model.GetAll()
	}
	return nil, config.ErrorNoAuth
}

func (uc *UsersController) GetSingle(userInfo models.User, c *gin.Context, admin bool) (interface{}, error) {
	if admin {
		id, ok := c.Params.Get("uid")
		if !ok {
			return nil, config.ErrorMissingID
		}
		return uc.Model.Get(id)
	}
	user, err := uc.Model.Get(userInfo.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UsersController) Store(userData models.User, c *gin.Context) (interface{}, error) {
	var ReqBody models.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	rawBytes, err := jws.DecryptJWE(userData.ID, ReqBody.Payload)
	if err != nil {
		return nil, config.ErrorDecryptJWE
	}
	var newUserData models.User
	err = json.Unmarshal(rawBytes, &newUserData)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	oldUserData, err := uc.Model.Get(newUserData.ID)
	if err != nil {
		return nil, err
	}
	updateUserData := models.User{
		ID:       oldUserData.ID,
		Email:    oldUserData.Email,
		KYCData:  userData.KYCData,
		Role:     userData.Role,
		Shifts:   oldUserData.Shifts,
		Vouchers: oldUserData.Vouchers,
		Deposits: oldUserData.Deposits,
		Cards:    oldUserData.Cards,
	}
	err = uc.Model.Update(updateUserData)
	if err != nil {
		return nil, err
	}
	return true, nil
}
