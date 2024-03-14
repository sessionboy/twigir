package models

type ConversationInput struct {
	Uid            string       `json:"uid"`
	Dtype          string       `json:"dgraph.type"`
	ConversationId string       `json:"conversation_id"`
	Creater        UserWithId   `json:"creater"`
	Users          []UserWithId `json:"users"`
	CreatedAt      string       `json:"created_at"`
}

// 新私信
type MessageInput struct {
	Uid            string     `json:"uid"`
	Dtype          string     `json:"dgraph.type,omitempty"`
	Serialid       string     `json:"serialid,omitempty"` // 消息编号，方便前端判断消息是否送达
	ConversationId string     `json:"conversation_id"`
	IsRead         bool       `json:"isRead"`
	Sender         UserWithId `json:"sender"`
	Recipient      UserWithId `json:"recipient"`
	Msg            string     `json:"msg"`
	MsgType        int        `json:"msg_type"` // 0:文本、1:图片、2:视频
	MediaUrl       string     `json:"media_url,omitempty"`
	CreatedAt      string     `json:"created_at"`
}

type UserWithId struct {
	Uid string `json:"uid"`
}
