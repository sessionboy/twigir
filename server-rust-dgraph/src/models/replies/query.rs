use serde::{Deserialize, Serialize};
use crate::models::users::item::UserWithUid;
use crate::models::statuses::entity_query::Entity;
use crate::models::statuses::item::SlimStatus;
use crate::models::users::item::{ GeneralUser };
use crate::lib::dgraph::pagination::{ ExtractUid };

// 回复的信息
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ReplyInfo {
    pub uid: String,
    pub user: UserWithUid,
    pub status: Option<ReplyStatus>,
    pub replies_count: u32,
    pub forwards_count: u32,
    pub favorites_count: u32
}
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ReplyStatus {
    pub uid: String,
    pub replies_count: u32,
    pub forwards_count: u32,
    pub favorites_count: u32
}


// 点赞的帖子或回复信息
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Favorite {
    pub uid: String,
    pub user: UserWithUid,
    pub is_favorite: bool,
    pub favorites_count: u32
}


// 回复详情
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct Reply { 
  pub uid: String,
  pub user: GeneralUser,
  pub text: String,
  pub entities: Option<Entity>,
  pub status: Option<SlimStatus>,
  pub is_to_reply: bool,
  
  pub is_replied: bool,
  pub replies_count: u32,

  pub is_favorited: bool,
  pub favorites_count: u32,

  pub to_reply: Option<SlimReply>,

  pub created_at: String
}

impl ExtractUid for Reply {
  fn get_id(&self) -> String{
    self.uid.clone()
  }
}

#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct HeaderReply { 
  pub uid: String,
  pub user: GeneralUser,
  pub text: String,
  pub entities: Option<Entity>,
  pub created_at: String
}

impl ExtractUid for HeaderReply {
  fn get_id(&self) -> String{
    self.uid.clone()
  }
}

// 回复的原回复
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct SlimReply { 
  pub uid: String,
  pub user: GeneralUser,
  pub text: String,
  pub entities: Option<Entity>,
  pub created_at: String
}

/*
  帖子的回复列表
*/ 
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct Replies {
  pub uid: String,
  pub edges: Option<Vec<StatusReply>>
}
impl Replies {
  pub fn get_edges(&self) -> Vec<StatusReply> {
    match &self.edges {
      Some(edges)=> edges.clone(),
      None=> Vec::new()
    }
  }
}
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct StatusReply { 
  pub uid: String,
  pub user: GeneralUser,
  pub text: String,
  pub entities: Option<Entity>,
  pub is_to_reply: bool,
  
  pub is_replied: bool,
  pub replies_count: u32,

  pub is_favorited: bool,
  pub favorites_count: u32,

  pub to_reply: Option<ToReply>,

  pub created_at: String
}
impl ExtractUid for StatusReply {
  fn get_id(&self) -> String{
    self.uid.clone()
  }
}

#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct ToReply { 
  pub uid: String,
  pub user: ToReplyUser
}
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct ToReplyUser { 
  pub uid: String,
  pub name: String,
  pub username: String
}
