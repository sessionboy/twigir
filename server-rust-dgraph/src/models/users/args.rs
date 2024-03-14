use serde::{Deserialize, Serialize};
use validator::{ Validate };
use crate::utils::regex::*;

// 用户登录参数
#[derive(Debug,Clone,Default,Serialize,Deserialize,Validate)]
pub struct LoginArg {
  pub username: Option<String>,
  pub phone_code: Option<String>,
  pub phone_number: Option<String>,
  #[validate(email(message = "请输入合法的邮箱地址"))]
  pub email: Option<String>,
  pub password: String,
}


// 验证手机号
#[derive(Debug,Clone,Default,Serialize,Deserialize,Validate)]
pub struct VerifyPhoneCode {
  #[validate(length(equal = 6, message = "请输入6位验证码"))]
  pub code: Option<String>,
  pub phone_code: String,
  #[validate(regex(path = "RE_PHONE", message = "请输入合法的手机号"))]
  pub phone_number: String
}

// 验证邮箱地址
#[derive(Debug,Clone,Default,Serialize,Deserialize,Validate)]
pub struct VerifyEmailCode {
  #[validate(length(equal = 6, message = "请输入6位验证码"))]
  pub code: Option<String>,
  #[validate(email(message = "请输入合法的邮箱地址"))]
  pub email: String
}

// 更新名字
#[derive(Debug,Clone,Default,Serialize,Deserialize,Validate)]
pub struct AccountName { 
  pub uid: Option<String>,
  #[validate(regex(path = "RE_NAME", message = "请输入合法的名字啊"))]
  pub name: String,
  pub updated_at: Option<String>
}

// 更新主名
#[derive(Debug,Clone,Default,Serialize,Deserialize,Validate)]
pub struct AccountUsername { 
  pub uid: Option<String>,
  #[validate(regex(path = "RE_USERNAME", message = "请输入合法的主名"))]
  pub username: String,
  pub updated_at: Option<String>
}

// 更新手机号
#[derive(Debug,Clone,Default,Serialize,Deserialize,Validate)]
pub struct AccountPhone { 
  pub uid: Option<String>,
  #[validate(regex(path = "RE_PHONE", message = "请输入合法的手机号"))]
  pub phone_number: String,
  pub phone_code: String,
  pub updated_at: Option<String>
}

// 更新邮箱
#[derive(Debug,Clone,Default,Serialize,Deserialize,Validate)]
pub struct AccountEmail { 
  pub uid: Option<String>,
  #[validate(email(message = "请输入合法的邮箱地址"))]
  pub email: String,
  pub updated_at: Option<String>
}

// 更新密码
#[derive(Debug,Clone,Default,Serialize,Deserialize,Validate)]
pub struct AccountPassword { 
  pub uid: Option<String>,
  #[validate(regex(path = "RE_PASSWORD", message = "请输入合法的密码"))]
  pub password: String,
  #[validate(regex(path = "RE_PASSWORD", message = "请输入合法的旧密码"))]
  pub old_password: String,
}
