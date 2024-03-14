package status

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/db"
	"server/db/dgraph"
	dql "server/dql/status"
	nmodels "server/models/notification"
	models "server/models/status"
	"server/routes/notify"
	"server/shares"
	res "server/shares/response"
	"server/utils"
	"strconv"

	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
)

/*
  1，功能：发布贴子
  2，post: /status
*/
func Status(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetString("user_id")
	user_verified := c.GetBool("user_verified")

	// 1，参数校验
	var status models.StatusInput
	if err := c.ShouldBind(&status); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	// 2，查询用户最近发帖信息
	r, err := dgraph.QueryWithVars(ctx, txn, dql.QueryStatusInfo, map[string]string{
		"$loggedUserid": loggedUserid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	items := r["user"].([]interface{})
	if len(items) == 0 {
		c.JSON(http.StatusNotFound, res.Err(i18n.Tr(lang, "account_not_found")))
		return
	}
	user := items[0].(map[string]interface{})
	last_publish_at := user["last_publish_at"]
	status_count := 0
	if user["status_count"] != nil {
		status_count = int(user["status_count"].(float64))
	}

	// 3，比较当前时间和最近一次发文时间是否小于60秒，是则报错，防止恶意刷帖
	if last_publish_at != nil {
		d, err := utils.CompareNowTimeWithSeconds(last_publish_at.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
			return
		}
		if d < 60 {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "err_frequently_pub")))
			return
		}
	}

	// 4，创建贴文
	publish_at := utils.GetUtcNowRFC3339()
	context := NewStatusContext{
		StatusType:   0,
		LoggedUserid: loggedUserid,
		Verified:     user_verified,
		PublishAt:    publish_at,
		StatusInput:  status,
	}
	pb := NewStatusJson(c, context)
	newstatus, err := json.Marshal(pb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_status_pub")))
		return
	}
	mu1 := &api.Mutation{
		SetJson: newstatus,
	}
	mu2 := &api.Mutation{
		SetNquads: []byte(fmt.Sprintf(`
			<%s> <last_publish_at> "%s" .
			<%s> <statuses> <_:status> .
			<%s> <status_count> "%v" .			
			`,
			// 更新用户最近发帖时间
			loggedUserid, publish_at,
			// 添加[用户->贴文]的statuses边
			loggedUserid,
			// 用户的贴文数量+1
			loggedUserid, status_count+1,
		)),
	}
	mus := []*api.Mutation{mu1, mu2}
	req := &api.Request{
		Mutations: mus,
		CommitNow: true,
	}
	resp, err := txn.Do(ctx, req)
	if err != nil {
		shares.SugarLogger.Errorf("%v publish error: %s", loggedUserid, err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_status_pub")))
		return
	}
	uids := resp.GetUids()
	id := uids["status"]
	if len(id) == 0 {
		shares.SugarLogger.Errorf("the id is not found after mutate return")
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_status_pub")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

	data := map[string]interface{}{
		"status": map[string]string{"id": id},
	}
	c.JSON(
		http.StatusOK, res.Ok(i18n.Tr(lang, "success"), data))
}

/*
  1，功能：发布引用贴文
  2，post: /status/:id/quote
*/
func Quote(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetString("user_id")
	user_verified := c.GetBool("user_verified")
	statusid := c.Param("id")
	_, err := strconv.ParseInt(statusid, 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "invalid_id")))
		return
	}

	// 1，参数校验
	var status models.StatusInput
	if err := c.ShouldBind(&status); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	// 2，查询用户最近发帖信息
	resp, err := txn.QueryWithVars(ctx, dql.QueryQuoteInfo, map[string]string{
		"$loggedUserid": loggedUserid,
		"$statusid":     statusid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var _quoteInfo dql.QuoteInfo
	err = json.Unmarshal(resp.Json, &_quoteInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	// 贴文不存在，或已删除
	if len(_quoteInfo.Status) == 0 {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "not_found_status")))
		return
	}

	logged_user := _quoteInfo.User[0]
	to_quote := _quoteInfo.Status[0]
	if len(to_quote.CreatedAt) == 0 {
		// 防止出现只有uid的空节点
		c.JSON(http.StatusNotFound, res.Err(i18n.Tr(lang, "not_found_status")))
		return
	}

	// 3，比较当前时间和最近一次发文时间是否小于60秒，是则报错，防止恶意刷帖
	if len(logged_user.LastPublishAt) > 0 {
		d, err := utils.CompareNowTimeWithSeconds(logged_user.LastPublishAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
			return
		}
		if d < 60 {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "err_frequently_pub")))
			return
		}
	}

	// 4，创建贴文
	publish_at := utils.GetUtcNowRFC3339()
	context := NewStatusContext{
		StatusType:   1,
		ToQuote:      statusid,
		LoggedUserid: loggedUserid,
		Verified:     user_verified,
		PublishAt:    publish_at,
		StatusInput:  status,
	}
	pb := NewStatusJson(c, context)
	newstatus, err := json.Marshal(pb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_status_pub")))
		return
	}
	mu1 := &api.Mutation{
		SetJson: newstatus,
	}
	mu2 := &api.Mutation{
		SetNquads: []byte(fmt.Sprintf(`
			<%s> <status_count> "%v" .
			<%s> <last_publish_at> "%s" .
			<%s> <quote_count> "%v" .			
			`,
			// 更新用户贴文数量
			loggedUserid, logged_user.StatusCount+1,
			// 更新用户最近发帖时间
			loggedUserid, publish_at,
			// 被引用的贴文的引用数量+1
			to_quote.Id, to_quote.QuoteCount+1,
		)),
	}
	mus := []*api.Mutation{mu1, mu2}
	req := &api.Request{
		Mutations: mus,
		CommitNow: true,
	}
	r, err := txn.Do(ctx, req)
	if err != nil {
		shares.SugarLogger.Errorf("%v publish error: %s", loggedUserid, err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_status_pub")))
		return
	}
	uids := r.GetUids()
	id := uids["status"]
	if len(id) == 0 {
		shares.SugarLogger.Errorf("the id is not found after mutate return")
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_status_pub")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

	// 发送通知
	if to_quote.User.Id != loggedUserid {
		n := nmodels.NotificationInput{
			Action:      3, // 引用
			Sender:      nmodels.WithId{Uid: loggedUserid},
			Recipient:   []nmodels.WithId{{Uid: to_quote.User.Id}},
			TargetType:  0,                                 // 贴文
			Target:      nmodels.WithId{Uid: id},           // 当前引用贴文
			TargetOwner: nmodels.WithId{Uid: loggedUserid}, // 当前引用贴文作者
		}
		go notify.SendNotification(c, n)
	}

	data := map[string]interface{}{
		"status": map[string]string{"id": id},
	}
	c.JSON(
		http.StatusOK, res.Ok(i18n.Tr(lang, "success"), data))
}

/*
  1，功能：转帖/撤销转帖
  2，POST: /status/:id/restatus
*/
func ReStatus(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetString("user_id")
	statusid := c.Param("id")
	_, err := strconv.ParseInt(statusid, 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "invalid_id")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	q := fmt.Sprintf(`
		query {
			status(func: uid("%s")){
				id:uid
				user{
					id:uid
				}
				restatus_count
				cnt as cnt:count(~restatuses @filter(uid("%s")))
				restatused: math(cnt == 1)			
			}
		}`, statusid, loggedUserid)
	r, err := txn.QueryWithVars(ctx, q, map[string]string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var _rstatus models.QueryReStatus
	err = json.Unmarshal(r.Json, &_rstatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	if len(_rstatus.Status) == 0 {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "not_found_status")))
		return
	}
	rstatus := _rstatus.Status[0]

	mus := []*api.Mutation{}
	if rstatus.Restatused {
		// 取消转帖
		mu1 := &api.Mutation{
			DelNquads: []byte(fmt.Sprintf(`
				<%s> <restatuses> <%s> .
			`, loggedUserid, statusid)),
		}
		mu2 := &api.Mutation{
			SetNquads: []byte(fmt.Sprintf(`
				<%s> <restatus_count> "%v" .
			`, statusid, rstatus.RestatusCount-1)),
		}
		mus = append(mus, mu1, mu2)
	} else {
		// 转帖
		mu1 := &api.Mutation{
			SetNquads: []byte(fmt.Sprintf(`
				<%s> <restatuses> <%s> .
				<%s> <restatus_count> "%v" .
			`, loggedUserid, statusid, statusid, rstatus.RestatusCount+1)),
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
		shares.SugarLogger.Errorf("%s restatus[%s] error: %s", loggedUserid, statusid, err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_restatus")))
		return
	}

	// 给贴文作者发送转帖通知
	if !rstatus.Restatused && rstatus.User.Id != loggedUserid {
		n := nmodels.NotificationInput{
			Action:      2, // 转帖
			Sender:      nmodels.WithId{Uid: loggedUserid},
			Recipient:   []nmodels.WithId{{Uid: rstatus.User.Id}},
			TargetType:  0, // 贴文
			Target:      nmodels.WithId{Uid: statusid},
			TargetOwner: nmodels.WithId{Uid: rstatus.User.Id},
		}
		go notify.SendNotification(c, n)
	}

	c.JSON(http.StatusOK, res.Ok("", nil))
}

/*
  1，功能：发布回复
  2，post: /status/:id/reply
	@规则和逻辑
		1，为防止恶意刷帖，增加60秒的发帖限制
		2，创建回复后，更新用户last_reply_at，以及to_reply的replies和reply_count
*/
func Reply(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetString("user_id")
	user_verified := c.GetBool("user_verified")
	statusid := c.Param("id")
	_, err := strconv.ParseInt(statusid, 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "invalid_id")))
		return
	}

	// 1，参数校验
	var status models.StatusInput
	if err := c.ShouldBind(&status); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	// 2，查询回复目标和登录用户最近回复等信息
	resp, err := txn.QueryWithVars(ctx, dql.GetReplyInfo, map[string]string{
		"$loggedUserid": loggedUserid,
		"$statusid":     statusid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var _replyInfo dql.ReplyInfo
	err = json.Unmarshal(resp.Json, &_replyInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	// 贴文不存在，或已删除
	if len(_replyInfo.Status) == 0 {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "not_found_status")))
		return
	}

	logged_user := _replyInfo.User[0]
	to_reply := _replyInfo.Status[0]
	if len(to_reply.CreatedAt) == 0 {
		// 防止出现只有uid的空节点
		c.JSON(http.StatusNotFound, res.Err(i18n.Tr(lang, "not_found_status")))
		return
	}

	// 如果目标贴文作者已注销账号，导致该贴文无效，则不能再回复
	if len(to_reply.User.Id) == 0 {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "reply_invalid_status")))
		return
	}

	// 3，比较当前时间和最近一次回复时间是否小于60秒，是则报错，防止恶意刷帖
	if len(logged_user.LastReplyAt) > 0 {
		d, err := utils.CompareNowTimeWithSeconds(logged_user.LastReplyAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
			return
		}
		if d < 60 {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "err_frequently_pub")))
			return
		}
	}

	// 4，创建贴文
	publish_at := utils.GetUtcNowRFC3339()
	context := NewStatusContext{
		StatusType:   2, // 默认是回复
		LoggedUserid: loggedUserid,
		Verified:     user_verified,
		PublishAt:    publish_at,
		ToUser:       to_reply.User.Id,
		ToStatus:     statusid, // 祖帖，默认是回复目标贴文的uid
		ToReply:      statusid, // 回复目标贴文
		StatusInput:  status,
	}
	if to_reply.StatusType == 2 {
		// 如果目标to_reply是回复(2)，则将当前status_type设为 3(回复的回复)
		// 同时将to_status设置为to_reply的to_status
		// 确保该回复和目标回复都属于同一个贴文的回复
		// 如果该贴文已删除，则to_status为空，允许这种情况存在
		context.StatusType = 3
		context.ToStatus = to_reply.ToStatus.Id
	}

	pb := NewStatusJson(c, context)
	newstatus, err := json.Marshal(pb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_status_pub")))
		return
	}
	mu1 := &api.Mutation{
		SetJson: newstatus,
	}
	mu2 := &api.Mutation{
		SetNquads: []byte(fmt.Sprintf(`
			<%s> <last_reply_at> "%s" .
			<%s> <reply_count> "%v" .			
			`,
			// 更新用户最近回复时间
			loggedUserid, publish_at,
			// 回复目标的回复数量+1
			statusid, to_reply.ReplyCount+1,
		)),
	}
	req := &api.Request{
		Mutations: []*api.Mutation{mu1, mu2},
		CommitNow: true,
	}
	r, err := txn.Do(ctx, req)
	if err != nil {
		shares.SugarLogger.Errorf("%v reply error: %s", loggedUserid, err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_status_reply")))
		return
	}
	uids := r.GetUids()
	id := uids["status"]
	fmt.Println("uids", uids, id)
	if len(id) == 0 {
		shares.SugarLogger.Errorf("the id is not found after mutate return")
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_status_pub")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

	// 给原文作者发送回复通知
	if to_reply.User.Id != loggedUserid {
		n := nmodels.NotificationInput{
			Action:      0, // 回复
			Sender:      nmodels.WithId{Uid: loggedUserid},
			Recipient:   []nmodels.WithId{{Uid: to_reply.User.Id}},
			TargetType:  0,                                 // 贴文
			Target:      nmodels.WithId{Uid: id},           // 当前回复贴文
			TargetOwner: nmodels.WithId{Uid: loggedUserid}, // 当前回复贴文作者
		}
		go notify.SendNotification(c, n)
	}

	data := map[string]interface{}{
		"status": map[string]string{"id": id},
	}
	c.JSON(http.StatusOK, res.Ok(i18n.Tr(lang, "success"), data))
}

/*
  1，功能：喜欢/撤销喜欢贴子
  2，post: /status/:id/favorite
*/
func Favorite(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetString("user_id")
	statusid := c.Param("id")
	_, err := strconv.ParseInt(statusid, 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "invalid_id")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	q := fmt.Sprintf(`
		query {
			status(func: uid("%s")){
				id:uid
				user{
					id:uid
				}
				favorite_count
				cnt as cnt:count(~favorites @filter(uid("%s")))
				favorited: math(cnt == 1)			
			}
		}`, statusid, loggedUserid)
	r, err := txn.QueryWithVars(ctx, q, map[string]string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var _status models.QueryFavoriteStatus
	err = json.Unmarshal(r.Json, &_status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	if len(_status.Status) == 0 {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "not_found_status")))
		return
	}
	status := _status.Status[0]

	mus := []*api.Mutation{}
	if status.Favorited {
		// 取消点赞
		mu1 := &api.Mutation{
			DelNquads: []byte(fmt.Sprintf(`
				<%s> <favorites> <%s> .
			`, loggedUserid, statusid)),
		}
		mu2 := &api.Mutation{
			SetNquads: []byte(fmt.Sprintf(`
				<%s> <favorite_count> "%v" .
			`, statusid, status.FavoriteCount-1)),
		}
		mus = append(mus, mu1, mu2)
	} else {
		// 点赞
		mu1 := &api.Mutation{
			SetNquads: []byte(fmt.Sprintf(`
				<%s> <favorites> <%s> .
				<%s> <favorite_count> "%v" .
			`, loggedUserid, statusid, statusid, status.FavoriteCount+1)),
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
		shares.SugarLogger.Errorf("%s favorite status[%s] error: %s", loggedUserid, statusid, err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_favorite")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

	// 给原文作者发送点赞通知
	if status.User.Id != loggedUserid {
		n := nmodels.NotificationInput{
			Action:      1, // 点赞
			Sender:      nmodels.WithId{Uid: loggedUserid},
			Recipient:   []nmodels.WithId{{Uid: status.User.Id}},
			TargetType:  0,                                   // 贴文
			Target:      nmodels.WithId{Uid: status.Id},      // 原文
			TargetOwner: nmodels.WithId{Uid: status.User.Id}, // 当原文作者
		}
		go notify.SendNotification(c, n)
	}

	c.JSON(http.StatusOK, res.Ok("", nil))
}

/*
  1，功能：删除贴子
  2，delete: /status/:id/delete
*/
func Delete(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetString("user_id")
	statusid := c.Param("id")
	_, err := strconv.ParseInt(statusid, 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "invalid_id")))
		return
	}

	ctx := c.Request.Context()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	resp, err := txn.QueryWithVars(ctx, dql.QueryDeleteStatusInfo, map[string]string{
		"$loggedUserid": loggedUserid,
		"$statusid":     statusid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var _statusInfo dql.DStatusInfo
	err = json.Unmarshal(resp.Json, &_statusInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}

	// 贴文不存在，或已删除
	if len(_statusInfo.Status) == 0 {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "not_found_status")))
		return
	}
	if len(_statusInfo.User) == 0 {
		// 账号不存在或已注销
		c.JSON(http.StatusNotFound, res.Err(i18n.Tr(lang, "account_not_found")))
		return
	}
	logged_user := _statusInfo.User[0]
	_status := _statusInfo.Status[0]
	if len(_status.CreatedAt) == 0 {
		// 防止出现只有uid的空节点
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "not_found_status")))
		return
	}

	query := fmt.Sprintf(`
		query {
			var(func: uid(%s)) {
				S as uid
				U as user
				TU as to_user
				TQ as to_quote
				TS as to_status
				TR as to_reply
				RS as restatuses
				F as favorites
			}
		}`, statusid)

	// 删除status
	mu1 := &api.Mutation{
		DelNquads: []byte(`
			uid(S) * * .
			uid(S) <user> uid(U) .
			uid(U) <statuses> uid(S) .
			uid(S) <to_user> uid(TU) .
			uid(S) <to_status> uid(TS) .
			uid(S) <to_quote> uid(TQ) .
			uid(S) <to_reply> uid(TR) .
			uid(RS) <restatuses> uid(S) .
			uid(F) <favorites> uid(S) .
		`),
	}
	mus := []*api.Mutation{mu1}

	if _status.StatusType == 0 {
		// 如果是删除贴文，则用户的 status_count-1
		mu2 := &api.Mutation{
			SetNquads: []byte(fmt.Sprintf(`<%s> <status_count> "%v" .`,
				loggedUserid, logged_user.StatusCount-1,
			)),
		}
		mus = append(mus, mu2)
	} else if _status.StatusType == 1 {
		// 如果是删除引用，则用户的 status_count-1，引用贴文的 quote_count-1
		mu2 := &api.Mutation{
			SetNquads: []byte(fmt.Sprintf(`
				<%s> <status_count> "%v" .
				<%s> <quote_count> "%v" .
			`, loggedUserid, logged_user.StatusCount-1,
				_status.ToQuote.Id, _status.ToQuote.QuoteCount-1,
			)),
		}
		mus = append(mus, mu2)
	} else if _status.StatusType == 2 || _status.StatusType == 3 {
		fmt.Println("_status.StatusType", _status.StatusType)
		// 如果是删除回复，回复贴文的 reply_count-1
		mu2 := &api.Mutation{
			SetNquads: []byte(fmt.Sprintf(`
				<%s> <reply_count> "%v" .
			`, _status.ToReply.Id, _status.ToReply.ReplyCount-1,
			)),
		}
		mus = append(mus, mu2)
	}

	req := &api.Request{
		Query:     query,
		Mutations: mus,
		CommitNow: true,
	}

	dresp, err := txn.Do(ctx, req)
	_ = dresp
	if err != nil {
		shares.SugarLogger.Errorf("%s delete status[%s] error: %s", loggedUserid, statusid, err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_delete_status")))
		return
	}

	// 提交事务
	txn.Commit(ctx)
	txn.Discard(ctx)

	c.JSON(http.StatusOK, res.Ok("", nil))
}
