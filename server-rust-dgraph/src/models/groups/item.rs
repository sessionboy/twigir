use serde::{Deserialize, Serialize};
use crate::models::commons::verify::{ SlimVerified };
use crate::models::groups::member::{ MemberWithRole, MemberItem };
use crate::models::groups::_enum::{ Access, Visible };
use crate::models::statuses::entity_query::Entity;
use crate::lib::dgraph::pagination::{ ExtractUid };

// 小组详情
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct Group { 
  pub uid: String,
  pub group_name: String,
  pub group_description: Option<String>,
  pub group_entities: Option<Entity>,
  pub access: Access,
  pub visible: Visible,
  pub announcement: Option<String>,
  pub avatar_url: Option<String>,
  pub cover_url: Option<String>,
  pub default_cover: bool,
  // 登录用户是否已加入该小组
  pub is_joined: Option<bool>,
  pub is_verified: bool,
  pub verified: Option<SlimVerified>,
  pub members_count: u32,
  pub statuses_count: u32,
  pub created_at: String,
  pub updated_at: Option<String>
}
impl ExtractUid for Group {
  fn get_id(&self) -> String{
    self.uid.clone()
  }
}

// 可用于通知等场景
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct HeaderGroup { 
  pub uid: String,
  pub group_name: String,
  pub is_verified: bool,
  pub created_at: String
}
impl ExtractUid for HeaderGroup {
  fn get_id(&self) -> String{
    self.uid.clone()
  }
}

// 小组列表的小组信息
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct GroupItem { 
  pub uid: String,
  pub group_name: String,
  pub group_description: Option<String>,
  pub group_entities: Option<Entity>,  
  pub access: Access,
  pub avatar_url: Option<String>,
  // 登录用户是否已加入该小组
  pub is_joined: Option<bool>,
  // 登录用户在该小组的成员信息
  pub member_self: Option<Vec<MemberWithRole>>,
  pub is_verified: bool,
  pub verified: Option<SlimVerified>,
  pub members_count: u32,
  pub statuses_count: u32,
  pub created_at: String
}
impl ExtractUid for GroupItem {
  fn get_id(&self) -> String{
    self.uid.clone()
  }
}

// 用户的小组列表
#[derive(Debug,Clone,Serialize,Deserialize)]
pub struct Groups { 
  pub uid: String,
  pub edges: Option<Vec<GroupItem>>
}
impl Groups {
  pub fn get_edges(&self) -> Vec<GroupItem> {
    match &self.edges {
      Some(edges)=> edges.clone(),
      None=> Vec::new()
    }
  }
}

// 所属小组
#[derive(Debug, Default,Clone, Serialize, Deserialize)]
pub struct GroupWithName {
    pub uid: String,
    pub group_name: String
}

// 查询小组名称是否已存在 
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct GroupWithUid { 
  pub uid: String
}


#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct GroupWithMe { 
  pub uid: String,
  pub members_me: Option<Vec<MemberWithRole>>,
  pub members_count: u32
}

#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct GroupWithMeAndMember { 
  pub uid: String,
  pub members_me: Option<Vec<MemberWithRole>>,
  pub members: Option<Vec<MemberWithRole>>,
  pub members_count: u32
}