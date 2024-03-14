
### 1，数据结构
```js
type Notification {
	isRead: bool,
	sender: int64, // 发送者id
	recipient: int64, // 接收者id
	action: int, // 0:回复，1:喜欢，2:转帖，3:引用，4:关注，5:提及，6:公告
	msg: string, // 额外的信息
	target: int64, // 动作目标id ，比如贴文id、回复id
	target_owner: int64, // 动作目标拥有者id ，比如贴文作者id、回复作者id
	target_type: int, // 0:贴文，1:回复，2:私信
	created_at datetime NOT NULL DEFAULT datetime(),
	updated_at datetime
}
```

### 2，在nebula中创建tag、edge、index
```go
const STATUS = `
CREATE TAG notify(
	isRead bool NOT NULL DEFAULT false,
	action int NOT NULL DEFAULT 0,
	msg string,
	target_type int,
	created_at datetime NOT NULL DEFAULT datetime()
);
create edge sender();
create edge recipient();
create edge target();
create edge target_owner();
`
```
