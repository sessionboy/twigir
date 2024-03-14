package dgraph

/*
  @msg 私信内容
  @msg_type 0:文本、1:图片、2:视频
  @media_url 发送的媒体url，比如图片链接
  @sender 发送者
  @recipient 接收者
  @isRead  接收方是否已读该消息
  @conversation_id 房间id，格式：userid_userid (左边是初次私聊发起人id)
  @last_publish_at 最近一次消息时间，用于排序
*/
var MessageSchema = `
	type Message {
    isRead
    conversation_id
    msg
    msg_type
    media_url
    sender
    recipient
    created_at
  }
  type Conversation {
    conversation_id
    creater
    users
    messages
    created_at
    last_publish_at
  }

  conversation_id: string @index(term) .
  creater: uid .
  users: [uid] @reverse @count .
  messages: [uid] @reverse @count .
  msg_type: int .
  media_url: string .
`
