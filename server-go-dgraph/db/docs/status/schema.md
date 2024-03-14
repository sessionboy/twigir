#### 1，数据模型
```js
type Status {
  id: ID!
  text: String! @search(by: [fulltext])
  status_type: Int! // 0：一般贴文；1：转帖；2：引用；3：回复；4：回复的回复
  media_type: Int! // 0:photo, 1:video
  reviewed: Boolean! // 是否已通过审核，默认为false
  ip String
	platform String
	device String

  mentions: [User]
  urls: [Url]
  hashtags: [Hashtag]
  photos: [Photo]
  videos: [Video]

  user: User
  to_user: User       // 回复谁
  to_status: Status   // 父贴文/回复
  favorite: [User]

  replies_count: Int,
	favorites_count: Int,
	quotes_count: Int,
	restatuses_count: Int,
  
  created_at: DateTime!
  update_at: DateTime
}
type Url {
	url String
	url_key String
}
type Hashtag {
  user User
	name String
	description String
}
type Photo {
  ip: String
	url: String
	source: String
  created_at: DateTime!
}
type Video {
  ip: String
	url: String
	play_count: Int  // 观看次数
	duration: Int    // 时长
	source: String
	platform: string // 存储平台
  created_at: DateTime!
}
```

#### 2，数据schema
```go
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
    photos
    videos

    user
    to_user
    to_status
    favorites

    replies_count
    favorites_count
    quotes_count
    restatuses_count

    created_at
    updated_at
  }

  type Url {
    url
	  url_key
  }
  type Hashtag {
    user
	  name
    description
  }
  type Photo {
    url
    source
    ip
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
  status_type: int .
  media_type: int .
  reviewed: bool .
  
  mentions: [uid] @reverse .
  urls: [uid] .
  hashtags: [uid] @reverse @count .
  photos: [uid] .
  videos: [uid] .

  user: uid @reverse .
  to_user: uid @reverse .
  to_status: uid @reverse .
  favorites: [uid] @reverse @count .

  replies_count: int .
  favorites_count: int .
  quotes_count: int .
  restatuses_count: int .

  url: string .
	url_key: string .
  description: string .
  source: string .
  play_count: int .
  duration: int .
`
```