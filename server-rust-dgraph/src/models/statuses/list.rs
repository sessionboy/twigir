use serde::{Deserialize, Serialize};
use crate::models::statuses::item::Status;

/*
  说明：帖子列表 -> 二级查询
  功能：用户的帖子列表、小组的帖子列表
  注意： dgraph查询中，需将列表命名为：edges， 比如：(edges: statuses ...)
*/ 
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct Statuses {
  pub uid: String,
  pub edges: Option<Vec<Status>>
}
impl Statuses {
  pub fn get_edges(&self) -> Vec<Status> {
    match &self.edges {
      Some(edges)=> edges.clone(),
      None=> Vec::new()
    }
  }
}