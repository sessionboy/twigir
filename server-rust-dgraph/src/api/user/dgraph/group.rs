
/*
  1，功能：用户{id}的小组列表 
*/ 
pub const USER_GROUPS: &'static str = r#"query all(
  $user_id: string,
  $logged_user_id: string,
  $first: int,
  $after: string
) {
  var(func: uid($user_id)) {
    ~member_user @filter(type(Member)) (first: $first, after: $after) {      
      G as ~members @filter(type(Group))
    }
  }
  data(func: type(Group)) @filter(uid(G)) {
    uid
    group_name
    group_description
    avatar_url
    access
    members_count
    statuses_count
    created_at

    member_self : members @filter(uid_in(member_user,$logged_user_id)){
      uid
      is_owner
      is_admin
    }

    cnt as count(members @filter(uid_in(member_user,$logged_user_id)))
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