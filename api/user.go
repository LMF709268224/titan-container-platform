package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

	"titan-container-platform/core"
	"titan-container-platform/core/dao"
	"titan-container-platform/core/token"
	"titan-container-platform/errors"
	"titan-container-platform/kubesphere"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
)

type userNonce struct {
	Code       string
	Expiration time.Time
}

var (
	defaultNonceExpiration = 5 * time.Minute
	userNonceMap           = make(map[string]*userNonce)
	mapLock                = new(sync.Mutex)
)

func putUserNonce(account string, info *userNonce) {
	mapLock.Lock()
	defer mapLock.Unlock()

	userNonceMap[account] = info
}

func getUserNonce(account string) string {
	mapLock.Lock()
	defer mapLock.Unlock()

	info := userNonceMap[account]
	if info == nil {
		return ""
	}

	if info.Expiration.Before(time.Now()) {
		return ""
	}

	return info.Code
}

func deleteUserNonce(account string) {
	mapLock.Lock()
	defer mapLock.Unlock()

	delete(userNonceMap, account)
}

func getUserInfoHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	id := claims[identityKey].(string)

	resp, err := dao.GetUserResponse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusOK, respError(errors.ErrNotFound))
		return
	}

	c.JSON(http.StatusOK, respJSON(resp))
}

func getNonceStringHandler(c *gin.Context) {
	account := c.Query("account")
	if account == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	nonce, err := generateNonceString(account)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JSONObject{
		"code": nonce,
	}))
}

func generateNonceString(account string) (string, error) {
	rand := generateRandomNumber(6)
	verifyCode := "Titan(" + rand + ")"

	putUserNonce(account, &userNonce{Code: verifyCode, Expiration: time.Now().Add(defaultNonceExpiration)})

	return verifyCode, nil
}

func generateRandomNumber(length int) string {
	seededRand := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	return fmt.Sprintf("%0*d", length, seededRand.Intn(1000000))
}

func addUserInfo(ctx context.Context, account, userName string) error {
	_, err := dao.GetUserByAccount(ctx, account)

	switch err {
	case sql.ErrNoRows:
		err = kubesphere.CreateUserAccount(account)
		if err != nil {
			log.Errorf("CreateUserAccount: %s", err.Error())
			return err
		}

		user := &core.User{
			Account:  account,
			Username: userName,
			// KubespherePwd: pwd,
		}

		err = dao.CreateUser(ctx, user)
		if err != nil {
			return err
		}
	case nil:
	default:
		return err
	}

	return nil
}

func getTokenHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	id := claims[identityKey].(string)

	code, err := token.ClaimTokens(id)
	if code > 0 {
		log.Errorf("getTokenHandler err:%v", err)
		c.JSON(http.StatusOK, respErrorCode(code, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JSONObject{
		"msg": "success",
	}))
}

func getBalanceHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	id := claims[identityKey].(string)

	balance, err := token.GetBalance(id)
	if err != nil {
		c.JSON(http.StatusOK, respError(errors.ErrNotFound))
		return
	}

	c.JSON(http.StatusOK, respJSON(JSONObject{
		"balance": balance,
	}))
}
