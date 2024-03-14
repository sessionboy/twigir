
## 一，账户相关
```js
// 登录
post : /auth/login    
// 注册
post : /auth/register
// 重置密码
post : /auth/reset

// 发送手机验证码
post : /auth/send_phone
// 验证验证码
post : /auth/verify_phone
// 发送邮箱验证码
post : /auth/send_email
// 验证邮箱
post : /auth/verify_email

// 更改资料
put : /account/profile
// 更改名字
put : /account/name
// 更改名字
put : /account/username
// 更改名字
put : /account/phone
// 更改名字
put : /account/email
// 更改名字
put : /account/password

// 我的个人主页信息
get : /users/me
// 推荐用户
get : /users/recommends
// 获取指定用户个人主页信息
get : /users/:id
// 我关注某人[id]
post : /users/:id/follow
// 我取关某人[id]
post : /users/:id/unfollow
// 我屏蔽某人[id]
post : /users/:id/shield
// 我解除屏蔽某人[id]
post : /users/:id/unshield
// 获取用户粉丝列表
get : /users/:id/followers
// 获取用户关注列表
get : /users/:id/followings
// 我和[id]的共同关注
get : /users/:id/same_followings
// 我和[id]的微关系户 ，即我关注的A、B、C等人也关注了[id]
get : /users/:id/relation_followings
// 获取用户屏蔽列表
get : /users/:id/shields
// 获取用户发布的帖子列表
get : /users/:id/status
// 获取用户小组列表
get : /users/:id/groups
// 获取用户所有小组的帖子列表
get : /users/:id/groups_status
// 获取用户的喜欢列表
get : /users/:id/favorites
// 获取用户的媒体(照片/视频等)帖子列表
get : /users/:id/media_status
```

## 二，贴子相关
```js
// 发布/转发帖子
post : /status
// 更新帖子 - 暂不做
put : /status/:id
// 点赞/取消点赞该帖子
post : /status/:id/favorite
// 删除帖子
delete : /status/:id/destroy
// 获取帖子
get : /status/:id
// 获取该帖子的回复列表
get : /status/:id/replies
// 获取该帖子的点赞列表
get : /status/:id/favorites

// 发布回复
post : /reply
// 获取回复详情
get : /reply/:id
// 获取回复的回复列表
get : /reply/:id/replies
// 删除回复
delete : /reply/:id/destroy
// 点赞/取消点赞回复
put : /reply/:id/favorite
```

## 三，小组相关
```js
// 创建小组
post : /group
// 获取帖子
get : /group/:groupId
// 更新帖子 - 暂不做
put : /group/:groupId
// 小组发帖
put : /group/:groupId/publish
// 解散小组
delete : /group/:groupId/destroy
// 加入小组
post : /group/:groupId/join
// 退出小组
post : /group/:groupId/leave
// 添加管理员
post : /group/:groupId/add_admin/:memberId
// 开除管理员
post : /group/:groupId/fire_admin/:memberId
// 开除成员
// post : /group/:groupId/fire_member/:memberId
// 禁言
post : /group/:groupId/forbidden/:memberId

```

## 四，通知相关
```js
// 获取通知列表
get : /notifications
```

## 四，私信相关
```js
// 获取私信列表
get : /messages/list
// 获取私信
get : /messages/:id
// 发私信
post : /messages/send
// 屏蔽私信
post : /messages/shield
```

## 五，搜索相关
```js
// 搜索
get : /search?q=type&k=keyword
```

## 六，系统相关
```js
// 关于
get : /about
// 下载
get : /download
```