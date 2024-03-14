package status

import (
	"fmt"
	"net/http"
	"server/db"
	"server/db/dgraph"
	dql "server/dql/status"
	res "server/shares/response"
	"server/utils"
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
	loggedUserid := utils.GetDefaultString(c, "user_id", "0")
	statusid := c.Param("id")
	_, err := strconv.ParseInt(statusid, 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "invalid_id")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	r, err := dgraph.QueryWithVars(ctx, txn, dql.QueryStatus, map[string]string{
		"$statusid":     statusid,
		"$loggedUserid": loggedUserid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	items := r["status"].([]interface{})

	if len(items) == 0 {
		c.JSON(http.StatusNotFound, res.Err(i18n.Tr(lang, "not_found_status")))
		return
	}

	status := items[0].(map[string]interface{})

	if status == nil {
		c.JSON(http.StatusNotFound, res.Err(i18n.Tr(lang, "not_found_status")))
		return
	}

	id := status["id"]
	if id == nil {
		c.JSON(http.StatusNotFound, res.Err(i18n.Tr(lang, "not_found_status")))
		return
	}

	delete(status, "fcnt")
	delete(status, "rcnt")

	c.JSON(http.StatusOK, res.Ok("", status))
}

/*
  1，功能：帖子的回复列表
  2，path: /status/{id}/replies
  3，status_id：帖子的id
*/
func QueryStatusReplies(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := utils.GetDefaultString(c, "user_id", "0")
	first := c.DefaultQuery("first", "12")
	after := c.DefaultQuery("after", "0")
	statusid := c.Param("id")
	_, err := strconv.ParseInt(statusid, 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "invalid_id")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	r, err := dgraph.QueryWithVars(ctx, txn, dql.QueryStatusReplies, map[string]string{
		"$statusid":     statusid,
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
	items := r["status"].([]interface{})
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
  1，功能：为你推荐/热门推荐
  2，GET: /status/recommends?t=image
	3，@t media_type类型
*/
func QueryRecommends(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := utils.GetDefaultString(c, "user_id", "0")
	first := c.DefaultQuery("first", "12")
	after := c.DefaultQuery("after", "0")

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	q := dql.QueryRecommendStatuses
	vars := map[string]string{
		"$loggedUserid": loggedUserid,
		"$first":        first,
		"$after":        after,
	}
	t := c.Query("t")
	if len(t) > 0 && (t == "image" || t == "video") {
		q = dql.QueryRecommendMediaStatuses
		if t == "image" {
			vars["$media_type"] = "1"
		}
		if t == "video" {
			vars["$media_type"] = "2"
		}
	}
	r, err := dgraph.QueryWithVars(ctx, txn, q, vars)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var list []interface{} = r["statuses"].([]interface{})
	if len(list) == 0 {
		list = make([]interface{}, 0)
	}

	c.JSON(http.StatusOK, res.Ok("", list))
}

/*
  1，功能：全局搜索
  2，GET: /search?t=all&q=react
	3，@t 搜索类型: all:所有，latest:最新，video:视频，image:图片
	4，@n sort new，按最新时间排序的方式
*/
func Search(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := utils.GetDefaultString(c, "user_id", "0")
	first := c.DefaultQuery("first", "12")
	after := c.DefaultQuery("after", "0")
	t := c.DefaultQuery("t", "all")
	keyword := c.Query("q")

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	q := dql.SearchQuery(t)
	vars := map[string]string{
		"$keyword":      keyword,
		"$loggedUserid": loggedUserid,
		"$first":        first,
		"$after":        after,
	}
	r, err := dgraph.QueryWithVars(ctx, txn, q, vars)
	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	c.JSON(http.StatusOK, res.Ok("", r))
}
