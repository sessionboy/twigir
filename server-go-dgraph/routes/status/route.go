package status

import (
	md "server/middleware"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.RouterGroup) {
	r.GET("/search", md.AuthNotRequired(), Search)                         // 全局搜索
	r.POST("/status", md.AuthRequired(), Status)                           // 发帖
	r.GET("/status/recommends", md.AuthNotRequired(), QueryRecommends)     // 热门推荐/视频/图片
	r.GET("/status/:id", md.AuthNotRequired(), QueryStatus)                // 帖文详情
	r.GET("/status/:id/replies", md.AuthNotRequired(), QueryStatusReplies) // 帖文回复列表
	r.POST("/status/:id/quote", md.AuthRequired(), Quote)                  // 引用发帖
	r.POST("/status/:id/reply", md.AuthRequired(), Reply)                  // 回复
	r.POST("/status/:id/restatus", md.AuthRequired(), ReStatus)            // 转帖
	r.POST("/status/:id/favorite", md.AuthRequired(), Favorite)            // 点赞贴子
	r.POST("/status/:id/delete", md.AuthRequired(), Delete)                // 删除帖子

}
