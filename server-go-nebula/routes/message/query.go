package message

import (
	"net/http"
	"server/service"
	res "server/shares/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
)

/*
  1，功能：私信列表
  2，GET: /messages
*/
func Messages(c *gin.Context) {
	loggedUserid := c.GetInt64("user_id")
	keyword := c.Query("keyword")
	limit := c.GetInt64("limit")
	offset := c.GetInt64("offset")

	result, _ := service.Message.GetMessages(loggedUserid, offset, limit, keyword)
	c.JSON(http.StatusOK, res.Ok("", result.Tables))
}

/*
  1，功能：私信详情，消息列表
  2，GET: /message/:id
*/
func Message(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetInt64("user_id")
	limit := c.GetInt64("limit")
	offset := c.GetInt64("offset")
	userid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	result, _ := service.Message.GetMessage(loggedUserid, userid, offset, limit)
	c.JSON(http.StatusOK, res.Ok("", result.Tables))
}
