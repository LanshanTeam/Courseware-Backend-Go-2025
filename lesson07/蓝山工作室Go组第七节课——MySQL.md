## 蓝山工作室Go组第七节课——MySQL

### 什么是数据库

> 我们上节课学习了网络编程，也学习了Gin框架，大家已经能够写出一个http服务了。以一个简单的http服务为例，你做了一个注册登陆功能。当用户在你的网站上注册了用户，那如何存储用户信息呢？有些同学会说我们可以使用map存储，但这带来了一个问题，一旦服务重启，map的内容自然而然会丢掉；那还有同学说可以存在文件里，这是一个很好的思路，即使我们重启也不会导致数据丢失。但同时也带来了一些问题，如果数据量过大，每次加载文件和查询数据时会异常缓慢。那我们本节课就是来解决这个问题的，数据库，它可以帮助我们存储信息并且进行快速查询。

数据库可以理解为一个**专门存放、管理数据的电子系统**。它类似一个功能强大的“信息仓库”，能把大量数据（例如用户信息、商品、订单、成绩等）**结构化（以一种特定的格式）地保存起来**，并且支持**快速查询、更新和删除**。相比手写记录或文件存储，数据库更安全、更高效，能保证多人同时操作时数据不冲突、不丢失。网站、APP、公司系统等几乎所有软件都依赖数据库来运作。简单来说，数据库就是现代软件世界中**用来存数据并让你随时能查到的核心工具**。

[数据库其实就那么回事](https://www.bilibili.com/video/BV1Hh411C7Up/?spm_id_from=333.337.search-card.all.click&vd_source=a57b854dceb9e43274d666a1811d0403)(后面是广告不用看，看前面科普就好)

[什么是MySQL](https://www.bilibili.com/video/BV1p5qhYsE4f/?spm_id_from=333.337.search-card.all.click&vd_source=a57b854dceb9e43274d666a1811d0403)

[一小时MySQL教程](https://www.bilibili.com/video/BV1AX4y147tA?spm_id_from=333.788.videopod.sections&vd_source=a57b854dceb9e43274d666a1811d0403&p=7)

### 数据库的分类

> 数据库的分类方式有很多很多，这里简单分为关系型数据库和非关系型数据库两种，有兴趣的同学可以自行了解其他分类

##### 关系型数据库

1. 用**表格（行列）**存数据，结构固定

2. 用 **SQL** 查询，能做复杂操作

3. 支持 **ACID 事务**，数据一致性强

4. 适合订单、用户、财务等 **复杂业务场景**  

   ***如MySQL，PostgreSQL***

##### 非关系型数据库

1. 数据结构**灵活**（键值、文档、列族等）

2. 查询简单，没有统一语言

3. 性能高、扩展容易，多为**最终一致性**

4. 常用于缓存、日志、海量数据等 **高并发场景**

   ***如Redis，MongoDB***

### MySQL数据库

#### 部署

可以直接部署在本机上或者使用Docker部署。更加推荐使用Docker

**直接安装**：[MySQL官网](https://dev.mysql.com/downloads/mysql/)

[MySQL安装、配置与卸载教程（Windows版）](https://juejin.cn/post/7337186171210170383)

**使用Docker安装**：

如何安装Docker，看这里：[Docker安装]()

```go
docker pull mysql
docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql
```

进入容器

```go
docker exec -it mysql /bin/bash
```

**登陆**

```go
Mysql -u root -p
```

然后输入密码即可





#### 介绍

MySQL是一种**关系型数据库**，关系型数据库的数据都是以**数据表**的形式进行存储的

> 看起来和excel差不多

```sql
+--------------+-----------------+------+-----+---------+----------------+
| Field        | Type            | Null | Key | Default | Extra          |
+--------------+-----------------+------+-----+---------+----------------+
| id           | bigint unsigned | NO   | PRI | NULL    | auto_increment |
| created_at   | datetime(3)     | YES  |     | NULL    |                |
| updated_at   | datetime(3)     | YES  |     | NULL    |                |
| nick_name    | varchar(42)     | YES  |     | NULL    |                |
| gender       | varchar(12)     | YES  |     | NULL    |                |
| password     | varchar(258)    | YES  |     | NULL    |                |
| avatar       | varchar(256)    | YES  |     | NULL    |                |
| introduction | longtext        | YES  |     | NULL    |                |
| email        | varchar(128)    | YES  |     | NULL    |                |
| qq           | bigint          | YES  |     | NULL    |                |
| tel          | varchar(18)     | YES  |     | NULL    |                |
| birthday     | varchar(128)    | YES  |     | NULL    |                |
| role         | tinyint         | YES  |     | 1       |                |
| user_name    | varchar(42)     | YES  |     | NULL    |                |
+--------------+-----------------+------+-----+---------+----------------+
```

上图就是一个user表，接下来我们简单来解析一下

1. **Field**（字段名）

- **含义**：表示表中列的名称，也就是数据表中每一列的名字。

2. **Type**（数据类型）

- **含义**：表示该字段的数据类型及其长度或精度。MySQL支持各种数据类型，如整数（INT）、浮动小数（FLOAT）、字符串（VARCHAR）、日期时间（DATETIME）等。详细内容可参考下表。

| **类别**         | **数据类型**          | **说明**                             | **范围/长度**                                                |
| ---------------- | --------------------- | ------------------------------------ | ------------------------------------------------------------ |
| **数字类型**     | `TINYINT`             | 存储小范围的整数                     | 有符号：-128 到 127；无符号：0 到 255                        |
|                  | `SMALLINT`            | 存储较小范围的整数                   | 有符号：-32,768 到 32,767；无符号：0 到 65,535               |
|                  | `MEDIUMINT`           | 存储中等范围的整数                   | 有符号：-8,388,608 到 8,388,607；无符号：0 到 16,777,215     |
|                  | `INT` / `INTEGER`     | 存储普通范围的整数                   | 有符号：-2,147,483,648 到 2,147,483,647；无符号：0 到 4,294,967,295 |
|                  | `BIGINT`              | 存储大范围的整数                     | 有符号：-9,223,372,036,854,775,808 到 9,223,372,036,854,775,807；无符号：0 到 18,446,744,073,709,551,615 |
|                  | `FLOAT`               | 存储单精度浮点数                     | 4 字节，精度通常为 7 位数字                                  |
|                  | `DOUBLE`              | 存储双精度浮点数                     | 8 字节，精度通常为 15 位数字                                 |
|                  | `DECIMAL` / `NUMERIC` | 存储定点数（高精度）                 | 定义时指定精度和小数位数，例如 `DECIMAL(10,2)`，表示最多 10 位，2 位小数 |
| **字符串类型**   | `CHAR`                | 固定长度字符串                       | 1 到 255 字符                                                |
|                  | `VARCHAR`             | 可变长度字符串                       | 1 到 65,535 字符                                             |
|                  | `TEXT`                | 变长文本数据                         | 最多 65,535 字符                                             |
|                  | `TINYTEXT`            | 变长文本数据                         | 最多 255 字符                                                |
|                  | `MEDIUMTEXT`          | 变长文本数据                         | 最多 16,777,215 字符                                         |
|                  | `LONGTEXT`            | 变长文本数据                         | 最多 4,294,967,295 字符                                      |
| **日期时间类型** | `DATE`                | 存储日期                             | `YYYY-MM-DD`，范围：1000-01-01 到 9999-12-31                 |
|                  | `DATETIME`            | 存储日期和时间                       | `YYYY-MM-DD HH:MM:SS`，范围：1000-01-01 00:00:00 到 9999-12-31 23:59:59 |
|                  | `TIMESTAMP`           | 存储时间戳                           | `YYYY-MM-DD HH:MM:SS`，范围：1970-01-01 00:00:01 到 2038-01-19 03:14:07 (UTC) |
|                  | `TIME`                | 存储时间                             | `HH:MM:SS`，范围：`-838:59:59` 到 `838:59:59`                |
|                  | `YEAR`                | 存储年份                             | `YYYY`，范围：1901 到 2155                                   |
| **布尔类型**     | `BOOLEAN` / `BOOL`    | 布尔值                               | 0（`FALSE`）或 1（`TRUE`）                                   |
| **二进制类型**   | `BINARY`              | 固定长度的二进制数据                 | 1 到 255 字节                                                |
|                  | `VARBINARY`           | 可变长度的二进制数据                 | 1 到 65,535 字节                                             |
|                  | `BLOB`                | 二进制大对象                         | 最多 65,535 字节                                             |
|                  | `TINYBLOB`            | 小型二进制大对象                     | 最多 255 字节                                                |
|                  | `MEDIUMBLOB`          | 中型二进制大对象                     | 最多 16,777,215 字节                                         |
|                  | `LONGBLOB`            | 大型二进制大对象                     | 最多 4,294,967,295 字节                                      |
| **JSON 类型**    | `JSON`                | 存储 JSON 格式数据                   | 最多 4GB 数据                                                |
| **集合类型**     | `ENUM`                | 枚举类型（限制为一组预定义值之一）   | 1 到 65,535 个预定义值                                       |
|                  | `SET`                 | 集合类型（可存储多个预定义值的组合） | 1 到 64 个预定义值的组合                                     |

3. **Null**（是否允许为NULL）

- **含义**：表示该字段是否允许存储 `NULL` 值，`NULL` 表示缺失或未知的数据。
- 值：

  - `YES`：字段允许为 `NULL`，即可以没有值。
- `NO`：字段不允许为 `NULL`，即该列必须有值。

4. **Key**（索引类型）

- **含义**：表示该字段是否作为索引的一部分，并指示索引的类型。MySQL支持多种类型的索引。
- **值**：
  - `PRI`：主键索引（Primary Key）。这是一个唯一且非空的索引，表中每个记录的主键值**必须唯一**。
  - `UNI`：唯一索引（Unique Key）。该字段的值**必须唯一**，但允许 `NULL` 值。
  - `MUL`：多重索引（Multiple）。表示该字段是普通索引的一部分，允许有重复值。
  - 如果该列没有任何索引，则该列显示为空。

5. **Default**（默认值）

- **含义**：表示该字段在没有指定值时使用的默认值。默认值可以是常量，也可以是 `NULL`。

- 示例：

  - 如果你在创建表时定义 `age INT DEFAULT 18`，当插入数据时，如果没有给 `age` 字段指定值，它会自动使用默认值 `18`。
- 如果一个字段没有默认值，`Default` 会显示为 `NULL`，表示没有默认值。

6. **Extra**（附加信息）

- **含义**：提供与字段相关的额外信息，例如是否自动递增（AUTO_INCREMENT）等。
- **值**：
  - `auto_increment`：字段是自动递增的，通常用于主键字段（如自增ID）。
  - `on update CURRENT_TIMESTAMP`：表示当记录被更新时，字段会自动更新为当前时间，通常用于记录最后更新时间的字段。
  - 如果没有额外信息，则该列为空

**插入数据后，表内容示例**（这里去掉了一些字段，看起来直观一些）

```sql
+----+------------+-----------+------------------------+---------+------+---------------------+
| id | user_name  | nick_name |         email          | gender  | role |     created_at      |
+----+------------+-----------+------------------------+---------+------+---------------------+
| 1  | xiaoming   | 小明      | xiaoming@example.com   | male    | 1    | 2025-01-15 10:12:33 |
| 2  | xiaohong   | 小红      | xiaohong@example.com   | female  | 1    | 2025-01-15 11:20:10 |
| 3  | admin      | Admin     | admin@example.com      | male    | 2    | 2025-01-10 08:33:21 |
| 4  | anon       | 匿名用户   | anon@example.com       | NULL    | 1    | 2025-01-17 15:02:55 |
| 5  | sonwwall   | 外城      | waicheng@example.com   | male    | 1    | 2025-01-18 21:33:59 |
+----+------------+-----------+------------------------+---------+------+---------------------+
```

### SQL

#### 介绍

SQL 是一种用于**查询和管理关系型数据库**的语言。它可以用来**增删改查数据**，以及创建或修改表结构。语法接近英文，是 MySQL、PostgreSQL、Oracle 等数据库都支持的标准语言。

<img src="./img.png" />

希望大家成为SQL高手.jpg

#### 语法

[sql教程-菜鸟教程](https://www.runoob.com/sql/sql-tutorial.html)

[MySQL入门到精通-黑马程序员](https://www.bilibili.com/video/BV1Kr4y1i7ru/?spm_id_from=333.337.search-card.all.click&vd_source=a57b854dceb9e43274d666a1811d0403)

**语法部分不再展示了，自己看文档教程即可，这里展示一些基本用法**

【创建数据库】

```sql
CREATE DATABASE school;
USE school;
```

---

【创建数据表】
 创建一个名为 students 的表，包含学号、姓名、年龄、专业和成绩等字段：

```sql
CREATE TABLE students (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50),
    age INT,
    major VARCHAR(50),
    score INT
);
```

------

【插入数据】
 向 students 表中插入多行示例数据：

```sql
INSERT INTO students (name, age, major, score)
VALUES 
('Alice', 20, 'Computer Science', 85),
('Bob', 21, 'Math', 90),
('Tom', 19, 'Physics', 78);
```

------

【查询数据】
 查询表中的所有数据：

```sql
SELECT * FROM students;
```

查询指定字段，如仅查看姓名与成绩：

```sql
SELECT name, score FROM students;
```

------

【条件查询】
 查询年龄大于 20 的学生：

```sql
SELECT * FROM students
WHERE age > 20;
```

查询数学专业并且成绩大于 80 的学生：

```sql
SELECT * FROM students
WHERE major = 'Math' AND score > 80;
```

使用 LIKE 进行模糊查询，例如查询名字以 A 开头的学生：

```sql
SELECT * FROM students
WHERE name LIKE 'A%';
```

------

【排序查询】
 按成绩从高到低排序：

```sql
SELECT * FROM students
ORDER BY score DESC;
```

------

【更新数据】
 更新 Alice 的成绩，将其修改为 95：

```sql
UPDATE students
SET score = 95
WHERE name = 'Alice';
```

------

【删除数据】
 删除指定行，例如删除名为 Tom 的学生：

```sql
DELETE FROM students
WHERE name = 'Tom';
```

删除表中所有数据但保留表结构：

```sql
DELETE FROM students;
```

使用 TRUNCATE 也可以快速清空表：

```sql
TRUNCATE TABLE students;
```

------

【删除表】
 当不再需要该数据表时，可以将其完全删除：

```sql
DROP TABLE students;
```

------

【删除数据库】
 若整个数据库不再使用，可以删除数据库 school：

```sql
DROP DATABASE school;
```

### GORM

> 上面介绍了如何使用SQL语句操作数据库，但是在实际写代码中一般不会直接写SQL语句，而是通过ORM库进行操作，Go的一个ORM库便是GORM

#### 什么是ORM

ORM 全称 **对象关系映射（Object–Relational Mapping）**

简单来说，ORM 让你可以像写普通代码一样访问数据库，不需要频繁手写 SQL。它既提高了开发效率，也减少了 SQL 拼接出错的风险。

#### GORM使用

首先首先，必须要看的就是gorm的官方文档[GORM指南](https://gorm.io/zh_CN/docs/index.html),所有关于gorm的用法都在里面提到了，建议多看

**下面给大家一个例子，告诉大家如何使用gorm，详细的使用方法一定要看文档！！！**

1. *安装*

```go
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

2. *连接到数据库*

```go
 dsn := "root:123456@tcp(127.0.0.1:3306)/lanshanteam?charset=utf8mb4&parseTime=True&loc=Local"
 db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
 if err != nil {
        log.Fatal("连接数据库失败: ", err)
 }
```

3. *建立模型*

```go
type Member struct {
	gorm.Model
	Name       string `gorm:"size:32;not null;comment:成员姓名"`
	Age        int    `gorm:"default:18;comment:年龄"`
	Department string `gorm:"size:64;index;comment:部门名称"`
}
```

4. *自动迁移*

```go
err = db.AutoMigrate(&Member{})
	if err != nil {
		log.Fatal("自动迁移失败: ", err)
	}
	log.Println("自动迁移成功")
```

5. *基本增删改查*

```go
// 准备数据
member1 := Member{Name: "kq", Age: 19, Department: "Backend-Go"}
member2 := &Member{Name: "cy", Age: 19, Department: "Backend-Go"}
members := []*Member{
    {Name: "grt", Age: 19, Department: "Backend-Go"},
    {Name: "wjk", Age: 18, Department: "Backend-Go"},
    {Name: "hym", Age: 19, Department: "Backend-Go"},
    {Name: "zbw", Age: 19, Department: "Backend-Python"},
    {Name: "lhy", Age: 20, Department: "Backend-Python"},
    {Name: "lk", Age: 19, Department: "Ops"},
}

// =======================
//       增（Create）
// =======================
db.Create(&member1)
db.Create(member2)
db.Create(members)
log.Println("数据创建完成")

// =======================
//       查（Read）
// =======================

// 1. 查第一条记录
var m1 Member
db.First(&m1)
fmt.Println("First 查询:", m1)

// 2. 根据 ID 查询
var m2 Member
db.First(&m2, 3) // 查 ID = 3 的记录
fmt.Println("根据 ID 查询:", m2)

// 3. 条件查询（查 Backend-Go 的成员）
var goMembers []Member
db.Where("department = ?", "Backend-Go").Find(&goMembers)
fmt.Println("Backend-Go 成员:")
for _, m := range goMembers {
    fmt.Println("成员:", m.Name)
}
var age19plusMembers []Member
db.Where("age >= ?", 20).Find(&age19plusMembers)
fmt.Println("年龄大于等于20的成员:")
for _, m := range age19plusMembers {
    fmt.Println("成员:", m.Name)
}

// =======================
//       改（Update）
// =======================

// 1. 更新单个字段
db.Model(&m1).Update("Age", 20)

// 2. 更新多个字段
db.Model(&m1).Updates(Member{Name: "kq-new", Department: "Backend-Python"})

fmt.Println("更新完成:", m1)

// =======================
//       删（Delete）
// =======================

// 1. 根据 id 删除
db.Delete(&Member{}, 2) // 删除 ID = 2 的成员

// 2. 按条件删除
db.Where("department = ?", "Ops").Delete(&Member{})

var opsMember Member
err = db.Where("department = ?", "Ops").First(&opsMember).Error
if err != nil {
    fmt.Println(err.Error())
}

log.Println("删除操作完成")
}
```



#### 事务

*事务（Transaction）指一组要么同时成功，要么同时失败的操作。
中间不能有部分成功、部分失败的情况。*

事务具有四大特性（ACID）：

- 原子性（Atomicity）
- 一致性（Consistency）
- 隔离性（Isolation）
- 持久性（Durability）

示例：

```go
func transactionExample(db *gorm.DB) {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 插入数据 1
		member1 := Member{Name: "md", Age: 20, Department: "Go", ID: 50001}
		if err := tx.Create(&member1).Error; err != nil {
			return err
		}

		// 插入数据 2（这里模拟失败）
		member2 := Member{Name: "md", Age: 21, Department: "Python", ID: 50001} // 假设唯一键冲突
		if err := tx.Create(&member2).Error; err != nil {
			return err
		}

		return nil // 没有错误则提交事务
	})

	if err != nil {
		log.Println("事务回滚：", err)
	} else {
		log.Println("事务提交成功")
	}
}
```

### 其他

#### 数据库工具

推荐使用[Navicat](https://www.navicat.com.cn/)

[DataGrip](https://www.jetbrains.com/zh-cn/datagrip/)

实在不想装可以使用Goland自带的数据库工具

## 作业

#### lv0

自己学习一下SQL语句，没事练习几个，不用提交

#### lv1

**完成一个ToDo小项目**

*基本要求*：1.利用本节课学习的gorm框架操作数据库，实现一个任务清单的增删改查

*进阶要求*：

1.利用上节课学习的gin框架或者hz框架，把增删改查四个任务包装为四个接口，对外提供http服务，使用api工具对任务的增删改查

2.结合上节课学习的注册登陆，注册时把用户信息存储到数据库中（密码实现加密存储），实现持久化存储。

3.使用jwt进行鉴权，只有注册登陆过的用户才可以对任务清单操作

---

希望大家都做一下这个小项目，基本要求并不难，理解好本节课的示例代码就能写出来。也鼓励大家做一些进阶要求，在寒假考核中这些都会用到。作业完成后把GitHub仓库地址发送到*guoruitong@lanshan.email*，有问题随时与我联系。