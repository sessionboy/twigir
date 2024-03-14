
#### 1，关注列表
```js
// 查询 1389805409983795200 的关注列表
go from 1389805409983795200 over follows yield follows._dst as id, 
case (1 == 1) when true then true else false end as following,
$$.user.name as name,
$$.user.username as username,
$$.user.auth_name as auth_name,
$$.user.bio as bio, 
$$.user.avatar_url as avatar_url,
$$.user.verified as verified, 
$$.user.followers_count as followers_count |
ORDER BY verified DESC, followers_count DESC |
LIMIT 0, 15
```

#### 2，粉丝列表 (含关注状态)
```js
// 查询 1389805409983795200 的粉丝列表
// 1388784544961794048为登录用户的id (查询粉丝的粉丝列表中是否有我，有则表示我已关注他)
// DISTINCT去重
go from 1389805409983795200 over follows reversely yield follows._dst as id |  
go from $-.id over follows reversely yield DISTINCT
$-.id as id, 
case (follows._dst == 1388784544961794048) when true then true else false end as following,
$^.user.name as name,
$^.user.username as username,
$^.user.auth_name as auth_name,
$^.user.bio as bio, 
$^.user.avatar_url as avatar_url,
$^.user.verified as verified, 
$^.user.followers_count as followers_count |
ORDER BY verified DESC, followers_count DESC |
LIMIT 0, 15
```

#### 3，好友列表 (含关注状态)
```js
go from 1389805409983795200 over follows yield follows._dst as id |
go from $-.id over follows WHERE follows._dst == 1389805409983795200 yield 
  $-.id as id, 
  case (1 == 1) when true then true else false end as following,
  $^.user.name as name,
  $^.user.username as username,
  $^.user.auth_name as auth_name,
  $^.user.bio as bio, 
  $^.user.avatar_url as avatar_url,
  $^.user.verified as verified, 
  $^.user.followers_count as followers_count |
  ORDER BY verified DESC, followers_count DESC |
  LIMIT 0, 15
```

#### 4，共同关注 (含关注状态)
说明：查询A和B的交集即可
```js
go from 1389805409983795200 over follows  yield 
 follows._dst as id, 
 case (1 == 1) when true then true else false end as following,
 $$.user.name as name, 
 $$.user.username as username, 
 $$.user.auth_name as auth_name, 
 $$.user.bio as bio, 
 $$.user.avatar_url as avatar_url, 
 $$.user.verified as verified, 
 $$.user.followers_count as followers_count
intersect
go from 1388784544961794048 over follows  yield 
 follows._dst as id, 
 case (1 == 1) when true then true else false end as following,
 $$.user.name as name, 
 $$.user.username as username, 
 $$.user.auth_name as auth_name, 
 $$.user.bio as bio, 
 $$.user.avatar_url as avatar_url, 
 $$.user.verified as verified, 
 $$.user.followers_count as followers_count |
 ORDER BY verified DESC, followers_count DESC |
 LIMIT 0, 15
```

#### 5，我(v)关注的A、B、C...也关注了他(v2)
```js
// 我(v) : 1389805409983795200
// 他(v2): 1388784544961794048
go from 1389805409983795200 over follows yield follows._dst as id |
go from $-.id over follows WHERE follows._dst == 1388784544961794048 yield 
  $-.id as id, 
  case (1 == 1) when true then true else false end as following,
  $^.user.name as name,
  $^.user.username as username,
  $^.user.auth_name as auth_name,
  $^.user.bio as bio, 
  $^.user.avatar_url as avatar_url,
  $^.user.verified as verified, 
  $^.user.followers_count as followers_count |
  ORDER BY verified DESC, followers_count DESC |
  LIMIT 0, 15
```

#### 6，可能喜欢的人
我关注的用户的关注(排除我，以及我关注的人)
```js
// 集合的减法，关注的关注 - 关注
GO 2 STEPS FROM 1388784544961794048 OVER follows 
WHERE follows._dst != 1388784544961794048 
YIELD DISTINCT 
  follows._dst as id,
  case (1 != 1) when true then true else false end as following,
  $$.user.name AS name,
  $$.user.username as username, 
  $$.user.auth_name as auth_name, 
  $$.user.bio as bio, 
  $$.user.avatar_url as avatar_url, 
  $$.user.verified as verified, 
  $$.user.followers_count as followers_count
MINUS 
GO 1 STEPS FROM 1388784544961794048 OVER follows 
WHERE follows._dst != 1388784544961794048 
YIELD DISTINCT 
  follows._dst as id,
  case (1 != 1) when true then true else false end as following,
  $$.user.name AS name,
  $$.user.username as username, 
  $$.user.auth_name as auth_name, 
  $$.user.bio as bio, 
  $$.user.avatar_url as avatar_url, 
  $$.user.verified as verified, 
  $$.user.followers_count as followers_count |
ORDER BY verified DESC, followers_count DESC |
LIMIT 0, 15
```

#### 7，推荐的用户 (待定)
认证用户、热门用户(排除我，以及我关注的人)
```js

```

#### 8，搜索功能：用户列表
```js
match (v:user) where v.name CONTAINS '陈' 
RETURN v.name as name, id(v) as vid | 
GO FROM $-.vid OVER follows REVERSELY YIELD 
  $-.name AS name, 
  follows._dst AS dst | 
YIELD any(d IN COLLECT(DISTINCT $-.dst) WHERE d==1388784544961794048) AS d, $-.name as name
```

#### 9，用户时间线 (待定)
```js
GO FROM 1388784544961794048 OVER follows YIELD follows._dst AS id |
GO FROM $-.id OVER statuses WHERE 
  $-.id > 1389805409983795200
YIELD 
  $$.status.text AS text
UNION
GO FROM $-.id OVER likes WHERE 
  $-.id > 1389805409983795200
YIELD 
  $$.status.text AS text
```


#### 10，相关用户 GetUserByIds
根据 ids 列表查询相关用户，比如贴文详情右侧的相关用户[贴文作者、提及用户]

#### 11，照片列表
用户主页右侧的照片列表

#### 12，用户贴文列表

#### 13，通知列表

#### 14，搜索功能：帖子列表