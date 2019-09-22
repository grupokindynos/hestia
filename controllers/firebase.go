package controllers

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/jws"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"strings"
)

type FirebaseController struct {
	App        *firebase.App
	UsersModel *models.UsersModel
}

func (fb *FirebaseController) CheckAuth(c *gin.Context, method func(userData models.User, context *gin.Context, admin bool) (res interface{}, err error), admin bool) {
	reqToken, ok := c.Request.Header["Authorization"]
	if !ok {
		config.GlobalResponseNoAuth(c)
		return
	}
	splitToken := strings.Split(reqToken[0], "Bearer ")
	token := splitToken[1]
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
	userData, err := fb.UsersModel.Get(uid)
	if admin {
		if userData.Role != "admin" {
			config.GlobalResponseError(nil, config.ErrorNoAuth, c)
			return
		}
	}
	if err != nil {
		fbUserData, err := fbAuth.GetUser(context.Background(), uid)
		if err != nil {
			config.GlobalResponseError(nil, config.ErrorNoUserInformation, c)
			return
		}
		newUserData := models.User{
			ID:       fbUserData.UID,
			Email:    fbUserData.Email,
			KYCData:  models.KYCInformation{},
			Role:     "user",
			Orders:   []string{},
			Shifts:   []string{},
			Vouchers: []string{},
			Deposits: []string{},
			Cards:    []string{},
		}
		err = fb.UsersModel.Update(newUserData)
		if err != nil {
			config.GlobalResponseError(nil, config.ErrorDBStore, c)
			return
		}
		goto user
	}
	res, err := method(userData, c, admin)
	if err != nil {
		config.GlobalResponseError(nil, err, c)
		return
	}
	switch res.(type) {
	case bool:
		config.GlobalResponseError(res, err, c)
		return
	default:
		jwe, err := jws.EncryptJWE(uid, res)
		config.GlobalResponseError(jwe, err, c)
		return
	}

}
