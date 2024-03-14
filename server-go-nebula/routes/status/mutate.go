package status

import (
	"fmt"
	"net/http"
	models "server/models/status"
	"server/service"
	"server/shares"
	res "server/shares/response"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
)

/*
  1，功能：发布贴子
  2，post: /status
  3，body：贴文信息
*/
func Status(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetInt64("user_id")
	var status models.NewStatus
	if err := c.ShouldBind(&status); err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	// status_type=> 0：一般贴文；1：转帖；2：引用；3：回复；4：回复的回复
	switch status.StatusType {
	case 1:
		if status.ToStatus == 0 {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
			return
		}
	case 2:
		if status.ToStatus == 0 {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
			return
		}
	case 3:
		if status.ToStatus == 0 || status.ToUser == 0 {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
			return
		}
	case 4:
		if status.ToStatus == 0 || status.ToUser == 0 {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
			return
		}
	}

	// 媒体类型，0：photo，1：videos
	if status.MediaType == 0 && len(status.Photos) == 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "required_photo")))
		return
	}
	if status.MediaType == 1 && status.Video == 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "required_video")))
		return
	}

	agent := utils.Parseua(c)
	status.Ip = agent.Ip
	status.Device = agent.Os
	status.Platform = agent.Platform
	status.Id = utils.GenerateId()
	result, err := service.Status.CreateStatus(status, loggedUserid)
	_ = result
	if err != nil {
		shares.SugarLogger.Errorf("error: create status fail : %v", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_status_pub")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能：喜欢/点赞该贴子
  2，post: /status/:id/favorite
*/
func Favorite(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetInt64("user_id")
	statusid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	// 检查是否已点赞
	has, err := service.Status.HasFavorite(loggedUserid, statusid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	if has {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "err_favorited")))
		return
	}

	var favorite_err = service.Status.FavoriteStatus(loggedUserid, statusid)
	if favorite_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_favorite")))
		return
	}
	c.JSON(http.StatusOK, res.Ok("", nil))
}

/*
  1，功能：删除贴子
  2，delete: /status/:id
*/
func Destroy(c *gin.Context) {
	lang := c.GetString("lang")
	statusid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	var delete_err = service.Status.DeleteStatus(statusid)
	if delete_err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "delete_favorite")))
		return
	}
	c.JSON(http.StatusOK, res.Ok("", nil))
}
