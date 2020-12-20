package main

import (
	"github.com/daodao97/egin"
	"github.com/daodao97/egin/middleware"
	"github.com/gin-gonic/gin"

	"oms/config"
	"oms/service/user"
)

// config/*** 存放启动时配置, 如路由, 中间件
// app.json 存放运行时配置, 如 db/redis 配置, 由 egin/config 持有
//go:generate /Users/mac/go/src/github.com/daodao97/egin-tools/egin-tools -route
//go:generate goimports -w config
func main() {
	boot := egin.Bootstrap{
		HttpMiddlewares: []func() gin.HandlerFunc{
			middleware.Cors,
			func() gin.HandlerFunc {
				return middleware.JWTMiddleware(user.New())
			},
			middleware.IPAuth,
			middleware.IpLimiter,
			middleware.HttpLog,
			middleware.Prometheus,
		},
		RegRoutes: config.RegRouter,
	}
	boot.Start()
}
