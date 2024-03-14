package notify

import (
	"net/http"
	"server/service"
	res "server/shares/response"

	"github.com/gin-gonic/gin"
)

/*
  1，功能：通知列表
  2，GET: /notifications
*/
func Notifications(c *gin.Context) {
	loggedUserid := c.GetInt64("user_id")
	keyword := c.Query("keyword")
	limit := c.GetInt64("limit")
	offset := c.GetInt64("offset")

	result, _ := service.Message.GetMessages(loggedUserid, offset, limit, keyword)
	c.JSON(http.StatusOK, res.Ok("", result.Tables))
}
