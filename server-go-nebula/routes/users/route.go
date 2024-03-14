package users

import (
	md "server/middleware"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.RouterGroup) {

	// 登录注册相关
	r.POST("/auth/login", Login)
	r.POST("/auth/register", Register)
	r.POST("/auth/send_phone", SendPhoneCode)
	r.POST("/auth/verify_phone", VerifyPhoneCode)
	r.POST("/auth/send_email", SendEmailCode)
	r.POST("/auth/verify_email", VerifyEmailCode)

	// 更新用户账号信息
	r.PUT("/account/profile", md.AuthRequired(), SetProfile)
	r.PUT("/account/name", md.AuthRequired(), SetName)
	r.PUT("/account/username", md.AuthRequired(), SetUsername)
	r.PUT("/account/phone", md.AuthRequired(), SetPhone)
	r.PUT("/account/email", md.AuthRequired(), SetEmail)
	r.PUT("/account/password", md.AuthRequired(), SetPassword)
	r.PUT("/account/avatar", md.AuthRequired(), SetAvatar)
	r.PUT("/account/cover", md.AuthRequired(), SetCover)
	r.PUT("/account/bio", md.AuthRequired(), SetBio)
	r.PUT("/account/setting", md.AuthRequired(), SetSetting)

	// 更新用户关系
	r.POST("/users/:id/follow", md.AuthRequired(), Follow)
	r.POST("/users/:id/unfollow", md.AuthRequired(), Unfollow)
	r.POST("/users/:id/shield", md.AuthRequired(), Shield)
	r.POST("/users/:id/unshield", md.AuthRequired(), Unshield)

	// 用户查询
	r.GET("/users/me", md.AuthRequired(), QueryMe)
	r.GET("/users/blacklists", md.Pagination(), QueryBlacklists)
	r.GET("/users/:id", md.AuthNotRequired(), QueryUser)
	r.GET("/users/:id/followings", md.Pagination(), QueryFollowings)
	r.GET("/users/:id/followers", md.Pagination(), QueryFollowers)
	r.GET("/users/:id/friends", md.Pagination(), QueryFriends)
	r.GET("/users/:id/status", QueryStatus)
	r.GET("/users/:id/status/medias", QueryStatusMedias)
	r.GET("/users/:id/status/favorites", QueryStatusFavorites)

}
