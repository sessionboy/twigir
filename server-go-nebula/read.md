
### 一，创建图空间
相关参数：
- partition_num  默认为100
- replica_factor 默认为1

##### 1，使用默认值创建
创建图空间：
```
create space default_value;
```
查看图空间：
```
describe space default_value;
```

##### 2，指定partition_num值创建
```
create space my_partition(partition_num=10);
```
```
describe space my_partition;
```

##### 3，指定replica_factor值创建
```
create space my_partition(replica_factor=3);
```
```
describe space my_partition;
```

##### 4，指定partition_num和replica_factor值创建
```
create space my_partition(partition_num=10,replica_factor=3);
```
```
describe space my_partition;
```