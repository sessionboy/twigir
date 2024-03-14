package notify

import (
	"fmt"
	"net/http"
	"server/db"
	"server/db/dgraph"
	dql "server/dql/notification"
	res "server/shares/response"

	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
)

/*
  1，功能：通知列表
  2，GET: /notifications?t=reply
*/
func Notifications(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetString("user_id")
	first := c.DefaultQuery("first", "12")
	after := c.DefaultQuery("after", "0")
	t := c.Query("t")

	var q string
	switch {
	case t == "reply":
		q = dql.Notification_Replies
	case t == "favorite":
		q = dql.Notification_Favorites
	case t == "restatus":
		q = dql.Notification_Restatus
	case t == "quote":
		q = dql.Notification_Quotes
	case t == "mention":
		q = dql.Notification_Mentions
	case t == "user":
		q = dql.Notification_Users
	case t == "system":
		q = dql.Notification_Broadcast
	}

	if len(q) == 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	r, err := dgraph.QueryWithVars(ctx, txn, q, map[string]string{
		"$loggedUserid": loggedUserid,
		"$first":        first,
		"$after":        after,
	})
	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var list = dgraph.GetList(r, "notifications")
	if list == nil {
		list = make([]interface{}, 0)
	}

	c.JSON(http.StatusOK, res.Ok("", list))
}
