use serde::{Deserialize, Serialize};
use crate::models::users::item::UserWithUid;

#[derive(Debug, Serialize, Deserialize)]
pub struct UserExist {
  pub name: Vec<UserWithUid>,
  pub username: Vec<UserWithUid>,
  pub phone: Vec<UserWithUid>
}

// 名字是否存在
#[derive(Debug, Serialize, Deserialize)]
pub struct ExistName {
  pub name: Vec<UserWithUid>
}

// 主名是否存在
#[derive(Debug, Serialize, Deserialize)]
pub struct ExistUsername {
  pub username: Vec<UserWithUid>
}

// 手机号是否存在
#[derive(Debug, Serialize, Deserialize)]
pub struct ExistPhone {
  pub phone: Vec<UserWithUid>
}

// 邮箱是否存在
#[derive(Debug, Serialize, Deserialize)]
pub struct ExistEmail {
  pub email: Vec<UserWithUid>
}