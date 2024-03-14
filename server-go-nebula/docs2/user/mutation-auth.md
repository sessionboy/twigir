// LOOKUP ON User WHERE User.name != "" yield User.name;

#### 一，注册用户

1，检查名字、手机号是否已注册
```js
LOOKUP ON User WHERE User.name == "知禾";
LOOKUP ON User WHERE User.phone_number == "13267029842";
```

2，检查主名是否已存在
```js
LOOKUP ON User WHERE User.username == "jack";
```

3，注册新用户
```js
insert vertex User(
  name,
  username,
  phone_code,
  phone_number,
  password,
  lang,
  profile_birthday,
  profile_gender
) values 
  "165171858868801539":("高颖浠","gaoyingxi","+86","18680373233","12345678a","zh-CN","1992-02-12","female")
```


#### 二，用户登录

```js
// 根据username/phone_number查询用户

match (v:User{ username:"zhihe"}) return id(v) AS id, v.name AS name | \;
// or
$var = LOOKUP ON User WHERE User.phone_number == "13267029842";
// or
$var = LOOKUP ON User WHERE User.username == "zhihe" yield 
  User._vid AS id,
  User.name,
  User.username,
  User.password,
  User.lang,
  User.avatar_url,
  User.description,
  User.role,
;
GO FROM $var.id OVER authenticated yield  authenticated.name;
```