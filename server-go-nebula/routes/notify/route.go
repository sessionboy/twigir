package notify

import (
	md "server/middleware"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.RouterGroup) {

	// 私信列表
	r.GET("/notifications", md.AuthRequired(), Messages)
	// 发送私信
	r.POST("/notify", md.AuthRequired(), SendMessage)

}
