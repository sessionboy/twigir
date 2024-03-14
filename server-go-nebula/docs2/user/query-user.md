
#### 1，查询个人资料
```js
fetch prop on User "145171858868801536" yield \
  User.name AS name, \
  User.username AS username, \
  User.password AS password, \
  User.role AS role, \
  User.lang AS lang, \
  User.description AS description, \
  User.avatar_url AS avatar_url, \
  User.profile_cover_url AS profile_cover_url, \
  User.profile_default_cover AS profile_default_cover, \
  User.profile_school AS profile_school, \
  User.profile_gender AS profile_gender, \
  User.profile_birthday AS profile_birthday, \
  User.profile_website AS profile_website, \
  User.profile_emotion AS profile_emotion, \
  User.followers_count AS followers_count, \
  User.followings_count AS followings_count, \
  User.friends_count AS friends_count, \
  User.created_at AS created_at;
```