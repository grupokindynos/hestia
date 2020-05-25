package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	e "errors"
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

type VouchersControllerV2 struct {
	Model           *models.VouchersModelV2
	UserModel       *models.UsersModel
	BitcouModel     *models.BitcouModel
	BitcouConfModel *models.BitcouConfModel
	CachedVouchers  VouchersCacheV2
}

func (vc *VouchersControllerV2) GetAll(userData hestia.User, params Params) (interface{}, error) {
	if params.Admin {
		f, _ := strconv.Atoi(params.Filter)
		return vc.Model.GetAll(f, "")
	}
	userInfo, err := vc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, errors.ErrorNoUserInformation
	}
	var Array []hestia.VoucherV2
	for _, id := range userInfo.VouchersV2 {
		obj, err := vc.Model.Get(id)
		if err != nil {
			continue
			/* return nil, errors.ErrorNotFound*/
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (vc *VouchersControllerV2) GetSingle(userData hestia.User, params Params) (interface{}, error) {
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

func (vc *VouchersControllerV2) GetVouchersByTimestampLadon(c *gin.Context) {
	// Check if the user has an id
	userId := c.Query("userid")
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
	var userVouchers []hestia.VoucherV2
	timestamp, _ := strconv.ParseInt(ts, 10, 64)

	for _, id := range userInfo.Vouchers {
		obj, err := vc.Model.Get(id)
		if err != nil {
			continue
		}

		if timestamp <= obj.CreatedTime {
			userVouchers = append(userVouchers, obj)
		}
	}

	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), userVouchers, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (vc *VouchersControllerV2) GetSingleLadon(c *gin.Context) {
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

func (vc *VouchersControllerV2) GetVoucherInfo(c *gin.Context) {
	// Check if the user has an id
	id, ok := c.Params.Get("product_id")
	country, okCountry := c.Params.Get("country")
	fmt.Println(country)
	if !ok || !okCountry {
		responses.GlobalResponseError(nil, errors.ErrorMissingID, c)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	_, err = mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	vouchers, err := vc.BitcouModel.GetCountryV2(country)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
	}
	for _, v := range vouchers.Vouchers {
		if v.ProductID == idInt {
			header, body, _ := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), v, os.Getenv("HESTIA_PRIVATE_KEY"))
			responses.GlobalResponseMRT(header, body, c)
			return
		}
	}
	if err != nil {
		responses.GlobalResponseError(nil, e.New("voucher id not found"), c)
		return
	}
}

func (vc *VouchersControllerV2) GetAllLadon(c *gin.Context) {
	filter := c.Query("filter")
	if filter == "" {
		filter = "all"
	}
	_, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	f, _ := strconv.Atoi(filter)
	vouchersList, err := vc.Model.GetAll(f, "")
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), vouchersList, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (vc *VouchersControllerV2) Store(c *gin.Context) {
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	// Try to unmarshal the information of the payload
	var voucherData hestia.VoucherV2
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
	err = vc.UserModel.AddVoucherV2(voucherData.UserId, voucherData.Id)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), voucherData.Id, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (vc *VouchersControllerV2) GetAvailableCountries(userData hestia.User, params Params) (interface{}, error) {
	if len(vc.CachedVouchers.CachedCountries) > 0 && vc.CachedVouchers.CachedCountriesUpdated+voucherCacheTimeFrame > time.Now().Unix() {
		return vc.CachedVouchers.CachedCountries, nil
	} else {
		countries, err := vc.BitcouModel.GetCountries(false)
		if err != nil {
			return nil, err
		}
		vc.CachedVouchers.CachedCountries = countries
		vc.CachedVouchers.CachedCountriesUpdated = time.Now().Unix()
		return countries, nil
	}
}

func (vc *VouchersControllerV2) GetTestAvailableCountries(userData hestia.User, params Params) (interface{}, error) {
	countries, err := vc.BitcouModel.GetCountries(true)
	if err != nil {
		return nil, err
	}
	return countries, nil
}

func (vc *VouchersControllerV2) GetVouchersV2(_ hestia.User, params Params) (interface{}, error) {
	cachedData, ok := vc.CachedVouchers.Vouchers[params.Country]
	if !ok {
		countryData, err := vc.BitcouModel.GetCountryV2(params.Country)
		if err != nil {
			return nil, err
		}
		vc.CachedVouchers.AddCountryVouchersV2(params.Country, countryData.Vouchers)
		return countryData.Vouchers, nil
	}
	if cachedData.LastUpdated+voucherCacheTimeFrame > time.Now().Unix() {
		return vc.CachedVouchers.Vouchers[params.Country].Vouchers, nil
	} else {
		countryData, err := vc.BitcouModel.GetCountryV2(params.Country)
		if err != nil {
			return nil, err
		}
		vc.CachedVouchers.AddCountryVouchersV2(params.Country, countryData.Vouchers)
		return countryData.Vouchers, nil
	}
}


func (vc *VouchersControllerV2) GetTestVouchers(userData hestia.User, params Params) (interface{}, error) {
	country := params.Country
	countryData, err := vc.BitcouModel.GetTestCountry(country)
	if err != nil {
		return nil, err
	}
	return countryData.Vouchers, nil
}

func (vc *VouchersControllerV2) AddFilters(c *gin.Context) {
	// Try to unmarshal the information of the payload
	var filterData models.ApiBitcouFilter
	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(reqBody, &filterData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	if filterData.Target != "dev" && filterData.Target != "prod" {
		responses.GlobalResponseError(nil, e.New("api can only have one of the following values (dev | prod)"), c)
		return
	}

	filter := models.BitcouFilter{
		ID:        filterData.Target,
		Providers: filterData.Providers,
		Vouchers:  filterData.Vouchers,
	}
	err = vc.BitcouConfModel.UpdateFilters(filter)
	if err != nil {
		responses.GlobalResponseError(nil, e.New("failed to update filter"), c)
		return
	}
	// Store voucher data to processs
	return
}

