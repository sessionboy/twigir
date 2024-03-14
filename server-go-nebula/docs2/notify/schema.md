
#### 1，创建 tag
```
CREATE TAG Notify(
  text string,
  type string,
  target_id string,
  target_type string,
  action_type string,
  sender_id string,
  sender_type string,
  is_read bool default false,
  created_at timestamp default now(),
  updated_at timestamp
);

```