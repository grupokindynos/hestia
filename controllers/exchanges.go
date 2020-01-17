package controllers

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
	"github.com/grupokindynos/hestia/models"
)

type ExchangesController struct {
	Model *models.ExchangesModel
}

func (ec *ExchangesController) GetOrders(c *gin.Context) {
	includeComplete := c.Query("include_complete")
	sinceTimestamp := c.Query("added_since")
	_, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	include, _ := strconv.ParseBool(includeComplete)
	timestamp, _ := strconv.Atoi(sinceTimestamp)
	orders, err := ec.Model.GetAll(include, timestamp)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), orders, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (ec *ExchangesController) StoreOrder(c *gin.Context) {
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	// Try to unmarshal the information of the payload
	var orderData hestia.AdrestiaOrder
	err = json.Unmarshal(payload, &orderData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	err = ec.Model.Update(orderData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), orderData.ID, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (ec *ExchangesController) UpdateOrder(c *gin.Context) {
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	// Try to unmarshal the information of the payload
	var updateData hestia.AdrestiaOrder
	err = json.Unmarshal(payload, &updateData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	_, err = ec.Model.Get(updateData.ID)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	err = ec.Model.Update(updateData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), updateData.ID, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (ec *ExchangesController) UpdateOrderStatus(c *gin.Context) {
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	// Try to unmarshal the information of the payload
	var updateData hestia.AdrestiaOrderUpdate
	err = json.Unmarshal(payload, &updateData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	orderInfo, err := ec.Model.Get(updateData.ID)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	orderInfo.Status = updateData.Status
	err = ec.Model.Update(orderInfo)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), orderInfo.ID, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}
