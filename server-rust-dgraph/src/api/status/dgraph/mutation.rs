
// 查询用户信息
pub const STATUS_CREATE: &'static str = r#"query all($uid: string) {
  data(func: uid($uid)) {
    uid 
    last_publish_at
    last_reply_at
    statuses_count
  }
}"#;

// 查询用户是否是该小组成员
pub const USER_IS_GROUP_MEMBER: &'static str = r#"query all(
  $uid: string,
  $logged_user_id: string
){
  data(func: uid($uid)) {
    uid 
    members @filter(uid_in(member_user, $logged_user_id)) {
      uid
      forbidden_date
      last_publish_at
    }
    statuses_count
  }
}"#;

// 查询帖子信息
pub const STATUS_INFO: &'static str = r#"query all(
  $uid: string
) {
  data(func: uid($uid)) {
    uid 
    user{
      uid
    }
    group{
      uid
      statuses_count
    }
    forward_to_status{
      uid
      forwards_count
    }
    replies_count
    forwards_count
    favorites_count
  }
}"#;

// 查询帖子 -> 数量
pub const STATUS_WITH_COUNT: &'static str = r#"query all(
  $uid: string
) {
  data(func: uid($uid)) {
    uid 
    replies_count
    forwards_count
    favorites_count
  }
}"#;

// 点赞的帖子信息
pub const STATUS_OR_REPLY_FAVORITE_INFO: &'static str = r#"query all(
  $uid: string,
  $logged_user_id: string
) {
  data(func: uid($uid)) {
    uid 
    user{
      uid
    }
    cnt as count(favorites @filter(uid($logged_user_id)))
    is_favorite: math(cnt == 1)
    favorites_count
  }
}"#;

// 查询回复信息
pub const REPLY_INFO: &'static str = r#"query all(
  $uid: string
) {
  data(func: uid($uid)) {
    uid 
    user{
      uid
    }
    status{
      uid
      replies_count
      forwards_count
      favorites_count
    }
    replies_count
    forwards_count
    favorites_count
  }
}"#;


// 删除帖子或回复
pub const STATUS_OR_REPLY_DELETE_QUERY: &'static str = r#"
  query all($uid: string){
    V as var(func: uid($uid)) {               
        Entity as entities {
            Url as urls
            Media as medias
        } 
    }
  }
"#;
pub const STATUS_OR_REPLY_DELETE: &'static str = r#"
  uid(V) * * .
  uid(Entity) * * .    
  uid(Url) * * .    
  uid(Media) * * .   
"#;