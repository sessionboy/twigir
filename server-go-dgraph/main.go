package main

import (
	"net/http"
	"server/config"
	"server/db"

	// "server/db/dgraph"
	"server/middleware"
	"server/routes"
	"server/service"
	"server/shares"
	"server/ws"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func chatPage(c *gin.Context) {
	c.HTML(http.StatusOK, "chat.html", gin.H{
		"title": "we chat",
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 初始化dgraph schema
	// dgraph.InitSchema()

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
	// 初始化dgraph数据库
	db.InitDb()
	service.InitService()
	// websocket初始化
	ws.InitStart()

	// 私信/通知测试
	r.LoadHTMLGlob("chat.html")
	r.GET("/chat", chatPage)

	// 路由
	v1 := r.Group("/api/v1")
	v1.GET("/ws", middleware.AuthRequiredForWs(), ws.Handler)
	routes.Routes(v1)

	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run(":8080")
	// r.RunTLS(":8080", "./cert/server.pem", "./cert/server.key")
}
