package controllers

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
)

type FirebaseController struct {
	App        *firebase.App
	UsersModel *models.UsersModel
}

func (fb *FirebaseController) CheckAuth(c *gin.Context, method func(userData models.User, context *gin.Context) (res interface{}, err error)) {
	token := c.GetHeader("token")
	// If there is no token on the header, return non-authed
	if token == "" {
		config.GlobalResponseNoAuth(c)
		return
	}
	// Verify token and get user information
	fbAuth, err := fb.App.Auth(context.Background())
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorFbInitializeAuth, c)
		return
	}
	tk, err := fbAuth.VerifyIDToken(context.Background(), token)
	if err != nil {
		config.GlobalResponseNoAuth(c)
		return
	}
	uid := tk.UID
user:
	userData, err := fb.UsersModel.GetUserInformation(uid)
	if err != nil {
		fbUserData, err := fbAuth.GetUser(context.Background(), uid)
		if err != nil {
			config.GlobalResponseError(nil, config.ErrorNoUserInformation, c)
			return
		}
		newUserData := models.User{
			Email:    fbUserData.Email,
			KYCData:  models.KYCInformation{},
			Role:     "user",
			Shifts:   []string{},
			Vouchers: []string{},
			Deposits: []string{},
			Cards:    []string{},
		}
		err = fb.UsersModel.UpdateUser(newUserData)
		if err != nil {
			config.GlobalResponseError(nil, config.ErrorUnableToStoreUser, c)
			return
		}
		goto user
	}
	res, err := method(userData, c)
	if err != nil {
		config.GlobalResponseError(nil, err, c)
		return
	}
	config.GlobalResponseError(res, nil, c)
	return
}

func (fb *FirebaseController) CheckAuthAdmin(c *gin.Context, method func(context *gin.Context) (res interface{}, err error)) {
	token := c.GetHeader("token")
	// If there is no token on the header, return non-authed
	if token == "" {
		config.GlobalResponseNoAuth(c)
		return
	}
	// Verify token and get user information
	fbAuth, err := fb.App.Auth(context.Background())
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorFbInitializeAuth, c)
		return
	}
	tk, err := fbAuth.VerifyIDToken(context.Background(), token)
	if err != nil {
		config.GlobalResponseNoAuth(c)
		return
	}
	uid := tk.UID
	userData, err := fb.UsersModel.GetUserInformation(uid)
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorNoUserInformation, c)
		return
	}
	if userData.Role != "admin" {
		config.GlobalResponseError(nil, config.ErrorNoAuth, c)
		return
	}
	res, err := method(c)
	if err != nil {
		config.GlobalResponseError(nil, err, c)
		return
	}
	config.GlobalResponseError(res, nil, c)
	return
}
