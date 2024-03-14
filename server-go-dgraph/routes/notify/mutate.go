package notify

import (
	"context"
	"fmt"
	"server/db"
	"server/db/dgraph"
	models "server/models/notification"
	"server/ws"

	"github.com/gin-gonic/gin"
)

/*
  1，功能：创建通知
*/
func SendNotification(c *gin.Context, n models.NotificationInput) {
	n.Dtype = "Notification"
	n.Uid = "_:notification"
	n.IsRead = false

	ctx := context.Background()
	txn := db.Dgraph.NewTxn()
	defer txn.Discard(ctx)

	r, err := dgraph.Mutate(ctx, txn, n)
	if err != nil {
		fmt.Println("create notification error: ", err)
	}
	id := dgraph.GetUid(r, "notification")
	n.Uid = id
	n.Dtype = ""

	txn.Commit(ctx)
	txn.Discard(ctx)

	// 推送通知
	for i := 0; i < len(n.Recipient); i++ {
		recipient := n.Recipient[i].Uid
		if ws.Hub.OnLine(recipient) {
			notify := ws.Notification{
				Type:      2, // 通知
				Sender:    n.Sender.Uid,
				Recipient: []string{n.Recipient[i].Uid},
				Data:      n,
			}
			ws.Hub.Send(notify)
		} else {
			// 如果不在线则缓存离线消息
		}
	}
}
