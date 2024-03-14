
### Vertex 点

基于Tag结构创建的节点(node)。

#### 一，创建Vertex
语法：
```
insert vertex tag_name(field1,field2) values ID:(value1, value2), ID:(value1, value2), ...;
```
例子1，插入一个user节点：
```
insert vertex User(name,age) values 100:("jack", 22);
```

例子2，插入多个user节点：
```
insert vertex User(name,age) values 100:("jack", 22), 101:("zhihe", 20), 103:("panda", 12);
```

#### 二，查询点Vertex

vid为Vertex点的id。

```
fetch prop on tag_name vertex_id, vid, vid, ...;
```

例子1，根据id=100查询单个用户：
```
fetch prop on User 100;
```

例子2，查询多个用户：
```
fetch prop on User 100, 200, 101, 300;
```


#### 三，删除点Vertex
```
delete vertex 100;
```