
// 查询粉丝列表 followers
pub const USERS_FOLLOWERS: &'static str = r#"query all(
  $uid: string, 
  $first: int,
  $after: string,
  $logged_user_id: string
) {
  data(func: uid($uid)) {
    uid 
    edges: ~follows (first: $first, after:$after) {
      uid
      name
      username
      cnt as count(~follows @filter(uid($logged_user_id)))
      following: math(cnt == 1)
      description
      avatar_url
      is_verified
      verified{
        uid
        name
        description
      }
    }
  }
}"#;

// 查询关注列表 followings
pub const USERS_FOLLOWINGS: &'static str = r#"query all(
  $uid: string, 
  $first: int,
  $after: string,
  $logged_user_id: string
) {
  data(func: uid($uid)) {
    uid 
    edges: follows (first: $first, after:$after) {
      uid
      name
      username
      cnt as count(~follows @filter(uid($logged_user_id)))
      following: math(cnt == 1)
      description
      avatar_url
      is_verified
      verified{
        uid
        name
        description
      }
    }
  }
}"#;


// 1，场景：查询好友列表
// 2，语义： 查询我关注的用户 -> 的关注列表里面 —> 是否有我 ，有则是彼此关注的好友
pub const USERS_FRIENDS: &'static str = r#"query all(
  $uid: string, 
  $first: int,
  $after: string,
  $logged_user_id: string
) {
  data(func: uid($uid)) {
    uid 
    edges: follows @filter(uid_in(follows, $uid)) (first: $first, after:$after) {
      uid
      name
      username
      cnt as count(~follows @filter(uid($logged_user_id)))
      following: math(cnt == 1)
      description
      avatar_url
      is_verified
      verified{
        uid
        name
        description
      }
    }
  }
}"#;

/*
  1，场景：查询 用户1 和 用户2 的共同关注
  2，语义：查询 我(用户1)关注的用户 -> 的粉丝列表里 -> 是否有uid_two ，
     如果有则表示uid_two也关注了我的这些粉丝，这些粉丝实为我(用户1)与用户2的共同关注
  3，展示：表示为“我(uid_one)和uid_two共同关注了user1、user2、user3等人”
*/ 
pub const USERS_SAME_FOLLOWINGS: &'static str = r#"query all(
  $uid_one: string, 
  $uid_two: string, 
  $first: int,
  $after: string
) {
  data(func: uid($uid_one)) {
    uid 
    edges: follows @filter(uid_in(~follows, $uid_two)) (first: $first, after:$after) {
      uid
      name
      username
      description
      avatar_url
      is_verified
      verified{
        uid
        name
        description
      }
    }
  }
}"#;

/*
  1，语义： 查询我(uid_one)的关注里面，有哪些用户关注了uid_two
  2，场景：通常用于查看某个用户的个人主页时，查询我的关注里有谁关注了这个用户
  3，展示：表示为“我关注的user1、user2、user3等人也关注了他”
  4，简述：uid_two是 [uid_one 和 uid_one的部分关注] 的共同关注
*/ 
pub const ONE_FOLLOWINGS_TOFOLLOW_TWO: &'static str = r#"query all(
  $uid_one: string, 
  $uid_two: string, 
  $first: int,
  $after: string
) {
  data(func: uid($uid_one)) {
    uid 
    edges: follows @filter(uid_in(follows, $uid_two)) (first: $first, after:$after) {
      uid
      name
      username
      description
      avatar_url
      is_verified
      verified{
        uid
        name
        description
      }
    }
  }
}"#;


// 查询屏蔽列表 shields
pub const USERS_SHIELDS: &'static str = r#"query all(
  $uid: string, 
  $first: int,
  $after: string,
  $logged_user_id: string
) {
  data(func: uid($uid)) {
    uid 
    edges: shields (first: $first, after:$after) {
      uid
      name
      username
      cnt as count(~follows @filter(uid($logged_user_id)))
      following: math(cnt == 1)
      description
      avatar_url
      is_verified
      verified{
        uid
        name
        description
      }
    }
  }
}"#;
