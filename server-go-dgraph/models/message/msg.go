package models

import (
	models "server/models/status"
	"time"
)

// 私信
type Message struct {
	Id        string       `json:"id"`
	Sender    models.Owner `json:"sender"`
	Recipient UserWithId   `json:"recipient"`
	Msg       string       `json:"msg"`
	MsgType   int          `json:"msg_type"` // 0:文本、1:图片、2:视频
	MediaUrl  int          `json:"media_url"`
	CreatedAt time.Time    `json:"created_at"`
}
