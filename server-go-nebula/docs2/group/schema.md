
#### 1，创建 tag
```
CREATE TAG Group(
  name string,
  description string,
  avatar_url string,
  cover_url string,
  default_cover bool default true,
  access string,
  visible string,
  announcement string,
  authenticated bool default false,
  created_at timestamp default now(),
  updated_at timestamp
);

CREATE TAG Member(
  role string default "MEMBER",
  level int default 1,
  anonymously bool default false,
  alias_name string,
  forbidden_date timestamp,
  last_query_date timestamp,
  last_publish_at timestamp,
  last_reply_at timestamp,
  created_at timestamp default now(),
  updated_at timestamp
);

create edge creater(count int default 0);
create edge owner(count int default 0);
create edge joins(count int default 0);
create edge members(count int default 0);
create tag index Group_index_name on Group(name);

```