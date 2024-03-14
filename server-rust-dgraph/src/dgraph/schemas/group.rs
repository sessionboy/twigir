
pub const GROUP_SCHEMA: &str = "
  type Group {
    group_creater
    group_name
    group_description
    group_entities
    announcement
    avatar_url
    cover_url
    default_cover
    access
    visible
    is_verified
    verified    
    owner
    members
    members_count
    statuses
    statuses_count
    
    created_at
    updated_at
  }

  type Member {
    member_user
    is_owner
    is_admin
    level
    is_anonymously
    alias_name
    forbidden_date
    last_query_date
    last_publish_at
    last_reply_at
    
    created_at
    updated_at
  }

  group_creater: uid .
  group_name: string @index(exact) .
  group_description: string @index(exact) .
  group_entities: uid .
  avatar_url: string .
  cover_url: string .
  default_cover: bool .
  access: string .
  visible: string .
  announcement: string .
  category: uid @reverse .
  is_verified: bool .
  verified: uid @reverse .
  owner: uid @reverse .
  members: [uid] @reverse .
  members_count: int .

  member_user: uid @reverse .
  is_owner: bool .
  is_admin: bool .
  is_anonymously: bool .
  level: int .
  alias_name: string .
  forbidden_date: dateTime .
  last_query_date: dateTime .
";
