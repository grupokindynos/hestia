package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
	"github.com/grupokindynos/common/utils"
	"github.com/grupokindynos/hestia/models"
)

/*

	VouchersController is a safe-access query for vouchers on Firestore Database
	Database Structure:

	vouchers/
		VoucherID/
			voucherData

*/

type VouchersController struct {
	Model       *models.VouchersModel
	UserModel   *models.UsersModel
	BitcouModel *models.BitcouModel
}

func (vc *VouchersController) GetAll(userData hestia.User, params Params) (interface{}, error) {
	if params.Admin {
		return vc.Model.GetAll(params.Filter)
	}
	userInfo, err := vc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, errors.ErrorNoUserInformation
	}
	var Array []hestia.Voucher
	for _, id := range userInfo.Vouchers {
		obj, err := vc.Model.Get(id)
		if err != nil {
			return nil, errors.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (vc *VouchersController) GetSingle(userData hestia.User, params Params) (interface{}, error) {
	if params.VoucherID == "" {
		return nil, errors.ErrorMissingID
	}
	if params.Admin {
		return vc.Model.Get(params.VoucherID)
	}
	userInfo, err := vc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, errors.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Vouchers, params.VoucherID) {
		return nil, errors.ErrorInfoDontMatchUser
	}
	return vc.Model.Get(params.VoucherID)
}

func (vc *VouchersController) GetVouchersByTimestampLadon(c *gin.Context) {
	// Check if the user has an id
	userId:= c.Query("userid")
	if userId == "" {
		responses.GlobalResponseError(nil, errors.ErrorMissingID, c)
		return
	}
	ts := c.Query("timestamp")
	if ts == "" {
		responses.GlobalResponseError(nil, errors.ErrorMissingID, c)
		return
	}
	_, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}

	userInfo, err := vc.UserModel.Get(userId)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorNoUserInformation, c)
		return
	}
	var userVouchers []hestia.Voucher
	timestamp, _ := strconv.ParseInt(ts, 10, 64)

	for _, id := range userInfo.Vouchers {
		obj, err := vc.Model.Get(id)
		if err != nil {
			responses.GlobalResponseError(nil, errors.ErrorNotFound, c)
			return
		}
		fmt.Println("obj:", obj)
		if timestamp <= obj.Timestamp {
			userVouchers = append(userVouchers, obj)
		}
	}
	fmt.Println("userVouchers", userVouchers)
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), userVouchers, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (vc *VouchersController) GetSingleLadon(c *gin.Context) {
	// Check if the user has an id
	id, ok := c.Params.Get("voucherid")
	if !ok {
		responses.GlobalResponseError(nil, errors.ErrorMissingID, c)
		return
	}
	_, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	voucher, err := vc.Model.Get(id)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), voucher, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (vc *VouchersController) GetAllLadon(c *gin.Context) {
	filter := c.Query("filter")
	if filter == "" {
		filter = "all"
	}
	_, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	vouchersList, err := vc.Model.GetAll(filter)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), vouchersList, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (vc *VouchersController) Store(c *gin.Context) {
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	// Try to unmarshal the information of the payload
	var voucherData hestia.Voucher
	err = json.Unmarshal(payload, &voucherData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	// Store voucher data to process
	err = vc.Model.Update(voucherData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	// Store ID on user information
	err = vc.UserModel.AddVoucher(voucherData.UID, voucherData.ID)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), voucherData.ID, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (vc *VouchersController) GetAvailableCountries(userData hestia.User, params Params) (interface{}, error) {
	usaVoucherData, err := vc.BitcouModel.GetCountry("usa")
	if err != nil {
		return nil, err
	}
	var countries []string
	for k := range usaVoucherData.Vouchers[0].Countries {
		countries = append(countries, k)
	}
	return countries, nil
}

func (vc *VouchersController) GetTestAvailableCountries(userData hestia.User, params Params) (interface{}, error) {
	usaVoucherData, err := vc.BitcouModel.GetTestCountry("usa")
	if err != nil {
		return nil, err
	}
	var countries []string
	for k := range usaVoucherData.Vouchers[0].Countries {
		countries = append(countries, k)
	}
	return countries, nil
}

func (vc *VouchersController) GetVouchers(userData hestia.User, params Params) (interface{}, error) {
	country := params.Country
	countryData, err := vc.BitcouModel.GetCountry(country)
	if err != nil {
		return nil, err
	}
	return countryData.Vouchers, nil
}

func (vc *VouchersController) GetTestVouchers(userData hestia.User, params Params) (interface{}, error) {
	country := params.Country
	countryData, err := vc.BitcouModel.GetTestCountry(country)
	if err != nil {
		return nil, err
	}
	return countryData.Vouchers, nil
}
