# XORM 功能更改记录
 
 1. 细化日志开关：支持分别对显示SQL、执行时间、缓存、事件、基础日记的开关控制
 2. 将日志中SQL的“?”占位符替换为实际的值，也就是日志中的SQL是完整的生成后的语句
 3. SQL执行时间日志中的耗时统一转为秒作为单位
 4. 改进表别名的生成规则:
 ```
    对于采用添加“extends”标签进行关联查询的操作，
    在不指定表别名的情况下默认采用结构体字段名作为别名。
 ```
例如（在这里我们的表前缀为“webx_”,下同）：
```go
type PostCollection struct {
        A     *D.Post     `xorm:"extends"`
        B     *D.User     `xorm:"extends"`
}
ms := []*PostCollection{}
engine.Where(`A.id=1`).Join(`LEFT`, `webx_user`, `A.uid=B.id`).Find(&ms)
```
会生成SQL：
```sql
SELECT `A`.`id`,`A`.`uid`,`B`.`id`,`B`.`passwd`
 FROM `webx_post` AS `A` 
 LEFT JOIN `webx_user` AS `B` ON A.uid=B.id WHERE A.id=1
```
再例如：
```go
type PostMore struct {
        Id    int64         `xorm:"not null pk autoincr INT(20)"`
        Uid   int64         `xorm:"not null default 0 INT(20)"`
        B     *D.User       `xorm:"extends"`
}
m := []*PostMore{}
engine.Where(`PostMore.id=1`).Join(`LEFT`, `webx_user`, `PostMore.uid=B.id`).Find(&m)
```
会生成SQL：
```sql
SELECT `PostMore`.`id`,`PostMore`.`uid`,`B`.`id`,`B`.`passwd`
 FROM `webx_post_more` AS `PostMore`
 LEFT JOIN `webx_user` AS `B` ON PostMore.uid=B.id WHERE PostMore.id=1
```
带关联条件标签的关联查询：
```go
type PostCollection2 struct {
	Post     *D.Post     `xorm:"extends" rel:""`
	User     *D.User     `xorm:"extends" rel:"LEFT:User.id=Post.uid"`
}
ms3 := []*PostCollection2{}
engine.Omit(`User.passwd`).Where(`Post.id=2`).Find(&ms3)
```
会生成SQL：
```sql
SELECT `Post`.`id`,`Post`.`uid`,`User`.`id` 
FROM `webx_post` AS `Post` 
LEFT JOIN `webx_user` AS `User` ON User.id=Post.uid WHERE Post.id=2
```

```
如果想要使用自己的别名，可以在结构体字段中增加tag(alias)来实现。比如：
type PostCollection struct {
        *D.Post     `xorm:"extends" alias:"a"`
        *D.User     `xorm:"extends" alias:"b"`
}
```


 5. 增加大量原生SQL执行接口(详见coscms.go)。
 6. 增加查询结果为`[]map[string]string`和`[]map[string]interface{}`类型的数据
    