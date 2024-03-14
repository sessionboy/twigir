package message

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"server/db"
	"server/db/dgraph"
	dql "server/dql/message"
	models "server/models/message"
	res "server/shares/response"
	"server/utils"
	"server/ws"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
)

/*
  1，功能：查找或创建对话
  2，GET: /conversation/:id
*/
func Conversation(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetString("user_id")
	loggedUserVerified := c.GetBool("user_verified")
	userid := c.Param("id")

	ctx := context.Background()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	// 1，检查对话是否已存在，存在则直接返回
	cid0 := loggedUserid + "-" + userid
	cid1 := userid + "-" + loggedUserid
	r, err := dgraph.QueryWithVars(ctx, txn, dql.HasConversation, map[string]string{
		"$cid0": cid0,
		"$cid1": cid1,
	})
	if err != nil {
		fmt.Println("err=>", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	cs := dgraph.GetItem(r, "conversation")
	if cs != nil {
		data := map[string]interface{}{
			"conversation": cs,
		}
		c.JSON(http.StatusOK, res.Ok("", data))
		return
	}

	// 2，检查是否已被对方禁止对话
	resp, err := txn.QueryWithVars(ctx, dql.ConversationAuth, map[string]string{
		"$userid":       userid,
		"$loggedUserid": loggedUserid,
	})
	if err != nil {
		fmt.Println("err2=>", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var conversation_user dql.ConversationUser
	conver_err := json.Unmarshal(resp.Json, &conversation_user)
	if conver_err != nil {
		fmt.Println("conver_err=>", conver_err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	if len(conversation_user.User) == 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "account_not_found")))
		return
	}
	user := conversation_user.User[0]
	// 私信权限过滤，
	// 如果对方关注了我，则无视其他规则，可对其发起私信
	if !user.Following {
		// 在对方的黑名单中，不能对其发起私信
		if user.Blacklist {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "chat_with_blacklist")))
			return
		}
		// 对方设置了他(她)未关注的人不能对其发起私信
		if user.ChatUnFollow {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "chat_with_following")))
			return
		}
		// 对方设置了未关注他(她)的人，不能对其发起私信
		if user.ChatUnFollowMe && user.Followme {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "chat_with_followMe")))
			return
		}
		// 对方设置了非认证用户不能对其发起私信
		if user.ChatUnVerified && !loggedUserVerified {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "chat_with_verified")))
			return
		}
	}

	// 3，允许对话，开始创建新的对话
	createdAt := utils.GetUtcNowRFC3339()
	m := map[string]interface{}{
		"uid":             "_:conversation",
		"dgraph.type":     "Conversation",
		"conversation_id": cid0,
		"creater":         map[string]string{"uid": loggedUserid},
		"users":           utils.UidsMap([]string{loggedUserid, userid}),
		"created_at":      createdAt,
		"last_publish_at": createdAt,
	}
	r, rerr := dgraph.Mutate(ctx, txn, m)
	if rerr != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	conversation := dgraph.GetUid(r, "conversation")

	data := map[string]interface{}{
		"conversation": map[string]string{
			"id":              conversation,
			"conversation_id": cid0,
		},
	}

	txn.Commit(ctx)
	txn.Discard(ctx)

	c.JSON(http.StatusOK, res.Ok("", data))
}

/*
  1，功能：发送私信
  2，POST: /message/:id
	@conversation_id  格式：userid1-userid2
*/
type BindInput struct {
	Msg string `json:"msg"`
}

func PostMessage(c *gin.Context) {
	lang := c.GetString("lang")
	loggedUserid := c.GetString("user_id")
	loggedUserVerified := c.GetBool("user_verified")
	conversationid := c.Param("id")

	var msg models.MessageInput
	if err := c.ShouldBind(&msg); err != nil {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}
	if len(msg.Msg) == 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
		return
	}

	ctx := context.Background()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	// 查询对话
	conversation_resp, err := txn.QueryWithVars(ctx, dql.ConversationWithId, map[string]string{
		"$conversationid": conversationid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var _conversation dql.Conversation
	conver_err := json.Unmarshal(conversation_resp.Json, &_conversation)
	if conver_err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	if len(_conversation.Conversation) == 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "conversation_not_found")))
		return
	}
	conversation := _conversation.Conversation[0]
	if len(conversation.ConversationId) == 0 {
		// 对话不存在或已被删除
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "conversation_not_found")))
		return
	}
	conversation_id := conversation.ConversationId
	if !strings.Contains(conversation_id, loggedUserid) {
		// 当前用户不在该对话中
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "conversation_not_match_user")))
		return
	}

	var userid string
	ids := strings.Split(conversation_id, "-")
	for i := 0; i < len(ids); i++ {
		if ids[i] != loggedUserid {
			userid = ids[i]
		}
	}
	if len(userid) == 0 {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}

	// 1，检查是否已被对方禁止对话
	resp, err := txn.QueryWithVars(ctx, dql.ConversationAuth, map[string]string{
		"$userid":       userid,
		"$loggedUserid": loggedUserid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	var _conversationuser dql.ConversationUser
	conver_user_err := json.Unmarshal(resp.Json, &_conversationuser)
	if conver_user_err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	if len(_conversationuser.User) == 0 {
		c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "account_not_found")))
		return
	}
	user := _conversationuser.User[0]

	// 私信权限过滤，
	// 如果对方关注了我，则无视其他规则，可对其发起私信
	if !user.Following {
		// 在对方的黑名单中，不能对其发起私信
		if user.Blacklist {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "chat_with_blacklist")))
			return
		}
		// 对方设置了他(她)未关注的人不能对其发起私信
		if user.ChatUnFollow {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "chat_with_following")))
			return
		}
		// 对方设置了未关注他(她)的人，不能对其发起私信
		if user.ChatUnFollowMe && user.Followme {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "chat_with_followMe")))
			return
		}
		// 对方设置了非认证用户不能对其发起私信
		if user.ChatUnVerified && !loggedUserVerified {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "chat_with_verified")))
			return
		}
	}

	// 2，允许对话，开始创建消息
	createdAt := utils.GetUtcNowRFC3339()
	msg.Uid = "_:message"
	msg.Dtype = "Message"
	msg.IsRead = false
	msg.Sender.Uid = loggedUserid
	msg.Recipient.Uid = userid
	msg.ConversationId = conversation_id
	msg.CreatedAt = createdAt
	var m struct {
		Uid           string              `json:"uid"`
		LastPublishAt string              `json:"last_publish_at"`
		Messages      models.MessageInput `json:"messages"`
	}
	m.Uid = conversationid
	m.LastPublishAt = createdAt
	m.Messages = msg
	r, err := dgraph.Mutate(ctx, txn, m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	id := dgraph.GetUid(r, "message")

	qres, err := dgraph.QueryWithVars(ctx, txn, dql.QueryMessage, map[string]string{
		"$messageid": id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "server_error")))
		return
	}
	message := dgraph.GetItem(qres, "message")

	txn.Commit(ctx)
	txn.Discard(ctx)

	// 推送通知
	notify := ws.Notification{
		Type:         1,
		Conversation: conversationid,
		Sender:       loggedUserid,
		Recipient:    []string{loggedUserid, userid},
		Data:         message,
	}
	go ws.Hub.Send(notify)
	// 离线通知

	// 客户端响应
	data := map[string]interface{}{
		"message": message,
	}
	c.JSON(http.StatusOK, res.Ok("", data))
}
