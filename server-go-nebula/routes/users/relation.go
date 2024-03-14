package users

import (
	"net/http"
	"server/service"
	"server/shares"
	res "server/shares/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
)

/*
  1，功能：关注某人
  2，path: /users/{id}/follow
  3，user_id：被关注的用户的id
  4，logged_user：当前登录用户信息
*/
func Follow(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetInt64("user_id")
	userid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	// 检查是否已关注
	has, err := service.User.HasFollow(loggedUserid, userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	if has {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "err_followed")))
		return
	}

	var follow_err = service.User.Follow(loggedUserid, userid)
	if follow_err != nil {
		shares.SugarLogger.Errorf("error: %v follow  %v is fail", loggedUserid, userid)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_follow")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能：取关某人
  2，path: /users/{id}/unfollow
  3，user_id：被取关的用户的id
  4，logged_user：当前登录用户信息
*/
func Unfollow(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetInt64("user_id")
	userid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	var unFollow_err = service.User.UnFollow(loggedUserid, userid)
	if unFollow_err != nil {
		shares.SugarLogger.Errorf("error: %v unfollow  %v is fail", loggedUserid, userid)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_unfollow")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能：屏蔽某人
  2，path: /users/{id}/shield
  3，user_id：被屏蔽的用户的id
  4，logged_user：当前登录用户信息
*/
func Shield(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetInt64("user_id")
	userid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	// 检查是否已关注
	has, err := service.User.HasPullBlack(loggedUserid, userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	if has {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "err_shielded")))
		return
	}

	var _err = service.User.Shield(loggedUserid, userid)
	if _err != nil {
		shares.SugarLogger.Errorf("error: %v shield  %v is fail", loggedUserid, userid)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_shield")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}

/*
  1，功能：解除屏蔽某人
  2，path: /users/{id}/unshield
  3，user_id：被解除屏蔽的用户的id
  4，logged_user：当前登录用户信息
*/
func Unshield(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetInt64("user_id")
	userid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	var _err = service.User.UnShield(loggedUserid, userid)
	if _err != nil {
		shares.SugarLogger.Errorf("error: %v unshield  %v is fail", loggedUserid, userid)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_unshield")))
		return
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), nil))
}
