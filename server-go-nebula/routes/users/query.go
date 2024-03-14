package users

import (
	"fmt"
	"net/http"
	"server/service"
	res "server/shares/response"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
)

/*
  1，功能：我的个人主页信息
  2，path: /users/me
  3，logged_user：当前登录用户信息
*/

func QueryMe(c *gin.Context) {
	data := map[string]interface{}{
		"name": "jack",
		"age":  22,
	}
	lang := c.GetString("lang")
	println("lang", lang)
	println("loggedUserid", c.GetInt64("user_id"))
	println("username", c.GetString("user_username"))

	utils.Parseua(c)
	c.JSON(200, data)
}

/*
  1，功能：{id}的个人主页信息
  2，path: /users/{id}
  3，user_id：目标用户{id}
  3，logged_user：当前登录用户信息，判断是否已关注(假如已登录的话)
*/
func QueryUser(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetInt64("user_id")
	userid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	result, err := service.User.GetUserHomepage(userid, loggedUserid)
	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	if len(result.Tables) == 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "account_not_found")))
		return
	}
	var _user = result.Tables[0]
	// fmt.Println("_user", _user)
	// user := models.User{}
	// var decode_err error = mapstructure.Decode(_user, &user)
	// if decode_err != nil {
	// 	fmt.Println("decode_err", decode_err)
	// 	c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
	// 	return
	// }
	// fmt.Println("user", user)

	c.JSON(http.StatusOK, res.Ok("", _user))
}

/*
  1，功能：获取{id}的关注列表
  2，path: /users/{id}/followings
  3，user_id：目标用户的id
  4，query：分页参数first、after等
*/
func QueryFollowings(c *gin.Context) {
	lang := c.GetString("lang")
	limit := c.GetInt64("limit")
	offset := c.GetInt64("offset")
	userid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	result, err := service.User.GetFollowings(userid, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var list = result.Tables
	c.JSON(http.StatusOK, res.Ok("", list))
}

/*
  1，功能：获取{id}的粉丝列表
  2，path: /users/{id}/followers
  3，user_id：目标用户的id
  4，query：分页参数first、after等
*/
func QueryFollowers(c *gin.Context) {
	lang := c.GetString("lang")
	limit := c.GetInt64("limit")
	offset := c.GetInt64("offset")
	userid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	result, err := service.User.GetFollowers(userid, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var list = result.Tables
	c.JSON(http.StatusOK, res.Ok("", list))
}

/*
  1，功能：获取{id}的好友列表
  2，path: /users/{id}/friends
  3，user_id：目标用户的id
  4，query：分页参数first、after等
*/
func QueryFriends(c *gin.Context) {
	lang := c.GetString("lang")
	limit := c.GetInt64("limit")
	offset := c.GetInt64("offset")
	userid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	result, err := service.User.GetFriends(userid, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var list = result.Tables
	c.JSON(http.StatusOK, res.Ok("", list))
}

/*
  1，功能：获取{id}的屏蔽列表
  2，path: /users/{id}/shields
  3，user_id：目标用户的id
  4，query：分页参数first、after等
*/
func QueryBlacklists(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetInt64("user_id")
	limit := c.GetInt64("limit")
	offset := c.GetInt64("offset")
	result, err := service.User.GetBlacklists(loggedUserid, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var list = result.Tables
	c.JSON(http.StatusOK, res.Ok("", list))
}

/*
  1，功能：获取loggedUserid(我)与{id}的共同关注
  2，path: /users/{id}/same_followings
  3，query：分页参数first、after等
*/
func QuerySameFollowings(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetInt64("user_id")
	limit := c.GetInt64("limit")
	offset := c.GetInt64("offset")
	userid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	result, err := service.User.GetSameFollowings(loggedUserid, userid, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var list = result.Tables
	c.JSON(http.StatusOK, res.Ok("", list))
}

/*
  1，功能：我关注的A、B等人也关注了[id]
  2，path: /users/{id}/relation_followings
  3，query：分页参数first、after等
*/
func QueryRelationFollowings(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetInt64("user_id")
	limit := c.GetInt64("limit")
	offset := c.GetInt64("offset")
	targetid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	result, err := service.User.GetRelationFollowings(loggedUserid, targetid, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var list = result.Tables
	c.JSON(http.StatusOK, res.Ok("", list))
}

/*
  1，功能：用户主页帖子列表
  2，path: /users/{id}/status
  3，user_id：用户id
*/
func QueryStatus(c *gin.Context) {
	data := map[string]interface{}{
		"name": "jack",
		"age":  22,
	}
	c.JSON(200, data)
}

/*
  1，功能：用户{id}的照片/视频等媒体帖子列表
  2，path: /users/{id}/status
  3，user_id：用户id
  4，query: media_type(媒体类型)、first、after
*/
func QueryStatusMedias(c *gin.Context) {
	data := map[string]interface{}{
		"name": "jack",
		"age":  22,
	}
	c.JSON(200, data)
}

/*
  1，功能：用户{id}喜欢的帖子列表
  2，path: /users/{id}/favorite
  3，user_id：用户id
  4，query: media_type(媒体类型)、first、after
*/
func QueryStatusFavorites(c *gin.Context) {
	data := map[string]interface{}{
		"name": "jack",
		"age":  22,
	}
	c.JSON(200, data)
}
