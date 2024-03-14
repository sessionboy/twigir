use serde::{Deserialize, Serialize};
use validator::{ Validate };
use crate::models::statuses::entity_input::EntityInput;
use crate::models::users::item::UserWithUid;

// 发布回复->前端传递的参数
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct CreateReplyInput { 
  pub text: String,
  // 回复某个回复
  pub to_reply: Option<String>,
  // 帖子的id
  pub status_id: String,
  // 实体
  pub entities: Option<EntityInput>,
}


// 实际写入数据库的回复
#[derive(Debug,Clone,Default,Serialize,Deserialize,Validate)]
pub struct ReplyInput { 
  // uid
  pub uid: String,
  // 创建者
  pub user: UserWithUid,
  // 贴文
  pub text: String,
  // 实体
  pub entities: Option<EntityInput>,
  // 原贴id
  pub forward_to_status: Option<String>,
  // 是否是转载
  pub is_forward: bool,
  
  pub replies_count: u32,
  pub forwards_count: u32,
  pub favorites_count: u32,

  pub created_at: String
}
