
pub const STATUS_SCHEMA: &str = "
  type Status {
    user
    text 
    is_forward
    forward_to_status
    replies
    replies_count
    forwards 
    forwards_count
    favorites
    favorites_count
    entities
    group
    
    created_at
    updated_at
  }

  type Reply {
    user
    status
    text
    is_to_reply
    to_reply
    replies
    replies_count
    entities

    reply_forwards 
    forwards_count

    reply_favorites
    favorites_count
  
    created_at
    updated_at
  }

  type Entity {
    urls
    mentions
    hashtags
    medias
    media_type
  }

  type Url {
    url
    url_key
  }

  type Hashtag {
    creater
    name
    description
    created_at
    updated_at
  }

  type Media {
    media_type
    url
    media_url
    source
  }

  user: uid @reverse .
  status: uid @reverse .
  text: string @index(fulltext) .
  is_forward: bool .
  forward_to_status: uid @reverse .
  replies: [uid] @reverse . 
  replies_count: int .
  forwards: [uid] @reverse . 
  forwards_count: int .
  favorites: [uid] @reverse .
  favorites_count: int .
  reply_forwards: [uid] @reverse . 
  reply_favorites: [uid] @reverse . 

  group: uid @reverse .
  entities: uid @reverse .
  creater: uid @reverse .
  urls: [uid] .
  mentions: [uid] .
  hashtags: [uid] .
  medias: [uid] .
  
  url: string .
  url_key: string .
  source: string .
  media_type: string .
  media_url: string .

  created_at: dateTime .
  updated_at: dateTime .

  is_to_reply: bool .
  to_reply: uid .
";
