
// 查询回复详情
pub const REPLY_DETAIL: &'static str = r#"query all(
  $uid: string,
  $logged_user_id: string
) {
  data(func: uid($uid)) {
    uid 
    text
    is_to_reply

    reply_cnt as count(replies @filter(uid_in(replies,$logged_user_id)))
    is_replied: math(reply_cnt == 1)
    replies_count

    favorite_cnt as count(favorites @filter(uid($logged_user_id)))
    is_favorited: math(favorite_cnt == 1)
    favorites_count

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
    status{
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
    to_reply{
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
  1，功能：帖子/回复的回复列表
*/ 
pub const STATUS_REPLIES: &'static str = r#"query all(
  $uid: string,
  $logged_user_id: string,
  $first: int,
  $after: string
) {
  data(func: uid($uid)) {
    uid
    edges: replies (first: $first, after:$after) {
      uid 
      text
      is_to_reply

      reply_cnt as count(replies @filter(uid_in(replies,$logged_user_id)))
      is_replied: math(reply_cnt == 1)
      replies_count

      favorite_cnt as count(favorites @filter(uid($logged_user_id)))
      is_favorited: math(favorite_cnt == 1)
      favorites_count

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
      to_reply{
        uid 
        user{
          uid
          name
          username
        }
      }
    }
  }
}"#;