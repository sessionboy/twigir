use serde::{Deserialize, Serialize};
use crate::models::users::item::{ GeneralUser, UserWithUid };
use crate::models::statuses::entity_query::Entity;
use crate::lib::dgraph::pagination::{ ExtractUid };

/*
  说明：单个帖子
  功能：帖子详情、帖子列表
*/ 
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct Status { 
  pub uid: String,
  pub user: GeneralUser,
  pub text: String,
  pub entities: Option<Entity>,
  pub forward_to_status: Option<SlimStatus>,
  pub is_forward: bool,
  
  pub is_replied: bool,
  pub replies_count: u32,

  pub is_forwarded: bool,
  pub forwards_count: u32,

  pub is_favorited: bool,
  pub favorites_count: u32,

  // 小组id
  pub group: Option<Group>,

  pub created_at: String,
  pub update_at: Option<String>
}

impl ExtractUid for Status {
  fn get_id(&self) -> String{
    self.uid.clone()
  }
}

#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct HeaderStatus { 
  pub uid: String,
  pub user: GeneralUser,
  pub text: String,
  pub entities: Option<Entity>,
  pub created_at: String
}

impl ExtractUid for HeaderStatus {
  fn get_id(&self) -> String{
    self.uid.clone()
  }
}

#[derive(Debug, Default,Clone, Serialize, Deserialize)]
pub struct StatusCreater {
    pub uid: String,
    pub last_publish_at: Option<String>,
    pub last_reply_at: Option<String>,
    pub statuses_count: u32,
}

#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct StatusInfo { 
  pub uid: String,
  pub replies_count: u32,
  pub forwards_count: u32,
  pub favorites_count: u32
}

#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct StatusWithCount { 
  pub uid: String,
  pub replies_count: u32,
  pub forwards_count: u32,
  pub favorites_count: u32
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct QueryStatus {
    pub uid: String,
    pub user: UserWithUid,
    pub group: Option<StatusGroup>,
    pub forward_to_status: Option<ForwardWithCount>,
    pub replies_count: u32,
    pub forwards_count: u32,
    pub favorites_count: u32
}
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ForwardWithCount {
    pub uid: String,
    pub forwards_count: u32
}
// 帖子所在小组
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StatusGroup {
    pub uid: String,
    pub statuses_count: u32
}

// 转发的原贴
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct SlimStatus { 
  pub uid: String,
  pub user: GeneralUser,
  pub text: String,
  pub entities: Option<Entity>,
  pub created_at: String
}

// 所属小组
#[derive(Debug, Default,Clone, Serialize, Deserialize)]
pub struct Group {
    pub uid: String,
    pub group_name: String,
    pub is_verified: bool
}

#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct StatusMedia {
  pub media_type: String
}

// 点赞的帖子或回复信息
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Favorite {
    pub uid: String,
    pub user: UserWithUid,
    pub is_favorite: bool,
    pub favorites_count: u32
}