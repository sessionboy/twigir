
/*
  1，功能：用户{id}发布的帖子列表
*/ 
pub const USER_STATUS: &'static str = r#"query all(
  $user_id: string,
  $logged_user_id: string,
  $first: int,
  $after: string
) {
  data(func: uid($user_id)) {
    uid
    edges: statuses (first: $first, after: $after) {
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

/*
  1，功能：用户{id}发布的照片帖子列表
*/ 
pub const USER_STATUS_MEDIA: &'static str = r#"query all(
  $user_id: string,
  $logged_user_id: string,
  $first: int,
  $after: string,
  $media_type: string
) {
  var(func: uid($user_id)) {
    statuses { 
      entities @filter(eq(media_type, $media_type)) {
        S as ~entities @filter(type(Status))
      }
    }
  }
  data(func: uid($user_id)) {
    uid
    edges: statuses @filter(uid(S)) (first: $first, after: $after) {
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

/*
  1，功能：用户{id}喜欢的帖子列表
*/ 
pub const USER_STATUS_FAVORITE: &'static str = r#"query all(
  $user_id: string,
  $logged_user_id: string,
  $first: int,
  $after: string
) {
  data(func: uid($user_id)) {
    uid
    edges: ~favorites @filter(type(Status)) (first: $first, after: $after) {
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


/*
  1，功能：用户{id}加入的小组的帖子列表
*/ 
pub const USER_GROUPS_STATUS: &'static str = r#"query all(
  $user_id: string,
  $logged_user_id: string,
  $first: int,
  $after: string
) {
  var(func: uid($user_id)) {
    ~member_user @filter(type(Member)) {
      G as ~members @filter(type(Group))
    }
  }
  data(func: type(Group)) @filter(uid(G)) {
    uid
    edges: statuses @filter(type(Status)) (first: $first, after: $after) {
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