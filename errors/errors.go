package errors

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Language represents a programming language.
type Language string

const (
	// LanguageEN represents the English language.
	LanguageEN = "en"
	// LanguageCN represents the Chinese language.
	LanguageCN = "cn"
)

// NotFound indicates that a requested resource was not found.
const (
	Success = 0

	NotFound = iota + 1000
	InvalidParams
	InternalServer
	VerifyCodeExpired
	InvalidSignature
	InvalidWalletAddress
	InvalidPublicKey
	QuotaIssued
	Received

	Unknown = -1
)

// ErrMap maps error codes to their corresponding error messages.
var ErrMap = map[int]string{
	Unknown:              "unknown error:未知错误",
	NotFound:             "not found:信息未找到",
	InternalServer:       "Server Busy:服务器繁忙，请稍后再试",
	InvalidParams:        "invalid params:参数有误",
	VerifyCodeExpired:    "verify code expired:验证码过期",
	InvalidSignature:     "invalid signature: 无效的签名",
	InvalidWalletAddress: "invalid wallet address: 无效的钱包地址",
	InvalidPublicKey:     "invalid public key: 无效的公钥地址",
	QuotaIssued:          "the quota has been issued: 额度已发完",
	Received:             "received: 已领取",
}

// ErrUnknown represents an unknown error.
var (
	ErrUnknown        = newError(Unknown, "Unknown Error")
	ErrNotFound       = newError(NotFound, "Record Not Found")
	ErrInvalidParams  = newError(InvalidParams, "Invalid Params")
	ErrInternalServer = newError(InternalServer, "Server Busy")
)

// APIError represents an error with a specific code and underlying error.
type APIError struct {
	code int
	err  error
}

// Code returns the error code associated with the APIError.
func (e APIError) Code() int {
	return e.code
}

func (e APIError) Error() string {
	return e.err.Error()
}

// APIError returns the error code and message for the APIError.
func (e APIError) APIError() (int, string) {
	return e.code, e.err.Error()
}

func newError(code int, message string) APIError {
	return APIError{code, errors.New(message)}
}

// New creates a new error with the provided message.
func New(message string) error {
	return errors.New(message)
}

// GenericError represents an error with a code and an underlying error.
type GenericError struct {
	Code int
	Err  error
}

func (e GenericError) Error() string {
	return e.Err.Error()
}

// NewErrorCode creates a new GenericError based on the provided code and context.
func NewErrorCode(Code int, c *gin.Context) GenericError {
	l := c.GetHeader("Lang")
	errSplit := strings.Split(ErrMap[Code], ":")
	var e string
	switch l {
	case LanguageCN:
		e = errSplit[1]
	default:
		e = errSplit[0]
	}
	return GenericError{Code: Code, Err: errors.New(e)}
}
