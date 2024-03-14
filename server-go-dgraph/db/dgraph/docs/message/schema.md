#### 1，数据模型
```js
type Message {
	msg string NOT NULL, // 私信内容
	msg_type: int, // 0:文本、1:图片、2:视频
	media_url: string, // 发送的媒体url，比如图片链接
	sender: Sender, // 发送者
	recipient: recipient,  // 接收者
	created_at datetime NOT NULL DEFAULT datetime(),
}
```

#### 2，数据schema
```go
var MessageSchema = `
type Message {
    msg
    msg_type
    media_url
    sender
    recipient
    created_at
    updated_at
  }

  msg_type: int .
  media_url: string .
`
```