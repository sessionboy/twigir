
### Edge 边类型

edge相当于关联字段，比如follow等。

#### 一，创建Edge
语法：
```
create edge edge_name(field_name field_type,field_name field_type,...);
```

例子，创建follow边：
```
create edge  follow(dgreee int);
```

也可以创建没有任何属性的Edge类型：
```
create edge edge_name();
```


#### 二，查看Edge

1，查看所有Edge
```
show edges;
```

1，查看指定Edge
```
describe edge edge_name;
```
例如:
```
describe edge follow;
```

#### 三，删除Edge
```
drop edge edge_name;
```