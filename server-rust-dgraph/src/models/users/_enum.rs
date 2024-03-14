use serde::{Deserialize, Serialize};

// 用户角色
#[derive(Debug,SmartDefault,Clone,Serialize,Deserialize)]
pub enum Role {
  #[default]
  User,       // 一般用户
  Admin,      // 管理员
  Super       // 超级管理员
}

// 性别
#[derive(Debug,SmartDefault,Clone,Serialize,Deserialize)]
pub enum Gender {
  #[default]
  Unknown,    // 未知
  Male,       // 男性
  Female      // 女性
}

// 感情状况
#[derive(Debug,SmartDefault,Clone,Serialize,Deserialize)]
pub enum Emotion {
  #[default]
  Single,     // 单身
  Loving,     // 恋爱中
  Married,    // 已婚
  Divorced    // 离异
}

