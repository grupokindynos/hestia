package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
	"github.com/grupokindynos/common/utils"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/services/bitcou"
	"os"
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

func (vc *VouchersController) GetCountries(userData hestia.User, params Params) (interface{}, error) {
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

func (vc *VouchersController) GetCategories(userData hestia.User, params Params) (interface{}, error) {
	country := params.Country
	countryData, err := vc.BitcouModel.GetCountry(country)
	if err != nil {
		return nil, err
	}
	cat := make(map[string]interface{})
	for _, voucher := range countryData.Vouchers {
		if voucher.Benefits["Mobile"] && voucher.Benefits["Minutes"] {
			_, ok := cat["Credit for calls"]
			if !ok {
				cat["Credit for calls"] = nil
			}
		}
		if voucher.Benefits["Mobile"] && voucher.Benefits["Data"] {
			_, ok := cat["Credit for internet"]
			if !ok {
				cat["Credit for internet"] = nil
			}
		}
		if voucher.Benefits["Mobile"] && voucher.Benefits["Data"] && voucher.Benefits["Minutes"] {
			_, ok := cat["Credit for calls and internet"]
			if !ok {
				cat["Credit for calls and internet"] = nil
			}
		}
		if voucher.Benefits["DigitalProducts"] {
			_, ok := cat["Gift Card"]
			if !ok {
				cat["Gift Card"] = nil
			}
		}
		if voucher.Benefits["Gaming"] {
			_, ok := cat["Gaming"]
			if !ok {
				cat["Gaming"] = nil
			}
		}
	}
	var catResponse []string
	for k, _ := range cat {
		catResponse = append(catResponse, k)
	}
	return catResponse, nil
}

func (vc *VouchersController) GetProviders(userData hestia.User, params Params) (interface{}, error) {
	country := params.Country
	category := params.Category
	countryData, err := vc.BitcouModel.GetCountry(country)
	if err != nil {
		return nil, err
	}
	var vouchersFiltered []bitcou.Voucher
	for _, voucher := range countryData.Vouchers {
		switch category {
		case "Credit for calls":
			if voucher.Benefits["Mobile"] && voucher.Benefits["Minutes"] {
				vouchersFiltered = append(vouchersFiltered, voucher)
			}
		case "Credit for internet":
			if voucher.Benefits["Mobile"] && voucher.Benefits["Data"] {
				vouchersFiltered = append(vouchersFiltered, voucher)
			}
		case "Credit for calls and internet":
			if voucher.Benefits["Mobile"] && voucher.Benefits["Data"] && voucher.Benefits["Minutes"] {
				vouchersFiltered = append(vouchersFiltered, voucher)
			}
		case "Gift Card":
			if voucher.Benefits["DigitalProducts"] {
				vouchersFiltered = append(vouchersFiltered, voucher)
			}
		case "Gaming":
			if voucher.Benefits["Gaming"] {
				vouchersFiltered = append(vouchersFiltered, voucher)
			}
		}

	}

	providers := make(map[string]interface{})
	for _, filtVoucher := range vouchersFiltered {
		_, ok := providers[filtVoucher.ProviderName]
		if !ok {
			if filtVoucher.ProviderName == "" {
				providers["Others"] = nil
			} else {
				providers[filtVoucher.ProviderName] = nil
			}
		}
	}
	var providerRes []string
	for k, _ := range providers {
		providerRes = append(providerRes, k)
	}
	return providerRes, nil
}

func (vc *VouchersController) GetVouchers(userData hestia.User, params Params) (interface{}, error) {
	country := params.Country
	category := params.Category
	provider := params.Provider
	countryData, err := vc.BitcouModel.GetCountry(country)
	if err != nil {
		return nil, err
	}
	var vouchersFiltered []bitcou.Voucher
	for _, voucher := range countryData.Vouchers {
		switch category {
		case "Credit for calls":
			if voucher.Benefits["Mobile"] && voucher.Benefits["Minutes"] {
				if voucher.ProviderName == provider {
					vouchersFiltered = append(vouchersFiltered, voucher)
				}
			}
		case "Credit for internet":
			if voucher.Benefits["Mobile"] && voucher.Benefits["Data"] {
				if voucher.ProviderName == provider {
					vouchersFiltered = append(vouchersFiltered, voucher)
				}
			}
		case "Credit for calls and internet":
			if voucher.Benefits["Mobile"] && voucher.Benefits["Data"] && voucher.Benefits["Minutes"] {
				if voucher.ProviderName == provider {
					vouchersFiltered = append(vouchersFiltered, voucher)
				}
			}
		case "Gift Card":
			if voucher.Benefits["DigitalProducts"] {
				if voucher.ProviderName == provider {
					vouchersFiltered = append(vouchersFiltered, voucher)
				}
			}
		case "Others":
			if voucher.Benefits["DigitalProducts"] {
				if voucher.ProviderName == provider {
					vouchersFiltered = append(vouchersFiltered, voucher)
				}
			}
		}
	}
	return vouchersFiltered, nil
}
