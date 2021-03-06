package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/jwt"
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

func (uc *UsersController) GetAll(userInfo hestia.User, c *gin.Context, admin bool, filter string) (interface{}, error) {
	if admin {
		return uc.Model.GetAll()
	}
	return nil, errors.ErrorNoAuth
}

func (uc *UsersController) GetSingle(userInfo hestia.User, params Params) (interface{}, error) {
	if params.Admin {
		if params.UserID == "" {
			return nil, errors.ErrorMissingID
		}
		return uc.Model.Get(params.UserID)
	}
	user, err := uc.Model.Get(userInfo.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UsersController) Store(userData hestia.User, c *gin.Context) (interface{}, error) {
	var ReqBody hestia.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		return nil, errors.ErrorUnmarshal
	}
	rawBytes, err := jwt.DecryptJWE(userData.ID, ReqBody.Payload)
	if err != nil {
		return nil, errors.ErrorDecryptJWE
	}
	var newUserData hestia.User
	err = json.Unmarshal(rawBytes, &newUserData)
	if err != nil {
		return nil, errors.ErrorUnmarshal
	}
	oldUserData, err := uc.Model.Get(newUserData.ID)
	if err != nil {
		return nil, err
	}
	updateUserData := hestia.User{
		ID:         oldUserData.ID,
		Email:      oldUserData.Email,
		KYCData:    userData.KYCData,
		Role:       userData.Role,
		Shifts:     oldUserData.Shifts,
		Vouchers:   oldUserData.Vouchers,
		Deposits:   oldUserData.Deposits,
		Cards:      oldUserData.Cards,
		ShiftV2:    oldUserData.ShiftV2,
		VouchersV2: oldUserData.VouchersV2,
	}
	err = uc.Model.Update(updateUserData)
	if err != nil {
		return nil, err
	}
	return true, nil
}
