
### 插入边 Edge

基于Edge结构创建边。

#### 一，插入边Edge
语法：
```
insert edge edge_name(field1,field2 ...) values vid1 -> vid2:(field1_value, field2_value ...);
```

例子1，插入一个follow边：
```
// 用户100 关注了 用户102
insert edge follow(dgreee) values 100 -> 102:(22);
```

例子2，插入多个follow边：
```
insert edge follow(dgreee) values 100 -> 102:(42), 100 -> 103:(52), 102 -> 103:(12), ...;
```

在两个点之间插入多条边，需要在终边指定"@ranking"值:
```
insert edge edge_name(field1,field2 ...) values vid1 -> vid2@2:(field1_value, field2_value ...);
```

#### 二，查询边Edge

```
fetch prop on edge_name vid1 -> vid2, vid1 -> vid2, ...;
```
将返回边的属性。

例子1，查询单条边：
```
fetch prop on follow 100 -> 102;
```

例子2，查询多条边：
```
fetch prop on follow 100 -> 102, 100 -> 103, 100 -> 104;
```


#### 三，删除边Edge
```
delete edge follow 100 -> 102;
```