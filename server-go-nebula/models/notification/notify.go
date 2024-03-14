package models

import (
	"time"
)

// 私信
type Notification struct {
	Id          int64     `json:"id"`
	IsRead      int64     `json:"isRead"`
	Sender      int64     `json:"sender"`
	Msg         string    `json:"msg"`
	Action      int       `json:"action"`
	Target      int64     `json:"target"`
	TargetOwner int64     `json:"target_owner"`
	TargetType  int       `json:"target_type"`
	CreatedAt   time.Time `json:"created_at"`
}
