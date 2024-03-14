use serde::{Deserialize, Serialize};

#[derive(Debug,SmartDefault,Clone,Serialize,Deserialize)]
pub enum Access {
  #[default]
  Public,       // 公开小组
  Private,      // 私人小组
  Closed        // 封闭小组
}

#[derive(Debug,SmartDefault,Clone,Serialize,Deserialize)]
pub enum Visible {
  #[default]
  Visible,       // 可被发现/搜索
  Invisible      // 不可发现
}
