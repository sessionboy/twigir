package service

import (
	"server/service/message"
	"server/service/status"
	"server/service/user"
)

var (
	User    *user.UserModel
	Status  *status.StatusModel
	Message *message.MessageModel
	CodeMap map[string]bool // 验证码(临时存储，后期要改为数据库存储)
)

func InitService() {
	User = &user.UserModel{}
	Status = &status.StatusModel{}
	Message = &message.MessageModel{}
}
