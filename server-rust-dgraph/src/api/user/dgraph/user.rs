
// 查询我的个人主页信息
pub const USER_ME: &'static str = r#"query all($uid: string) {
  data(func: uid($uid)) {
    uid 
    name
    username
    description
    lang
    avatar_url
    profile_cover_url
    profile_default_cover
    profile_school
    profile_gender
    profile_birthday
    profile_website
    profile_emotion
    is_verified
    verified{
      uid
      name
      description
    }
    followers_count
    followings_count
    friends_count
    created_at
  }
}"#;

// 查询指定用户个人主页信息
// $logged_user_id: string
// is_followed: count(follows @filter(uid($logged_user_id)))
pub const USER: &'static str = r#"query all(
  $uid: string,
  $logged_user_id: string
) {
  data(func: uid($uid)) {
    uid 
    name
    username
    description
    lang
    cnt as count(~follows @filter(uid($logged_user_id)))
    following: math(cnt == 1)
    avatar_url
    profile_cover_url
    profile_default_cover
    profile_school
    profile_gender
    profile_birthday
    profile_website
    profile_emotion
    is_verified
    verified{
      uid
      name
      description
    }
    followers_count
    followings_count
    friends_count
    created_at
  }
}"#;
