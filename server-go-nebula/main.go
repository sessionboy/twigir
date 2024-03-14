package main

import (
	"server/config"
	"server/conn"
	"server/middleware"
	"server/routes"
	"server/service"
	"server/shares"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {

	// 初始化日志
	shares.InitLogger()
	defer shares.SugarLogger.Sync()

	// 创建gin应用
	r := gin.Default()
	store := cookie.NewStore([]byte(config.StoreSecret))
	store.Options(sessions.Options{
		MaxAge: config.LoginCountStoreTimeMaxAge,
	})
	r.Use(sessions.Sessions(config.LoginSession, store))

	// lang 国际化语言中间件 (尽量放在前面，因为其他中间件需要用到lang)
	r.Use(middleware.LangMiddleware())
	// 错误处理中间件
	r.Use(middleware.GinRecovery(false))

	// 初始化数据库连接
	conn.InitDb()
	service.InitService()

	// 路由
	v1 := r.Group("/api/v1")
	routes.Routes(v1)

	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run()
}
