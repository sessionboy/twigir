
#### 1，白名单问题
```js
rpc error: code = Unknown desc = unauthorized ip address: 223.73.108.241
```
dgraph默认仅允许本机执行schema alter，本机以外的ip没有权限。需要设置白名单才可以拥有权限。
```js
dgraph alpha --security whitelist=127.0.0.1 ...

dgraph alpha --security whitelist=172.17.0.0:172.20.0.0,192.168.1.1 ...

// 允许所有ip访问
dgraph alpha --security whitelist=0.0.0.0/0 ...
```

官方文档：https://dgraph.io/docs/deploy/dgraph-administration/#whitelisting-admin-operations