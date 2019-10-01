package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/jwt"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
)

type FirebaseController struct {
	App        *firebase.App
	UsersModel *models.UsersModel
}

func (fb *FirebaseController) CheckAuth(c *gin.Context, method func(userData hestia.User, context *gin.Context, admin bool, filter string) (res interface{}, err error), admin bool) {
	token := c.GetHeader("token")
	if token == "" {
		responses.GlobalResponseError(nil, config.ErrorHeaderToken, c)
		return
	}
	// Verify token and get user information
	fbAuth, err := fb.App.Auth(context.Background())
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorFbInitializeAuth, c)
		return
	}
	tk, err := fbAuth.VerifyIDToken(context.Background(), token)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	uid := tk.UID
user:
	userData, err := fb.UsersModel.Get(uid)
	if admin {
		if userData.Role != "admin" {
			responses.GlobalResponseError(nil, config.ErrorNoAuth, c)
			return
		}
	}
	if err != nil {
		fbUserData, err := fbAuth.GetUser(context.Background(), uid)
		if err != nil {
			responses.GlobalResponseError(nil, config.ErrorNoUserInformation, c)
			return
		}
		newUserData := hestia.User{
			ID:       fbUserData.UID,
			Email:    fbUserData.Email,
			KYCData:  hestia.KYCInformation{},
			Role:     "user",
			Orders:   []string{},
			Shifts:   []string{},
			Vouchers: []string{},
			Deposits: []string{},
			Cards:    []string{},
		}
		err = fb.UsersModel.Update(newUserData)
		if err != nil {
			responses.GlobalResponseError(nil, config.ErrorDBStore, c)
			return
		}
		goto user
	}
	filter := c.Query("filter")
	if filter == "" {
		filter = "all"
	}
	res, err := method(userData, c, admin, filter)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	switch res.(type) {
	case bool:
		responses.GlobalResponseError(res, err, c)
		return
	default:
		jwe, err := jwt.EncryptJWE(uid, res)
		responses.GlobalResponseError(jwe, err, c)
		return
	}

}

func (fb *FirebaseController) CheckToken(c *gin.Context) {
	reqBody, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorNoBody, c)
		return
	}
	headerSignature := c.GetHeader("service")
	if headerSignature == "" {
		responses.GlobalResponseError(nil, config.ErrorNoHeaderSignature, c)
		return
	}

	decodedHeader, err := jwt.DecodeJWSNoVerify(headerSignature)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorSignatureParse, c)
		return
	}

	var serviceStr string
	err = json.Unmarshal(decodedHeader, &serviceStr)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	// Check which service the request is announcing
	var pubKey string

	switch serviceStr {
	case "ladon":
		pubKey = os.Getenv("LADON_PUBLIC_KEY")
	case "tyche":
		pubKey = os.Getenv("TYCHE_PUBLIC_KEY")
	case "adrestia":
		pubKey = os.Getenv("ADRESTIA_PUBLIC_KEY")
	default:
		responses.GlobalResponseError(nil, config.ErrorWrongMessage, c)
		return
	}

	var strBody string
	err = json.Unmarshal(reqBody, &strBody)

	valid, payload := mvt.VerifyMVTToken(headerSignature, strBody, pubKey, os.Getenv("MASTER_PASSWORD"))

	if !valid {
		responses.GlobalResponseError(nil, config.ErrorInvalidPassword, c)
		return
	}
	var fbToken string
	err = json.Unmarshal(payload, &fbToken)
	if err != nil {
		fmt.Println(err)
		responses.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	// Verify token and get user information
	fbAuth, err := fb.App.Auth(context.Background())
	if err != nil {
		fmt.Println(err)
		responses.GlobalResponseError(nil, config.ErrorFbInitializeAuth, c)
		return
	}
	user, err := fbAuth.VerifyIDToken(context.Background(), fbToken)
	if err != nil {
		fmt.Println(err)
		responses.GlobalResponseError(false, nil, c)
		return
	}
	responsePayload := hestia.TokenVerification{
		Valid: true,
		UID:   user.UID,
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), responsePayload, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}
