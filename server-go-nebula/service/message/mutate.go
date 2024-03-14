package message

import (
	"fmt"
	"server/conn"
	models "server/models/message"
)

type MessageModel struct{}

// 发送私信
func (u MessageModel) CreatedMessage(message models.NewMessage) (err error) {
	gql := fmt.Sprintf(`INSERT VERTEX message(
		msg,msg_type,media_url, created_at
		) VALUES %v:("%s",%v,"%s",datetime());`,
		message.Id, message.Msg, message.MsgType, message.MediaUrl,
	)
	res, err := conn.Execute(gql)
	_ = res
	return
}
