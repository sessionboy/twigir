use std::fmt::Debug;
use crate::config::constants;

// 分页参数
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct PaginationArgs {
  pub first: Option<u32>,
  pub after: Option<String>
}
impl PaginationArgs {
  pub fn get_first(&self) -> u32 {
    match self.first.clone() {
      Some(first) => first,
      None => constants::PAGINATION_FIRST
    }
  }
  // first+1
  pub fn get_add_first(&self) -> u32 {
    match self.first.clone() {
      Some(first) => {
        first + 1
      },
      None => constants::PAGINATION_FIRST + 1
    }
  }
  pub fn get_add_first_str(&self) -> String {
    self.get_add_first().to_string()
  }
  // after
  pub fn get_after(&self) -> String {
    match self.after.clone() {
      Some(after) => {
        if after.is_empty() {
          return String::from("0");
        }
        after
      },
      None => String::from("0")
    }
  }
}

// 分页返回数据结构
#[derive(Debug, Default, Serialize, Deserialize, Clone)]
pub struct PageInfo {
    // 用endCursor去查,查询创建日期有没有大于endCursor的
    pub hasNextPage: bool,
    // after存在则为true，否则为false
    pub hasPrevPage: bool,
    // edges[0]
    pub startCursor: Option<String>,
    // edges.length - 1
    pub endCursor: Option<String>
}
#[derive(Debug, Default, Serialize, Deserialize, Clone)]
pub struct PaginationResult<T> {
    pub count: u32,  
    pub edges: Vec<T>,
    pub pageInfo: PageInfo
}

// 分页参数解析 constants::PAGINATION_FIRST;
// 注意，查询时必须让 first+1 ，多查一条用于判断
// 这里的first是原始的first，没有+1

pub trait ExtractUid {
  fn get_id(&self) -> String;
}
#[derive(Debug)]
pub struct Expand {
  uid: String,
}
impl ExtractUid for Expand {
  fn get_id(&self) -> String {
      self.uid.clone()
  }
}

pub fn pagination<T:Debug+Clone+ExtractUid>(
  edges: &Vec<T>, 
  first: u32,
  after: String
)-> PaginationResult<T> {
  let _len = edges.len() as u32;
  let has_prev_page:bool = !after.eq(&"0".to_owned()) && !after.is_empty(); 
  let has_next_page:bool = _len > first;
  let mut start_cursor: Option<String> = None;
  let mut end_cursor: Option<String> = None;
  let end = match _len > 0 {
    false => 0,
    true => edges.len()-1
  };

  let _edges = match has_next_page {
    false => &edges,
    true => &edges[0..end]
  };

  if !_edges.is_empty() {
    start_cursor = Some(_edges.first().unwrap().get_id());
    end_cursor = Some(_edges.last().unwrap().get_id());
  }
  let count:u32 = _edges.len() as u32;
 
  let page_info = PageInfo {
    hasNextPage: has_next_page,
    hasPrevPage: has_prev_page,
    startCursor: start_cursor,
    endCursor: end_cursor
  };

  PaginationResult{
    count,
    edges: _edges.to_vec(),
    pageInfo: page_info
  }
  
}
