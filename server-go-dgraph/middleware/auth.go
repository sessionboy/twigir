package middleware

import (
	"net/http"
	"server/shares"
	res "server/shares/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
)

// 拦截并提取登录用户信息
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.JSON(
				http.StatusUnauthorized,
				res.Err(i18n.Tr(c.GetString("lang"), "unauthorized")),
			)
			c.Abort()
			return
		}
		token := strings.Replace(authorization, "Bearer ", "", 1)
		claims, err := shares.ParseToken(token)
		if err != nil || len(claims.Id) == 0 {
			c.JSON(
				http.StatusUnauthorized,
				res.Err(i18n.Tr(c.GetString("lang"), "invalid_authorized")),
			)
			c.Abort()
			return
		}

		c.Set("user_id", claims.Id)
		c.Set("user_role", claims.Role)
		c.Set("user_verified", claims.Verified)
		c.Set("user_username", claims.Username)
		c.Next()
	}
}

// 不拦截，仅提取登录用户信息，
func AuthNotRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if len(authorization) > 0 {
			token := strings.Replace(authorization, "Bearer ", "", 1)
			claims, err := shares.ParseToken(token)
			if err == nil && len(claims.Id) > 0 {
				c.Set("user_id", claims.Id)
				c.Set("user_role", claims.Role)
				c.Set("user_verified", claims.Verified)
				c.Set("user_username", claims.Username)
			}
		}
		c.Next()
	}
}

// websocket用户认证
func AuthRequiredForWs() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Query("token")
		if authorization == "" {
			c.JSON(
				http.StatusUnauthorized,
				res.Err(i18n.Tr(c.GetString("lang"), "unauthorized")),
			)
			c.Abort()
			return
		}
		token := strings.Replace(authorization, "Bearer ", "", 1)
		claims, err := shares.ParseToken(token)
		if err != nil || len(claims.Id) == 0 {
			c.JSON(
				http.StatusUnauthorized,
				res.Err(i18n.Tr(c.GetString("lang"), "invalid_authorized")),
			)
			c.Abort()
			return
		}

		c.Set("user_id", claims.Id)
		c.Set("user_role", claims.Role)
		c.Set("user_verified", claims.Verified)
		c.Set("user_username", claims.Username)
		c.Next()
	}
}
