package dgraph

/*
  @notify_unFollow  隐藏未关注的人的通知
  @chat_unFollow  禁止未关注的人的私信
  @unconversations 拒绝对话名单
*/

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
    status_count
    favorites
    restatuses
    follows
    followers_count
	  followings_count
    blacklists
    unconversations
    records
	
    notify_unFollow
    notify_unFollowMe
    notify_unVerified
    notify_blacklist
    chat_unFollow
    chat_unFollowMe
    chat_unVerified
    chat_blacklist

    last_publish_at
    last_reply_at
    created_at
    updated_at
  }

  type Record {
    user
    type
    ip
    platform
    device
    created_at
  }

  name: string @index(term) .
  username: string @index(term) .
  verify_name: string @index(term) .
	verify_level: int .
  verified: bool @index(bool) .
  phone_code: string .
  phone_number: string @index(hash) .
  phone_country: string .
  email: string @index(hash) .
  password: string .
  lang: string .
  avatar_url: string .
  bio: string @index(term) .
  role: int .
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
  status_count: int .
 
  follows: [uid] @reverse @count .
  followers_count: int .
  followings_count: int .
  blacklists: [uid] @reverse @count .
  records: [uid] @reverse @count .

  notify_unFollow: bool .
  notify_unFollowMe: bool .
  notify_unVerified: bool .
  notify_blacklist: bool .
  chat_unFollow: bool .
  chat_unFollowMe: bool .
  chat_unVerified: bool .
  chat_blacklist: bool .
 
  last_publish_at: dateTime .
  last_reply_at: dateTime .
  created_at: dateTime .
  updated_at: dateTime .
  
  type: int .
  ip: string .
  platform: string .
  device: string .
`
