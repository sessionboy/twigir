
#### 没有生成dgraph.type导致无法删除的问题？
set schema时没有生成dgraph.type，需要在创建数据时添加dgraph.type属性，比如
```js
{
  "set":{
    "name":"jack",
    "dgraph.type":"User"
  }
}
```

### Error("invalid type: map, expected a string", line: 1, column: 37)
可能是返回数据类型错误，比如:
```js
// 本应该是 Option<Vec<Follow>>
pub struct UserFollowings {
  pub uid: String,
  pub followings: Option<Vec<Follow>>
}

// 但却写成了 Option<Follow>
pub struct UserFollowings {
  pub uid: String,
  pub followings: Option<Follow>
}
```

### 部署

```js
cd ${docker-compose.yml 目录}
docker-compose up -d // 启动
docker-compose down // 关闭
```

#### 查看日志
```js
// 日志存储路径
/var/lib/docker/containers/<container id>/<container id>-json.log
```