
### 1，数据结构
```js
type Message {
	msg string NOT NULL,
	msg_type: int, // 0:文本、1:图片、2:视频
	media_url: string,
	sender: Sender, // 发送者
	recipient: recipient,  // 接收者
	created_at datetime NOT NULL DEFAULT datetime(),
}
```
```go
const STATUS = `
CREATE TAG message(
	msg string NOT NULL,
	msg_type int NOT NULL DEFAULT 0,
	media_url string,
	created_at datetime NOT NULL DEFAULT datetime()
);
create edge sender();
create edge recipient();
`
```
