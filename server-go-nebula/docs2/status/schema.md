
#### 1，创建 tag
```
CREATE TAG Status(
  text string,
  lang string,
  is_group bool default false,
  status_type string default "STATUS",
  media_type string, 
  reply_count default 0,
  quote_count default 0,
  favorite_count default 0,
  forward_count default 0,
  created_at timestamp default now(),
  updated_at timestamp
);

CREATE TAG Entity(
  media_type string, 
  created_at timestamp default now()
);

CREATE TAG Url(
  url string, 
  url_key string
);

CREATE TAG Hashtag(
  name string, 
  description string, 
  created_at timestamp default now(),
  updated_at timestamp
);

CREATE TAG Media(
  url string, 
  path string, 
  media_type string
);

create edge entities(count int default 0);
create edge user(count int default 0);
create edge reply_user(count int default 0);
create edge reply_permission(count int default 0);
create edge quote_from(count int default 0);
create edge from_status(count int default 0);
create edge reply_to(count int default 0);
create edge reply_to_user(count int default 0);
create edge quote_from(count int default 0);
create edge group(count int default 0);
create edge member(count int default 0);
create tag index Status_index_text on Status(text);

create edge urls(count int default 0);
create edge mentions(count int default 0);
create edge hashtags(count int default 0);
create edge medias(count int default 0);
create tag index Hashtag_index_name on Hashtag(name);
```