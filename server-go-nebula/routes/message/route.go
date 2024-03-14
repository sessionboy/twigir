package message

import (
	md "server/middleware"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.RouterGroup) {

	// 私信列表
	r.GET("/messages", md.AuthRequired(), Messages)
	// 发送私信
	r.POST("/message/:id", md.AuthRequired(), SendMessage)
	// 私信详情
	r.GET("/message/:id", md.AuthRequired(), Message)

}
