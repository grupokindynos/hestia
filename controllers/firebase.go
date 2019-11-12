package controllers

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/jwt"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
	"github.com/grupokindynos/hestia/models"
	"os"
)

type Params struct {
	Admin     bool
	Filter    string
	ShiftID   string
	OrderID   string
	UserID    string
	VoucherID string
	DepositID string
	CardCode  string
	Body      []byte
}

type FirebaseController struct {
	App        *firebase.App
	UsersModel *models.UsersModel
}

func (fb *FirebaseController) CheckAuth(c *gin.Context, method func(userData hestia.User, params Params) (res interface{}, err error), admin bool) {
	token := c.GetHeader("token")
	if token == "" {
		responses.GlobalResponseError(nil, errors.ErrorHeaderToken, c)
		return
	}
	// Verify token and get user information
	fbAuth, err := fb.App.Auth(context.Background())
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorFbInitializeAuth, c)
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
			responses.GlobalResponseError(nil, errors.ErrorNoAuth, c)
			return
		}
	}
	if err != nil {
		fbUserData, err := fbAuth.GetUser(context.Background(), uid)
		if err != nil {
			responses.GlobalResponseError(nil, errors.ErrorNoUserInformation, c)
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
			responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
			return
		}
		goto user
	}

	// Grab Params
	filter := c.Query("filter")
	if filter == "" {
		filter = "all"
	}
	var ReqBody hestia.BodyReq
	var bodyBytes []byte
	reqBody, _ := c.GetRawData()
	if len(reqBody) > 0 {
		err := json.Unmarshal(reqBody, &ReqBody)
		if err != nil {
			responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
			return
		}
		bodyBytes, err = jwt.DecryptJWE(userData.ID, ReqBody.Payload)
		if err != nil {
			responses.GlobalResponseError(nil, errors.ErrorDecryptJWE, c)
			return
		}
	}
	params := Params{
		Admin:     admin,
		Filter:    filter,
		ShiftID:   c.Param("shiftid"),
		OrderID:   c.Param("orderid"),
		UserID:    c.Param("uid"),
		VoucherID: c.Param("voucherid"),
		DepositID: c.Param("depositid"),
		CardCode:  c.Param("cardcode"),
		Body:      bodyBytes,
	}
	res, err := method(userData, params)
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
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	var fbToken string
	err = json.Unmarshal(payload, &fbToken)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	// Verify token and get user information
	fbAuth, err := fb.App.Auth(context.Background())
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorFbInitializeAuth, c)
		return
	}
	user, err := fbAuth.VerifyIDToken(context.Background(), fbToken)
	if err != nil {
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
