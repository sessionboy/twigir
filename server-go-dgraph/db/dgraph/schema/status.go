package dgraph

/*
  @status_type  0：一般贴文；1：转帖；2：引用；3：回复；4：回复的回复
  @media_type  0:无媒体文件；1:图片；2:视频
  @reviewed 是否已通过审核
  @user_id  贴文作者的id，用于搜索条件筛选(比如搜索我关注的人发布的贴文)
  @to_user 回复哪个用户，通常是to_reply(回复的目标贴文)的作者
  @to_status 回复所属的贴文，祖帖。
  @to_reply 回复目标贴文，可是是贴文，也可以是回复
  @verify_int 表示是否是认证用户发布的贴文 (用于排序/搜索)，0则表示否，1则表示认证
  @recommended 推荐指数(用于排序/搜索)
  @replies 查询 `~to_reply`即可
  @quotes 查询 `~to_quote`即可
*/
var StatusSchema = `
	type Status {
    text
    status_type
    media_type
	  reviewed
    ip
    platform
    device

    mentions
    urls
    hashtags
    images
    video

    user
    user_id
    to_user
    to_status
    to_quote
    to_reply
    restatuses
    favorites
    verify_int
    recommended

    reply_count
    quote_count
    restatus_count
    favorite_count

    created_at
    updated_at
  }
  type Hashtag {
    user
	  name
    description
  }
  type Image {
    url
    source
    ip
    platform
    created_at
  }
  type Video {
    url
    play_count
    duration
    source
    platform
    ip
    created_at
  }

  text: string @index(fulltext) .
  status_type: int @index(int) .
  media_type: int @index(int) .
  verify_int: int @index(int) .
  reviewed: bool .
  user_id: string @index(hash) .
  
  urls: [string] .
  mentions: [uid] @reverse .
  hashtags: [uid] @reverse @count .
  images: [uid] .
  video: uid .

  user: uid @reverse .
  to_user: uid @reverse .
  to_status: uid @reverse .
  to_reply: uid @reverse .
  to_quote: uid @reverse .
  restatuses: [uid] @reverse @count .
  favorites: [uid] @reverse @count .

  reply_count: int .
  favorite_count: int .
  quote_count: int .
  restatus_count: int .
  recommended: int .

  url: string .
  description: string .
  source: string .
  play_count: int .
  duration: int .
`
