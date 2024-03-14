
#### 1，查询关注列表
```js
MATCH (v)-[:follows]->(v2), p=(v2)<-[:follows]-(me) \
WHERE id(v) == "145171858868801536" AND id(me) == "165171858868801539" \ 
RETURN v2.name AS name, v2.username AS username, id(v2) AS id, exists(p) AS following
```

```js
go from "145171858868801536" over follows yield follows._dst as id | \
go from $-.id over follows reversely yield $-.id as dst, \
case (follows._dst == '165171858868801539') when true then 1 else 0 end as following | \
group by $-.dst yield $-.dst AS target,  BIT_OR($-.following) as following | \
fetch prop on User $-.target YIELD \
  $-.following AS following, \
  User.name AS name, \
  User.username AS username;
```

```js
GO FROM "145171858868801536" OVER follows YIELD follows._dst AS target_user | \
GO FROM $-.target_user OVER follows REVERSELY \
YIELD $-.target_user AS target_user, \
CASE follows._dst == "165171858868801539" WHEN true THEN 1 ELSE 0 END AS following | \
GROUP BY $-.target_user YIELD $-.target_user AS target_user, BIT_OR($-.following) AS following
```

```js
GO FROM "145171858868801536" OVER follows YIELD follows._dst AS target_user | \
GO FROM $-.target_user OVER follows REVERSELY \
YIELD $-.target_user AS target_user, \
follows._dst == "165171858868801539" AS following | \
MATCH (v) WHERE id(v) == "165171858868801539" RETURN v.name;

FETCH PROP ON User $-.target_user YIELD \
  $-.following AS following, \
  User.name AS name, \
  User.username AS username;
```

```js
GO FROM "145171858868801536" OVER follows YIELD follows._dst AS target_user | \
GO FROM $-.target_user OVER follows REVERSELY \
YIELD $-.target_user AS target_user, \
follows._dst == "145171858868801536" AS following | \
YIELD $-.target_user AS target_user, $-.following AS following | \
FETCH PROP ON User $-.target_user YIELD $-.following, \
  User.name AS name, \
  User.username AS username | \
YIELD $-.VertexID AS id, \
  $-.name AS name, \
  $-.username AS username, \
  $-.following AS following;
```

```js
GO FROM "145171858868801536" OVER follows YIELD follows._dst AS follow_user_id | \
GO FROM $-.follow_user_id OVER follows REVERSELY \
YIELD $-.follow_user_follow_id AS follow_user_follow_id, \
follows._dst == "145171858868801536" AS following | \
FETCH PROP ON User $-.target_user YIELD \
  $-.following AS following, \
  User.name AS name, \
  User.username AS username | \
YIELD $-.VertexID AS id, \
  $-.name AS name, \
  $-.username AS username, \
  $-.following AS following;
```


#### 2，查询粉丝列表

```js
go from 1388784544961794048 over follows yield follows._dst as id | 
go from $-.id over follows reversely yield $-.id as dst, 
case (follows._dst == 1389805409983795200) when true then 1 else 0 end as following | 
group by $-.dst yield $-.dst AS target,  BIT_OR($-.following) as following
```


#### 3，查询好友列表