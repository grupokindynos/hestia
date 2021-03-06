package controllers

import (
	"encoding/json"
	"github.com/grupokindynos/hestia/services/bitcou"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
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

const voucherCacheTimeFrame = 2 * 60 * 60 // 24 hours

type CachedVouchersData struct {
	LastUpdated int64
	Vouchers    []bitcou.LightVoucher
}

type CachedVouchersDataV2 struct {
	LastUpdated int64
	Vouchers    []bitcou.LightVoucherV2
}

type VouchersCache struct {
	lock                   sync.RWMutex
	Vouchers               map[string]CachedVouchersData
	CachedCountries        []string
	CachedCountriesUpdated int64
}

type VouchersCacheV2 struct {
	lock                   sync.RWMutex
	Vouchers               map[string]CachedVouchersDataV2
	CachedCountries        []string
	CachedCountriesUpdated int64
}

func (vc *VouchersCache) AddCountryVouchers(country string, vouchers []bitcou.LightVoucher) {
	vc.lock.Lock()
	vc.Vouchers[country] = CachedVouchersData{
		LastUpdated: time.Now().Unix(),
		Vouchers:    vouchers,
	}
	vc.lock.Unlock()
	return
}

func (vc *VouchersCacheV2) AddCountryVouchersV2(country string, vouchers []bitcou.LightVoucherV2) {
	vc.lock.Lock()
	vc.Vouchers[country] = CachedVouchersDataV2{
		LastUpdated: time.Now().Unix(),
		Vouchers:    vouchers,
	}
	vc.lock.Unlock()
	return
}

type VouchersController struct {
	Model           *models.VouchersModel
	UserModel       *models.UsersModel
	BitcouModel     *models.BitcouModel
	BitcouConfModel *models.BitcouConfModel
	CachedVouchers  VouchersCache
}

func (vc *VouchersController) GetAll(userData hestia.User, params Params) (interface{}, error) {
	if params.Admin {
		return vc.Model.GetAll(params.Filter, "")
	}
	userInfo, err := vc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, errors.ErrorNoUserInformation
	}
	var Array []hestia.Voucher
	for _, id := range userInfo.Vouchers {
		obj, err := vc.Model.Get(id)
		if err != nil {
			continue
			/* return nil, errors.ErrorNotFound*/
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
	_, _, err := mvt.VerifyRequest(c)
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
			continue
			/* responses.GlobalResponseError(nil, errors.ErrorNotFound, c)
			return*/
		}

		if timestamp <= obj.Timestamp {
			userVouchers = append(userVouchers, obj)
		}
	}

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
	_, _, err := mvt.VerifyRequest(c)
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
	_, _, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	vouchersList, err := vc.Model.GetAll(filter, "")
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), vouchersList, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (vc *VouchersController) Store(c *gin.Context) {
	payload, _, err := mvt.VerifyRequest(c)
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

func (vc *VouchersController) GetTestAvailableCountries(_ hestia.User, _ Params) (interface{}, error) {
	countries, err := vc.BitcouModel.GetCountries(true)
	if err != nil {
		return nil, err
	}
	return countries, nil
}

func (vc *VouchersController) GetVouchers(_ hestia.User, params Params) (interface{}, error) {
	cachedData, ok := vc.CachedVouchers.Vouchers[params.Country]
	if !ok {
		countryData, err := vc.BitcouModel.GetCountry(params.Country)
		if err != nil {
			return nil, err
		}
		vc.CachedVouchers.AddCountryVouchers(params.Country, countryData.Vouchers)
		return countryData.Vouchers, nil
	}
	if cachedData.LastUpdated+voucherCacheTimeFrame > time.Now().Unix() {
		return vc.CachedVouchers.Vouchers[params.Country].Vouchers, nil
	} else {
		countryData, err := vc.BitcouModel.GetCountry(params.Country)
		if err != nil {
			return nil, err
		}
		vc.CachedVouchers.AddCountryVouchers(params.Country, countryData.Vouchers)
		return countryData.Vouchers, nil
	}
}

func (vc *VouchersController) GetTestVouchers(userData hestia.User, params Params) (interface{}, error) {
	country := params.Country
	countryData, err := vc.BitcouModel.GetTestCountry(country)
	if err != nil {
		return nil, err
	}
	return countryData.Vouchers, nil
}

func (vc *VouchersController) AddFilters(c *gin.Context) {
	// Try to unmarshal the information of the payload
	var filterData models.ApiBitcouFilter
	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(reqBody, &filterData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	log.Println(filterData)
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