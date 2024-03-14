package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/db"
	nmodels "server/models/notification"
	models "server/models/users"
	"server/routes/notify"
	"server/shares"
	res "server/shares/response"
	"strconv"

	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
)

/*
  1，功能：关注/取消关注某人
  2，path: /users/:id/follow (id为雪花id)
*/
func Follow(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetString("user_id")
	userid := c.Param("id")
	_, err := strconv.ParseInt(userid, 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "invalid_id")))
		return
	}

	// 不能关注自己
	if loggedUserid == userid {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_follow_self")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	q := fmt.Sprintf(`
		query {
			user(func: uid("%s")) {
				id:uid
				followers_count
			}
			loggedUser(func: uid("%s")) {
				id:uid
				followings_count
				cnt as cnt:count(follows @filter(uid("%s")))
				following: math(cnt == 1)			
			}			
		}`, userid, loggedUserid, userid)
	r, err := txn.QueryWithVars(ctx, q, map[string]string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var _user models.QueryUser
	err = json.Unmarshal(r.Json, &_user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	if len(_user.User) == 0 {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "not_found_status")))
		return
	}
	user := _user.User[0]
	loggedUser := _user.LoggedUser[0]

	mus := []*api.Mutation{}
	if loggedUser.Following {
		// 取消关注
		mu1 := &api.Mutation{
			DelNquads: []byte(fmt.Sprintf(`
				<%s> <follows> <%s> .
			`, loggedUserid, userid)),
		}
		mu2 := &api.Mutation{
			SetNquads: []byte(fmt.Sprintf(`
				<%s> <followers_count> "%v" .
				<%s> <followings_count> "%v" .
			`,
				userid, user.FollowersCount-1,
				loggedUserid, loggedUser.FollowingsCount-1,
			)),
		}
		mus = append(mus, mu1, mu2)
	} else {
		// 关注
		mu1 := &api.Mutation{
			SetNquads: []byte(fmt.Sprintf(`
				<%s> <follows> <%s> .
				<%s> <followers_count> "%v" .
				<%s> <followings_count> "%v" .
			`,
				loggedUserid, userid,
				userid, user.FollowersCount+1,
				loggedUserid, loggedUser.FollowingsCount+1,
			)),
		}
		mus = append(mus, mu1)
	}

	req := &api.Request{
		Mutations: mus,
		CommitNow: true,
	}
	resp, err := txn.Do(ctx, req)
	_ = resp
	if err != nil {
		shares.SugarLogger.Errorf("%v follow  %v is fail: %s", loggedUserid, userid, err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_follow")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

	if !loggedUser.Following {
		n := nmodels.NotificationInput{
			Action:    5, // 关注
			Sender:    nmodels.WithId{Uid: loggedUserid},
			Recipient: []nmodels.WithId{{Uid: userid}},
		}
		go notify.SendNotification(c, n)
	}

	data := map[string]interface{}{
		"id":        userid,
		"following": !loggedUser.Following,
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), data))
}

/*
  1，功能：屏蔽某人，加入黑名单
  2，path: /users/:id/blacklist
*/
func Blacklist(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetString("user_id")
	userid := c.Param("id")
	_, err := strconv.ParseInt(userid, 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "invalid_id")))
		return
	}

	// 不能屏蔽自己
	if loggedUserid == userid {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "fail_shield_self")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	q := fmt.Sprintf(`
		query {
			user(func: uid("%s")) {
				id:uid
				cnt as cnt:count(blacklists @filter(uid("%s")))
				blacklisted: math(cnt == 1)			
			}			
		}`, loggedUserid, userid)
	r, err := txn.QueryWithVars(ctx, q, map[string]string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var _user models.QueryBlacklistUser
	err = json.Unmarshal(r.Json, &_user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	if len(_user.User) == 0 {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "not_found_status")))
		return
	}
	user := _user.User[0]

	mus := []*api.Mutation{}
	if user.Blacklisted {
		// 取消黑名单
		mu1 := &api.Mutation{
			DelNquads: []byte(fmt.Sprintf(`
				<%s> <blacklists> <%s> .
			`, loggedUserid, userid)),
		}
		mus = append(mus, mu1)
	} else {
		// 加入黑名单
		mu1 := &api.Mutation{
			SetNquads: []byte(fmt.Sprintf(`
				<%s> <blacklists> <%s> .
			`,
				loggedUserid, userid,
			)),
		}
		mus = append(mus, mu1)
	}

	req := &api.Request{
		Mutations: mus,
		CommitNow: true,
	}
	resp, err := txn.Do(ctx, req)
	_ = resp
	if err != nil {
		shares.SugarLogger.Errorf("error: %v shield  %v is fail: %s", loggedUserid, userid, err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_shield")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

	data := map[string]interface{}{
		"id":          userid,
		"blacklisted": !user.Blacklisted,
	}

	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), data))
}
