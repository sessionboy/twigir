use serde::{Deserialize, Serialize};

// 查询实体
#[derive(Debug, Default,Clone, Serialize, Deserialize)]
pub struct Entity {
    pub uid: String,
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

// url
#[derive(Debug, Default,Clone, Serialize, Deserialize)]
pub struct Url {
  pub uid: String,
  pub url: String,
  pub url_key: String,
}

// 提及
#[derive(Debug, Default,Clone, Serialize, Deserialize)]
pub struct User {
    pub uid: String,
    pub name: String,
    pub username: String
}

// 主题
#[derive(Debug, Default,Clone, Serialize, Deserialize)]
pub struct Hashtags {
    pub uid: String,
    pub name: String
}


// 媒体类型
#[derive(Debug,SmartDefault,Clone,Serialize,Deserialize)]
pub enum MediaType {
  #[default]
  Photo,
  Video,
  Music,
  Live,
  Vote
}

// 查询媒体
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct Media {
  pub uid: String,
  pub media_type: MediaType,
  pub url: String,
  pub media_url: Option<String>,
  pub source: Option<String>
}