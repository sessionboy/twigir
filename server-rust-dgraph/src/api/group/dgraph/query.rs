
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


/*
  1，功能：小组详情
*/ 
pub const GROUP_DETAIL: &'static str = r#"query all(
  $uid: string,
  $logged_user_id: string
) {
  data(func: uid($uid)) {
    uid
    group_name
    group_description
    announcement
    avatar_url
    cover_url
    default_cover
    access
    visible
    members_count
    statuses_count
    created_at

    cnt as count(members @filter(uid_in(user,$logged_user_id)))
    is_joined: math(cnt == 1)

    is_verified
    verified{
      uid
      name
      description
    }    
    group_entities{
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
}"#;


/*
  1，功能：小组{id}的帖子列表
*/ 
pub const GROUP_STATUS: &'static str = r#"query all(
  $uid: string,
  $logged_user_id: string,
  $first: int,
  $after: string
){
  data(func: uid($uid)) {
    uid
    edges: ~group @filter(type(Status)) (first: $first, after: $after) {
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
  }
}"#;