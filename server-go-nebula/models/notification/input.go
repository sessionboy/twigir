package models

// 新通知
type NewNotification struct {
	Id          int64  `json:"id"`
	IsRead      int64  `json:"isRead"`
	Sender      int64  `json:"sender"`
	Receiver    int64  `json:"receiver"`
	Action      int    `json:"action"` // 0:回复，1:喜欢，2:转帖，3:引用，4:关注，5:提及，6:公告
	Msg         string `json:"msg"`
	Target      int64  `json:"target"`
	TargetOwner int64  `json:"target_owner"`
	TargetType  int    `json:"target_type"` // 0:贴文，1:回复，2:私信
}
