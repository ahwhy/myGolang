# Golang-Database  Golang的数据库

## 一、SQL语法简介

### 1. MySQL初始配置
- SQL(i/ˈsiːkwəl/; Structured Query Language)是一套语法标准，不区分大小写

- MySQL、sql-server和Oracle都是关系型数据库，在一些高级语法上跟标准SQL略有出入

- MySQL事务
	- 事务就是由单独单元的⼀个或多个sql语句组成，在这个单元中，每个sql语句都是相互依赖的；⽽整个单独单元是作为⼀个不可分割的整体存在，类似于物理当中的原⼦(⼀种不可分割的最⼩单位)
		- 事务就是⼀个整体，⾥⾯的内容要么都执⾏成功，要么都不成功
		- 不存在部分执⾏成功，⽽部分执⾏不成功的情况
		- 如果单元中某条sql语句⼀旦执⾏失败或者产⽣错误，那么整个单元将会回滚(返回最初状态)，所有受到影响的数据将返回到事务开始之前的状态
		- 如果单元中的所有sql语句都执⾏成功的话，那么该事务也就被顺利执⾏
	- 事务的四个特性(ACID)
		- 原⼦性(Atomicity)：指事务是⼀个不可分割的最⼩⼯作单位，事务中的操作只有都发⽣和都不发⽣两种情况
		- ⼀致性(Consistency)：事务必须使数据库从⼀个⼀致状态变换到另外⼀个⼀致状态
		- 隔离性(Isolation)：⼀个事务的执⾏不能被其他事务⼲扰，即⼀个事务内部的操作及使⽤的数据对并发的其他事务是隔离的，并发执⾏的各个事务之间不能互相⼲扰
		- 持久性(Durability)：⼀个事务⼀旦提交成功，它对数据库中数据的改变将是永久性的，接下来的其他操作或故障不应对其有任何影响
	- 事物的分类
		- 隐式事务：该事务没有明显的开启和结束标记，它们都具有⾃动提交事务的功能，如DML语句(insert、update、delete)等
		- 显式事务：该事务具有明显的开启和结束标记
	- 事务并发时出现的问题
		- 因为某⼀刻不可能总只有⼀个事务在运⾏，可能出现A在操作t_account表中的数据，B也同样在操作t_account表，那么就会出现并发问题，对于同时运⾏的多个事务，当这些事务访问数据库中相同的数据时，如果没有采⽤必要的隔离机制，就会发⽣以下各种并发问题
			- 脏读：对于两个事务T1,T2，T1读取了已经被T2更新但还没有被提交的字段之后，若T2回滚，T1读取的内容就是临时且⽆效的
			- 不可重复读：对于两个事务T1,T2，T1读取了⼀个字段，然后T2更新了该字段之后，T1在读取同⼀个字段，值就不同了
			- 幻读：对于两个事务T1,T2，T1在A表中读取了⼀个字段，然后T2⼜在A表中插⼊了⼀些新的数据时，T1再读取该表时，就会发现多出⼏⾏了
	- 事物的隔离级别：
		- RU：read uncommitted(读未提交数据)，允许事务读取未被其他事务提交的变更。(脏读、不可重复读和幻读的问题都会出现)。
		- RC：read committed(读已提交数据)，只允许事务读取已经被其他事务提交的变更。(可以避免脏读，但不可重复读和幻读的问题仍然可能出现)
		- RR：repeatable read(可重复读)：确保事务可以多次从⼀个字段中读取相同的值，在这个事务持续期间，禁⽌其他事务对这个字段进⾏更新(update)。(可以避免脏读和不可重复读，但幻读仍然存在，Mysql默认的事务隔离级别)
		- Serializable(串⾏化)：确保事务可以从⼀个表中读取相同的⾏，在这个事务持续期间，禁⽌其他事务对该表执⾏插⼊、更新和删除操作，所有并发问题都可避免，但性能⼗分低下(因为不完成就都不可以进行下一步，效率太低)
	- [MySQL事务(transaction)](https://blog.csdn.net/qq_56880706/article/details/122653735)

- MVCC(multi-version-concurrent-control)
	- MVCC，即多版本并发控制，MVCC是⼀种并发控制的⽅法，⼀般在 数据库管理系统 中，实现对数据库的并发访问，在编程语⾔中实现事务内存
		- MVCC 在 MySQL InnoDB 中的实现主要是为了提⾼数据库的并发性能，⽤更好的⽅式去处理读-写冲突，做到**即使有读写冲突时，也能做到不加锁，⾮阻塞并发读**
		- 快照读
			- 像不加锁的 select 操作就是快照读，即不加锁的⾮阻塞读；
			- 快照读的前提是隔离级别不是串⾏级别，串⾏级别下的快照读会退化成当前读；
			- 之所以出现快照读的情况，是基于提⾼并发性能的考虑，快照读的实现是基于多版本并发控制
		- 当前读
			- select lock in share mode（共享锁）
			- select for update ； update，insert，delete（排他锁）
			- 这些操作都是⼀种当前读，它读取的记录都是⽬前数据库中最新的版本，读取时还要保证其它并发事务不能修改当前记录，所以会对读取数据加锁
	- 关于幻读的补充
		- 官⽅定义：当同⼀个查询在不同的时间产⽣不同的结果集时，事务中就会出现所谓的幻象问题
			- 例如，如果 SELECT 执⾏了两次，但第⼆次返回了第⼀次没有返回的⾏，则该⾏是“幻像”⾏
		- 针对MySQL的幻读，在默认隔离级别是 可重读RR 下，有以下两种解决⽅案：
			- 快照读，通过MVCC ⽅式解决了幻读，在默认隔离级别可重复读下，⼀个事务启动后 快照读的数据和 启动时刻 所有的数据是⼀致的，即使中途有其他事务插⼊了⼀条数据，也查不出来这⼀条数据，所有避免了幻读问题
			- 当前读，虽然上⾯说了幻读发⽣在当前读，但MySQL 通过next-key lock（记录锁+间隙锁）解决了，当执⾏select … for update 语句的时候，会加上next-key lock （在⼀个区间内加锁），如果有其它事务在这个范围内插⼊⼀条数据，那么这个插⼊就会被阻塞，⽆法插⼊，所以就避免了幻读问题
		- [MySQL 幻读](https://blog.csdn.net/m0_51809035/article/details/127565218)

- MySQL存储引擎
	- `show engines;` 查看支持的存储引擎
	- InnoDB 支持事务，具有提交，回滚和崩溃恢复能力，事务安全，适合大量insert或update操作
	- MyISAM 不支持事务和外键，访问速度快，提供高速存储和检索，适合大量的select查询操作
	- Memory 利用内存创建表，访问速度非常快，因为数据在内存，而且默认使用Hash索引，但是一旦关闭，数据就会丢失
	- Archive 归档类型引擎，仅能支持insert和select语句
	- Csv 以CSV文件进行数据存储，由于文件限制，所有列必须强制指定not null，另外CSV引擎也不支持索引和分区，适合做数据交换的中间表
	- BlackHole 黑洞，只进不出，进来消失，所有插入数据都不会保存
	- Federated 可以访问远端MySQL数据库中的表。一个本地表，不保存数据，访问远程表内容
	- MRG_MyISAM 一组MyISAM表的组合，这些MyISAM表必须结构相同，Merge表本身没有数据，对Merge操作可以对一组MyISAM表进行操作
	- [MySQL - 存储引擎MyISAM和Innodb](https://blog.csdn.net/weixin_42201180/article/details/125696197)

- InnoDB
	- InnoDB 的数据是按「数据⻚」为单位来读写的，默认数据⻚⼤⼩为 16 KB；每个数据⻚之间通过双向链表的形式组织起来，物理上不连续，但是逻辑上连续
	- 数据⻚内包含⽤户记录，每个记录之间⽤单项链表的⽅式组织起来，为了加快在数据⻚内⾼效查询记录，设计了⼀个⻚⽬录，⻚⽬录存储各个槽(分组)，且主键值是有序的，于是可以通过⼆分查找法的⽅式进⾏检索从⽽提⾼效率
	- 为了⾼效查询记录所在的数据⻚，InnoDB 采⽤ b+ 树作为索引，每个节点都是⼀个数据⻚；如果叶⼦节点存储的是实际数据的就是聚簇索引，⼀个表只能有⼀个聚簇索引；如果叶⼦节点存储的不是实际数据，⽽是主键值则就是⼆级索引，⼀个表中可以有多个⼆级索引
	- 在使⽤⼆级索引进⾏查找数据时，如果查询的数据能在⼆级索引找到，那么就是「索引覆盖」操作，如果查询的数据不在⼆级索引⾥，就需要先在⼆级索引找到主键值，需要去聚簇索引中获得数据⾏，这个过程就叫作「回表」
	- 回表：回到主键索引树搜索的过程，称为回表
	- 最左前缀原则：B+Tree这种索引结构，可以利⽤索引的"最左前缀"来定位记录；
		- 只要满⾜最左前缀，就可以利⽤索引来加速检索；
		- 最左前缀可以是联合索引的最左N个字段，也可以是字符串索引的最左M个字符
		- 第⼀原则是：如果通过调整顺序，可以少维护⼀个索引，那么这个顺序往往就是需要优先考虑采⽤的。
	- 索引下推
		- 在MySQL5.6之前，只能从根据最左前缀查询到ID开始⼀个个回表；到主键索引上找出数据⾏，再对⽐字段值
		- MySQL5.6引⼊的索引下推优化，可以在索引遍历过程中，对索引中包含的字段先做判断，直接过滤掉不满⾜条件的记录，减少回表次数

- 聚簇索引
	- 指索引项的排序⽅式和表中数据记录排序⽅式⼀致的索引
	- 聚簇索引并不是⼀种单独的索引类型，⽽是⼀种数据存储⽅式
	- 术语“聚簇”表示数据⾏和相邻的键值紧凑的存储在⼀起
	- 也就是说聚集索引的顺序就是数据的物理存储顺序，它会根据聚集索引键的顺序来存储表中的数据，即对表的数据按索引键的顺序进⾏排序，然后重新存储到磁盘上
	- 因为数据在物理存放时只能有⼀种排列⽅式，所以⼀个表只能有⼀个聚集索引

- MySQL外键
	- 外键是某个表中的一列，它包含在另一个表的主键中
		- 外键也是索引的一种，是通过一张表中的一列指向另一张表中的主键，来对两张表进行关联
		- 一张表可以有一个外键，也可以存在多个外键，与多张表进行关联
	- 外键的主要作用是保证数据的一致性和完整性，并且减少数据冗余
		- 阻止执行
			- 从表插入新行，其外键值不是主表的主键值便阻止插入
			- 从表修改外键值，新值不是主表的主键值便阻止修改
			- 主表删除行，其主键值在从表里存在便阻止删除(要想删除，必须先删除从表的相关行)
			- 主表修改主键值，旧值在从表里存在便阻止修改(要想修改，必须先删除从表的相关行)
		- 级联执行
			- 主表删除行，连带从表的相关行一起删除
			- 主表修改主键值，连带从表相关行的外键值一起修改
		- 外键创建限制
			- 父表必须已经存在于数据库中，或者是当前正在创建的表；如果是后一种情况，则父表与子表是同一个表，这样的表称为自参照表，这种结构称为自参照完整性
			- 必须为父表定义主键
			- 外键中列的数目必须和父表的主键中列的数目相同
			- 两个表必须是InnoDB表，MyISAM表暂时不支持外键。
			- 外键列必须建立了索引，MySQL 4.1.2以后的版本在建立外键时会自动创建索引，但如果在较早的版本则需要显式建立
			- 外键关系的两个表的列必须是数据类型相似，也就是可以相互转换类型的列，比如int和tinyint可以，而int和char则不可以

- MySQL索引
	- 存储方式区分
		- BTREE 索引，MySQL索引默认使用B+树
		- HASH 索引，散列表(Hash table，也叫哈希表)
				- Hash table 查询时间为 O(1)，但是其对范围查询的支持不如 B+树
				- 即Hash table只支持等于或不等于，不支持关键词检索
	- 逻辑区分
		- 普通索引
		- 唯一索引
		- 主键索引
		- 空间索引
		- 全文索引
	- 主键默认会加索引
		- 按主键构建的B+树里包含所有列的数据，而普通索引的B+树里只存储了主键，还需要再查一次主键对应的B+树(回表)
		- 使用 `explain`命令 查看一个SQL语句的执行计划，如使用的索引，是否做全表扫描等
	- 联合索引的前缀同样具有的索引的效果
	- sql语句前加 `explain`可以查看索引使用情况
	- 如果MySQL没有选择最优的索引方案，可以在 where前 `force index (index_name)`
	- 规避慢查询
		- 大部分的慢查询都是因为没有正确地使用索引
		- 一次select不要超过1000行
		- 分页查询limit m,n 会检索前m+n行，只是返回后n行，通常用id>x来代替这种分页方式
		- 批量操作时最好一条sql语句搞定；其次打包成一个事务，一次性提交(高并发情况下减少对共享资源的争用)
		- 不要使用连表操作，join逻辑在业务代码里完成
```sql
// 创建索引
// PRIMARY KEY（主键索引）
ALTER TABLE `table_name` ADD PRIMARY KEY ( `column` )

// UNIQUE (唯⼀索引)
ALTER TABLE `table_name` ADD UNIQUE (`column` )

// INDEX(普通索引)
ALTER TABLE `table_name` ADD INDEX index_name ( `column` )

// FULLTEXT(全⽂索引)
ALTER TABLE `table_name` ADD FULLTEXT ( `column` )

// 多列索引(组合索引)
 mysql > ALTER TABLE `table_name` ADD INDEX index_name
( `column1`, `column2`, `column3` )
CREATE TABLE table_name ( ID INT NOT NULL, username VARCHAR(16) NOT NULL, city
VARCHAR(50) NOTNULL, age INT NOT NULL );

// 删除索引的语法
DROP INDEX [ indexName ] ON mytable;
```

- B+树
	- B即Balance，对于m叉树每个节点上最多有m个数据，最少有m/2个数据(根节点除外)
	- 叶节点上存储了所有数据，把叶节点链接起来可以顺序遍历所有数据
	- 每个节点设计成内存页的整倍数；MySQL的 m=1200，树的前两层放在内存中
```
	            2  28 65
	            P1 P2 P3
	    |          |         |
	2  13 20   28 35 56   65 80 90
	P1 P2 P3   P1 P2 P3   P1 P2 P3
	|  |  |    |  |  |    |  |  | 
	2  13 20   28 35 56   65 80 90
	8  15 23   30 38 60   73 85 96
	9  19 27   32 50 63   79 88 99
	Q  Q  Q    Q  Q  Q    Q  Q  Q   // data 
```

- MySQL配置
	- Linux安装MySQL客户端   `yum install mysql`
	- 安装MySQL服务端        `yum install mysql-server`
	- 启动MySQL服务端        `systemctl start mysqld.service`
	- 以root登录             `mysql -uroot`
		- SQL管理员创建账号      `create user 'tester' identified by '123456';`
		- 查看账号创建是否成功    `select host, user from mysql.user where user='tester';`
		- 赋予账号对应权限
```sql
	// grant <privileges> on <database>.<table> to 'tester'@'localhost';
	grant create, insert on *.* to 'tester'; 
	// 这里表示赋予该用户所有数据库所有表(*.*表示所有表)，%表示所有IP地址
	grant all privileges on *.* to '用户名'@'%' identified by '密码' with grant option;
```
	
- MySQL管理	
```sql
	// 以tester登录
	`mysql -utester -p'123456' -h121.40.150.39 -P9528`

	// 创建database
	`create database test;`
	
	// 使用database  
	use test;
	show tables;

	// 创建表
	create table if not exists student(
	id int not null auto_increment comment '主键自增id',
	name char(4) not null comment '姓名',
	province char(6) not null comment '省',
	city char(10) not null comment '城市',
	addr varchar(100) default '' comment '地址',
	score float not null default 0 comment '考试成绩',
	enrollment date not null comment '入学时间',
	primary key (id),  unique key idx_name (name),  
	key idx_location (province,city)
	)default charset=utf8 comment '学员基本信息';

	show variables like 'innodb_large_prefix';
	show variables like 'innodb_file_format';
	set global innodb_large_prefix=1;
	set global innodb_file_format=BARRACUDA;
	CREATE TABLE `resource` (
	`id` char(64) CHARACTER SET latin1 NOT NULL,
	`vendor` tinyint(1) NOT NULL,
	`region` varchar(64) CHARACTER SET latin1 NOT NULL,
	`zone` varchar(64) CHARACTER SET latin1 NOT NULL,
	`create_at` bigint(13) NOT NULL,
	`expire_at` bigint(13) DEFAULT NULL,
	`category` varchar(64) CHARACTER SET latin1 NOT NULL,
	`type` varchar(120) CHARACTER SET latin1 NOT NULL,
	`instance_id` varchar(120) CHARACTER SET latin1 NOT NULL,
	`name` varchar(255) NOT NULL,
	`description` varchar(255) DEFAULT NULL,
	`status` varchar(255) CHARACTER SET latin1 NOT NULL,
	`update_at` bigint(13) DEFAULT NULL,
	`sync_at` bigint(13) DEFAULT NULL,
	`sync_accout` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
	`public_ip` varchar(64) CHARACTER SET latin1 DEFAULT NULL,
	`private_ip` varchar(64) CHARACTER SET latin1 DEFAULT NULL,
	`pay_type` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
	`describe_hash` varchar(255) NOT NULL,
	`resource_hash` varchar(255) NOT NULL,
	PRIMARY KEY (`id`),
	KEY `name` (`name`) USING BTREE,
	KEY `status` (`status`) USING HASH,
	KEY `private_ip` (`private_ip`) USING BTREE,
	KEY `public_ip` (`public_ip`) USING BTREE,
	KEY `instance_id` (`instance_id`) USING HASH
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

	// 查看索引 [Specified key was too long](https://blog.csdn.net/shaochenshuo/article/details/51064685)
	show index from student\G;

	// 新增记录
	// 必须给`not null`且无 default值的列赋值
	insert into student (name,province,city,enrollment) values
	('张三','北京','北京','2021-03-05'),
	('李四','河南','郑州','2021-04-25'),
	('小丽','四川','成都','2021-03-10');

	// 查询记录
	select id,name from student where id>0;
	select province,avg(score) as avg_score from student 
		where score>0 
		group by province having avg_score>50 
		order by avg_score desc;
	select * from student order by c1 asc, c2 desc;    // 先按列 c1 升序排列，再按 c2 降序排列

	// 修改记录
	update student set score=score+10,addr='海淀' where province='北京';
	update student set
		score=case province
			when '北京' then score+10     
			when '四川' then score+5 
			else score+7
		end,
		addr=case province
			when '北京' then '东城区'        
			when '四川' then '幸福里'        
			else '朝阳区'    
		end
	where id>0;

	// 删除记录
	delete from student where city='郑州';
	// 删除表里的所有行
	delete from student;
	// 删除表
	drop table student;     
```

### 2. MySQL初始配置
- 注意事项
	- 注意存储引擎的选择(InnoDB)
	- 主键选择 和 唯一键考虑清楚
	- 考虑数据类型与长度，选择合适的类型，避免空间浪费
	- 字符串注意确认字符集，如果需要存入中文，请选择utf8编码
	- 为过滤条件的字段 添加索引
	- 高频组合查询可以考虑 联合索引
	- 注意选择使用索引的方法: `Hash` `Btree` `Normal`
	- 写sql时一律使用小写
	- 建表时先判断表是否已存在 `if not exists` 
	- 所有的列和表都加 `comment`
		- `comment` 是备注、注释的意思，写上 `comment 'id'` 之后，在建表信息里可以看到添加的备注信息
	- 字符串长度比较短时尽量使用 `char`，定长有利于内存对齐，读写性能更好，而 `varchar`字段频繁修改时容易产生内存碎片
	- 满足需求的前提下尽量使用短的数据类型，如`tinyint` vs `int`, `float` vs `double`, `date` vs `datetime`
	- null
		- default null有别于 default '' 和 default 0
		- is null, is not null有别于 != '', !=0
		- 尽量设为 not null
			- 有些DB索引列不允许包含null
			- 对含有null的列进行统计，结果可能不符合预期
			- null值 有时候会严重拖慢系统性能

### 3. Go语言中SQL驱动接口
- database/sql
	- sql包提供了保证SQL或类SQL数据库的泛用接口
	- Go官方没有提供数据库驱动，而是为开发数据库驱动定义了一些标准接口(即database/sql)，开发者可以根据定义的接口来开发相应的数据库驱动
	- Go语言中支持MySQL的驱动比较多，如
		- `github.com/go-sql-driver/mysql` 支持 database/sql
		- `github.com/ziutek/mymysql` 支持 database/sql，支持自定义的接口
		- `github.com/Philio/GoMySQL` 不支持 database/sql，支持自定义的接口

- database/sql/driver
	- driver包定义了应被数据库驱动实现的接口，这些接口会被sql包使用，绝大多数代码应使用sql包绝大多数代码应使用sql包

- 详情参考
  - [Golang-Sql](./27_Golang-Sql.md)

### 4. Go语言中 数据库的操作
- 增删改查
	- Go语言中的第三方库 `go get github.com/go-sql-driver/mysql`
	- 连接数据库 `db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/test?charset=utf8")`
	- 增删改 `func (*sql.DB).Exec(sql string) (sql.Result, error)`
	- 查 `func (*sql.DB).Query(sql string) (*sql.Rows, error)`

- Stmt
	- 参数化查询 `db.Where("merchant_id = ?", merchantId)`
	- 拼接sql `db.Where(fmt.Sprintf("merchant_id = %s", merchantId))`
	- 定义一个sql模板 `stmt, err := db.Prepare("update student set score=score+? where city=?")`
	- 多次使用模板 `res, err := stmt.Exec(10, "上海")`

### 5. 数据库使用问题
- SQL注入
	- 问题现象
		- `sql = "select username,password from user where username='" + username + "' and password='" + password + "'";`
			- 变量username和password从前端输入框获取，如果用户输入的username为lily， password为aaa' or '1'='1
			- 则完整的sql为select username,password from user where username='lily' and password='aaa' or '1'='1'
			- 会返回表里的所有记录，如果记录数大于0就允许登录，则lily的账号被盗
		- `sql = "insert into student (name) values ('"+username+"')";`
			- 变量username从前端输入框获取，如果用户输入的username为lily'); drop table student;--
			- 完整sql为insert into student (name) values ('lily'); drop table student;--')
			- 通过注释符--屏蔽掉了末尾的')，删除了整个表
	- 预防措施
		- 前端输入要加正则校验、长度限制
		- 对特殊符号`(<>&*; '"等)`进行转义或编码转换，Go的 text/template包里面的HTMLEscapeString函数可以对字符串进行转义处理
		- 不要将用户输入直接嵌入到sql语句中，而应该使用参数化查询接口，如Prepare、Query、Exec(query string, args ...interface{})
		- 使用专业的SQL注入检测工具进行检测，如sqlmap、SQLninja
		- 避免网站打印出SQL错误信息，以防止攻击者利用这些错误信息进行SQL注入

- SQL预编译
	- DB执行SQL分为3步
		- 词法和语义解析
		- 优化 SQL 语句，制定执行计划
		- 执行并返回结果
	- SQL 预编译技术是指将用户输入用占位符?代替，先对这个模板化的sql进行预编译，实际运行时再将用户输入代入
	- 除了可以防止 SQL 注入，还可以对预编译的SQL语句进行缓存，之后的运行就省去了解析优化SQL语句的过程

### 6. ORM与NoSQL技术
- SQLBuilder
	- Go语言中的第三方库
		- Go-SQLBuilder 是一个用于创建SQL语句的工具函数库，提供一系列灵活的、与原生SQL语法一致的链式函数，归属于艾润物联公司 `go get -u github.com/parkingwang/go-sqlbuilder`
		- Gendry 是一个用于辅助操作数据库的Go包，基于go-sql-driver/mysql，它提供了一系列的方法来为调用标准库database/sql中的方法准备参数 `go get –u github.com/didi/gendry`
	- 自行封装SQL构建器
		- 写一句很长的sql容易出错，且出错后不好定位
		- 函数式编程可以直接定位到是哪个函数的问题
		- 函数式编程比一长串sql更容易编写和理解
		- [custom sqlbuilder](../database/03_sql_builder/custom_sqlbuilder/main.go)
```go
	// Go-SQLBuilder 函数链
	sql := gsb.NewContext().
		Select("id", "name", "score", "city").
		From("student").
		OrderBy("score").DESC().                  // 按"score"降序
		Column("name").ASC().                     // 当"score"相同，按"name"升序
		Limit(10).Offset(20).                     // 从第20个开始，读10个 limit 20, 10
		ToSQL()

	// Gendry
	//  map
	where := map[string]interface{}{
		"city": []string{"北京", "上海", "杭州"},
		"score<": 30,
		"addr": builder.IsNotNull,
		"_orderby": "score desc",
	}
	fields := []string{"id", "name", "city", "score"}
	_, _, err := builder.BuildSelect("student", where, fields)
```

- ORM技术与GORM
	- ORM
		- ORM 即 Object Relational Mapping，对象关系映射
		- Relational 指各种sql类的关系型数据为
		- Object 指面向对象编程(object-oriented programming)中的对象
		- ORM在数据库记录和程序对象之间做一层映射转换，使程序中不用再去编写原生SQL，而是面向对象的思想去编写类、对象、调用相应的方法来完成数据库操作
	- GORM
		- `go get -u gorm.io/gorm`
		- `go get -u gorm.io/driver/mysql`
		- GORM是一个全能的、友好的、基于golang的ORM库
		- GORM 倾向于约定，而不是配置
			- 默认情况下，GORM 使用 ID 作为主键，使用结构体名的 蛇形复数 作为表名，字段名的 蛇形 作为列名，并使用 `CreatedAt`、`UpdatedAt` 字段追踪创建、更新时
```go
		// 完全是在操作struct，看不到sql的影子
		dsn := "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"   // data source name
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		student := Student{Name: "光绪", Province: "北京", City: "北京", Score: 38, Enrollment: time.Now()}
		db.Create(&student)
```

### 7. MongoDB
- MongoDB相关概念
	- NoSQL泛指非关系型数据库，如 mongodb, redis, HBase
	- mongo使用高效的二进制数据存储，文件存储格式为BSON(一种json的扩展，比json性能更好，功能更强大)
	- MySQL中表的概念在mongo里叫集合(collection)， MySQL中行的概念在mongo中叫文档(document)，一个文档看上去像一个json

- MongoDB的安装配置
```shell
	// 安装MongoDB
	$ vim /etc/yum.repos.d/mongodb-org-4.2.repo
		[mongodb-org-4.2] 
		name=MongoDB Repository
		baseurl=https://repo.mongodb.org/yum/redhat/$releasever/mongodb-org/4.2/x86_64/
		gpgcheck=1
		enabled=1
		gpgkey=https://www.mongodb.org/static/pgp/server-4.2.asc
	$ yum install -y mongodb-org
	$ systemctl start mongod.service
	
	// MongoDB常用命令
	$ mongo
	use test;                                                                                // 切换到test库，如果没有则创建
	db.createUser({user: "tester", pwd: "123456", roles: [{role: "dbAdmin", db: "test"}]});  // 创建用户
	db.createCollection("student");                                                          // 创建collection
	db.student.createIndex({"name":1});                                                      // 在name上创建索引,不是唯一索引
	db.student.insertOne({name:"张三",city:"北京"});
	db.student.find({name:"张三"});
	db.student.update({name:"张三"},{name:"张三",city:"上海"});
	db.student.deleteOne({name:"张三"});
```

- Go语言中 MongoDB的操作
```go
	// 安装mongo-driver
	go get go.mongodb.org/mongo-driver
	go get go.mongodb.org/mongo-driver/x/bsonx/bsoncore@v1.7.1
	go get go.mongodb.org/mongo-driver/x/mongo/driver@v1.7.1
	go get go.mongodb.org/mongo-driver/mongo/options@v1.7.1
	go get go.mongodb.org/mongo-driver/x/mongo/driver/topology@v1.7.1
	go get go.mongodb.org/mongo-driver/mongo@v1.7.1
	
	// 连接MongoDB
	option := options.Client().ApplyURI("mongodb://127.0.0.1:27017").
	SetConnectTimeout(time.Second).                                                         // 连接超时时长
	SetAuth(options.Credential{Username: "tester", Password: "123456", AuthSource: "test"}) // 指定用户名和密码，AuthSource代表Database
	client, err := mongo.Connect(context.Background(), option)
	// 注意: Ping成功才代表连接成功
	err = client.Ping(ctx, nil) 

	// 查询MongoDB
	sort := bson.D{{"name", 1}}                                 // 1升序，-1降序
	filter := bson.D{{"score", bson.D{{"$gt", 3}}}} //score>3   // greater than
	findOption := options.Find()
	findOption.SetSort(sort)                                    // 按name排序
	findOption.SetLimit(10)                                     // 最多返回10个
	findOption.SetSkip(3)                                       // 跳过前3个
	cursor, err := collection.Find(ctx, filter, findOption)
```