
// 查询用户信息
pub const STATUS_DETAIL: &'static str = r#"query all(
  $uid: string,
  $logged_user_id: string
) {
  data(func: uid($uid)) {
    uid 
    text
    is_forward
    replies_count
    forwards_count
    favorites_count
    created_at
    group{
      uid
      name
      is_verified
    }
    user{
      uid
      name
      username
      avatar_url
      description
      is_verified
      verified{
        uid
        name
        description
      }
      cnt as count(~follows @filter(uid($logged_user_id)))
      following: math(cnt == 1)
    }
    entities:{
      urls{
        url
        url_key
      }
      mentions{
        uid
        name
        username
      },
      hashtags{
        uid
        name
      }
    }
    medias{
      photos{
        url
      }
      video{
        url
      }
      music{
        url
      }
    }
  }
}"#;

