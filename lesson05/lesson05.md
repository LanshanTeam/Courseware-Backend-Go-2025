# 前言

大家在之前的学习中，熟悉了 go 的基本语法，基本的数据结构，面向对象编程以及常用的包，那么这节课则是介绍 Go 在并发编程方面的一些特性，同时会浅浅引入一些操作系统的知识（操作系统是个很宏大的话题，初学不理解很正常，之后有兴趣可以系统地学习），主要都是理解性的内容，代码含量很少😋，实在不理解也没关系，目前只要会用 goroutine 这些，知道它们是干啥的就可以了。

这节课主要讲的就是 Goroutine，channel，以及一些同步原语，但是在此之前，还需要介绍一下操作系统是个什么东西。

# 1. 操作系统是什么？

## 1.1 抽象

你们前面学过抽象的概念，抽象，类比一下就是将一个复杂的对象或操作给包装起来，而只暴露关键部分，比如说我们希望使用 typora 或者记事本对一个文档进行编辑，此时我们只需要考虑**输入文本内容**这个行为，就能够编辑我们期望的内容。

对比一下，如果我们使用编程语言对文件内容进行编辑，需要做的事情就多了，我们要先 read 这个文件的内容，然后观察，我们应该在哪里进行文本信息的插入，然后计算这个执行位置的偏移量，然后对字符串进行操作随后进行覆盖写，或者我们希望对文本文档进行追加操作，我们也需要调用 write 启用 append 模式：

```Go
    f, err := os.OpenFile("myfile", os.O_APPEND | os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
        return
    }
    _, err = f.Write([]byte("追加内容"))
    fmt.Println(err)
```

相比之下，文本编辑器就少了很多复杂的东西，我们并不需要去考虑每个文字的偏移量，只需要用鼠标点到想要插入的区域然后进行输入就行了。

这个时候你可能会联想到其他的一些应用，比如在线协作文档等一系列应用，如果没有在线协作文档，你可能需要通过其他的一些通讯软件将文档发送给别人，让别人修改一下，然后反馈给你，这样一来一往非常麻烦，于是有了在线协作文档这种抽象帮我们屏蔽了一些细节，可见，在我们的日常生活中，**抽象无处不在。**

## 1.2 操作系统

现在你可能对抽象有了一定了解，那么回到我们的计算机。

我们的电脑主要由存储器（磁盘和内存），CPU，输入输出设备等组成。而**操作系统**正是对这些底层的硬件进行封装，主要以系统调用的形式将这些硬件操作暴露给用户，比如我们常用的 open，write，read 都是常见的系统调用（当然，go 里面的文件读写都是经过一些封装的，但是依旧是使用的这些系统调用）。现在我们所看见的 QQ，微信或者游戏一系列应用，都是建立在操作系统提供的这些系统调用大厦之上的，可知操作系统的重要性……

不同操作系统提供的系统调用基本都差不多，但也有一些区别，也正是有这一层原因，有的应用程序会有多个操作系统版本。

在我们蓝山之后的学习中，很少会有直接用系统调用去编程的机会，大多数情况都是用封装好的包来编程，但是了解这些底层知识还是比较重要的。当然如果你是自学小能手，这样的机会有很多😍

## 1.3 并发

前面只是热身，操作系统的本事远不止提供系统调用。回到本节的主角——**并发**。

**先说进程。** 进程是一个正在运行的程序实例。操作系统会为每个进程分配彼此隔离的用户空间，并把它们交给 CPU 调度执行。一个应用可以只开一个进程，也可以开很多个，这取决于应用自身；对操作系统来说，进程之间没有啥区别，它只负责把要跑的活儿交给 CPU。

**那一颗 CPU 怎么“同时”跑这么多东西？** 关键在于“时间片”。单个 CPU 核心同一时刻只能执行一个任务，但操作系统会把时间切成很细的片段，轮流把这些时间片分给不同的可运行实体。每个时间片到了就切换到下一个，如此快速轮转，让我们主观上感觉“所有程序都在同时跑”。

**再说线程。线程是操作系统调度的最小单位**。一个进程至少包含一个线程；所谓“调度进程”本质上是在**调度进程里的线程**。同一进程中的线程共享进程的资源，但各自拥有独立的栈（栈主要是负责存局部变量和函数调用栈等数据的）。

线程大致分两类：
- **内核级线程**：创建、销毁、调度与上下文切换由内核完成，需要系统调用，开销相对更高。
- **用户级线程**：在线程库/运行时于用户态完成调度与切换，通常只在必要时与少量内核线程关联，因而切换开销更小。

**而我们今天的主角 goroutine 就是一种用户级线程** 它由 Go 运行时在用户态调度，创建成本低、切换轻量，还配有按需增长的栈。运行时采用 M:N 调度（后面可以了解一下 **GMP 模型**），把大量 goroutine 映射到少量内核线程上，从而在保持高并发的同时控制系统调用与上下文切换的成本。

**PS：** 虽然讲了这么多，你们可能都没听懂，但是没关系，上完这节课，以上的全部都没有掌握也没关系！

下面正式介绍我们怎么用这玩意！

# 2. Go 并发编程
## 2.1 Goroutine

前面花了大量的篇幅介绍了背景知识，终于到我们的代码环节了，其实我们的 goroutine 使用起来非常简单，比如我们想要实现一下用 10 个协程去累加一个数字，每个协程将这个变量自增 10w 次，我们期望得到 100w 的结果

```Go
package main

import (
    "fmt"
    "time"
)

func main() {
    sum := 0
    for range 10 {
        go func() {
            for range 100000 {
                sum += 1
            }
        }()
    }
    time.Sleep(1 * time.Second)
    fmt.Println(sum)
}
```

结果呢？

```Go
root@RINAI-SWORD:/project/lanshan/count# go run ./
569738
root@RINAI-SWORD:/project/lanshan/count# go run ./
503786
root@RINAI-SWORD:/project/lanshan/count# go run ./
595125
```

得到是不确定的数字，可能有人觉得奇怪，但是在现实中，这才是符合预期的，你如果把数值调小，让多个线程累加到 100 或者 1000，你的结果可能正是你期望的，但是仅仅是因为样本数量太少而没有引发多线程出现的问题。

那么为什么会出现这种结果？

这涉及到我们 cpu 的工作原理，别担心，没那么复杂🥰，我们的一条自增 ++ 代码，在 cpu 层面并不是原子执行的，实际上，一个简单的自增在计算机层面会是这样的顺序：

```Bash
1. cpu 从指定内存读取 sum 的值
2. cpu 将 sum + 1。
3. cpu 将 sum 写回指定内存。
```

大概是这三个步骤，那么在多线程场景下思考，只要在 1 -> 2 的间隙发生线程切换，这里线程切换会将当前线程中的某些上下文数据保存下来到内存中，等到切换回线程1的时候就会从内存中恢复这些数据到 cpu 中，在另一个线程中执行了 sum + 1，就会导致我们当前 cpu 读取的数值是旧的，比如：

```Bash
线程1：cpu 读取到 sum = 10。
线程2：cpu 读取到 sum = 10。
线程2：cpu 将 sum + 1。
线程2：cpu 回写 sum，此时 sum = 11。
线程1：cpu 根据读取到的 sum + 1。
线程1：cpu 回写 sum，此时 sum = 11.
```

由此可见，多线程虽然利用多核的优势，但是也给我们带来了一些麻烦。

针对于赋值这些操作，go 标准库提供了原子操作库 `sync/atomic` 对于一些变量，可以原子性的进行更改：

```Go
package main

import (
    "fmt"
    "sync/atomic"
    "time"
)

func main() {
    sum := atomic.Int64{}
    sum.Store(0)
    for range 10 {
        go func() {
            for range 100000 {
                sum.Add(1)
            }
        }()
    }
    time.Sleep(1 * time.Second)
    fmt.Println(sum.Load())
}
```

结果：

```Bash
root@RINAI-SWORD:/project/lanshan/count# go run ./
1000000
```

符合预期，这个 `sync/atomic` 库还支持 CompareAndSwap 操作，也就是让一个变量的值可以实现原子的切换，由于我们的指针基本可以表示任何数据，所以**从理论上**我们可以通过切换指针来实现对自定义数据结构的原子替换，保证了线程安全，但是实际上开发日常基本不会用，但是 CAS 在一些 go 官方库的源码中用的会比较多，所以还是得注意一下。

  

## 2.2 互斥锁

当然，我们无法依靠原子操作打遍天下，所以这里还需要介绍一个常用的标准库 `sync`，在这个包里面，我们可以通过一个叫锁的东西去保护我们的变量：

```Go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    sum := 0
    // 互斥锁
    lock := sync.Mutex{}
    for range 10 {
        go func() {
            for range 100000 {
                lock.Lock()
                sum += 1
                lock.Unlock()
            }
        }()
    }
    time.Sleep(1 * time.Second)
    fmt.Println(sum)
}
```

这里通过锁来保护我们的 sum 变量在自增过程中不会进行切换，那么锁是什么？为什么能够保护 sum 这个变量？

同一时刻只有一个线程能够拿到同一把锁，然后才能继续执行接下来的代码，这就是一个临界区域，当锁被其他线程占有的时候我们调用 Lock，他无法拿到这把锁，也就无法向下继续执行代码，然后会自旋几次然后将这个线程挂起，不消耗 cpu 资源（这部分详细见 Lock 源码和调度器），当这个锁被释放时，陷入睡眠的线程就会被唤醒然后重新去持有这把锁，这里涉及到一个概念叫**惊群效应**，可以去了解一下。

锁的设计保证了关于临界区的操作执行是原子的，并且同一时刻只有一个线程能够持有，所以可以很好的保证我们 sum ++ 的安全性！

引申一下，在淘宝购物，双11秒杀场景下，我们需要保证大量的用户同时抢购一个商品的准确的库存扣减应该怎么做？在我们 go 语言内部，模拟一下这个场景很简单，加个锁或者原子操作就行了：

```Go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type Product struct {
    mu sync.Mutex
    Name string
    Price int
    Stock int
}

func (p *Product) DecrStock(count int) (bool, int) {
    p.mu.Lock()
    if p.Stock >= count {
        p.Stock -= count
        p.mu.Unlock()
        return true, p.Stock
    }
    p.mu.Unlock()
    return false, p.Stock

}

type User struct {
    mu sync.Mutex
    Name string
    buyed int
}

func (u *User) Buy(p *Product, count int) {
    ok, reserved := p.DecrStock(count)
    if ok {
        fmt.Printf("%s 成功购买了 %v 个 %s，剩余 %v\n", u.Name, count, p.Name, reserved)
        // 这里是为了保护 buyed
        u.mu.Lock()
        u.buyed += count
        u.mu.Unlock()
    } else {
        fmt.Printf("%s 想要购买 %v 个 %s 失败，剩余 %v\n", u.Name, count, p.Name, reserved)
    }
}

func main() {
    pro := &Product{
        mu: sync.Mutex{},
        Name: "《卧樱》",
        Price: 10,
        Stock: 2000000,
    }

    u1 := &User{
        Name: "草薙直哉",
    }

    u2 := &User{
        Name: "大饼老师",
    }

    for range 20000 {
        b1 := rand.Int() % 3 + 1
        go u1.Buy(pro, b1)
        b2 := rand.Int() % 3 + 1
        go u2.Buy(pro, b2)
    }
    // 这里推荐改成 waitGroup，方便演示直接 sleep
    time.Sleep(5 * time.Second)
    fmt.Println("总计购买" , u1.buyed + u2.buyed, "当前库存", pro.Stock)
}
```

这里仅作演示，有问题欢迎指出。

除了 `sync.Mutex` 还有一种另一种锁 `sync.RWmutex` 也就是读写锁，因为在只有并发读，没有并发写的时候，再加锁并没有什么用，因为并发读始终是安全的，因此读写锁就是针对于互斥锁的一个优化，在读多写少的场景下我们可以选择读写锁来替代互斥锁。

## 2.3 Channel

典中典之"**不要通过共享内存来通信，而要通过通信来共享内存**"，这句话便是 Go 语言的并发哲学，什么意思？就是说，如果两个 goroutine 希望共享一个变量，不应该通过一个外部的全局变量来进行加锁读写，而是应该通过 channel 将 goroutine A 中的变量传递给 goroutine B。

其实个人觉得这玩意只能算是一种写代码的方法，并不是说所有 go 程序都应该这么写，共享内存也有其优点，不能一棒子打死。

看看 channel 的使用，下面是一个任务调度器：

```Go
package main

import (
    "fmt"
    "time"
)

type Task struct {
    // 函数也可以是结构体的成员
    Runnable func(workerId int)
}

func main() {
    // 一个负责任务分发的管道
    ch := make(chan Task, 10)

    // 启动几个 worker 负责处理任务
    for id := range 10 {
        go func(workerId int) {
            for t := range ch {
                t.Runnable(workerId)
            }
        }(id)
    }

    // 任务分发
    for i := range 20 {
	    j := i
        t1 := Task{
            Runnable: func(workerId int) {
                fmt.Printf("workerId%v：task%v做一件事情\n", workerId, j)
            },
        }
        ch <- t1
    }
    time.Sleep(1 * time.Second)
    close(ch)
}
```

Channel 这东西说常用，我其实也不是经常用，但是依旧有他的使用场景，比如协程池，上面的任务调度器的 worker 其实就是一个最简单的协程池，他可以限制 goroutine 的数量，防止突发情况下，有大量的购买请求到来，产生大量的 goroutine，即便我们的协程比较轻量，但也扛不住过大的量，而且也会影响调度的效率（以后可以了解 GMP 模型）

## 2.4 其他同步原语

刚刚讲的 goroutine，互斥锁和 channel，你如果今天能理解其实也就差不多了。下面是扩展部分

### 2.4.1 WaitGroup

这个其实还是挺重要的，他的作用就跟他的名字一样，可以等待一组 goroutine 执行完成，怎么做到的？来看一段代码：

```Go
package main

import (
    "fmt"
    "sync"
)

func main() {
    wg := sync.WaitGroup{}
    for range 10 {
        // wg.Add 代表添加一个执行任务
        wg.Add(1)
        go func() {
            fmt.Println(1)
            // wg.Done() 代表执行完成
            wg.Done()
        }()
    }
    // 等待 Add 的任务全部 Done
    wg.Wait()
}
```

其实就是通过计数器来保证了等待 goroutine 执行完成，虽然我们自己也能够用锁来实现这个计数器，但是 wg 确实是一个方便的工具，有了它，我们就不需要再执行很多 goroutine 的时候，为了方便查看结果，总是使用 time.Sleep 了。

### 2.4.2 sync.Once

程序运行时，懒加载初始化常用（其实平时也不咋用这个），保证了整个生命周期，一个 sync.Once 对象仅会执行一次 Do 方法，多余的都会跳过。

比如：

```Go
package main

import (
    "fmt"
    "sync"
)

func main() {
    init := sync.Once{}
    num := 0
    for range 10 {
        init.Do(func() {
            num ++
        })
    }
    fmt.Println(num)
}
```

结果为 1，无需多言。

### 2.4.3 sync.Map

这个在 go 的八股文里面还是比较常考的吧，所以我也一起讲了

众所不周知，我们 go 中原生的 `map` 并不是并发安全的，在并发读写的情况下会 panic，导致程序崩溃，于是 sync 标准库推出了 `sync.Map` ，这是一个并发安全的 `map`，严格意义来说，我们直接在原生 `map` 的基础上加互斥锁或者读写锁也能够解决问题，但是 `sync.Map` 是做了一定基础的优化的，go 1.24 之前的版本和之后的版本有不同的实现方式，感兴趣可以去了解一下。总之，有了它，我们就不需要为 `map` 维护一个互斥锁了。

### 2.4.4 sync.Cond

它可以实现一种信号通知的功能（很少用，但是蓝妹里面用到了😋）
例子：
```go
package main

import (
	"sync"
	"time"
)

func main() {
	l := sync.Mutex{}
	c := sync.NewCond(&l)

	go func() {
		// 调用 wait 前需要先加锁
		c.L.Lock()
		c.Wait()
		println("Hello, world!")
		// 调用 wait 后需要解锁
		c.L.Unlock()
	}()

	time.Sleep(time.Second)
	println("唤醒")
	// Signal 方法唤醒等待的 goroutine
	c.Signal()

	time.Sleep(time.Second)
	println("end")
}

```

剩下的不常用，用到的时候再说吧。

参考链接：https://draven.co/golang/docs/part3-runtime/ch06-concurrency/golang-sync-primitives/

# 3. 基础语法补充内容

到了上面其实基础语法和 go 的特性就讲完了，但是有些内容没讲到，覆盖并不全，因此我塞了个补充内容部分用于扫盲，仅作了解即可，前期写 web 的话基本上没啥用。

## 3.1 Context 上下文

context 可以在函数传播链中用于存储一些 kv 键值对信息用于下游获取，也可以在并发控制中发挥作用，比如我们可以为这个函数调用设计一个超时时间，我们的 context 就可以通过这个超时时间来取消这个函数调用链，除了超时控制，我们还可以手动地去 cancel 这个上下文，取消这次调用，但是需要注意的是，你即便取消了这个 context,已经执行的代码并不能撤回，谨慎设置超时时间。
例子：
```go
package main

import (
	"context"
	"fmt"
	"time"
)

func doSomething(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("doSomething stopped:", ctx.Err())
			return // 退出 goroutine
		default:
			fmt.Println("hello")
			time.Sleep(time.Second)
		}
	}
}

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)

	defer cancel()

	doSomething(ctx)

}
```

详细讲解
https://draven.co/golang/docs/part3-runtime/ch06-concurrency/golang-context/

## 3.2 timer 定时器

time 包里面有个比较好玩的没讲，就是 Timer 和 Ticker：
```go
package main

import "time"

func main() {
	t := time.NewTicker(1 * time.Second)
	i := 0
	for {
		// t.C 按照每秒一次的频率向 channel 发送时间信号
		// 而 timer 的 t.C 方法只会发送一次时间信号，如果希望持续发送，则需要通过 Reset 方法重新设置时间
		<-t.C
		println("Hello, world!")
		i++
		if i == 3 {
			break
		}
	}
	t.Stop()
}

```

包括 timer 的底层及其 API 的详细讲解
https://draven.co/golang/docs/part3-runtime/ch06-concurrency/golang-timer/

## 3.3 select 关键字

我们 go 里面的 select 和操作系统里面的 select 系统调用比较相似，我们 go 里面的 select 可以同时等待多个 channel 传来信号，但是最终只会执行一个：
```go
package main

import "time"

func main() {
	chan1 := make(chan int)
	chan2 := make(chan int)

	go func() {
		time.Sleep(2 * time.Second)
		chan1 <- 1
	}()

	go func() {
		time.Sleep(1 * time.Second)
		chan2 <- 2
	}()

	// select 关键字用于在多个通信操作中进行选择，如果多个通道同时接收到数值，会随机选择一个进行接收。
	select {
	case msg1 := <-chan1:
		println("Received from chan1:", msg1)
	case msg2 := <-chan2:
		println("Received from chan2:", msg2)
	case <-time.After(3 * time.Second):
		println("Timeout occurred")
	}
}

```
select 还有一个 default 分支，没有接收到数据的时候会直接走这个。
select 详解：
https://draven.co/golang/docs/part2-foundation/ch05-keyword/golang-select/

## 3.4 defer 关键字

我发现前面没讲过这个，在这里补充一下
defer 常见的使用就是在函数结束后执行一些收尾工作，比如我们在打开文件之后，在程序结束的时候需要关闭这个文件，以及在我们创建一个临时管道的时候需要在函数退出的时候能够及时关闭这个管道。

另外 defer 也是一个比较常考的八股文，比如，在以下两个函数中，函数的返回值会是什么呢？为什么会有区别，这就涉及到 defer 和 return 的执行过程：
```go
package main

import "fmt"

func returnValue1() int {
	a := 1
	defer func() {
		a += 1
	}()
	return a
}

func returnValue2() (a int) {
	a = 1
	defer func() {
		a += 1
	}()

	return a
}

func main() {
	fmt.Println(returnValue1()) // 1 
	fmt.Println(returnValue2()) // 2
}

```
在 return 的时候，如果函数签名没有定义返回的变量，那就会定义一个临时变量，或者说帮你填充了一个你不知道的变量名，然后将返回的数据赋予这个临时变量，如果你已经定义好了变量名，那就将返回的数据赋予给这个定义的变量；第二步，我们会执行 defer 的操作；最后，我们真正将这个临时变量/已经定义的变量返回回去。
也就是说，我们在函数签名已经定义好了变量名的时候，在最后会将这个变量名的值返回，如果没有，那么就会首先定义一个临时变量，然后赋予它返回值，再执行 defer 的函数，最后将这个临时变量返回。
大概是这样的：
```go
// 如果定义了返回值变量名
a = 返回的数据
defer xxxx // 此时对 a 执行修改操作会反馈到 return 值上面
return a

// 如果没有定义变量名
tmp := a
defer xxx // 此时如果对 a 进行修改无法反馈到 tmp 上面
return tmp
```
简单来说分为三步，就是赋值，执行defer，返回值。

# 4. 作业

### Lv1 

  动手敲一下课件里面的代码，尝试理解一下（不用提交）

### Lv2 

分几个简单的阶段完成这个作业：

- 阶段1：使用指定数量的协程完成对某个变量的指定数量的自增操作，同时保证并发安全，在自增一定次数之后，需要得到期望的结果。
- 阶段2：将一个自增操作封装成 task ，将自增任务通过一些策略**分发**给 goroutine 进行执行（课件里面有类似的例子）
- 阶段3：把你现在实现的 task -> goroutine 的模式封装成一个工具包，使得任何外部的库都可以通过你的包提供的函数/方法进行 task 分发。（PS：考虑一下提供什么样的 API 让用户使用起来更方便）

进阶要求（**选做**）：
- 性能优化：进行一些性能优化，我们的协程池还可以如何改造？（一些提示：sync.Pool，多管道轮询分发 task，属于是课外扩展的内容，不懂可以问 AI）
- 思考：做完性能优化之后，尝试思考或者百度这几个问题：
    - 为什么需要协程池？
    - sync.Pool 能拿来干什么，为什么需要池？
    - 原生的 map （前面课有讲）不是并发安全的，了解 1.24 之前及之后的 sync.Map 的实现，他是怎么实现并发安全的？

### Lv3 写个小项目
**项目一：多线程文件搜索工具**
实现一个多文件内容搜索的命令行工具，运行格式：
```bash
go build -o catch ./
./catch [需要检索的目录] [需要检索的关键词]
```
基本要求：
1. 针对的是多个文件的搜索，所以需要先递归搜索所有文件的路径，并把每个待搜索的文件抽象成一个 task，然后交给 goroutine 进行搜索。
2. 使用你在 Lv2 中实现的协程池来进行搜索，而不是一个文件一个 goroutine。
3. 要求匹配到关键字段的时候，将对应整行的内容**规范地**输出出来。

进阶要求（**选做**）：
4. 你可能是在 goroutine 里面直接读取文件的内容然后打印出来，我希望你能够通过 waitgroup 等待所有的任务完成，然后通过 channel 接受每个任务的搜索结果，最后统一打印出来。
5. 了解一下 [cobra](https://github.com/spf13/cobra) ，使用 cobra 将这个命令行工具升级一下！
6. linux 里面有一个搜索命令叫做 grep，这里不需要你去实现，简单了解一下 grep 有什么更高级的功能，思考一下怎么实现。

**PS：Lv3 不一定要写这个文件搜索工具，如果你有其他好点子也可以写😋**

### Lv4 （选做）

**了解**并发导致的常见多线程模型，比如：哲学家吃饭问题，生产者消费者模型，读者写者问题，理发师问题……（网上资料很多，到处都能搜到这种常见的并发模型，了解一下这些问题可以帮助你更好的理解并发以及并发带来的问题）

### Lv5 （选做）

将你在本节课学到的知识通过 AI 或者别人的博客巩固一下，整理成一篇博客发布到平台上（也可以发到个人博客上，重要的是思考了，然后动手写，然后发布出去能有一个正反馈）

### Lv6 （选做）

未来一两年的长期目标了这是，操作系统的知识对于后端仔还是比较重要的，所以你们可以系统的学习一下操作系统，但是在此之前请先有一定的计算机系统基础，学习操作系统个人推荐 [MIT 6.S081: Operating System Engineering - CS自学指南](https://csdiy.wiki/%E6%93%8D%E4%BD%9C%E7%B3%BB%E7%BB%9F/MIT6.S081/) 这门课程，并实现相关 lab，或者其他的课程也行，这个作业你们自由发挥，不做强制要求。

----
# 4. 结语
  完成前三个作业**非选做部分**就算做完了，作业提交至 chenyue@lanshan.email，前三个完成基本要求的同学奖励奶茶一杯，可以向 AI 寻求帮助，但是请不要直接 copy AI 的代码，请自己思考理解之后再写。

作业部分虽然很多，但是仅仅是给你们一个相对长期的学习方向，并不是说要在下一节课或者本学期之前就得做完，前三个作业的主要目的就是考验你们关于 goroutine 和 channel 的应用，如果能够完成的话，其实这节课也就学的差不多了，希望更进一步了解的可以看看选做部分，也可以接着预习下一节课 web 的内容，每年人走的最多的就是 web 这一节课，所以不要忽视了它～