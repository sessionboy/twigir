use serde::{Deserialize, Serialize};
use crate::models::users::item::UserWithAccount;

// 小组的成员详情
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct Member { 
  pub uid: String,
  pub member_user: UserWithAccount,
  pub is_owner: bool,
  pub is_admin: bool,
  pub level: u32,
  pub is_anonymously: bool,
  pub alias_name: Option<String>,
  pub created_at: Option<String>
}


#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct MemberWithUid { 
  pub uid: String
}

#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct MemberWithPublish { 
  pub uid: String,
  pub forbidden_date: Option<String>,
  pub last_publish_at: Option<String>
}

// 小组的成员 - 展示
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct MemberItem { 
  pub uid: String,
  pub member_user: UserWithAccount,
  pub is_owner: bool,
  pub is_admin: bool
}


#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct MemberWithRole { 
  pub uid: String,
  pub is_owner: bool,
  pub is_admin: bool
}

// 查询是否是该小组的成员
#[derive(Debug, Default,Clone, Serialize, Deserialize)]
pub struct IsGroupMember {
    pub uid: String,
    pub members: Vec<MemberWithPublish>,
    pub statuses_count: u32
}