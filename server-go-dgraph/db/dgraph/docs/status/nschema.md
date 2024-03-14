
### 1，数据结构
```js
type Status {
	text string NOT NULL,
	media_type: int, // 媒体类型，0：photo，1：videos
	reviewed: bool, // 是否通过审核
	status_type: int, // 0：一般贴文；1：转帖；2：引用；3：回复；4：回复的回复
  ...

	// edge边
	owner: user,
	to_user: user,  // 回复哪个user
	to_status: status, // 父贴文/父回复
	urls: [url],
	hashtags: [hashtag],
	mentions: [user],
	photos: [photo],
	media_video: video

	...
	created_at datetime NOT NULL DEFAULT datetime(),
	updated_at datetime
}

type url {
	url string
	url_key string
}
type hashtag {
	title string
	description string
}
type photo {
	url string
	source string
}
type video {
	url string
	play_count int  // 观看次数
	duration int    // 时长
	source string
	platform string // 存储平台
}
```

后端组织 entities
```go
type entities {
	media_type: int,
	urls: [url],
	hashtags: [hashtag],
	mentions: [user],
	photos: [photo],
	media_video: video
}
```

CREATE TAG User(
	name string
);
CREATE TAG Profile(
	age int
);
CREATE TAG Photo(
	url string
);

### 2，在nebula中创建tag、edge、index
```go
// AND $$.status.reviewed == true
// media_type -> 0:photo, 1:video
// to_status -> 父贴文，可能是回复，也可能是贴文
const STATUS = `
CREATE TAG status(
	text string NOT NULL,
	status_type int NOT NULL DEFAULT 0,
	media_type int,
	reviewed bool NOT NULL DEFAULT false,
	replies_count int NOT NULL DEFAULT 0,
	favorites_count int NOT NULL DEFAULT 0,
	quotes_count int NOT NULL DEFAULT 0,
	restatuses_count int NOT NULL DEFAULT 0,
	ip string,
	platform string,
	device string,
	created_at datetime NOT NULL DEFAULT datetime(),
	updated_at datetime
);
CREATE TAG url(
	url string,
	url_key string
);
CREATE TAG hashtag(
	title string,
	description string
);
CREATE TAG photo(
	url string,
	source string,
	platform string
);
CREATE TAG video(
	url string,
	play_count int, 
	duration int,
	source string,
	platform string
);
create edge owner();
create edge to_status();
create edge to_user();
create edge mentions();
create edge urls();
create edge hashtags();
create edge photos();
create edge videos();
create edge favorites();
create tag index index_text on status(text(256));
`
```
