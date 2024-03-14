use serde::{Deserialize, Serialize};
use validator::{ Validate };
use crate::models::users::_enum::{ Gender,Role, Emotion };
use crate::utils::regex::*;

// 用户注册参数
#[derive(Debug,Clone,Default,Serialize,Deserialize,Validate)]
pub struct RegisterInput {
  #[validate(regex(path = "RE_NAME", message = "请输入合法的名字啊"))]
  pub name: String,
  #[validate(regex(path = "RE_USERNAME", message = "请输入合法的主名"))]
  pub username: String,
  #[validate(regex(path = "RE_PASSWORD", message = "请输入合法的密码"))]
  pub password: String,
  pub phone_code: String,
  #[validate(regex(path = "RE_PHONE", message = "请输入合法的手机号"))]
  pub phone_number: String,
  pub avatar_url: String,
  pub description: Option<String>,
 
  pub statuses_count: Option<u32>,
  pub replies_count: Option<u32>,
  pub role: Option<Role>,
  pub lang: Option<String>,
  pub profile_default_cover: Option<bool>,
  pub followers_count: Option<u32>,
  pub followings_count: Option<u32>,
  pub friends_count: Option<u32>,
  pub is_verified: Option<bool>,
  pub receive_chat: Option<bool>,
  pub receive_like_notification: Option<bool>,
  pub receive_reply_notification: Option<bool>,
  pub created_at: Option<String>
}

// 更新用户资料
#[derive(Debug,Clone,Default,Serialize,Deserialize,Validate)]
pub struct ProfileInput { 
  pub uid: Option<String>,
  pub avatar_url: Option<String>,
  pub description: Option<String>,

  pub profile_default_cover: Option<bool>,
  pub profile_cover_url: Option<String>,
  pub profile_gender: Option<Gender>,
  pub profile_emotion: Option<Emotion>,
  pub profile_birthday: Option<String>,
  pub profile_school: Option<String>,
  pub profile_isgraduation: Option<bool>,
  pub profile_job: Option<String>,
  pub profile_website: Option<String>,
  pub profile_country: Option<String>,
  pub profile_province: Option<String>,
  pub profile_city: Option<String>,
  pub receive_chat: Option<bool>,
  pub receive_like_notification: Option<bool>,
  pub receive_reply_notification: Option<bool>,
  pub updated_at: Option<String>
}
