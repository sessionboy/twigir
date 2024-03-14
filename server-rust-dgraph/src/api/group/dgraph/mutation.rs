
// 查询小组名称是否已存在
pub const GROUP_NAME_EXIST: &'static str = r#"query all($group_name: string) {
  data(func: eq(group_name,$group_name)) {
    uid
  }
}"#;

// 查询当前登录用户在该小组的信息
pub const GROUP_WITH_ME: &'static str = r#"query all(
  $uid: string,
  $logged_user_id: string
) {
  data(func: uid($uid)) {
    uid,
    members_count
    members_me: members @filter(uid_in(member_user,$logged_user_id)) {
      uid
      is_admin
      is_owner
    }
  }
}"#;

// 查询我和指定成员的小组信息
pub const GROUP_WITH_ME_AND_MEMBER: &'static str = r#"query all(
  $uid: string,
  $logged_user_id: string,
  $member_id: string
) {
  data(func: uid($uid)) {
    uid,
    members_count
    members_me: members @filter(uid_in(member_user,$logged_user_id)) {
      uid
      is_admin
      is_owner
    },
    members @filter(uid($member_id)) {
      uid
      is_admin
      is_owner
    }
  }
}"#;
