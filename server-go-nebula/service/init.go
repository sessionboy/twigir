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
)

func InitService() {
	User = &user.UserModel{}
	Status = &status.StatusModel{}
	Message = &message.MessageModel{}
}
