package models

// 新私信
type NewMessage struct {
	Id       int64  `json:"id"`
	Sender   int64  `json:"sender"`
	Receiver int64  `json:"receiver"`
	Msg      string `json:"msg"`
	MsgType  int    `json:"msg_type"` // 0:文本、1:图片、2:视频
	MediaUrl string `json:"media_url"`
}
