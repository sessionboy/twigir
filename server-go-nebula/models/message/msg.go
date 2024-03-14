package models

import (
	models "server/models/status"
	"time"
)

// 私信
type Message struct {
	Id        int64        `json:"id"`
	Sender    models.Owner `json:"sender"`
	Receiver  models.Owner `json:"receiver"`
	Msg       string       `json:"msg"`
	MsgType   int          `json:"msg_type"` // 0:文本、1:图片、2:视频
	MediaUrl  int          `json:"media_url"`
	CreatedAt time.Time    `json:"created_at"`
}
