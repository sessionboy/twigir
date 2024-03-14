
pub const USER_SCHEMA: &str = "
  type User {
    name
    username
    phone_code
    phone_number
    password
    email
    lang 
    avatar_url
    description
    user_entities
    role
    location

    profile_cover_url
    profile_default_cover
    profile_gender
    profile_birthday
    profile_school
    profile_isgraduation
    profile_job
    profile_website
    profile_emotion
    profile_country
    profile_province
    profile_city

    is_verified
    verified

    statuses
    statuses_count

    follows
    <~follows>
    followers_count
    followings_count
    friends_count
    shields
    sign_records

    receive_chat
    receive_like_notification
    receive_reply_notification

    last_sign_at
    last_publish_at
    last_reply_at
    created_at
    updated_at
  }

  type Verified {
    type
    name
    is_group 
    Verified.category
    level
    description 
    created_at 
    updated_at 
  }

  type SignRecord {
    platform
    ip
    device
    browser
    is_register
    created_at
  }

  name: string @index(hash) .
  username: string @index(hash) .
  phone_code: string .
  phone_number: string @index(hash) .
  password: string .
  email: string @index(hash) .
  lang: string .
  avatar_url: string .
  description: string @index(hash) .
  user_entities: uid .
  role: string .
  location: geo @index(geo) .

  profile_cover_url: string .
  profile_default_cover: bool .
  profile_gender: string .
  profile_birthday: dateTime .
  profile_school: string .
  profile_isgraduation: bool .
  profile_job: string .
  profile_website: string .
  profile_emotion: string .
  profile_country: string .
  profile_province: string .
  profile_city: string .

  is_verified: bool .
  verified: uid @reverse .

  statuses: [uid] @count .
  statuses_count: int .

  follows: [uid] @reverse @count .
  followers_count: int .
  followings_count: int .
  friends_count: int .
  shields: [uid] @reverse @count .
  sign_records: [uid] .
 
  receive_chat: bool .
  receive_like_notification: bool .
  receive_reply_notification: bool .

  last_sign_at: dateTime .
  last_publish_at: dateTime .
  last_reply_at: dateTime .
  created_at: dateTime .
  updated_at: dateTime .
  
  type: string .
  is_group: bool .
  Verified.category: string .
  level: int .
  
  platform: string .
  ip: string .
  device: string .
  browser: string .
  is_register: bool .

  created_at: dateTime .
  updated_at: dateTime .
";