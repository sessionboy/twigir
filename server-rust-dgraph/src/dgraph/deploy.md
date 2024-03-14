
## 1，dgraph的安装

安装地址：
```
https://dgraph.io/docs/deploy/download/
```

## 2，dgraph基本配置
#### (1) 配置文件：
```js
// config.json
{
  "my": "localhost:7080",    
  "zero": "localhost:5080", 
  "lru_mb": 4096, 
  "postings": "/path/to/p",  // 数据存储的目录路径
  "wal": "/path/to/w"  // 预写日志记录的目录路径
}
```
应用该配置文件：
```
dgraph zero --config ./config.json
```

## 使用docker-compose部署

#### (1) 官方部署文件：
```
https://github.com/dgraph-io/dgraph/blob/master/contrib/config/docker/docker-compose.yml
```

#### (2) 注意设置白名单
在白名单内的ip客户端可以执行数据库的管理操作，比如设置schema等等。
意味着如果不设置白名单，你将无法在开发中设置schema。

白名单设置：
```
// 可以指定ip。这里设置为0.0.0.0/0表示允许任意ip操作

dgraph alpha --whitelist 0.0.0.0/0 --my=alpha:7080 --lru_mb=2048 --zero=zero:5080
```


## ratel 的使用

1，ratel访问地址：ip:8000

2，ratel需要连接alpha地址：ip:8080