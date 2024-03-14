package status

import (
	"net/http"
	"server/service"
	res "server/shares/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
)

/*
  1，功能：帖子详情
  2，GET: /status/:id
*/
func QueryStatus(c *gin.Context) {
	lang := c.GetString("lang")
	statusid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	result, _ := service.Status.GetStatus(statusid)
	if len(result.Tables) == 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "not_found_status")))
		return
	}
	c.JSON(http.StatusOK, res.Ok("", result.Tables[0]))
}

/*
  1，功能：帖子的回复列表
  2，path: /status/{id}/replies
  3，status_id：帖子的id
*/
func QueryStatusReplies(c *gin.Context) {
	lang := c.GetString("lang")
	limit := c.GetInt64("limit")
	offset := c.GetInt64("offset")
	statusid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	result, _ := service.Status.GetReplies(statusid, offset, limit)
	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), result.Tables))
}

/*
  1，功能：为你推荐/热门推荐
  2，GET: /status/recommends
*/
func QueryRecommends(c *gin.Context) {

	c.JSON(http.StatusOK, res.Ok("", nil))
}
