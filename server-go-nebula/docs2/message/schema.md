
#### 1，创建 tag
```
CREATE TAG Message(
  msg string,
  type string,
  platform string,
  created_at timestamp default now(),
  updated_at timestamp
);

create edge sender(count int default 0);
create edge to(count int default 0);
create edge msg_entities(count int default 0);
```
