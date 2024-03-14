use serde::{Deserialize, Serialize};

pub mod verify;

/*
  1，查询单条记录 
  2，从response的data列表中提取第一个结果 T
  3，注意：查询包装名必须是 data ，例如: data(func: uid($uid)) {...
*/ 
#[derive(Debug,Clone,Serialize,Deserialize)]
pub struct Item<T> {
  pub data: Option<Vec<T>>
}
impl<T> Item<T> {
  pub fn to_result(self) -> Option<T>{
    self.data?.into_iter().next()
  }
}

/*
  1，查询一级列表
  2，返回response的 data列表 Vec<T>
  3，注意：查询包装名必须是 data ，例如: data(func: uid($uid)) {...
*/ 
#[derive(Debug,Clone,Serialize,Deserialize)]
pub struct List<T> {
  pub data: Option<Vec<T>>
}
impl<T> List<T> {
  pub fn to_result(self) -> Vec<T>{
    if self.data.is_none() {
      let result: Vec<T> = Vec::new();
      return result;
    }
    self.data.unwrap()
  }
}


/*
  1，查询二级列表
  2，返回response的 data列表 Vec<T>
  3，注意：查询包装名必须是 data ，例如: data(func: uid($uid)) {...
*/ 
#[derive(Debug,Clone,Serialize,Deserialize)]
pub struct SubList<T> {
  pub data: Option<Vec<SubEdges<T>>>
}
impl<T> SubList<T> {
  pub fn to_result(self) -> Vec<T> {
    if self.data.is_none() {
      let result: Vec<T> = Vec::new();
      return result;
    }
    let list = self.data.unwrap();
    let result = list.into_iter().next();
    let sub_edges: SubEdges<T> = result.unwrap();
    sub_edges.get_edges()
  }
}


/*
  1，说明：二级列表查询 ，比如查询用户下的statuses列表
  2，注意：二级列表需要以"edges"作为别名，比如 
   {
     uid
     edges: statuses {
       uid
       ...
     }
   }
*/ 
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct SubEdges<T> {
  pub uid: String,
  pub edges: Option<Vec<T>>
}
impl<T> SubEdges<T> {
  pub fn get_edges(self) -> Vec<T> {
    match self.edges {
      Some(edges)=> edges,
      None=> Vec::new()
    }
  }
}