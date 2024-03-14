$var = lookup on User WHERE User.name == "知禾"; \
FETCH PROP ON User $var._vid;

#### 1，创建 User tag
```
CREATE TAG User(
  name string,
  username string,
  phone_code string,
  phone_number string,
  password string,
  email string,
  lang string,
  avatar_url string,
  description string,
  role string default "USER",
  location string,
  verified bool default false,
  authenticated bool default false,
  authenticated_name string,

  profile_cover_url string,
  profile_default_cover bool default true,
  profile_gender string,
  profile_birthday string,
  profile_school string,
  profile_isgraduation string,
  profile_job string,
  profile_website string,
  profile_emotion string,
  profile_country string,
  profile_province string,
  profile_city string,

  statuses_count int default 0,
  followers_count int default 0,
  followings_count int default 0,
  friends_count int default 0,
  last_sign_at timestamp,
  last_publish_at timestamp,
  last_reply_at timestamp,
  created_at timestamp default now(),
  updated_at timestamp
);

CREATE TAG Authenticate(
  type string,
  name string,
  level int default 1,
  description string,
  created_at timestamp default now(),
  updated_at timestamp
);

CREATE TAG Record(
  platform string,
  ip string,
  device string,
  browser string,
  is_register bool default false,
  created_at timestamp
);

create edge entities(count int default 0);
create edge authenticated(count int default 0);
create edge follows(count int default 0);
create edge statuses(count int default 0);
create edge blacklists(count int default 0);
create edge desc_entity(count int default 0);

create tag index User_index_name on User(name(20));
create tag index User_index_username on User(username(20));
create tag index User_index_phone_number on User(phone_number(50));
create tag index Authenticate_index_name on Authenticate(name(20));
```

```
GO FROM "145171858868801536" OVER follows, statuses YIELD follows.count, statuses.count;
```