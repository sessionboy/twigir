package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func LangMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptLanguages := c.GetHeader("Accept-Language")
		languages, _, err := language.ParseAcceptLanguage(acceptLanguages)
		if err != nil || len(languages) == 0 {
			// 设置zh-CN为默认语言
			c.Set("lang", "zh-CN")
			c.Next()
			return
		}
		c.Set("lang", languages[0].String())
		c.Next()
	}
}
