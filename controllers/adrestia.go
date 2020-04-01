package controllers

import (
	"log"
	"os"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
	"github.com/grupokindynos/hestia/models"
	"strconv"
)

type AdrestiaController struct {
	Model *models.AdrestiaModel
}

func (ac *AdrestiaController) StoreDeposit(c *gin.Context) {
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	ac.storeSimpleTx(c, payload, "deposits")
}

func (ac *AdrestiaController) StoreWithdrawal(c *gin.Context) {
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	ac.storeSimpleTx(c, payload, "withdrawals")
	return
}

func (ac *AdrestiaController) storeSimpleTx(c *gin.Context, payload []byte, txType string) {
	// Try to unmarshal the information of the payload
	var simpleTx hestia.SimpleTx
	err := json.Unmarshal(payload, &simpleTx)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	log.Println("Pasa primer error")
	err = ac.Model.UpdateSimpleTx(simpleTx, txType)
	if err != nil {
		log.Println(err)
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), simpleTx.Id, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	log.Println("Sale de controller")
	return
}

func (ac *AdrestiaController) StoreBalancerOrder(c *gin.Context) {
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	// Try to unmarshal the information of the payload
	var orderData hestia.BalancerOrder
	err = json.Unmarshal(payload, &orderData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	err = ac.Model.UpdateBalancerOrder(orderData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), orderData.Id, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (ac *AdrestiaController) StoreBalancer(c *gin.Context) {
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	// Try to unmarshal the information of the payload
	var balancerData hestia.Balancer
	err = json.Unmarshal(payload, &balancerData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	err = ac.Model.UpdateBalancer(balancerData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), balancerData.Id, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (ac *AdrestiaController) UpdateDeposit(c *gin.Context) {
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	ac.updateSimpleTx(c, payload, "deposits")
	return
}

func (ac *AdrestiaController) UpdateWithdrawal(c *gin.Context) {
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	ac.updateSimpleTx(c, payload, "withdrawals")
	return
}

func (ac *AdrestiaController) updateSimpleTx(c *gin.Context, payload[]byte, txType string) {
	// Try to unmarshal the information of the payload
	var updateData hestia.SimpleTx
	err := json.Unmarshal(payload, &updateData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	_, err = ac.Model.GetSimpleTx(updateData.Id, txType)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	err = ac.Model.UpdateSimpleTx(updateData, txType)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), updateData.Id, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (ac *AdrestiaController) UpdateBalancerOrder(c *gin.Context) {
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	// Try to unmarshal the information of the payload
	var updateData hestia.BalancerOrder
	err = json.Unmarshal(payload, &updateData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	_, err = ac.Model.GetBalancerOrder(updateData.Id)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	err = ac.Model.UpdateBalancerOrder(updateData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), updateData.Id, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (ac *AdrestiaController) UpdateBalancer(c *gin.Context) {
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	// Try to unmarshal the information of the payload
	var updateData hestia.Balancer
	err = json.Unmarshal(payload, &updateData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	_, err = ac.Model.GetBalancer(updateData.Id)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	err = ac.Model.UpdateBalancer(updateData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), updateData.Id, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (ac *AdrestiaController) GetDeposits(c *gin.Context) {
	includeComplete := c.Query("include_complete")
	sinceTimestamp := c.Query("added_since")
	_, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	include, _ := strconv.ParseBool(includeComplete)
	timestamp, _ := strconv.Atoi(sinceTimestamp)
	ac.getSimpleTx(c, "deposits", include, timestamp)
	return
}

func (ac *AdrestiaController) GetWithdrawals(c *gin.Context) {
	includeComplete := c.Query("include_complete")
	sinceTimestamp := c.Query("added_since")
	_, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	include, _ := strconv.ParseBool(includeComplete)
	timestamp, _ := strconv.Atoi(sinceTimestamp)
	ac.getSimpleTx(c, "withdrawals", include, timestamp)
	return
}

func (ac *AdrestiaController) getSimpleTx(c *gin.Context, txType string, includeComplete bool, timestamp int) {
	simpleTxs, err := ac.Model.GetAllSimpleTx(includeComplete, timestamp, txType)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), simpleTxs, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (ac *AdrestiaController) GetBalancerOrders(c *gin.Context) {
	includeComplete := c.Query("include_complete")
	sinceTimestamp := c.Query("added_since")
	_, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	include, _ := strconv.ParseBool(includeComplete)
	timestamp, _ := strconv.Atoi(sinceTimestamp)
	orders, err := ac.Model.GetAllBalancerOrder(include, timestamp)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), orders, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (ac *AdrestiaController) GetBalancers(c *gin.Context) {
	includeComplete := c.Query("include_complete")
	sinceTimestamp := c.Query("added_since")
	_, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	include, _ := strconv.ParseBool(includeComplete)
	timestamp, _ := strconv.Atoi(sinceTimestamp)
	balancers, err := ac.Model.GetAllBalancerOrder(include, timestamp)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), balancers, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}