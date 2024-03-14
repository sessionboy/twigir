package notify

import (
	md "server/middleware"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.RouterGroup) {

	// 通知列表
	r.GET("/notifications", md.AuthRequired(), Notifications)

}
