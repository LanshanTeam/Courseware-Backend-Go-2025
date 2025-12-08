# 非关系型数据库概述

在学习过 MySQL 等关系型数据库后，我们可以更好地理解非关系型数据库（NoSQL）的设计思路。关系型数据库的优势在于数据的规范性和严格的表格结构，但这也使得它在某些场景下变得不够灵活，尤其是在面对大规模数据、高并发请求和复杂数据结构时，会面临性能瓶颈。

与此相对，非关系型数据库摒弃了传统的表格结构，通常具有更灵活的数据存储方式，并且能够根据实际需求进行扩展，尤其适合处理大数据量和快速读写的场景。虽然有的非关系型数据库可能不完全遵循 ACID 原则，但通常遵循 CAP 理论，通常有着灵活的数据结构以及更高的性能。

常见的非关系型数据库包括：Redis、MongoDB 等。接下来我们向大家会介绍 Redis 和 MongoDB。

---
# Redis

## 为什么需要 Redis？

同学们已经学习了 web 开发以及 MySQL，实际上，你们已经有能力去完整地写出一个网站比如电商，个人博客，知乎，微信这些项目的后端部分，只不过可能承载不了过大的访问量，为什么呢？有过一些计算机知识的同学可能会知道我们的 CPU 并不是说会直接操作磁盘，也不会直接去访问内存去获取数据，而是有着好几级缓存的，如果每次 cpu 计算都要去访问磁盘或者内存的话，那么就太慢了。

在 web 开发里面，为了解决这个问题，我们也引入了**缓存**的概念。缓存的作用是将数据存储在访问速度更快的内存中，当请求访问这些数据时，直接从内存中读取，避免频繁访问磁盘和数据库，从而大大提高性能。

Redis 是一种高性能的内存数据存储系统，它是一个 NoSQL 数据库，广泛用于缓存场景。通过 Redis，我们可以将数据库中频繁访问的数据缓存起来，从而减少数据库压力，提高响应速度。

---

## 基础概念

Redis 有着以下的特性，决定了他天生适合当作分布式缓存中间件使用：

- **键值对存储**：Redis 是一个键值对存储数据库，可以将数据以键值对的形式存储在内存中，支持各种类型的数据结构，如字符串、哈希、列表、集合等。
    
- **内存存储**：Redis 将数据存储在内存中，这使得它的读写速度极快，远超传统的磁盘数据库。
    
- **持久化选项**：虽然 Redis 主要作为内存数据库，但它也提供了持久化机制，可以定期将数据保存到磁盘上，以防数据丢失。
    
- **高可用性**：Redis 支持主从复制、分区、哨兵等机制，能够实现高可用的分布式架构。

### Redis 的常见数据结构和命令

先连接 redis
```bash
redis-cli -h localhost -p 6379 -a <password>
```
如果你是 docker 部署，你要先进入要 redis 容器里面：
```bash
docker exec -it redis /bin/bash
```

#### 1. String

字符串，最基本的类型，可以存储任何数据，例如文本或数字。命令如下（不必死记，用到再查即可）：

- `SET [key] [value]`：添加或修改一个已有的 String 类型的键值对。
    
- `GET [key]`：根据 key 获取 String 类型的 value。
    
- `MSET [key1 value1] [key2 value2] ...`：批量添加多个 String 类型的键值对。
    
- `MGET [key1] [key2] ...`：根据多个 key 获取多个 String 类型的 value。
    
- `INCR [key]`：让一个整型的 key 自增 1。（对应的当然有 DECR）
    
- `INCRBY [key] [increment]`：让一个整型的 key 自增并指定步长，例如 `INCRBY num 2` 让 num 值自增 2。
    
- `INCRBYFLOAT [key] [increment]`：让一个浮点类型的数字自增并指定步长。
    
- `SETNX [key] [value]` 或 `SET [key] [value] NX`：添加一个 String 类型的键值对，前提是这个 key 不存在，否则不执行。（分布式锁）
    
- `SETEX [key] [seconds] [value]` 或 `SET [key] EX [value]`：添加一个 String 类型的键值对，并且指定有效期（单位：秒）。

你到这里可能会想，既然是缓存，那么我们应该在这非关系型数据库里面如何存储 MySQL 里面的结构化数据呢？其实很简单，我们 MySQL 里面的一行数据其实对应到 Go 里面就是一个结构体，那么我们可以对这个结构体进行序列化成一个字符串，这样就能当作 json 字符串存储到 string 中了，当然这样也有一个弊端，那就是无法灵活修改 string 的成员，下面的数据类型可以帮我们克服这一点。

---

#### 2. Hash

哈希类型，可以理解为我们 Go 语言里面的 map，虽然你们可能会觉得 redis 本身就是一个 kv 数据库，里面还有一个 kv 类型很奇怪，但是这是必要的，一个 key 的 value 就是所谓的 map，可以理解为key里面又存储了多个key的键值对，相较于上面 json 字符串形式存储数据有着一定的优势，那就是对 json 字符串中的单个数据进行修改很不方便，而 hash 类型则可以对单个字段进行 CRUD。

```
key
├── field1: value1
├── field2: value2
└── field3: value3
```

- `HSET [key] [field1] [value1] [field2] [value2] ...`：添加或修改 hash 类型 key 的 field 的值。注：hmset也行，不过已经弃用了.
    
- `HGET [key] [field]`：获取 hash 类型 key 的 field 的值。
    
- `HMGET [key] [field1] [field2] ...`：批量获取多个 field 的值。
    
- `HGETALL [key]`：获取 key 中的所有 field 和 value。
    
- `HKEYS [key]`：获取 key 中的所有 field。
    
- `HVALS [key]`：获取 key 中的所有 value。
    
- `HINCRBY [key] [field] [increment]`：让指定 field 值增加指定步长。
    
- `HSETNX [key] [field] [value]`：添加 field 的值，前提是 field 不存在，否则不执行。
    

---

#### 3. List

可以看作是一个双向队列，但是查询速度 O(N)，原因在他的底层设计，为了节省内存，所以不支持下标查询，下面是他的命令：

- `LPUSH [key] [element] ...`：向列表左侧插入一个或多个元素。
    
- `LPOP [key]`：移除并返回列表左侧第一个元素，没有则返回 nil。
    
- `RPUSH [key] [element] ...`：向列表右侧插入一个或多个元素。
    
- `RPOP [key]`：移除并返回列表右侧第一个元素。
    
- `LRANGE [key] [start] [end]`：返回指定范围内的所有元素。
    
- `BLPOP [key] [timeout]`：在没有元素时等待指定时间，然后返回列表左侧元素。
    
- `BRPOP [key] [timeout]`：在没有元素时等待指定时间，然后返回列表右侧元素。
    

---

#### 4. Set

相当于C++的 `unordered_set` 或者 Java 的 `HashSet`，可以用于查看共同好友等，底层使用哈希表实现，特点是无序，元素不可重复，查找快，支持交集并集这些功能，没错，理解成你们初中学的集合就行了😁：

- `SADD [key] [member] ...`：向 set 中添加一个或多个元素。
    
- `SREM [key] [member] ...`：移除 set 中的指定元素。
    
- `SCARD [key]`：返回 set 中元素的个数。
    
- `SISMEMBER [key] [member]`：判断元素是否存在于 set 中。
    
- `SMEMBERS [key]`：获取 set 中的所有元素。
    
- `SINTER [key1] [key2] ...`：求 key1 与 key2 的交集。
    
- `SDIFF [key1] [key2] ...`：求 key1 与 key2 的差集。
    
- `SUNION [key1] [key2] ...`：求 key1 和 key2 的并集。
    

---

#### 5. SortedSet

有序集合，其实就是给上面 Set 的 member 换成 Key-Value 的形式，也就是 Key 为元素，Value 为数值，按照 Value 进行排序，可以实现排行榜之类的功能，常见命令如下：

- `ZADD [key] [score] [member]`：添加或更新元素的 score 值。
    
- `ZREM [key] [member]`：删除元素。
    
- `ZSCORE [key] [member]`：获取元素的 score 值。
    
- `ZRANK [key] [member]`：获取元素的升序排名。
    
- `ZCARD [key]`：获取元素数量。
    
- `ZCOUNT [key] [min] [max]`：统计 score 在指定范围内的元素个数。
    
- `ZINCRBY [key] [increment] [member]`：让元素 score 增加指定值。
    
- `ZRANGE [key] [min] [max]`：按升序获取指定排名范围的元素。
    
- `ZREVRANGE [key] [min] [max]`：按降序获取指定排名范围的元素。
    
- `ZRANGEBYSCORE [key] [min] [max]`：按 score 获取指定范围的元素。
    
- `ZDIFF`、`ZINTER`、`ZUNION`：求差集、交集、并集。



#### 其他

其他的数据结构并不算基本的数据结构类型，统一介绍一下：
1. **bitmap**：底层使用 String 类型实现，其实就是一个位图，每一位只需要 1 bit，占用非常小，很适合用于存储文章阅读，签到状况这些数据。
2. **Geo**：底层使用 SortedSet 实现，用于表示经纬度，如果想要找方圆距离多少的数据，就可以用他，本质是将经纬度对应成分数来打分排行的。
3. **Stream**：算是一个新加入的数据结构，用于实现消息队列，但是本身功能并不多，所以大多数人还是会使用专门的消息队列。

#### Redis 中键的设计规范

由于Redis中没有表这一结构，于是我们会需要key按照 `项目名:业务名:类型:主键id` 的方式命名，但并不固定，比如mysql里面的shopping库中的goods表的id为1的数据的 key 可以表示为 `shopping:goods:1`，而这一个 key 对应的 value 可以是结构体（对象）序列化后的 json 字符串,这里值得一提的是，如果你用的 `RDM` 的 redis 图形化界面，这样的命名在图形化界面里面会以**树**的形式出现，显示很清晰，但是 `Datagrip` 这类软件貌似并不支持这个功能。

### go-redis

go 里面既然有 gorm 可以操作 MySQL，那么 go 里面也有一个库可以帮助我们去操作 redis。我们可以通过
```bash
go get github.com/redis/go-redis/v9
```
去获取这个包，然后就可以在 go 里面愉快的操作 redis 了，API使用[教程](https://redis.ac.cn/docs/latest/develop/clients/go/)，下面是一个结构体 -> redis 数据的例子，我们可以使用这个方法来将 MySQL 中的结构化数据缓存到 Redis 中：
```go
package main

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
)

type Role struct {
	Role   string
	Gender string
	Age    int
}

func main() {
	cli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx := context.Background()
	cli.HSet(ctx, "魔裁", "樱羽艾玛", "粉色小狗")

	res := cli.HGet(ctx, "魔裁", "樱羽艾玛")

	fmt.Println(res.Val())

	cli.HDel(ctx, "魔裁", "樱羽艾玛")

	Role := Role{
		Role:   "樱羽艾玛",
		Gender: "粉色小狗",
		Age:    18,
	}

	cli.HSet(ctx, "樱羽艾玛", map[string]interface{}{
		"Role":   Role.Role,
		"Gender": Role.Gender,
		"Age":    Role.Age,
	})

	re := cli.HGetAll(ctx, "樱羽艾玛")

	mp, err := re.Result()
	if err != nil {
		fmt.Println(err)
	}
	// 一个将 map 对象反序列化为结构体的工具包
	mapstructure.Decode(mp, &Role)

	fmt.Printf("%+v\n", Role)
}
```
Redis 的客户端本身也提供了许多功能，比如 opt 里面的参数可以进行许多优化和扩展，但是这里不过多讲解（因为我也没咋用过），更多 API 还请自己探索。

我们日常使用 Redis 大多是用作缓存，问题来了，我们第一次读取 MySQL 数据的时候，希望将这个数据缓存到 Redis 中，之后的读取就全是使用 Redis 来读取了，那么当有没有 Redis 中缓存的数据与 MySQL 的数据不一致的时候呢？当我们需要写数据的时候，我们肯定需要去写 MySQL 了，那么 Redis 中的数据应该如何处理呢？你们可能觉得很简单，直接修改 Redis 里面的数据不就行了？大多数人最开始做缓存策略可能都是这样做的，但是这样其实还是会导致不一致的问题：
```
请求1：写 MySQL
请求2：写 MySQL
请求2：更新 redis 缓存
请求1：更新 redis 缓存
```
如果出现上面这种情况，那么就会导致缓存和数据库的数据不一致，在平常自己测试的时候可能很难发现，但是并发度上来了，这种问题就会非常明显。

所以最常见，最简单的做法是写后删除，也就是写完 MySQL 的数据之后，删除 redis 中的数据，缺点是之后又需要重新访问数据库来重建缓存（但是可以用 [singleflight](https://zhuanlan.zhihu.com/p/382965636) 来优化一下）这种情况也有很小的情况会引起数据不一致，如果当时 Redis 没有缓存数据时：
```
请求1：读请求，发现 Redis 没有数据，读 MySQL
请求2：写 MySQL
请求2：删除 Redis 缓存
请求1：重建 redis 缓存
```
虽然也会导致不一致，但是概率很小，所以实际上大多数人都是直接用的写后删除策略 + 合适的 TTL，那么有没有完全避免数据不一致的策略呢？是有的，据我所知：
1. 延迟双删：写后先删除一遍，然后过一段时间又删除一遍，当然，实现起来比较复杂，依赖了延迟这段时间，语义上防止比当前写请求还旧的读请求还残留在请求路径上。
2. 版本号机制：实现起来比较简单，使用单调增的版本号来防止读旧数据，每次写 MySQL 都将版本号写回到 Redis，这个版本号由于单调增的机制，不会出现旧数据覆盖新数据的情况，所以可以放心写，这样就不用担心旧数据的情况了。

上面还提到了重建缓存是有一定开销的，可以想象一下，在短时间内访问量比较大的时候，如果此时 Redis 中还没有缓存，那么就会有多个请求尝试去重建缓存，这里就会引起不必要的开销，因为实际上我们只需要有一个请求去重建缓存就可以了，剩下的只需要等待，这里就是典型的狗堆效应的问题，最简单的方法就是加锁，分布式场景下就用分布式锁，我之前看见个博客讲的挺好的，现在找不到了，还有个优化策略就是用上述的 singleflight。

当然，除了缓存，我们的 Redis 还可以干其他很多事情，比如 bitmap 可以用来做签到/阅读/点赞数据的统计，可以很好的节省空间，geo 可以用来做地图中的“附近的店”的功能，`SET NX` 命令可以用来做分布式锁，redis + lua 脚本可以做一个分布式限流器，但就目前而言，基本不会用到这些，有兴趣可以自己了解。

更多 Redis 的知识，我的建议是看[黑马点评](https://www.bilibili.com/video/BV1cr4y1671t)，**别看实战篇，有 Java 你看不懂**，这些应用你直接看别人总结的博客就行了，里面你还可以自己配一个 Redis 集群，了解一下  Redis 的一些底层设计（有些涉及操作系统，看不懂就算了），总之，黑马点评这个视频确实不错，别看实战篇就行了，除此之外没有一点 java 的内容。


---

# MongoDB

## MongoDB 的优势

MongoDB 是一个开源的文档数据库，它采用 BSON 格式存储数据，数据存储以文档为单位。与传统的关系型数据库相比，MongoDB 在数据存储方面具有更强的灵活性和扩展性。

MongoDB 不需要事先定义数据模式（schema），也就是说，存储的数据可以是不同结构的，这使得它特别适合存储大量半结构化或非结构化的数据，比如日志文件、社交网络数据等。

---

## 基础概念

- **文档存储**：MongoDB 将数据以文档的形式存储，文档是由键值对组成的对象，通常采用 BSON 格式。这种格式可以非常灵活地存储复杂的数据结构，如嵌套数组和子文档。
    
- **集合**：MongoDB 中的文档是存储在集合（collection）中的，类似于关系型数据库中的表格。不同于传统数据库，MongoDB 的集合不需要事先定义结构，因此在数据存储上更为灵活。
    
- **高可扩展性**：MongoDB 具有内建的分片支持，能够自动将数据分布到不同的服务器上，实现横向扩展。
    
- **强大的查询能力**：MongoDB 提供了丰富的查询语言，支持复杂的查询操作，如聚合查询、文本搜索、地理位置查询等。
    
### Go 操作 MongoDB

mongoDB 我不太熟，我也不咋用，不好直接给你们讲，其实 MySQL + Redis 够用了，本身这节课主要讲的是 Redis ，这里就只讲如何在 go 里面操作 mongoDB 了，例子：
```go
package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		panic(err)
	}
	collection := cli.Database("test").Collection("users")
	// 增
	result, err := collection.InsertOne(ctx, bson.M{"name": "John Doe", "age": 30})
	if err != nil {
		panic(err)
	}
	// 改
	resp, err := collection.UpdateOne(ctx, bson.M{"_id": result.InsertedID}, bson.M{"$set": bson.M{"age": 31}})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.ModifiedCount)
	// 查
	var user bson.M
	err = collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&user)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	// 删
	res, err := collection.DeleteOne(ctx, bson.M{"_id": result.InsertedID})
	if err != nil {
		panic(err)
	}
	fmt.Println(res.DeletedCount)

}
```
默认是使用的 `_id` 当作主键，当然，默认生成的 `_id` 虽然像一个无规则字符串，但是他其实是有规律的，他是是一个 16 进制的数字，默认为主键索引，你可以理解为比较自由的 MySQL 即可，而且他的 collection 里面的不同对象的字段可以不一样。

mongo 相对来说很灵活，不需要你提前建表啥的，同时也有 ACID 的特性，也有事务，bson 的字段也符合我们的 go 里面结构体的 json 直觉的，总之操作起来算是比较方便的，在比如文档，结构多变的数据可以选择 mongo 存储，但是我其实用得都不多。

但是其实本质这些关系型数据库和非关系型数据库你们现在掌握的也是增删查改，对你们来说应该只有 API 和数据结构操作不一样的区别，API 用到的时候再去问 AI 都可以，重要的是怎么用他们，用他们干什么。

## 作业

1. 尝试在 linux 上面部署 Redis（docker，虚拟机，物理机都可以），尝试用 go-redis 敲一下各种数据结构的增删查改（无需提交）
2. 给你们之前写的 todolist 或者其他的 web 项目加一个 Redis 缓存，注意你的缓存策略，学习一下缓存击穿，缓存雪崩，缓存穿透这几个常见问题和解决办法。
3. （选做）了解如何用 redis 实现分布式锁和分布式限流。
4. （长期任务）看完[黑马点评](https://www.bilibili.com/video/BV1cr4y1671t)的非实战篇部分，同时看别人的实战篇的博客总结，比如[这个](https://cyborg2077.github.io/2022/10/22/RedisPractice/)，但是这个人的博客挺卡的，可以换一个看，**重点在应用而不是代码**。