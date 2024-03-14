
// 查询帖子详情
pub const STATUS_DETAIL: &'static str = r#"query all(
  $uid: string,
  $logged_user_id: string
) {
  data(func: uid($uid)) {
    uid 
    text
    is_forward

    reply_cnt as count(replies @filter(uid_in(replies,$logged_user_id)))
    is_replied: math(reply_cnt == 1)
    replies_count

    forward_cnt as count(forwards @filter(uid($logged_user_id)))
    is_forwarded: math(forward_cnt == 1)
    forwards_count

    favorite_cnt as count(favorites @filter(uid($logged_user_id)))
    is_favorited: math(favorite_cnt == 1)
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
    entities{
      uid
      urls{
        uid
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
      medias{
        uid
        url
        media_type
      }
    }
    forward_to_status{
      uid
      text
      created_at
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
      }
      entities{
        uid
        urls{
          uid
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
        medias{
          uid
          url
          media_type
        }
      }
    }
  }
}"#;

/*
  1，功能：用户主页帖子列表
  2，规则1：我关注的用户的帖子
  3，规则2：我关注的用户喜欢的帖子
  4，规则3：推荐的认证用户的帖子
  5，规则4：我发布的帖子
*/ 
pub const USER_HOME_STATUS: &'static str = r#"query all(
  $logged_user_id: string,
  $first: int,
  $after: string
) {
  var(func: uid($logged_user_id))  {
    follows { 
      p1 as statuses
      p2 as ~favourites
    }
    p3 as statuses
  }
  data(
    func: type(Status), 
    first: $first, 
    after:$after,
    orderdesc: created_at 
  ) @filter(uid(p1,p2,p3)) {
    uid 
    text
    is_forward

    reply_cnt as count(replies @filter(uid_in(replies,$logged_user_id)))
    is_replied: math(reply_cnt == 1)
    replies_count

    forward_cnt as count(forwards @filter(uid($logged_user_id)))
    is_forwarded: math(forward_cnt == 1)
    forwards_count

    favorite_cnt as count(favorites @filter(uid($logged_user_id)))
    is_favorited: math(favorite_cnt == 1)
    favorites_count

    created_at
    group{
      uid
      group_name
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
    entities{
      uid
      urls{
        uid
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
      medias{
        uid
        url
        media_type
      }
    }
    forward_to_status{
      uid
      text
      created_at
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
      }
      entities{
        uid
        urls{
          uid
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
        medias{
          uid
          url
          media_type
        }
      }
    }
  }
}"#;

