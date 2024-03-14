package ws

// 私信消息
type ChatMessage struct {
	Userid   string `json:"userid"`   // 用户id
	Roomid   string `json:"roomid"`   // 房间id
	Serialid string `json:"serialid"` // 消息编号，用于错误返回，方便前端判断
	Type     string `json:"type"`     // 聊天类型，register:注册，single:单聊，broadcast:广播(群聊)
	MsgType  int    `json:"msg_type"` // 消息类型，0:文本、1:图片、2:视频
	MediaUrl string `json:"media_url"`
	Data     string `json:"data"` // 消息内容
	To       string `json:"to"`   // 发送给谁
}

// 通知
type Notification struct {
	Type         int         `json:"type"`                   // 响应类型， 1: 私信，2:通知，3: 系统通知
	Conversation string      `json:"conversation,omitempty"` // 消息发送者
	Sender       string      `json:"sender"`                 // 消息发送者
	Recipient    []string    `json:"recipient"`              // 消息接收方
	Appid        []string    `json:"appid,omitempty"`        // 如果有，则表示指定发送给该appid客户端
	Msg          string      `json:"msg,omitempty"`          // 附带消息
	Data         interface{} `json:"data"`                   // 响应内容
}
