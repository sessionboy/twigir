
### Tag(标签)

相当于schema中的type。

#### 一，创建Tag
语法：
```
create tag tag_name(
  field_name field_type,
  field_name field_type,
  field_name field_type,
  ...
 );
```

例子，创建User标签：
```
create tag User(name string, age int, password string);
```

也可以创建没有添加属性的空标签：
```
create tag tag_name();
```


#### 二，查看标签

1，查看所有标签
```
show tags;
```

1，查看指定标签
```
describe tag tag_name;
```
例如:
```
describe tag User;
```

#### 三，删除标签
```
drop tag tag_name;
```