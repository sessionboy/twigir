use serde::{Deserialize, Serialize};
use crate::models::statuses::entity_query::{ MediaType };

// 写入实体
#[derive(Debug, Default,Clone, Serialize, Deserialize)]
pub struct EntityInput {
    //  链接
    pub urls: Option<Vec<Url>>,
    /// 提及用户{id}
    pub mentions: Option<Vec<User>>,
    //  主题{id}
    pub hashtags: Option<Vec<Hashtags>>,
    // 帖子媒体类型
    pub media_type: Option<MediaType>,
    // 媒体: 图片、视频、音乐、直播
    pub medias: Option<Vec<Media>>
}

// 帖子的url
#[derive(Debug, Default,Clone, Serialize, Deserialize)]
pub struct Url {
    pub url: String,
    pub url_key: String
}

// 提及用户
#[derive(Debug, Default,Clone, Serialize, Deserialize)]
pub struct User {
    pub uid: String
}

// 主题
#[derive(Debug, Default,Clone, Serialize, Deserialize)]
pub struct Hashtags {
    pub uid: String
}

// 查询媒体
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct Media {
  pub media_type: MediaType,
  pub url: String,
  pub media_url: Option<String>,
  pub source: Option<String>
}

