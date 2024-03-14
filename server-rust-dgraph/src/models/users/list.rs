use serde::{Deserialize, Serialize};
use crate::models::users::item::Follow;

// 查询关系列表
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct Follows {
  pub uid: String,
  pub edges: Option<Vec<Follow>>
}
impl Follows {
  pub fn get_edges(&self) -> Vec<Follow> {
    match &self.edges {
      Some(edges)=> edges.clone(),
      None=> Vec::new()
    }
  }
}