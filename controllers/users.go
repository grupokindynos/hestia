package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/utils"
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

// GetSelfInfo is a protected function to help the user get their own information
func (uc *UsersController) GetSelfInfo(userInfo models.User, c *gin.Context) (interface{}, error) {
	user, err := uc.Model.GetUserInformation(userInfo.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserInfo is a protected function to help any admin to get the information of any user.
func (uc *UsersController) GetUserInfo(userData models.User, c *gin.Context) (interface{}, error) {
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

// UpdateUserInfo is a protected function to update role and kyc information for any user.
// This function will only work with already registered users, if the user is not registered, it will return error.
// Shift, Vouchers, Deposits and Card information should be handled on other controllers.
// ID and Email must never be changed.
func (uc *UsersController) UpdateUserInfo(userData models.User, c *gin.Context) (interface{}, error) {
	var ReqBody models.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	rawBytes, err := utils.DecryptJWE(userData.ID, ReqBody.Payload)
	if err != nil {
		return nil, config.ErrorDecryptJWE
	}
	var newUserData models.User
	err = json.Unmarshal(rawBytes, &newUserData)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	oldUserData, err := uc.Model.GetUserInformation(newUserData.ID)
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
	err = uc.Model.UpdateUser(updateUserData)
	if err != nil {
		return nil, err
	}
	return true, nil
}
