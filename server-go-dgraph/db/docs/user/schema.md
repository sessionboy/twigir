#### 1，数据模型
```js
type User {
  id: ID!
  name: String! @search(by: [term])
  username: String! @id @search(by: [term])
  verify_name String @id @search(by: [term])
	verify_level Int
  verified: Boolean!
  phone_code: String!
  phone_number: String! 
  phone_country: String!
  password: String!
  email: String
  lang: String 
  avatar_url: String
  bio: String @search(by: [term])
  role: Int 
  location: String
  
  // 个人资料
  cover_url: String
  gender: Int
  birthday: String
  school: String 
  isgraduation: Boolean
  job: String 
  website: String
  emotion: String
  country: String
  province: String
  city: String

  statuses: [Status]
  statuses_count: Int
  follows: [User]
  followers_count: Int
  followings_count: Int
  blacklists: [User]
  records: [Activity]

  // 用户设置
  notify_hide_unFollow: Boolean!    // 我未关注的人
	notify_hide_unFollowMe: Boolean!  // 未关注我的人
	notify_hide_unVerified: Boolean!  // 非认证用户
	notify_hide_blacklist: Boolean!   // 黑名单用户
	nochat_unFollow: Boolean!        
	nochat_unFollowMe: Boolean!
	nochat_unVerified: Boolean!
	nochat_blacklist bool: Boolean!

  created_at: DateTime!
  update_at: DateTime
}

type Activity {
  type: Int
	ip: String
	Platform: String
	Device: String
	created_at: DateTime!
}
```

#### 2，数据schema
```go
var UserSchema = `
type User {
    name
    username
    verify_name
	  verify_level
    verified
    phone_code
    phone_number
    phone_country
    email
    password
    lang 
    avatar_url
    bio
    role
    location

    cover_url
    gender
    birthday
    school
    isgraduation
    job
    website
    country
    province
    city

    statuses
    statuses_count
    follows
    followers_count
	  followings_count
    blacklists
    records
	
    notify_hide_unFollow
    notify_hide_unFollowMe
    notify_hide_unVerified
    notify_hide_blacklist
    nochat_unFollow
    nochat_unFollowMe
    nochat_unVerified
    nochat_blacklist

    created_at
    updated_at
  }

  type Activity {
    type
    ip
    platform
    device
    created_at
  }

  id: int .
  name: string @index(term) .
  username: string @index(term) .
  verify_name: string @index(term) .
	verify_level: int .
  verified: bool .
  phone_code: string .
  phone_number: string .
  phone_country: string .
  email: string .
  password: string .
  lang: string .
  avatar_url: string .
  bio: string @index(term) .
  role: string .
  location: string .

  cover_url: string .
  gender: int .
  birthday: dateTime .
  school: string .
  isgraduation: bool .
  job: string .
  website: string .
  country: string .
  province: string .
  city: string .

  statuses: [uid] @reverse @count .
  statuses_count: int .

  follows: [uid] @reverse @count .
  followers_count: int .
  followings_count: int .
  blacklists: [uid] @reverse @count .
  records: [uid] @reverse @count .

  notify_hide_unFollow: bool .
  notify_hide_unFollowMe: bool .
  notify_hide_unVerified: bool .
  notify_hide_blacklist: bool .
  nochat_unFollow: bool .
  nochat_unFollowMe: bool .
  nochat_unVerified: bool .
  nochat_blacklist: bool .
 
  created_at: dateTime .
  updated_at: dateTime .
  
  type: string .
  ip: string .
  platform: string .
  device: string .
`
```