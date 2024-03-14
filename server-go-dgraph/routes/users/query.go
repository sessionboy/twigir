package users

import (
	"fmt"
	"net/http"
	"server/db"
	"server/db/dgraph"
	dql "server/dql/user"
	res "server/shares/response"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
)

/*
  1，功能：{username}的个人主页信息
  2，path: /users/:username
*/
func QueryUser(c *gin.Context) {
	lang := c.GetString("lang")
	// 如果loggedUserid不存在则设置为"0"，防止查询出错
	loggedUserid := utils.GetDefaultString(c, "user_id", "0")
	username := c.Param("username")

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	r, err := dgraph.QueryWithVars(ctx, txn, dql.GetUserHomepage, map[string]string{
		"$username":     username,
		"$loggedUserid": loggedUserid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	items := r["item"].([]interface{})

	if len(items) == 0 {
		c.JSON(http.StatusNotFound, res.Err(i18n.Tr(lang, "account_not_found")))
		return
	}

	c.JSON(http.StatusOK, res.Ok("", items[0]))
}

/*
  1，功能：获取:username的关注列表
  2，GET: /users/:username/followings?first=12&after=id
*/
func QueryFollowings(c *gin.Context) {
	lang := c.GetString("lang")
	username := c.Param("username")
	loggedUserid := utils.GetDefaultString(c, "user_id", "0")
	first := c.DefaultQuery("first", "15")
	after := c.DefaultQuery("after", "0")

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	r, err := dgraph.QueryWithVars(ctx, txn, dql.Followings, map[string]string{
		"$username":     username,
		"$loggedUserid": loggedUserid,
		"$first":        first,
		"$after":        after,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}

	var list interface{}
	items := r["item"].([]interface{})
	if len(items) == 0 {
		list = make([]interface{}, 0)
	} else {
		item := items[0].(map[string]interface{})
		if item["edges"] == nil {
			list = make([]interface{}, 0)
		} else {
			list = item["edges"]
		}
	}

	c.JSON(http.StatusOK, res.Ok("", list))
}

/*
  1，功能：获取:username的粉丝列表
	2，GET: /users/:username/followers?first=12&after=id
*/
func QueryFollowers(c *gin.Context) {
	lang := c.GetString("lang")
	username := c.Param("username")
	loggedUserid := utils.GetDefaultString(c, "user_id", "0")
	first := c.DefaultQuery("first", "15")
	after := c.DefaultQuery("after", "0")

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	r, err := dgraph.QueryWithVars(ctx, txn, dql.Followers, map[string]string{
		"$username":     username,
		"$loggedUserid": loggedUserid,
		"$first":        first,
		"$after":        after,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}

	var list interface{}
	items := r["item"].([]interface{})
	if len(items) == 0 {
		list = make([]interface{}, 0)
	} else {
		item := items[0].(map[string]interface{})
		if item["edges"] == nil {
			list = make([]interface{}, 0)
		} else {
			list = item["edges"]
		}
	}

	c.JSON(http.StatusOK, res.Ok("", list))
}

/*
  1，功能：获取:username的好友列表
	2，GET: /users/:username/friends?first=12&after=uid
*/
func QueryFriends(c *gin.Context) {
	lang := c.GetString("lang")
	username := c.Param("username")
	loggedUserid := utils.GetDefaultString(c, "user_id", "0")
	first := c.DefaultQuery("first", "15")
	after := c.DefaultQuery("after", "0")

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	r, err := dgraph.QueryWithVars(ctx, txn, dql.Friends, map[string]string{
		"$username":     username,
		"$loggedUserid": loggedUserid,
		"$first":        first,
		"$after":        after,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}

	var list interface{}
	items := r["item"].([]interface{})
	if len(items) == 0 {
		list = make([]interface{}, 0)
	} else {
		item := items[0].(map[string]interface{})
		if item["edges"] == nil {
			list = make([]interface{}, 0)
		} else {
			list = item["edges"]
		}
	}

	c.JSON(http.StatusOK, res.Ok("", list))
}

/*
  1，功能：获取登录用户的黑名单列表
  2，path: /users/blacklists?first=12&after=uid
*/
func QueryBlacklists(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetString("user_id")
	first := c.DefaultQuery("first", "15")
	after := c.DefaultQuery("after", "0")

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	fmt.Println(loggedUserid, first, after)
	r, err := dgraph.QueryWithVars(ctx, txn, dql.Blacklists, map[string]string{
		"$loggedUserid": loggedUserid,
		"$first":        first,
		"$after":        after,
	})
	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}

	var list interface{}
	items := r["item"].([]interface{})
	if len(items) == 0 {
		list = make([]interface{}, 0)
	} else {
		item := items[0].(map[string]interface{})
		if item["edges"] == nil {
			list = make([]interface{}, 0)
		} else {
			list = item["edges"]
		}
	}

	c.JSON(http.StatusOK, res.Ok("", list))
}

/*
  1，功能：获取loggedUserid(我)与:username的共同关注
  2，GET: /users/:username/same_followings?first=12&after=uid
*/
func QuerySameFollowings(c *gin.Context) {
	lang := c.GetString("lang")
	username := c.Param("username")
	loggedUserid := utils.GetDefaultString(c, "user_id", "0")
	first := c.DefaultQuery("first", "15")
	after := c.DefaultQuery("after", "0")

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	r, err := dgraph.QueryWithVars(ctx, txn, dql.SameFollowings, map[string]string{
		"$username":     username,
		"$loggedUserid": loggedUserid,
		"$first":        first,
		"$after":        after,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}

	var list interface{}
	items := r["item"].([]interface{})
	if len(items) == 0 {
		list = make([]interface{}, 0)
	} else {
		item := items[0].(map[string]interface{})
		if item["edges"] == nil {
			list = make([]interface{}, 0)
		} else {
			list = item["edges"]
		}
	}

	c.JSON(http.StatusOK, res.Ok("", list))
}

/*
  1，功能：我关注的A、B等人也关注了:username
	2，GET: /users/:username/relation_followings?first=12&after=uid
*/
func QueryRelationFollowings(c *gin.Context) {
	lang := c.GetString("lang")
	username := c.Param("username")
	loggedUserid := utils.GetDefaultString(c, "user_id", "0")
	first := c.DefaultQuery("first", "15")
	after := c.DefaultQuery("after", "0")

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	r, err := dgraph.QueryWithVars(ctx, txn, dql.RelationFollowings, map[string]string{
		"$username":     username,
		"$loggedUserid": loggedUserid,
		"$first":        first,
		"$after":        after,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}

	var list interface{}
	items := r["item"].([]interface{})
	if len(items) == 0 {
		list = make([]interface{}, 0)
	} else {
		item := items[0].(map[string]interface{})
		if item["edges"] == nil {
			list = make([]interface{}, 0)
		} else {
			list = item["edges"]
		}
	}

	c.JSON(http.StatusOK, res.Ok("", list))
}

/*
  1，功能：用户时间线
  2，path: /users/timeline?first=12&after=0
*/
func QueryTimeline(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := utils.GetDefaultString(c, "user_id", "0")
	first := c.DefaultQuery("first", "15")
	after := c.DefaultQuery("after", "0")

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	r, err := dgraph.QueryWithVars(ctx, txn, dql.QueryTimeline, map[string]string{
		"$loggedUserid": loggedUserid,
		"$first":        first,
		"$after":        after,
	})
	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	c.JSON(http.StatusOK, res.Ok("", r))
}

/*
  1，功能：用户贴文列表[个人主页]
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
  1，功能：用户主页图片列表
  2，path: /users/{id}/images
*/
func QueryUserImages(c *gin.Context) {
	data := map[string]interface{}{
		"name": "jack",
		"age":  22,
	}
	c.JSON(200, data)
}

/*
  1，功能：用户{id}的图片/视频等媒体帖子列表
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
