# 包

### 包介绍

首先新建go的包管理文件go.mod(Goland新建项目的时候自动做了)

```shell
go mod init lesson03
```

Go语言中支持模块化的开发理念，在Go语言中使用`包（package）`来支持代码模块化和代码复用。

包的引入使得我们可以去调用自己或者别人（在Github等代码托管平台上）的模块代码，方便了我们的开发。

例如，在之前的课件中，我们引入了 fmt 这个包。这样使得我们可以调用 fmt 包内部的函数和变量。

```go
package main

import "fmt"

func main(){
  fmt.Println("Hello world!")
}
```

接下来详细介绍一下

### 定义包

我们可以根据自己的需要创建自定义包。一个包可以简单理解为一个存放`.go`文件的**文件夹**。

该文件夹下面的所有`.go`文件都要在非注释的第一行添加如下声明，声明该文件归属的包。

| 常用包名             | 用途示例                                     |
| -------------------- | -------------------------------------------- |
| **config**           | 配置加载与管理                               |
| **logger**           | 日志初始化与封装                             |
| **middleware**       | HTTP / gRPC 中间件                           |
| **common** / **pkg** | 公共常量、错误码、工具函数（可能导致高耦合） |
| **utils**            | 工具函数（字符串处理、时间格式化等）         |
| **consts**           | 常量分组，如错误码、Redis key、表名等        |

另外需要注意一个文件夹下面直接包含的文件只能归属一个包，同一个包的文件不能在多个文件夹下。

包名为`main`的包是应用程序的入口包，这种包编译（`go build`）后会得到一个可执行文件，而编译不包含`main`包的源代码则不会得到可执行文件。

~~以后比较复杂的作业记得分包哦~~

### 可见性

在同一个包内部声明的标识符都位于同一个命名空间下，在不同的包内部声明的标识符就属于不同的命名空间。想要在包的外部使用包内部的标识符就需要添加包名前缀，例如`fmt.Println("GenShin Impact")`。

如果想让一个包中的标识符（如变量、常量、类型、函数等）能被外部的包使用，那么标识符必须是对外可见的。在Go语言中是通过标识符的首字母大/小写来控制标识符的对外可见/不可见的。在一个包内部只有**首字母大写**的标识符才是对外可见的。

这也会在我们后面的**封装**中体现

例如我们定义一个名为`demo`的包，在其中定义了若干标识符。在另外一个包中并不是所有的标识符都能通过`demo.`前缀访问到，因为只有那些首字母是大写的标识符才是对外可见的。

```go
var	Name  string // 可在包外访问的方法
var	class string // 仅限包内访问的字段
```

### 包的引入

要在当前包中使用另外一个包的内容就需要使用`import`关键字引入这个包，并且import语句通常放在文件的开头，`package`声明语句的下方。完整的引入声明语句格式如下:

```go
import importname "path/to/package"
```

其中：

- importname：引入的包名，通常都省略。默认值为引入包的包名。（如果很多相似的包可以通过自定义包名来提高**区分度**）

- path/to/package：引入包的相对路径，必须使用双引号包裹起来。

- *`Go语言中禁止循环导入包`* 这个一定要注意，很有可能因为这个要**重新架构**

一个Go源码文件中可以同时引入多个包，例如：

```go
import "fmt"
import "net/http"
import "os"
```

当然可以使用批量引入的方式。

```go
import (
    "fmt"
  	"net/http"
    "os"
)
```

如果引入一个包的时候为其设置了一个特殊`_`作为包名，那么这个包的引入方式就称为匿名引入。

一个包被匿名引入的目的主要是为了加载这个包，从而使得这个包中的资源得以初始化。 被匿名引入的包中的`init`函数将被执行并且仅执行一遍。（多个将会按照声明顺序执行），`init`函数在main时会在匿名引入包的`init`后面执行

```go
import _ "github.com/go-sql-driver/mysql"
```

匿名引入的包与其他方式导入的包一样都会被编译到可执行文件中。



听到现在你可能还会有点疑惑，我们现在不还是在引用自己本地的包吗？左手倒右手有什么意思(x

可以用下面这个命令拉网上第三方仓库的依赖,然后就可以用别人造好的轮子啦～

```shell
go get "github.com/J1407B-K/buff"
```

也要记得定期`go mod tidy`一下，可以清理未用/补缺的包依赖

# 指针

> “指针就做指针的事”

任何程序数据载入内存后，在内存都有他们的**地址**，这就是指针。而为了保存一个数据在内存中的地址，我们就需要指针变量

Go语言中的指针操作非常**简单**，我们只需要记住两个符号：`&`（取地址）和`*`（根据地址取值）。

## 指针地址和指针类型

每个变量在运行时都拥有一个地址，这个地址代表变量在内存中的位置。Go语言中使用`&`字符放在变量前面对变量进行“取地址”操作。 Go语言中的值类型（int、float、bool、string、array、struct）都有对应的**指针类型**，如：`*int`、`*int64`、`*string`等。

**注意，指针类型不是原类型，指针类型只是地址**

取变量指针的语法如下：

```go
ptr := &v    // v的类型为T
```

其中：

- v:代表被取地址的变量，类型为`T`
- ptr:用于接收地址的变量，ptr的类型就为`*T`，称做T的指针类型。*代表指针。

举个例子：

```go
func main() {
	a := 10
	b := &a
	fmt.Printf("a:%d ptr:%p\n", a, &a) // a:10 ptr:0xc00001a078
	fmt.Printf("b:%p type:%T\n", b, b) // b:0xc00001a078 type:*int
}
```

我们来看一下`b := &a`的图示，~~一图传三代~~：

![](https://camo.githubusercontent.com/5920f54c36a7bcf807c6a5e87354b5610f828449bc179933f040f40037418c59/68747470733a2f2f7777772e6c6977656e7a686f752e636f6d2f696d616765732f476f2f706f696e7465722f7074722e706e67)

## 指针取值

在对普通变量使用&操作符取地址后会获得这个变量的指针，然后可以对指针使用*操作，也就是指针取值，代码如下。

```go
func main() {
	a := 10
	b := &a
  c := &b
	fmt.Printf("a:%d ptr:%p\n", a, &a) // a:10 ptr:0xc00001a078
	fmt.Printf("b:%p type:%T\n", b, b) // b:0xc00001a078 type:*int
  fmt.Printf("c:%p type:%T val:%d \n",c,c,**c) //c:0x14000052038 type:**int val:10 
}
```

**总结：** 取地址操作符`&`和取值操作符`*`是一对**互补**操作符，`&`取出地址，`*`根据地址取出地址指向的值。

变量、指针地址、指针变量、取地址、取值的相互关系和特性如下：

- 对变量进行取地址（&）操作，可以获得这个变量的指针变量。
- 指针变量的值是指针**地址**。
- *`对指针变量进行取值操作，可以获得指针变量指向的原变量的值`*。

so，有了指针之后我们就可以愉快的进行传递地址了(注意空指针引发的panic)，正常函数传值只会传递值的**副本**，但我们地址不怕啊，我们地址就算被copy一份，但还是指向那个值，**指针只是地址**

```go
// 值复制
func Update(a int){
	a++
}

func main() {
	a := 0 
	Update(a)
	fmt.Print(a)
}
```

```go
// 指针共用数据
func Update(a *int) {
	*a++		// 注意取值再自增
}

func main() {
	a := 0
	Update(&a) //传指针
	fmt.Print(a)
}

```

## “指针就做指针的事”

**Go语言中的指针不能进行偏移和运算**，是为了安全、可移植、便于 GC/优化，并鼓励用切片这种更高层的**抽象**

**内存安全 & 边界检查**

- C 风格的 `p+1`, `p+4` 很容易越界；Go 可以通过切片的 `len/cap + 边界检查` 把这类错误提前拦住。
- 指针算术会绕开这些保护，容易造成未定义行为。不符合Go的设计理念

**让 GC 更可靠**

- Go 的垃圾回收器需要**精确地知道哪些是指针**。任意算出来的“整数化指针”会让 GC 难以追踪，甚至漏标。

- Go 运行时可能**移动栈**（stack growth/shrink）。如果你手里拿的是“计算出来的地址”，而不是受运行时追踪的真正指针，栈移动后它就可能失效。

  ```go
  func f() {
      x := 123
      p := &x                 // p 是受 GC 追踪的合法指针
      addr := uintptr(unsafe.Pointer(p)) + 8 // 自己算地址（GC 不会追踪）
      // Your Work……
      // 这里可能触发栈扩容
      q := (*int)(unsafe.Pointer(addr))      // addr 可能已经失效！
      *q = 456                               // ❌ 未定义行为
  }
  
  ```

- 因此 Go 要求：**指针就做“指针的事”，算地址这类事尽量别做**；必须做时，用受控的 `unsafe` 入口。

**优化与别名分析**

- 指针随意乱跑会让编译器很难判断两块内存是否别名，从而抑制很多优化（内联、消除边界检查、逃逸分析等）。
- 限制指针运算能让编译器更自信地做优化，提升部分性能。

**可移植性与简洁性**

- Go 目标是“一份源码多平台”。不同架构字长/对齐差异下，裸指针算术更容易踩坑。
- 语法层面不提供指针算术，迫使大家用统一的切片/索引方式，代码风格与可读性更一致。

# 结构体，方法，接口——面向对象思想在Go中的浪漫诠释

## 什么是“类”

类（Class）是面向对象里的**自定义类型模板**，描述一组对象共享的**属性（字段）和行为（方法）**。关键要素：

- 字段/属性
- 方法/成员函数
- 可见性
- 构造/析构
- 静态成员
- 继承/多态（is-a）

```java
//Java
public class Person {
    // 字段
    private String name;  //可见性private    
    private int age;

    // 构造器
    public Person(String name, int age) {
        this.name = name;
        this.age = age;
        population++;
    }

    // 方法
    public void greet() {
        System.out.println("Hi, I'm " + name);
    }
  
    public static int population = 0; // 类变量（属于类，全局）

    // 封装
    public String getName() { return name; }
    public void setName(String name) { this.name = name; }
}

// 继承
public class Student extends Person {
    private String school;
    
    public Student(String name, int age, String school) {
        super(name, age);       // 调用父类构造器
        this.school = school;   // 子类自己的初始化
    }
    
    // 多态
  	@Override
    public void greet() {
        System.out.println("Hi, I'm a student named " + getName());
    }
}
```



Go语言中没有**类**的概念，也不支持**类**的继承等面向对象的概念,

但我们可以让结构体 + 方法 + 接口组成一个‘类系统’，具有更高的**扩展性和灵活性**。

## 前置知识

### 自定义类型

Go语言中可以使用`type`关键字来定义自定义如`string`、`int`、`bool`等的数据类型。

自定义类型是定义了一个***全新***的类型。我们可以基于内置的基本类型定义，也可以通过struct定义。例如：

```go
//将MyInt定义为int类型
type MyInt int
```

通过`type`关键字的定义，`MyInt`就是一种**新的类型**，它具有`int`的特性

### 类型别名

类型别名规定：本质上是同一个类型。就像一个孩子小时候有小名，这和他的名字都指向同一个人。

```go
type 类型的别名 = 类型名
```

### 类型定义和类型别名的区别

类型别名和原类型是同一种类型。自定义类型是一种全新的类型。

类型别名的类型只会在代码中存在，编译完成时并不会存在(发生在类型检查阶段)

> 类型别名 和 自定义类型 的意义

自定义类型：举个例子，我们想给 int 类型定义一个方法，但是又不想改变int本身的性质。可以基于内置的`int`类型使用type关键字可以定义新的自定义类型，然后为我们的自定义类型添加方法。

类型别名：想象一下你有一个非常长的类型名字，比如`map[int]string`，如果在代码中反复使用这个类型，那将会变得很啰嗦。

但是如果你使用类型别名来代替它，比如`Data`，那么你只需使用`Data`这个简短的名字就可以代替长长的类型名字了

## 结构体

Go语言中的基础数据类型可以表示一些事物的基本属性（int,float,string），但是当我们想表达一个事物的全部或部分属性时，这时候再用单一的基本数据类型明显就无法满足需求了，Go语言提供了一种自定义数据类型，可以封装多个基本数据类型，这种数据类型叫结构体，英文名称`struct`。

```go
alice_name := "alice"
alice_age := 16

// or 

type People struct{
		Name string
		Age  int
}

alice := People{
  	Name:"alice",
  	Age:16,
}
// 哪个更直观一目了然
```

同时，结构体只是**数据容器**（data structure）

### 结构体的定义

使用`type`和`struct`关键字来定义结构体，具体代码格式如下：

```go
type 类型名 struct {
    字段名 字段类型
    字段名 字段类型
    //…
}
```

其中：

- 类型名：标识自定义结构体的名称，在**同一个包**内不能重复。

- 字段名：表示结构体字段名。结构体中的字段名必须**唯一**。

- 字段类型：表示结构体字段的**具体类型**（int,float,string）。

### 初始化

 没有初始化的结构体，其成员变量都是对应其类型的零值,如string为空字符串，int为0，bool为false

  ```go
  type People struct {  
  	Name string    // "" 空字符串
  	Age  int8			 // 0
  }
  ```

 定义好结构体后，我们就可以初始化了，初始化有蛮多形式的，这里列举三种

```go
// 直接使用键值对赋值
bob := People{
		Name: "bob",
		Age: 16,
}

//或者用var声明
var bob People
bob.Name = "bob"
bob.Age = 16

// 字典顺序初始化(严格按照声明字段顺序)
alice := People{
  	"Alice", 		//Name
  	18,					//Age
}
```

### 嵌套结构体

一个结构体中可以嵌套包含另一个**结构体**或**结构体指针**,正因如此，声明也要嵌套式声明

```go
type Info struct{
		Name string
		Age int
}

type	Lanshaner struct{
		Info 			// Info Info 
		Group  string
}

//初始化时也要初始化内层结构体
func main(){
  kq := &Lanshaner{
    Info:Info{
      Name:"J1407B",
      Age:19,
    },
    Group:"后端Go组"
  }
  
  fmt.Println(kq.Info.Name) // "J1407B"
	fmt.Println(kq.Group)     // "后端Go组"
}
```

## 方法

Go语言中的`方法（Method）`是一种作用于特定类型变量的函数。这种特定类型变量叫做`接收者（Receiver）`。

**只有**特定的接收者变量才可以调用对应的方法，该方法会写入该类型的**方法集**（Method Set）

Go 的方法本质上就是一个带有“**隐式接收者变量参数**”的普通**函数**，只是语法上**绑定**到了某个类型上。所以方法的声明与函数类似，又有些许不同，见下

```go
func (接收者变量 接收者类型) 方法名(参数列表) (返回参数) {
    // Your Work
}
```

其中，

- 接收者变量：接收者中的参数变量名在命名时，官方建议使用接收者类型名称首字母的小写。例如，`Person`类型的接收者变量应该命名为 `p`。
- 接收者类型：接收者类型和参数类似，可以是`指针类型`和`非指针类型`。
- 方法名、参数列表、返回参数：具体格式与**函数**定义相同。

```go
//创建一个名为love的方法
func (l *Lanshaner)Love(){
  fmt.Printf("%v:Dont forget your love",l.Name)
}
```

如果想要修改接受者变量的某个字段，最常用的是通过定义它的接收者类型为**指针类型**，根据指针的特性，在方法结束后修改仍是有效的

```go
func (l *Lanshaner)Set(age int){
  l.Age = age // 语法糖 ==> (*l).Age = age
}
```

相反，如果接收者类型是**值类型**，Go语言会在代码运行时将**接收者**的值复制一份，同样，修改操作也只会针对该**副本**，不会修改本身

```go
func (l Lanshaner)Set(age int){
  l.Age = age
  p := &l
  p.Age = 0   //副本地址，一样修改不了
}
```

### 什么时候应该使用指针类型接收者

1. 需要**修改**接收者中的值
2. 接收者是**拷贝代价比较大**的大对象
3. 保证**一致性**，如果有某个方法使用了指针接收者，那么其他的方法也应该使用指针接收者。

**注意事项：** **非本地类型**不能定义方法，也就是说我们不能给别的包的类型定义方法，保证了隔离

## 接口

> 接口是一种**动态类型**
>
> 接口不是**继承**，而是**约定**

接口是一种由程序员来定义的类型，一个接口类型就是一组**方法的集合**（或类型约束），它规定了需要实现的所有方法。

相较于使用结构体类型，当我们使用接口类型说明相比于它是什么更关心**它能做什么**，而不是它是什么。

有兴趣的可以下来了解一下[鸭子类型](https://zh.wikipedia.org/zh-hans/%E9%B8%AD%E5%AD%90%E7%B1%BB%E5%9E%8B)

### 接口的定义

每个接口类型由任意个方法签名组成，接口的定义格式如下：

```go
type 接口类型名 interface{
    方法名1( 参数列表1 ) 返回值列表1
    方法名2( 参数列表2 ) 返回值列表2
    // …
}
```

其中：

- 接口类型名：Go语言的接口在命名时，一般会在单词后面添加`er`，如有写操作的接口叫`Writer`，有关闭操作的接口叫`closer`，管理数据库的接口叫`Manager`等。接口名最好要能突出该接口的类型含义。
- 方法名：当方法名首字母是大写且这个接口类型名首字母也是大写时，这个方法可以被接口所在的包（package）之外的代码访问。
- 参数列表、返回值列表：参数列表和返回值列表中的参数变量名可以省略。

```go
// Writer is the interface that wraps the basic Write method.
//
// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns n < len(p).
// Write must not modify the slice data, even temporarily.
//
// Implementations must not retain p.
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

这是Go语言源码中`io.Writer`接口的定义，我们看到这里面要实现一个`Write`方法，同时，还对它的返回值做了逻辑规范，即`n`只能是`len(p)`，调用需要`io.Writer`类型的标准库方法的话会检查返回值，`n!=len(p)`就返回一个**error**

当你看到一个`io.Writer`类型的值，我可以不知道它具体是什么，唯一知道的就是可以调用它的`Write`方法来做一些**具体实现不同，但抽象出来都是`Write`的事情(如网络IO、磁盘IO、打印等等)**

### 实现接口的条件

接口就是规定了一个**需要实现的方法列表**，在 Go 语言中一个类型只要实现了接口中规定的**所有方法**，那么我们就称它实现了这个接口，可以赋值到这个接口类型上

我们定义的`Singer`接口类型，它包含一个`Sing`方法。

```go
// Singer 接口 可以称之为 “会唱歌的东西”
type Singer interface {
		Sing()
}
```

我们有一个`Bird`结构体类型如下。

```go
type Bird struct {}
```

因为`Singer`接口只包含一个`Sing`方法，所以只需要给`Bird`结构体添加一个`Sing`方法就可以满足`Singer`接口的要求。

```go
// Sing Bird类型的Sing方法
func (b Bird) Sing() {
		fmt.Println("汪汪汪")
}
```

这样就称为`Bird`实现了`Singer`接口。bird 就属于`singer`类型了。

```go
var singer Singer
singer = Bird{}   // 动态类型，接口是Singer不变，具体类型从nil变成Bird
singer.Sing()
```

同时我们也可以通过**参数**传递接口类型

```go
// These routines end in 'f' and take a format string.

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
func Fprintf(w io.Writer, format string, a ...any) (n int, err error) {
	p := newPrinter()
	p.doPrintf(format, a)
	n, err = w.Write(p.buf)
	p.free()
	return
}
```

这里的`w`就是`io.Writer`，我们可以传入任意实现了`io.Writer`的具体类型，因为它可以是那只**鸭子**的值（注意会对`n==len(p)`进行检查，自定义类型的时候要注意）

### 接口的意义

1️⃣ 实现**多态性**（Polymorphism）：通过接口，可以实现多个不同类型的对象以相同的方式(如`.Write(p)`)进行操作，这增强了代码的灵活性和可复用性。

2️⃣ **解耦合**（Decoupling）：接口使得模块之间的依赖关系更松散，模块只需要关注接口定义的方法，不关注具体的实现细节。这有助于降低代码的耦合度，增加代码的可维护性和可测试性。

3️⃣ **扩展性**（Extensibility）：通过接口定义通用的行为，可以方便地对系统进行扩展和修改，而不需要改动已有的代码。当需要添加新的功能时，只需要实现接口定义的方法即可。

4️⃣ **接口断言**（Interface Assertion）：使用接口断言可以在**运行时**检查一个对象是否实现了某个接口，并根据情况进行处理。这样可以进行更灵活的类型转换和错误处理。



你可能还没怎么看懂这些，没事，我们会在下一个部分统合起来好好讲讲👍

# 面向对象思想，OOP

> ~~wc,OP~~

刚刚我们讲了**结构体、方法、接口**，尤其是接口，大家可能听的云里雾里的

Q：多写几个函数把函数串联起来、面向执行过程不就可以了吗？为什么要搞这么复杂

```go
func addUser(name string) {
	fmt.Println("添加用户:", name)
}
func deleteUser(name string) {
	fmt.Println("删除用户:", name)
}
func main() {
	addUser("陈越")
	deleteUser("郭瑞彤")
  addUser("王家宽")
  deleteUser("贺一鸣")
}
```

A：👀看起来确实没问题，但你的应用不可能都只做得这么简单，For example:

- **用户不止有名字，还有权限、角色、群组；**
- **还要扩展好多好多功能**
- **每个功能都得传一大堆参数；**
- **修改一处逻辑可能影响一堆函数传参等**

```go
package main

import "fmt"

// 面向过程：所有数据分散，每个函数都要传一堆参数
func addUser(name string, age int, info string, email string) {
	fmt.Println("添加用户:", name, age, info, email)
}

func deleteUser(name string, info string) {
	fmt.Println("删除用户:", name, "个人信息:", info)
}

func updateUser(name string, newInfo string, newEmail string) {
	fmt.Println("修改用户信息:", name, "→", newInfo, newEmail)
}

func main() {
  // 调用链路太💩
	addUser("康桥", 19, "面硬加咸蔬菜加倍蒜末和油多多", "1768832245@qq.com")
	updateUser("康桥", "0d000721", "kangqiao@lanshan.com")
	deleteUser("康桥", "0d000721")
}
```

这种所谓“面向过程”的代码会逐渐变得：

> 💣 难维护、难复用、难测试。

所以，引入了我们今天的主角，面向对象思想（object oriented programming，OOP）

## OOP 的出现就是为了解决这些“复杂度问题”

> 用**结构体**承载状态，用**方法**定义行为，用**接口**抽象能力

OOP 不是为了炫技，而是为了让：

1. **数据和行为绑定**（= 封装）；
2. **逻辑结构清晰、职责单一**（= 组合/继承）；
3. **可扩展、不改旧代码也能新增功能**（= 多态）。

将上述代码稍稍改一下:

```go
type User struct {
	Name   string
	Age    int
	Info  string
	Email  string
}


type UserManager struct{}

func (m UserManager) Add(u User) {
	fmt.Println("添加用户:", u.Name, u.Age, u.Info, u.Email)
}

func (m UserManager) Update(u *User, newInfo, newEmail string) {
	u.Info = newInfo
	u.Email = newEmail
	fmt.Println("更新用户:", u.Name, "→", u.Info, u.Email)
}

func (m UserManager) Delete(u User) {
	fmt.Println("删除用户:", u.Name)
}

func main() {
	manager := UserManager{} // 管理者 => 之后会用数据库进行持久化
  // 对象实例，直观修改
	user := User{
		Name:  "康桥",
		Age:   19,
		Info: "面硬加咸蔬菜加倍蒜末和油多多",
		Email: "kq@example.com",
	}

	manager.Add(user)
	manager.Update(&user, "0d000721", "kangqiao@lanshan.com")
	manager.Delete(user)
}
```

这样子的逻辑看起来就很直观了。我们先定义了一个对用户的**Manager**，用它来管理用户，再创建一个标准的用户结构体，之后在`Add`等方法里面用`Manager`去处理

因此，我们很多时候要把代码从面向过程的**跑起来**，改成面向对象的**长久维护、范式开发**

## 如何做OOP

OOP具有三要素：

1. **封装（Encapsulation）**
2. **继承** => 在Go中我们使用**组合代替继承（Composition over Inheritance）**
3. **多态（Polymorphism）**

### 一、封装

> **定义**：隐藏内部实现，只暴露必要的接口。
>  Go 通过**字段导出规则（大小写）+ 方法**实现封装。

```go
type Account struct {
	name   string // 小写 → 私有，包外不可见
	balance int64 
}

// 提供公开方法
func (a *Account) Deposit(amount int64) {
	a.balance += amount
}
func (a *Account) Balance() int64 {
	return a.balance
}

func main() {
	a := &Account{name: "二等分"}
	a.Deposit(100)
	fmt.Println("当前余额:", a.Balance()) // ✅ 通过方法访问内部状态
	// fmt.Println(a.balance)  // ❌ 外部直接访问会报错
}
```

具体来说就是通过**方法**控制读写，做到包与包之间私有字段可以**隔离**开来

### 二、组合代替继承

> **定义**：Go 没有类继承，用**结构体内嵌**实现复用与“方法提升”。

```go
type Info struct {
	Name string
}

func (i *Info) Introduce() {
	fmt.Println("我是", i.Name)
}

type Lanshaner struct {
	Info     // 内嵌结构体 → 继承 Info 的字段与*方法*
	Group string
}

func (l *Lanshaner) Work() {
	fmt.Println(l.Name, "正在", l.Group, "工作")
}

func main() {
	kq := Lanshaner{
		Info:  Info{Name: "康桥"},
		Group: "后端Go组",
	}

	kq.Introduce() // ✅ 方法提升：等价于 kq.Info.Introduce()
	kq.Work()
}
```

可以看到，我们的`lanshaner`组合了`info`，导致`kq`具有了`info`的字段和方法

**ps:若外层定义同名方法，可覆盖内层实现**



### 三、多态

> **定义**：不同类型实现相同接口，表现出相同行为(抽象角度下)
>  Go 通过 **interface + 隐式实现** 实现多态。

```go
type Notifier interface {
	Notify(msg string)
}

// 不同类型实现接口
type Email struct{}
func (Email) Notify(msg string) { fmt.Println("📧 邮件通知:", msg) }

type SMS struct{}
func (SMS) Notify(msg string) { fmt.Println("📱 短信通知:", msg) }

// 接口参数实现多态
func Send(n Notifier, msg string) {
	n.Notify(msg)
}

func main() {
	var e Email
	var s SMS

	Send(e, "欢迎加入蓝山团队！") // Email.Notify
	Send(s, "系统更新提醒")         // SMS.Notify
}
```

其实就是我在接口那里讲的，**实现的该接口的具体类型可以作为该接口值传入**，大部分时候我们只需要关注所谓抽象行为



# 泛型

> Go 的泛型在设计上就是**够用就好**——它不是要变成 C++ Template 或 Rust Trait 那种强泛化系统，而是希望你在有些时候能做得更安全、更整洁。
>
> 泛型对于 Go ：**是“语法糖级别的工具”，不是“类型系统革命”**

泛型 = **写一份，适配多种类型**，**编译期类型检查**，**不牺牲性能**



### 类型参数列表/泛型函数

每个泛型函数都具有他的**类型变量**，一般为`T`,代表着>=1个的类型

如下:

```go
func Max[T int|float32](a,b T)T{
		if a > b{
				return a
		}else{
				return b
		}
}

//Max[int](a,b)
```

`int | float32` 是类型约束，意思是,T **必须**是 int 或 float32 其中之一

然后我们指定`a/b/返回值`类型为当前类型变量`T`

那泛型只能在函数里面用吗？当然不是！！

### 泛型结构体

```go
type List[T any] struct {
	data []T
}

// l := List[float32]{
//   data: make([]float32,1024),
//}
```

显而易见，我们的`list`结构体里面可以是任意类型的切片



很简单对吧？我们来看看下一个

```go
type Map[K comparable, V any] struct {
	hashMap map[K]V
}
```

我们发现出现了一个我们没有见过的`comparable`，点进去我们发现是一个`interface`,注意我们刚刚讲接口的时候说的，接口也可以做**类型约束**，因为map中会涉及到K的比较，所以K的类型要约束到支持`==`and`!=`

```go
type myType interface{
	~int | ~int8 | ~float32		//～ 代表支不支持衍生类型，比如type Myint int，只写int就支持不了Myint
}
```

和刚刚结合起来试试

```go
func Max[T myType](a,b T)T{ // myType 替代了刚刚的 int | float32,支持的类型就是我们在myType中写的类型
	if a > b{
		return a
	}else{
		return b
	}
}
```

**注意:接口可以同时具有类型约束和方法，但只要存在类型约束就只能作为一个类型变量使用了**



# 底层简单解析

### 接口（iface）

每个进程（running的应用程序）都会维护一个`itabTableType`全局表，里面核心类型是`itab`

```go
type itabTableType struct {
	size    uintptr             // entries的长度
	count   uintptr             // 实际被填充的entries
	entries [itabInitSize]*itab // 实际为变长数组（really [size] large），容量为2的幂
}
```

```go
type itab struct {
	inter *interfacetype 	//接口抽象类型
	_type  *_type						//底层具体类型
	hash  uint32     // 类型hash，用于类型转换
	fun   [1]uintptr // 方法表，按接口“声明顺序”存放函数指针，变长尾随数组
}

// 当 fun[0] == 0 时，表示该具体类型未实现该接口（负缓存）,会报错
```

这个`entries`,它是整个**进程级（runtime 级）**的单例缓存表，用来缓存「接口类型 × 具体类型」对应的 itab 指针。

`itab`里面有接口类型签名，具体类型签名，和按照声明顺序的**方法表**

Q：什么时候创建`itab`？

A：第一次需要**接口赋值或断言**时，才在运行时创建 itab，例如:

```go
type Reader interface {
	Read()
}

type File struct{}

func (File) Read() { fmt.Println("reading...") }

func DoRead(r Reader) {
	r.Read()
}

func main() {
	f := File{}

	// 这里第一次触发 itab 构建
	DoRead(f)
}
```

```go
func main(){
	var r Reader
	// 这里也会触发
	r = File{}
}
```

触发的根本在于`getitab`，这是一个检查是否存在对应itab的函数，没有就会Init一个然后Add到entries上，最后创建具体`iface`结构体

```go
type iface struct {
    tab  *itab					// 指向从全局缓存表中查到/创建的itab
    data unsafe.Pointer	// 数据指针
}

// 空接口
type eface struct {
    _type *_type
    data  unsafe.Pointer
}
```

#### 编译期时：

- 记录每个类型的“方法集”——接收类型为T是T的方法集，接收类型为 *T 是 T 和 *T的方法集
- 记录每个接口的方法签名；

#### 运行期：

​     如果有类型要包装成接口，就通过`getitab`去查`itabTableType`,如果不存在就构造`itab`然后扔进`iface`

**总结**：编译期决定“接口签名集”和“类型方法集”；运行期通过 `getitab` 将 `(接口类型, 具体类型)` 映射成 `itab` 并缓存；接口赋值本质上就是构造 `iface{ tab, data }`，实现动态类型绑定

### 泛型(看看就行

**相同形状可复用（shape stenciling）**
 当类型参数的**内存布局 + GC 指针位图**一致（简称同一*shape*），编译器只生成**一份机器码**并复用，效果≈“手写几次类似函数”的优化。
 例：`Max[T constraints.Ordered]` 对 `int/int32/float64` 等标量，多数可共享同一形状实现。

**不同形状 + 需要具体算法时会带“字典”**
 若用到依赖具体类型的操作（如 **map 的哈希/相等性**、接口方法调用、`==`、`<=` 在某些约束下的实现），编译器会为实例**传入算法字典/专用 helper**，或生成额外实例，因此**底层略复杂**，但仍保持**类型安全**与**接近手写的性能**。
 例：`map[K]V` 的 `K` 需要专属的 hash/eq 函数



# 手搓time

看时间吧，时间还够就带大家手搓一个用到本课大部分内容的demo

# 作业

1. 好好消化，本节课内容有点多
2. 实现一个任意你想实现的东西，例如原神用户管理系统，`MyGo`防炸团系统等等；最好用到部分本节课的内容，体现一下本节课的思想hhh，**做得好的有奶茶QAQ** ~~老东西早该爆金币了~~

完成后提交到邮箱kangqiao@lanshan.com

有什么问题也可以发邮箱/私聊我问哦



# 扩展阅读

**Go 源码：`src/runtime/iface.go`**

- 定义了 `getitab`、`itabTableType`、`convT2I`、`convI2I` 等核心函数。
- [在线阅读](https://cs.opensource.google/go/go/+/master:src/runtime/iface.go)

**Go 源码：`src/runtime/runtime2.go`**

- 包含了 `eface`、`iface`、`itab`、`interfacetype` 等结构体定义。
- [在线阅读](https://cs.opensource.google/go/go/+/master:src/runtime/runtime2.go)

**官方博客**：[The Go Blog: Interfaces](https://go.dev/blog/laws-of-reflection)

**描述了 Go 1.18 泛型引入后，接口的“类型约束”与传统接口在 runtime 层的关系**：https://github.com/golang/go/issues/43651

**Go的OOP(讲的不赖）**:https://blog.gypsydave5.com/posts/2024/4/12/go-is-an-object-oriented-programming-language/





---
<u>Don't Forget Your Love 💙</u>
