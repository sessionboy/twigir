package status

import (
	md "server/middleware"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.RouterGroup) {

	// 帖文/回复的增删改操作
	r.POST("/status", md.AuthRequired(), Status)
	r.POST("/status/:id/favorite", md.AuthRequired(), Favorite)
	r.POST("/status/:id/destroy", md.AuthRequired(), Destroy)

	// 热门推荐/视频/照片
	r.GET("/status/recommends", QueryRecommends)
	// 帖文详情
	r.GET("/status/:id", QueryStatus)
	// 回复列表
	r.GET("/status/:id/replies", QueryStatusReplies)

}
