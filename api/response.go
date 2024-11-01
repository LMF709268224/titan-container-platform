package api

import (
	"strings"

	err "titan-container-platform/errors"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// JSONObject represents a JSON object.
type JSONObject map[string]interface{}

func respJSON(v interface{}) gin.H {
	return gin.H{
		"success": true,
		"data":    v,
		"code":    0,
	}
}

func respError(e error) gin.H {
	var apiError err.APIError
	if !errors.As(e, &apiError) {
		apiError = err.ErrUnknown
	}

	return gin.H{
		"success": false,
		"code":    apiError.Code(),
		"message": apiError.Error(),
	}
}

func respErrorCode(code int, c *gin.Context) gin.H {
	lang := c.GetHeader("Lang")

	var msg string

	messages := strings.Split(err.ErrMap[code], ":")
	if len(messages) == 0 {
		msg = err.ErrMap[code]
	} else {
		if lang == err.LanguageCN {
			msg = messages[1]
		} else {
			msg = messages[0]
		}
	}

	return gin.H{
		"code": -1,
		"err":  code,
		"msg":  msg,
	}
}
