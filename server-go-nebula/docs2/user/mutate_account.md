
#### 一，更新资料

```js
update vertex "145171858868801536" set \
  User.description = "我是歌手", \
  User.profile_school = "清华大学" \
  yield $^.User.name AS name, \
  $^.User.description AS description, \
  $^.User.profile_school AS profile_school;
```
