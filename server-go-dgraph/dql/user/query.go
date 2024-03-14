package user

// 查询名字/主名/手机号是否已注册
var UserExist = `query all(
  $name: string,
  $username: string, 
  $phone_number: string
) {
  name(func: eq(name,$name)) {
    id:uid        
  }
  username(func: eq(username,$username)) {
    id:uid
  }
  phone_number(func: eq(phone_number,$phone_number)) {
    id:uid
  }
}`

// 登录：根据username/email/phone_number查询账号信息
var FindAccount = `query all(
  $email: string,
  $username: string, 
  $phone_number: string
) {
  user(func: type(User)) @filter(eq(phone_number,$phone_number) OR eq(username,$username) OR eq(email,$email) ) {
    id:uid
    name
    username
    verify_name
    password
    role
    lang
    avatar_url
    verified
  }
}`

// 登录：根据id查询账号信息
var FindUserById = `query all($id: string) {
  user(func: eq(id,$id)) {
    id:uid
    name
    username
    verify_name
    password
    role
    lang
    avatar_url
    verified
  }
}`

// 查询名字是否已注册
var NameExist = `query all($name: string) {
  name(func: eq(name,$name)) {
    id:uid        
  }
}`

// 查询主名是否已注册
var UsernameExist = `query all($username: string) {
  username(func: eq(username,$username)) {
    id:uid
  }
}`

// 查询手机号是否已注册
var PhoneNumberExist = `query all($phone_number: string) {
  phone_number(func: eq(phone_number,$phone_number)) {
    id:uid
  }
}`

// 查询手机号是否已注册
var EmailExist = `query all($email: string) {
  email(func: eq(email,$email)) {
    id:uid
  }
}`

// 查询用户密码
var GetUserPassword = `query all($loggedUserid: string) {
  user(func: uid($loggedUserid)) {
    id:uid   
    password     
  }
}`

// 查询用户个人主页
var GetUserHomepage = `query all(
  $username: string,
  $loggedUserid: string
) {
  item(func: eq(username,$username)) {
    id:uid
    name
    username
    verify_name
    verified
    bio
    cnt as cnt:count(~follows @filter(uid($loggedUserid)))
    following: math(cnt == 1)
    same_follows: ~follows @filter(uid_in(~follows,$loggedUserid)) (first: 5) {
      id:uid
      name
      username
      verify_name
    }
    avatar_url
    cover_url
    school
    gender
    birthday
    website
    followers_count
    followings_count
    created_at
  }
}`

// 查询关注列表
var Followings = `query all(
  $username: string, 
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  item(func: eq(username, $username)) {
    id:uid 
    edges: follows (first: $first, after: $after) {
      id:uid
      name
      username
      verify_name
      verified
      cnt as cnt:count(~follows @filter(uid($loggedUserid)))
      following: math(cnt == 1)
      bio
      avatar_url
    }
  }
}`

// 查询粉丝列表
var Followers = `query all(
  $username: string, 
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  item(func: eq(username, $username)) {
    id:uid 
    edges: ~follows (first: $first, after:$after) {
      id:uid
      name
      username
      verify_name
      verified
      cnt as cnt:count(~follows @filter(uid($loggedUserid)))
      following: math(cnt == 1)
      bio
      avatar_url
    }
  }
}`

// 查询好友列表
var Friends = `query all(
  $username: string, 
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  item(func: eq(username, $username)) {
    id:uid 
    edges: follows @filter(uid_in(follows, $loggedUserid)) (first: $first, after:$after) {
      id:uid
      name
      username
      verify_name
      verified
      cnt as cnt:count(~follows @filter(uid($loggedUserid)))
      following: math(cnt == 1)
      bio
      avatar_url
    }
  }
}`

// 查询黑名单列表
var Blacklists = `query all(
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  item(func: uid($loggedUserid)) {
    id:uid 
    edges: blacklists (first: $first, after:$after) {
      id:uid
      name
      username
      verify_name
      verified
      cnt as cnt:count(~follows @filter(uid($loggedUserid)))
      following: math(cnt == 1)
      bio
      avatar_url
    }
  }
}`

// 查询共同关注
var SameFollowings = `query all(
  $username: string, 
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  item(func: eq(username, $username)) {
    id:uid 
    edges: follows @filter(uid_in(~follows, $loggedUserid)) (first: $first, after:$after) {
      id:uid
      name
      username
      verify_name
      verified
      cnt as cnt:count(~follows @filter(uid($loggedUserid)))
      following: math(cnt == 1)
      bio
      avatar_url
    }
  }
}`

// 查询关系关注=> 我关注的A、B等人也关注了:username
// 反向思维=> :username的粉丝的粉丝里有我
var RelationFollowings = `query all(
  $username: string, 
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  item(func: eq(username, $username)) {
    id:uid 
    edges: ~follows @filter(uid_in(~follows, $loggedUserid)) (first: $first, after:$after) {
      id:uid
      name
      username
      verify_name
      verified
      following: math(1 == 1)
      bio
      avatar_url
    }
  }
}`

// 用户时间线
var QueryTimeline = `query all(
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  var(func:uid($loggedUserid)) {
    US as statuses
    UFo as follows{
      F as favorites
      S as statuses
    }
  }
  statuses(func:uid(US,S,F)) {
    ...StatusFragment
  }
}
fragment StatusFragment {
  id:uid        
  text
  media_type
  user{
    name
    username
    verify_name
    avatar_url
    verified
  }
  to_user{
    name
    username
    verify_name
  }
  urls
  images{
    url
    source
    platform
  }
  hashtags{
    id:uid
    name
    description
  }
  video{
    url
    source
    platform
  }
  reply_count
  quote_count
  restatus_count
  favorite_count
  platform

  favorite_users:~favorites @filter(uid(UFo)) (first:3) {
    id:uid
    name
    count:count(name)
  }
  restatus_users:~restatuses @filter(uid(UFo)) (first:3) {
    id:uid
    name
    count:count(name)
  }

  favorite_cnt as fcnt:count(~favorites @filter(uid($loggedUserid)))
  favorited: math(favorite_cnt == 1)

  restatus_cnt as rcnt:count(~restatuses @filter(uid($loggedUserid)))
  restatused: math(restatus_cnt == 1)

  created_at
  to_quote{
    id:uid        
    text
    status_type
    media_type
    user{
      name
      username
      verify_name
      avatar_url
      verified
    }
    to_user{
      name
      username
      verify_name
    }
    urls
    images{
      url
      source
      platform
    }
    hashtags{
      id:uid
      name
      description
    }
    video{
      url
      source
      platform
    }
    created_at
  }
}
`
