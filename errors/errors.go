package errors

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Language string

const (
	LanguageEN = "en"
	LanguageCN = "cn"
)

const (
	NotFound = iota + 1000
	InvalidParams
	UserNotFound
	InvalidPassword
	InternalServer
	PermissionNotAllowed
	VerifyCodeExpired
	InvalidVerifyCode
	UnsupportedVerifyCodeType
	GetVCFrequently
	TimeoutCode
	InvalidReferralCode
	ReferralCodeBound

	PassWordNotAllowed
	UnauthorizedTwitter
	UnauthorizedDiscord
	UnauthorizedTelegram

	SocialMediaAccountIsAlreadyInUse

	MissionComplete
	MissionUnComplete
	NoImplement

	NotEnoughTagUsers

	InvalidSignature
	InvalidWalletAddress
	InvalidPublicKey
	WalletBound
	CaptchaError

	Unknown = -1
)

var ErrMap = map[int]string{
	Unknown:                          "unknown error:未知错误",
	NotFound:                         "not found:信息未找到",
	InternalServer:                   "Server Busy:服务器繁忙，请稍后再试",
	InvalidParams:                    "invalid params:参数有误",
	VerifyCodeExpired:                "verify code expired:验证码过期",
	InvalidVerifyCode:                "invalid verify code:无效的验证码",
	UnsupportedVerifyCodeType:        "unsupported verify code type:不支持的验证码类型",
	GetVCFrequently:                  "frequently request not allowed. please try again later.:请勿频繁获取验证码。请等待一段时间后再试。",
	TimeoutCode:                      "request timeout, please try again later: 请求超时, 请稍后再试",
	InvalidReferralCode:              "invalid referral code: 无效的邀请码",
	ReferralCodeBound:                "referral code bound: 已绑定邀请码",
	PassWordNotAllowed:               "password not allowed:密码错误",
	UnauthorizedTwitter:              "Unauthorized Twitter: 未授权 Twitter",
	UnauthorizedDiscord:              "Unauthorized Discord: 未授权 Discord",
	UnauthorizedTelegram:             "Unauthorized Telegram: 未授权 Telegram",
	SocialMediaAccountIsAlreadyInUse: "Binding Failed, The social media account has already been linked to another account: 绑定失败, 社交媒体账号已被其他账号绑定",
	MissionComplete:                  "mission completed: 任务已完成",
	MissionUnComplete:                "Please complete mission first: 请先完成任务",
	NoImplement:                      "No Implement: 正在开发中",
	NotEnoughTagUsers:                "Not Enough Tag Users: Tag 用户数量不满足要求",
	InvalidSignature:                 "invalid signature: 无效的签名",
	InvalidWalletAddress:             "invalid wallet address: 无效的钱包地址",
	InvalidPublicKey:                 "invalid public key: 无效的公钥地址",
	WalletBound:                      "wallet has been bound: 钱包已被绑定",
	CaptchaError:                     "Slide verification failed: 滑块校验失败",
}

var (
	ErrUnknown              = newError(Unknown, "Unknown Error")
	ErrNotFound             = newError(NotFound, "Record Not Found")
	ErrInvalidParams        = newError(InvalidParams, "Invalid Params")
	ErrUserNotFound         = newError(UserNotFound, "user not found")
	ErrInvalidPassword      = newError(InvalidPassword, "invalid password")
	ErrInternalServer       = newError(InternalServer, "Server Busy")
	ErrPermissionNotAllowed = newError(PermissionNotAllowed, "Permission Not Allowed")
)

type ApiError struct {
	code int
	err  error
}

func (e ApiError) Code() int {
	return e.code
}

func (e ApiError) Error() string {
	return e.err.Error()
}

func (e ApiError) APIError() (int, string) {
	return e.code, e.err.Error()
}

func newError(code int, message string) ApiError {
	return ApiError{code, errors.New(message)}
}

func New(message string) error {
	return errors.New(message)
}

type GenericError struct {
	Code int
	Err  error
}

func (e GenericError) Error() string {
	return e.Err.Error()
}

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
