use serde::{Deserialize, Serialize};
use crate::models::groups::_enum::{ Access, Visible };

// 创建小组参数
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct CreateGroupInput { 
  pub group_name: String,
  pub group_description: String,
  pub access: Access,
  pub visible: Visible
}

// 更新小组参数
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct UpdateGroupInput { 
  pub group_name: Option<String>,
  pub group_description: Option<String>,
  pub access: Option<Access>,
  pub visible: Option<Visible>,
  pub announcement: Option<String>,
  pub avatar_url: Option<String>,
  pub cover_url: Option<String>
}


