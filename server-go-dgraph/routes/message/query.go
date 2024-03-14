package message

import (
	"fmt"
	"net/http"
	"server/db"
	"server/db/dgraph"
	dql "server/dql/message"
	res "server/shares/response"

	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
)

/*
  1，功能：私信对话列表
  2，GET: /conversations?s=keyword
*/
func Conversations(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetString("user_id")
	first := c.DefaultQuery("first", "12")
	after := c.DefaultQuery("after", "0")
	keyword := c.Query("s")
	isSearch := len(keyword) > 0

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	dqlstr := dql.GetConversationsDql(loggedUserid, first, after, keyword)
	vars := map[string]string{
		"$loggedUserid": loggedUserid,
		"$first":        first,
		"$after":        after,
	}
	if isSearch {
		vars["$keyword"] = keyword
	}
	r, err := dgraph.QueryWithVars(ctx, txn, dqlstr, vars)

	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var list interface{}
	if isSearch {
		list = dgraph.GetSubList(r, "users", "conversations")
	} else {
		list = dgraph.GetList(r, "conversations")
	}

	data := map[string]interface{}{
		"conversations": list,
	}

	c.JSON(http.StatusOK, res.Ok("", data))
}

/*
  1，功能：私信详情，对话内容
  2，GET: /conversation/:id
*/
func Message(c *gin.Context) {
	lang := c.GetString("lang")
	conversationid := c.Param("id")
	first := c.DefaultQuery("first", "12")
	after := c.DefaultQuery("after", "0")

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	r, err := dgraph.QueryWithVars(ctx, txn, dql.QueryMessages, map[string]string{
		"$conversationid": conversationid,
		"$first":          first,
		"$after":          after,
	})

	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	list := dgraph.GetSubList(r, "conversations", "messages")

	data := map[string]interface{}{
		"conversations": list,
	}

	c.JSON(http.StatusOK, res.Ok("", data))
}
