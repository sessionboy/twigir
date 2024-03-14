
### 1，数据结构
```js
type user {
	name string NOT NULL,
	phone_code // 区号：+86 
	phone_number // 完整号码：+86132...
	phone_country // 国家码：CN、US...
	location // 目前所在的位置 ，比如北京
  ...
	// edge
	follows: [user],  // 包含 followers、followings
	statuses: [status],
	blacklists: [user],
	records: [activity]
	...
}
```

### 2，在nebula中创建tag、edge、index
```go
const USER = `
CREATE TAG user(
	name string NOT NULL,
	username string NOT NULL,
	auth_name string NOT NULL,
	phone_code string NOT NULL,
	phone_number string NOT NULL,
	phone_country string NOT NULL,
	email string,
	password string NOT NULL,
	bio string,
	lang string NOT NULL DEFAULT "zh",
	avatar_url string,
	role int NOT NULL DEFAULT 0,
	location string,
	verified bool NOT NULL DEFAULT false,
	verify_name string,
	verify_level int,
	followers_count int NOT NULL DEFAULT 0,
	followings_count int NOT NULL DEFAULT 0,
	friends_count int NOT NULL DEFAULT 0,
	statuses_count int NOT NULL DEFAULT 0,
	replies_count int NOT NULL DEFAULT 0,
	created_at datetime NOT NULL DEFAULT datetime(),
	updated_at datetime
);
create edge follows(start_at datetime);
create edge statuses();
create edge blacklists();
create edge records();
create tag index index_name on user(name(50));
create tag index index_username on user(username(50));
create tag index index_phone_number on user(phone_number(50));
create tag index index_email on user(email(50));
create tag index index_bio on user(bio(200));
`
// record边: 活动记录 activity
```

#### 用户资料
```go
const SETTING = `
CREATE TAG profile(
	location string,
	cover_url string,
	gender int,
	birthday date,
	school string,
	isgraduation bool,
	job string,
	website string,
	emotion int,
	country string,
	province string,
	city string,
	updated_at datetime
);
`
```

#### 用户设置
```go
// unFollow:我未关注的人
// unFollowMe:未关注我的人
// unVerified:非认证的用户
// shield:黑名单用户
const SETTING = `
CREATE TAG setting(
	notify_hide_unFollow bool NOT NULL DEFAULT false,
	notify_hide_unFollowMe bool NOT NULL DEFAULT false,
	notify_hide_unVerified bool NOT NULL DEFAULT false,
	notify_hide_shield bool NOT NULL DEFAULT false,
	nochat_unFollow bool NOT NULL DEFAULT false,
	nochat_unFollowMe bool NOT NULL DEFAULT false,
	nochat_unVerified bool NOT NULL DEFAULT false,
	nochat_shield bool NOT NULL DEFAULT false,
	updated_at datetime
);
`
```

#### 用户登录/注册记录
```go 
const AUTHENTICATE = `
CREATE TAG activity(
	type int NOT NULL DEFAULT 1,
	ip string NOT NULL,
	platform string NOT NULL,
	device string NOT NULL,
	created_at datetime NOT NULL DEFAULT datetime()
);
`

// type: 0：注册，1：登录
```
