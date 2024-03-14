package dgraph

/*
  @sender User 发送者
  @recipient []User 接收者
  @action Int 0:回复，1:喜欢，2:转帖，3:引用，4:提及，5:关注，6:广播
  @target uid 资源目标
  @target_owner uid 资源所有者
  @target_type Int 资源类型 0: 贴文，1:私信，2:用户
*/

var NotificationSchema = `
	type Notification {
    isRead
    sender
    recipient
    action
    msg
    target
    target_owner
    target_type
    created_at
    updated_at
  }

  isRead: bool .
  sender: uid @reverse .
  recipient: []uid @reverse .
  action: int .
  msg: string .
  target: uid .
  target_owner: uid .
  target_type: int .
`
