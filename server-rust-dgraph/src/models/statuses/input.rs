use serde::{Deserialize, Serialize};
use crate::models::statuses::entity_input::EntityInput;

// 发帖
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct CreateStatusInput { 
  pub text: String,
  // 转发的原贴id
  pub forward_to_status: Option<String>,
  // 帖子实体
  pub entities: Option<EntityInput>,
}
