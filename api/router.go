package api

import (
	"bytes"
	"io"
	"strings"

	"titan-container-platform/config"

	"github.com/TestsLing/aj-captcha-go/service"
	"github.com/gin-gonic/gin"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("api")

// 行为校验初始化
var (
	factory *service.CaptchaServiceFactory
)

// func InitCaptcha() {
// 	// 水印配置
// 	clickWordConfig := &config2.ClickWordConfig{
// 		FontSize: 25,
// 		FontNum:  4,
// 	}
// 	// 点击文字配置
// 	watermarkConfig := &config2.WatermarkConfig{
// 		FontSize: 12,
// 		Color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
// 		Text:     "",
// 	}
// 	// 滑动模块配置
// 	blockPuzzleConfig := &config2.BlockPuzzleConfig{Offset: 200}
// 	configcap := config2.BuildConfig(constant.RedisCacheKey, config.Cfg.ResourcePath, watermarkConfig,
// 		clickWordConfig, blockPuzzleConfig, 2*60)
// 	factory = service.NewCaptchaServiceFactory(configcap)
// }

// ServerAPI initializes the server API with the provided configuration.
func ServerAPI(cfg *config.Config) {
	gin.SetMode(cfg.Mode)
	r := gin.Default()
	r.Use(cors())
	r.Use(RequestLoggerMiddleware())

	apiV1 := r.Group("/api/v1")
	authMiddleware, err := jwtGinMiddleware(cfg.SecretKey)
	if err != nil {
		log.Fatalf("jwt auth middleware: %v", err)
	}

	err = authMiddleware.MiddlewareInit()
	if err != nil {
		log.Fatalf("authMiddleware.MiddlewareInit: %v", err)
	}

	user := apiV1.Group("/user")
	user.GET("/login_before", getNonceStringHandler)
	user.POST("/login", authMiddleware.LoginHandler)
	user.POST("/logout", authMiddleware.LogoutHandler)
	user.GET("/refresh_token", authMiddleware.RefreshHandler)

	user.Use(authMiddleware.MiddlewareFunc())
	user.GET("/info", getUserInfoHandler)
	user.POST("/faucet", getTokenHandler)
	user.GET("/balance", getBalanceHandler)

	order := apiV1.Group("/order")
	order.Use(authMiddleware.MiddlewareFunc())
	order.GET("/price", getPriceHandler)
	order.POST("/create", createOrderHandler)
	order.GET("/history", getOrderHistoryHandler)

	if err := r.Run(cfg.Listen); err != nil {
		log.Fatalf("starting server: %v\n", err)
	}
}

// RequestLoggerMiddleware logs requests to the server.
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "storage") {
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ := io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)
			if string(body) != "" {
				log.Debug(string(body))
			}
		}
		// log.Debug(c.Request.Header)
		c.Next()
	}
}
