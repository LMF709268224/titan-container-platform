package api

import (
	"net/http"
	"strconv"

	"titan-container-platform/core"
	"titan-container-platform/core/dao"
	"titan-container-platform/core/order"
	"titan-container-platform/errors"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getPriceHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	account := claims[identityKey].(string)

	log.Infoln("createOrderHandler account:", account)

	cpu, _ := strconv.Atoi(c.Query("cpu"))
	ram, _ := strconv.Atoi(c.Query("ram"))
	duration, _ := strconv.Atoi(c.Query("duration"))
	storage, _ := strconv.Atoi(c.Query("storage"))

	params := &core.OrderReq{CPUCores: cpu, RAMSize: ram, StorageSize: storage, Duration: duration}
	if checkOrderParams(params) > 0 {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	cost := order.CalculateTotalCost(params)

	c.JSON(http.StatusOK, respJSON(JSONObject{
		"cost": cost,
	}))
}

func getOrderHistoryHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	account := claims[identityKey].(string)

	log.Infoln("createOrderHandler account:", account)

	size, _ := strconv.Atoi(c.Query("size"))
	page, _ := strconv.Atoi(c.Query("page"))
	// lang := c.GetHeader("Lang")

	list, n, err := dao.LoadAccountOrders(c, account, page, size)
	if err != nil {
		log.Errorf("getOrderHistoryHandler: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JSONObject{
		"list":  list,
		"total": n,
	}))
}

func createOrderHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	account := claims[identityKey].(string)

	var params *core.OrderReq
	if err := c.BindJSON(params); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	if checkOrderParams(params) > 0 {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	orderID := uuid.NewString()
	order := &core.Order{
		Account:     account,
		CPUCores:    params.CPUCores,
		RAMSize:     params.RAMSize,
		StorageSize: params.StorageSize,
		Duration:    params.Duration,
		Status:      core.OrderStatusCreated,
		ID:          orderID,
	}

	err := dao.CreateOrder(c.Request.Context(), order)
	if err != nil {
		log.Errorf("CreateOrder: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JSONObject{
		"id": orderID,
	}))
}

func checkOrderParams(order *core.OrderReq) int {
	if order.CPUCores > 32 || order.CPUCores < 1 {
		return errors.InvalidParams
	}

	if order.RAMSize > 64 || order.CPUCores < 1 {
		return errors.InvalidParams
	}

	if order.StorageSize > 4000 || order.CPUCores < 40 {
		return errors.InvalidParams
	}

	if order.Duration > 30*24 || order.CPUCores < 1 {
		return errors.InvalidParams
	}

	return 0
}
