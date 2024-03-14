use serde::{Deserialize, Serialize};
use crate::models::users::item::{ HeaderUser };
use crate::models::statuses::item::{ HeaderStatus };
use crate::models::replies::item::{ HeaderReply };
use crate::models::groups::item::{ HeaderGroup };
use crate::models::notifications::_enum::{ NotificationType };
use crate::lib::dgraph::pagination::{ ExtractUid };

// 通知详情
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct Message { 
  pub uid: String,
  pub notification_type: NotificationType,
  pub message: Option<String>,

  pub sender: HeaderUser,
  pub receiver: Option<Vec<HeaderUser>>,
  pub target: Option<HeaderUser>,

  pub group: Option<HeaderGroup>,
  pub status: Option<HeaderStatus>,
  pub reply: Option<HeaderReply>,
  pub user: Option<HeaderUser>,

  pub created_at: String
}
impl ExtractUid for Notification {
  fn get_id(&self) -> String{
    self.uid.clone()
  }
}
