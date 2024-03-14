package models

// 新通知
type NotificationInput struct {
	Uid         string   `json:"uid"`
	Dtype       string   `json:"dgraph.type,omitempty"`
	IsRead      bool     `json:"isRead"`
	Sender      WithId   `json:"sender"`
	Recipient   []WithId `json:"recipient"`
	Action      int      `json:"action"` // 0:回复，1:喜欢，2:转帖，3:引用，4:提及，5:关注，6:广播
	Msg         string   `json:"msg"`
	Target      WithId   `json:"target"`
	TargetOwner WithId   `json:"target_owner"`
	TargetType  int      `json:"target_type"` // 0:贴文, 2:私信
}
type WithId struct {
	Uid string `json:"uid"`
}
