package api

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"titan-container-platform/core"
	"titan-container-platform/errors"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/gin-gonic/gin"
)

const (
	loginStatusFailure = iota
	loginStatusSuccess
)

const titanWalletPrefix = "titan"

type login struct {
	Account  string `form:"account" json:"account"`
	UserName string `form:"user_name" json:"user_name"`
	Sign     string `form:"sign" json:"sign"`

	PublicKey string `form:"publicKey" json:"publicKey"`
}

type loginResponse struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
}

var (
	identityKey = "id"
	roleKey     = "role"
	tenantID    = "tenant_id"
	tenantName  = "tenant_name"
)

func jwtGinMiddleware(secretKey string) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:             "User",
		Key:               []byte(secretKey),
		Timeout:           24 * time.Hour,
		IdentityKey:       identityKey,
		SendAuthorization: true,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*core.User); ok {
				return jwt.MapClaims{
					identityKey: v.Account,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			return &core.User{
				Account: claims[identityKey].(string),
			}
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": loginResponse{
					Token:  token,
					Expire: expire.Format(time.RFC3339),
				},
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
			})
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginParams login
			if err := c.BindJSON(&loginParams); err != nil {
				return "", fmt.Errorf("invalid input params")
			}

			// if loginParams.Username == "" {
			// 	return "", jwt.ErrMissingLoginValues
			// }
			// if loginParams.VerifyCode == "" && loginParams.Password == "" && loginParams.Sign == "" {
			// 	return "", jwt.ErrMissingLoginValues
			// }

			if loginParams.Sign != "" {
				return loginBySignature(c, loginParams.UserName, loginParams.Account, loginParams.Sign, loginParams.PublicKey)
			}

			// if loginParams.VerifyCode != "" {
			// 	return loginByVerifyCode(c, loginParams.Username, loginParams.VerifyCode)
			// }

			// if loginParams.Password != "" {
			// 	return loginByPassword(c, loginParams.Username, loginParams.Password)
			// }

			return nil, nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			if strings.Contains(message, "Token is expired") {
				msg := "Session expired, please log in again"

				if c.GetHeader("Lang") == "cn" {
					msg = "会话已过期, 请重新登陆"
				}

				message = msg
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"msg":     message,
				"success": false,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		// TokenLookup: "header: Authorization, query: token, cookie: jwt",
		TokenLookup: "header: Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,

		RefreshResponse: func(c *gin.Context, code int, token string, t time.Time) {
			c.Next()
		},
	})
}

func loginBySignature(c *gin.Context, userName, account, msg, publicKey string) (interface{}, error) {
	nonce := getUserNonce(account)
	if nonce == "" {
		return nil, errors.NewErrorCode(errors.VerifyCodeExpired, c)
	}

	// 小狐狸
	// recoverAddress, err := verifyMessage(nonce, msg)
	// if strings.ToUpper(recoverAddress) != strings.ToUpper(account) {
	// 	return nil, errors.NewErrorCode(errors.PassWordNotAllowed, c)
	// }

	// 开普勒
	success, err := verifyCosmosAddr(account, publicKey, titanWalletPrefix)
	if err != nil || !success {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidWalletAddress, c))
		return nil, errors.NewErrorCode(errors.InvalidWalletAddress, c)
	}

	bytePubKey, err := hex.DecodeString(publicKey)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidPublicKey, c))
		return nil, errors.NewErrorCode(errors.InvalidPublicKey, c)
	}

	pubKey := secp256k1.PubKey{Key: bytePubKey}

	byteSignature, err := hex.DecodeString(msg)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidSignature, c))
		return nil, errors.NewErrorCode(errors.InvalidSignature, c)
	}

	success, err = verifyArbitraryMsg(account, nonce, byteSignature, pubKey)
	if err != nil || !success {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidSignature, c))
		return nil, errors.NewErrorCode(errors.InvalidSignature, c)
	}

	err = addUserInfo(c.Request.Context(), account, userName)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return nil, errors.NewErrorCode(errors.InternalServer, c)
	}

	return &core.User{Account: account, Username: userName}, nil
}

// func authRequired(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		claims, e := authMiddleware.GetClaimsFromJWT(ctx)
// 		if e == nil {
// 			switch v := claims["exp"].(type) {
// 			case nil:
// 				authMiddleware.Unauthorized(ctx, http.StatusUnauthorized, authMiddleware.HTTPStatusMessageFunc(jwt.ErrMissingExpField, ctx))
// 				return
// 			case float64:
// 				if int64(v) < authMiddleware.TimeFunc().Unix() {
// 					authMiddleware.Unauthorized(ctx, http.StatusUnauthorized, authMiddleware.HTTPStatusMessageFunc(jwt.ErrExpiredToken, ctx))
// 					return
// 				}
// 			case json.Number:
// 				n, err := v.Int64()
// 				if err != nil {
// 					authMiddleware.Unauthorized(ctx, http.StatusUnauthorized, authMiddleware.HTTPStatusMessageFunc(jwt.ErrWrongFormatOfExp, ctx))
// 					return
// 				}
// 				if n < authMiddleware.TimeFunc().Unix() {
// 					authMiddleware.Unauthorized(ctx, http.StatusUnauthorized, authMiddleware.HTTPStatusMessageFunc(jwt.ErrExpiredToken, ctx))
// 					return
// 				}
// 			default:
// 				authMiddleware.Unauthorized(ctx, http.StatusUnauthorized, authMiddleware.HTTPStatusMessageFunc(jwt.ErrWrongFormatOfExp, ctx))
// 				return
// 			}
// 			ctx.Set("JWT_PAYLOAD", claims)
// 			identity := authMiddleware.IdentityHandler(ctx)

// 			if identity != nil {
// 				ctx.Set(authMiddleware.IdentityKey, identity)
// 			}

// 			if !authMiddleware.Authorizator(identity, ctx) {
// 				authMiddleware.Unauthorized(ctx, http.StatusUnauthorized, authMiddleware.HTTPStatusMessageFunc(jwt.ErrForbidden, ctx))
// 				return
// 			}
// 			if int64(claims["exp"].(float64)-authMiddleware.Timeout.Seconds()/2) < authMiddleware.TimeFunc().Unix() {
// 				tokenString, _, e := authMiddleware.RefreshToken(ctx)
// 				if e == nil {
// 					ctx.Header("new-token", tokenString)
// 				}
// 			}
// 		}
// 		ctx.Next()
// 	}
// }

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, lang")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}

// func verifyMessage(message string, signedMessage string) (string, error) {
// 	// Hash the unsigned message using EIP-191
// 	hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message)) + message)
// 	hash := crypto.Keccak256Hash(hashedMessage)
// 	// Get the bytes of the signed message
// 	decodedMessage := hexutil.MustDecode(signedMessage)
// 	// Handles cases where EIP-115 is not implemented (most wallets don't implement it)
// 	if decodedMessage[64] == 27 || decodedMessage[64] == 28 {
// 		decodedMessage[64] -= 27
// 	}
// 	// Recover a public key from the signed message
// 	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), decodedMessage)
// 	if sigPublicKeyECDSA == nil {
// 		log.Errorf("Could not get a public get from the message signature")
// 	}
// 	if err != nil {
// 		return "", err
// 	}

// 	return crypto.PubkeyToAddress(*sigPublicKeyECDSA).String(), nil
// }
