package message

import (
	md "server/middleware"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.RouterGroup) {

	// 获取或创建对话, id为对话目标的userid
	r.GET("/conversation/:id", md.AuthRequired(), Conversation)
	// 私信对话列表
	r.GET("/conversations", md.AuthRequired(), Conversations)
	// 对话详情，消息列表
	r.GET("/messages/:id", md.AuthRequired(), Message)
	// 发送私信
	r.POST("/message/:id", md.AuthRequired(), PostMessage)

}
