1，创建三个球员
insert vertex player(name,age) values "145171858868801536":("陈超",25),"145171858861801538":("郭超",21),"145171858811801578":("林超",24);

2，创建2个球队
insert vertex team(name) values "145171858868801333":("火箭部落"),"145171858868801444":("归心俱乐部");

3，创建球员与球队的边 (球员属于哪个球队的)
insert edge follow(dgree) values "145171858868801536" -> "145171858861801538":(95),"145171858868801536" -> "145171858811801578":(90),"145171858811801578" -> "145171858861801538":(93);

3，球员与球队的合同服务 边
insert edge serve(start_year, end_year) values "145171858868801536" -> "145171858868801333":(1997, 2016),"145171858861801538" -> "145171858868801444":(1999,  2018);


二，查询

yield: 指定要返回的结果
where: 指定筛选条件
$$: 表示目标顶点
$-: 在管道符号之前代表查询的输出
$^: 表示边的源顶点

1，查询球员 145171858868801536(陈超) 所关注的球员

先通过go查id，再通过fetch查具体的数据
```js
go from "145171858868801536" over follow yield follow._dst AS id | fetch prop on player $-.id;
```

2，查询球员 145171858868801536(陈超) 所关注的年龄大于23的球员

通过where进行筛选。
```js
go from "145171858868801536" over follow where $$.player.age >= 23 \ 
yield follow._dst AS id | fetch prop on player $-.id;
```

3，使用$var变量查询
```js
$var = go from "145171858868801536" over follow yield follow._dst AS id;  
fetch prop on player $var.id;
```

4，多级查询
  查询球员 145171858868801536(陈超) 所关注的球员
  查询这些球员所服务的球队
```js
go from "145171858868801536" over follow yield follow._dst AS id | \ 
go from $-.id over serve yield $$.team.name AS Team;
```

三，创建索引

1，创建单索引
```js
create tag index player_index_0 on player(name(10)); // 取name属性的前10个字符
```

2，创建复合索引
多字段唯一，例如 name+username 作为索引依据
```js
create tag index player_index_1 on player(name(10), username);
```

3，重建索引
索引新建或修改后，可以重建索引，使索引立即生效。
```js
REBUILD TAG INDEX player_index_0;
```

4，使用索引 
```js
LOOKUP ON player WHERE player.name == "Tony Parker";
```

四，变更(更新、删除)
