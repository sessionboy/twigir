use serde::{Deserialize, Serialize};
use crate::models::commons::verify::{ SlimVerified };
use crate::lib::dgraph::pagination::{ ExtractUid };

#[derive(Debug, Default, Clone, Serialize, Deserialize)]
pub struct UserWithUid {
  pub uid: String
}

#[derive(Debug, Clone, Serialize, Default, Deserialize)]
pub struct UserWithAccount {
  pub uid: String,
  pub name: String,
  pub username: String
}

// 用户个人主页
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct User {
  pub uid: String,
  pub name: String,
  pub username: String,
  pub description: Option<String>,
  pub lang: String,
  // "我"是否已关注该用户
  pub following: Option<bool>,
  pub avatar_url: Option<String>,
  pub profile_cover_url: Option<String>,
  pub profile_default_cover: Option<bool>,
  pub profile_school: Option<String>,
  pub profile_gender: Option<String>,
  pub profile_birthday: Option<String>,
  pub profile_website: Option<String>,
  pub profile_emotion: Option<String>,

  pub is_verified: Option<bool>,
  pub verified: Option<SlimVerified>,
  
  pub followers_count: u32,
  pub followings_count: u32,
  pub friends_count: u32,

  pub created_at: String,
}

// 普通用户信息，用于展示，比如用户列表、帖子的作者等
#[derive(Debug, Default, Serialize, Deserialize, Clone)]
pub struct GeneralUser {
    pub uid: String,
    pub name: String,
    pub username: String,
    pub avatar_url: Option<String>,
    pub description: Option<String>,
    pub is_verified: bool,
    pub verified: Option<SlimVerified>,
    // “我(当前登录用户)” 是否已关注该用户
    pub following: Option<bool>
}


#[derive(Debug, Default, Serialize, Deserialize, Clone)]
pub struct HeaderUser {
    pub uid: String,
    pub name: String,
    pub username: String,
    pub avatar_url: Option<String>,
    pub is_verified: bool
}

// 更新密码接口，查询用户的密码
#[derive(Debug,Serialize, Deserialize)]
pub struct UserPassword {
  pub user: Vec<Password>
}
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct Password { 
  pub uid: String,
  pub password: String,
  pub updated_at: Option<String>
}


/*
  功能： 粉丝/关注/共同关注等
*/ 
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct Follow {
  pub uid: String,
  pub name: String,
  pub username: String,
  // (我:登录用户)是否已关注该用户
  pub following: Option<bool>,
  pub description: Option<String>,
  pub avatar_url: Option<String>,
  pub is_verified: Option<bool>,
  pub verified: Option<SlimVerified>,
}
impl ExtractUid for Follow {
  fn get_id(&self) -> String{
    self.uid.clone()
  }
}

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