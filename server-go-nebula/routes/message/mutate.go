package message

import (
	"net/http"
	models "server/models/message"
	"server/service"
	res "server/shares/response"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
)

/*
  1，功能：发送私信
  2，POST: /message/:id
*/
func SendMessage(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetInt64("user_id")
	userid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	var message models.NewMessage
	if err := c.ShouldBind(&message); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	message.Id = utils.GenerateId()
	message.Sender = loggedUserid
	message.Receiver = userid
	var send_err = service.Message.CreatedMessage(message)
	if send_err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_sendMessage")))
		return
	}

	c.JSON(http.StatusOK, res.Ok("", nil))
}
