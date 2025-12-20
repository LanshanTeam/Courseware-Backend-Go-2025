# Docker、Linux基础

**Docker**是一个开放源代码的开放平台软件，用于开发应用、交付（shipping）应用和运行应用。Docker允许用户将基础设施（Infrastructure）中的应用单独分割出来，形成更小的颗粒（容器），从而提高交付软件的速度。

~~流行，工作要用，podman也不错其实~~



## Docker的特点——容器化，可移植

**容器化**将复杂的应用拆分为多个独立的服务。Docker容器为每个微服务提供独立的运行环境，确保服务之间的**隔离性**和**独立性**，便于开发、部署和维护。

**可移植**符合了应用在不同环境中的一致性，由于容器的启动部署的自动化和标准化，可以实现应用的快速部署和扩展。



## Docker容器、虚拟机？

### 虚拟机

虚拟机会虚拟化整个**操作系统**，运行在真正的物理机器上，最终导致的结果（哪怕当下已经存在较为成熟的**硬件辅助虚拟化技术**），但仍然显得太重、太慢，最大的好处其实是提供**强隔离性、跨操作系统内核运行**

### Docker容器（Containers）

Docker容器通过**共享宿主机内核**提供运行能力，Docker通过Linux内核中的资源分离机制，**容器化**（操作系统层虚拟化）建成一个**足够隔离的**的容器，避免再启动一个操作系统导致额外性能负担，相较于虚拟机开销小，更**轻**，更**快**。

Docker **没有**：

- ❌ 虚拟硬件
- ❌ 虚拟内核
- ❌ VMM / Hypervisor
- ❌ Guest OS



## Docker这么强，是怎么实现的呢？

主要的核心就是：控制群组（**Cgroups**）和命名空间（**NameSpace**）、分层文件系统(**Union File System**)

### Cgroups

> 决定容器能用多少资源

- **资源限制（Resource limiting）：组**可以被设置不超过设定的内存限制（这也包括页面缓存）、IO带宽限制、CPU配额限制、CPU集合限制或最大打开文件数。
- **优先级设定（Prioritization）：**一些**组**可能会得到更大的CPU或磁盘IO吞吐量。
- **结算（Accounting）：**衡量一个**组**的资源使用情况，可以用于计费等目的。
- **控制（Control）：**冻结**组**中的进程，运行检查点和重新启动。

### NameSpace

> 决定容器内的可视视图、容器能看到什么

- **PID**：只看到容器内进程
- **Mount**：看到自己的文件系统
- **Network**：独立于Host的IP/端口
- UTS/IPC/User：主机名、通信、权限隔离

### Union File System

> 容器为什么轻、快、省空间

Docker 镜像是**分层文件系统**：

- 底层只读、可复用
- 容器启动时加一层可写层
- 多容器共享同一镜像层

**镜像小、启动快、部署一致**



### DockerFile

由开发者编写，具有基本环境、依赖等，**从上到下执行**，~~一般给GPT写~~

```dockerfile
FROM ubuntu:18.04 #Ubuntu18.04风格的用户态文件系统
COPY . /app	#拷贝文件到镜像中
RUN make /app #依赖用户态文件系统提供的make工具
CMD python /app/app.py #用户态下python程序开始跑，其中系统调用落在宿主机内核
```

写好之后，可以`docker build`打包我们的镜像，甚至可以上传到docker registry给其他人使用



### Docker Compose

> 刚刚 Dockerfile 解决的是：“一个容器怎么构建、怎么运行”
>
> 那如果我们的应用不止一个容器，我们应该怎么建议的跑起来呢？

docker-compose是用一个配置文件，一次性启动、管理多个相关容器的工具，符合**YAML**的格式

```yaml
services:				#以下都是你需要的服务
  app:					#程序容器名
    build: .				#构建目录
    ports:
      - "8080:8080" #暴露端口
    depends_on:
      - db					#依赖

  db:
    image: mysql:latest  #拉取最新镜像
    environment:
      MYSQL_ROOT_PASSWORD: root	#环境变量
```



### Docker网络

#### 桥接网络（Bridge）--常用

默认网络模式，容器通过桥接网络进行通信，适用于单主机环境。

- 虚拟网卡
- 容器之间可以互相通信
- 容器访问外网，通过宿主机转发

```yaml
services:
  app:
    image: nginx
    ports:
      - "8080:80"

```



#### 主机网络（Host）

容器共享宿主机的网络栈，适用于对网络性能要求较高的场景。

- 共用网卡
- 容器**直接使用宿主机的网络**
- 没有端口映射的概念
- 网络性能最好，但隔离最弱

```yaml
services:
  app:
    image: nginx
    network_mode: host

```



#### 覆盖网络（Overlay）——多机器通信

跨主机的网络模式，适用于集群环境，通过 Docker Swarm 或 Kubernetes 实现容器间通信。



####  无网络（None）

容器不连接任何网络，适用于高度隔离的场景。

```yaml
services:
  app:
    image: nginx
    network_mode: none

```



#### 自定义网络（Custom Network）

根据需求创建自定义网络，支持 DNS 服务发现和网络隔离。

- Docker 允许创建自定义网络
- 容器之间可以直接用容器名通信
- 网络隔离更清晰

```yaml
services:
  web:
    image: nginx
    ports:
      - "8080:80"
    networks:
      - mynet

networks:
  mynet:
    driver: bridge

```



### Docker存储

#### 数据卷（Volumes）

宿主机上的目录或文件，挂载到容器中，适用于**持久化**数据和跨容器共享。

- 数据由 Docker 管理
- 和容器生命周期分离
- 容器删了，数据还在

```yaml
services:
  db:
    image: redis
    volumes:
      - data:/data

volumes:
  data:

```



#### 绑定挂载（Bind Mounts）

将宿主机的文件或目录直接挂载到容器，适用于开发环境中的实时代码更新。

- 直接把宿主机目录挂到容器
- 修改代码，容器里立刻生效 -- 适合调试、开发

```yaml
services:
  app:
    image: nginx
    volumes:
      - ./html:/usr/share/nginx/html

```



#### tmpfs 挂载

将数据存储在内存中，适用于需要高性能和短期存储的场景。

数据存在内存中，容器没了就消失，但是很快。

```yaml
services:
  app:
    image: nginx
    tmpfs:
      - /tmp

```



#### 存储驱动（Storage Driver）

管理镜像和容器文件系统，支持多种存储后端，如 AUFS、OverlayFS、Btrfs、ZFS 等。

- 镜像分层
- 容器文件系统



## Docker上手

- **镜像（Image）**：用于创建容器的只读模板。可以从 Docker Hub 拉取现有镜像或通过 Dockerfile 自行构建。
- **容器（Container）**：镜像的运行实例，具有独立的文件系统、网络和进程空间。
- **仓库（Repository）**：存储镜像的地方，可以是本地仓库或远程仓库（如 Docker Hub）。

> 人话讲就是，容器是跑起来的镜像，仓库是存镜像的“github”



### 常用命令

拉取镜像

```bash
docker pull [name]
```

把单一镜像跑起来作为容器（没有会自动拉取）

```bash
docker run [name]
```

查看容器状态

```bash
docker ps
```

查看所有的镜像

```bash
docker image ls
```

删除容器

```bash
docker rm [name]
```

------

启动多个容器服务

```bash
docker compose up -d #-d 代表后台运行
```

暂停

```bash
docker compose stop
```

查看服务状态

```bash
docker compose ps
```

重启服务

```bash
docker compose restart
```

构建服务

```bash
docker-compose build
```



### 常见问题

#### Docker镜像站位于国外，经常会遇到镜像拉不下来

1. 换镜像站

   一般来说，daemon.json位于/etc/docker下，然后sudo systemctl restart docker

   ```json
   {
     "registry-mirrors": [
       "https://docker.m.daocloud.io"
     ]
   }
   ```

2. 科学上网

   

## Linux基础

作为绝大部分服务器的首要选择，linux对后端来讲是十分重要的，对很多公司来说也会在面试中考察

常见命令大全如下：man,ls,cd,pwd,mkdir,touch,rm,rmdir,mv,cp,cat,which,whereis,find,tar,gzip,chmod,chown,top,free,lsof,ifconfig,ping,traceroute,netstat,telnet,ln,grep,ps,mount...



### 简易基础概念

#### 一切皆进程（Process）

Docker 容器的本质，就是一个**被限制和隔离的 Linux 进程**。

#### 用户态 & 内核态

Linux 把程序分成两层：
用户态负责跑程序，内核态负责管资源和硬件。

```bash
应用程序（用户态）
    ↓ 系统调用
Linux 内核
    ↓
		硬件
```

#### 文件系统 & 目录树

Linux 只有一个 `/`

程序看到的文件来自文件系统

不同进程**可以看到不同的文件系统视图**

> Docker 正是利用 Linux 的这个能力，
> 让不同容器看到不同的 `/`。



## 作业

这节课概念相关的非常简单，也就讲了讲怎么用，重要的是大家下来往深了去理解

### **容器到底是不是虚拟机？请用“证据”说明**

#### 要求

1. 启动一个容器：

   ```
   docker run -it ubuntu:22.04 bash
   ```

2. 在 **宿主机** 和 **容器内**，分别执行：

   ```
   ps aux | head -n 5
   ```

3. 回答问题：

   - 容器内看到的 PID=1 是谁？
   - 宿主机是否能看到容器里的进程？

### **用 docker-compose 构建一个「最小可用系统」**

#### 要求

使用 `docker-compose.yml` 启动 **两个服务**：

- `web`：your app
- `db`：任意

必须满足：

1. 使用 **自定义 bridge 网络**
2. `db` 使用 **volume** 持久化数据
3. `web` 可以通过 **服务名** 访问 `db`

### **Docker 容器在 Linux 里留下了什么痕迹？**

#### 要求

1. 启动任意一个容器
2. 在宿主机上找出以下内容之一即可（任选）：

**任选 A：进程链路**

- 找到与该容器相关的：
  - `dockerd`
  - `containerd`
  - `containerd-shim`
  - 容器进程本体

**任选 B：网络证据**

- 找到容器对应的：
  - IP 地址
  - 端口映射的证据（如 docker-proxy / ss / netstat）

1. 用一句话解释：

> Docker 是如何把一个普通进程，变成‘看起来像一台机器’的？

> 为什么Mysql容器里 **不能随便 kill PID 1**？
> PID 1 在Linux中有什么特殊含义？

