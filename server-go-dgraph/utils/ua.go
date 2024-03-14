package utils

import (
	models "server/models/users"

	"github.com/gin-gonic/gin"
	ua "github.com/mileusna/useragent"
)

func Parseua(c *gin.Context) models.UserAgent {
	ua_str := c.GetHeader("user-agent")
	ua := ua.Parse(ua_str)
	userAgent := models.UserAgent{
		Ip:       c.ClientIP(),
		Os:       ua.OS,
		Platform: ua.Name,
		Mobile:   ua.Mobile,
		App:      ua.OS,
	}
	return userAgent
}
