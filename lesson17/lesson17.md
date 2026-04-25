# 第 17 节课 · 杂项

---

## 0. 开场

前面 16 节课,从 Go 语法一路讲到 AI Agent,你已经能从 0 写一个项目、把它放进容器、拆成微服务、加上观测和网关,这节课主要是帮助你了解比较完善的生产线是如何工作的,以及相关的组件。

我们现在个人开发的工作流大概是——本地开发->Dockerfile 打包->推送镜像->服务拉镜像，当然也可能有人是直接推到 github 在服务器上打包运行，而对于关键的配置文件，我们可能是手动搬运的，因为 yaml 总不可能写镜像里面或者放到 github 上，当然如果没有在服务器上运行的需求那就不用说了，只需要本地打包调试完事。然而在团队多人开发，涉及到生产测试环境的时候，我们就不能干这么原始的事情了，这个时候就多了两个东西：

- 写完的代码怎么自动跑测试、自动打包、自动上线 —— **CI/CD**
- 上线之后想改个开关、调个阈值,不想重启服务 —— **配置中心**

这就是我们这节课要讲的其中之二。

---

## 1. CI/CD

### 1.1 先把名字搞清楚

CI = Continuous Integration,持续集成。说人话就是:你 push 代码,机器自动帮你跑 lint、跑测试、跑构建,跑挂了立刻通知你。目的是**保证主分支永远是能发布的状态**,而不是等到要上线那天才发现有人三周前合了个挂的代码。

CD 有两种意思,挺多人搞混:

- Continuous Delivery,持续交付。机器自动把可发布的产物(镜像、二进制)准备好,**人点按钮**才上线。
- Continuous Deployment,持续部署。准备好之后**直接自动**上线,人都不用管。

自动是啥意思，其实就是自己写个一个类似脚本的东西在服务器上面给你的项目执行打包，然后推送到服务器上运行，国内大多数公司其实是 Delivery,因为没人敢让代码自己往生产推。Deployment 在云原生、SaaS、那种迭代极快的产品里更常见。

为什么非要搞这套?因为不搞的话,日常对话就是这样:

> "我本地能跑啊"
> "啊我忘记跑测试了"
> "上线脚本第 7 步是不是又漏了"
> "上次发布是谁动的来着"

CI/CD 的价值很简单:**把容易出错、需要纪律的事交给机器**。人是不可靠的,机器是可靠的,就这么回事。

### 1.2 业界都用啥

主流方案大概这些:

| 方案 | 在哪跑 | 配置文件 | 一句话评价 |
| --- | --- | --- | --- |
| GitHub Actions | github.com | `.github/workflows/*.yml` | 开源项目首选,生态最猛 |
| GitLab CI | GitLab | `.gitlab-ci.yml` | 跟 MR 集成深,公司用 GitLab 就用它 |
| Jenkins | 自己搭 | `Jenkinsfile` | 老牌,插件多,UI 一股 2010 年的味道 |
| Drone | 自己搭 | `.drone.yml` | 轻量,容器原生 |
| Gitea Actions | 自己搭 | `.gitea/workflows/*.yml` | 语法跟 GitHub Actions 一样 |

其实这些压根不用学，因为现在都是让 AI 去写，同时在企业里面，大多数用的是 gitlab CI、jenkins、drone 这种开源的解决方案，还有一些是自己有一套基础设施，所以其实你压根学不完，所以放心交给 AI 吧👍

### 1.3 GitHub Actions

概念只有几个：**Workflow**（一个 yml 文件）、**Job**（默认并行）、**Step**（串行步骤）、**Action**（别人封装好的 Step，`uses:` 一下就用）、**Runner**（跑这些东西的机器，GitHub 免费给）。

最小能跑的 yml：

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4         # 拉代码
      - uses: actions/setup-go@v5         # 装 Go
        with:
          go-version: '1.22'
      - run: go test ./...                # 跑测试
```

提交这个文件到 main，下次 push 就在 GitHub 仓库的 Actions 标签页能看到结果。第一次跑成功的时候那个绿勾会让你上瘾。

### 1.4 实战：Go 项目的 lint + 测试

拿来就能用的版本：

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main]
  pull_request:

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true                    # 自动缓存 go mod
      - uses: golangci/golangci-lint-action@v6
        with:
          version: v1.59
          args: --timeout=5m

  test:
    runs-on: ${{ matrix.os }}
    strategy:
      fail-fast: false                   # 一个挂了别立刻终止其他的
      matrix:
        os: [ubuntu-latest, macos-latest]
        go: ['1.21', '1.22']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
      - run: go test -race -coverprofile=coverage.out ./...
      - if: matrix.os == 'ubuntu-latest' && matrix.go == '1.22'
        uses: codecov/codecov-action@v4
        with:
          files: ./coverage.out
          token: ${{ secrets.CODECOV_TOKEN }}
```

讲两个细节。**`matrix`** 是矩阵构建，可以同时跑多个 OS × 多个 Go 版本的组合，开源库基本得加，公司项目通常用一个固定环境就行，加 matrix 反而浪费 runner 资源。**`-race`** 一定要开，Go 的竞态检测器在 CI 里跑一次能帮你抓到很多本地发现不了的并发 bug，代价是慢 2-3 倍，但值。

### 1.5 实战：打镜像推到 GHCR

GHCR 是 GitHub Container Registry，跟 Docker Hub 一样是放镜像的地方，但跟 GitHub 账号天然打通。

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags: ['v*']                         # 打 tag 才发

permissions:
  contents: read
  packages: write                        # 写 GHCR 必须

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: docker/setup-qemu-action@v3      # 多架构构建
      - uses: docker/setup-buildx-action@v3

      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}    # 内置的,不用自己配

      - id: meta
        uses: docker/metadata-action@v5
        with:
          images: ghcr.io/${{ github.repository }}
          tags: |
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=sha,format=short

      - uses: docker/build-push-action@v6
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

用法就是 `git tag v1.0.0 && git push origin v1.0.0`，几分钟后 `ghcr.io/你的用户名/仓库名:1.0.0` 就有了。

有个坑：`secrets.GITHUB_TOKEN` 是 GitHub 自动给你的，但默认权限不够写 GHCR，所以 yml 顶上那段 `permissions: packages: write` 是必须的，漏了会报 403，排查半天找不到原因。

关于 Secrets：密码别写进 yml，放到仓库 Settings → Secrets 里，用 `${{ secrets.XXX }}` 引用。GitHub 在日志里会自动打码成 `***`，但别在 step 里 `echo` 它，字符串变换之后打码可能失效。

本地调试 yml 可以装个 [`act`](https://github.com/nektos/act)，`brew install act && act -j test`，用 docker 模拟 GitHub Runner，不是 100% 一致但抓 80% 的低级错误够用。

---

### 1.6 Gitea：你不能每个项目都放 GitHub

公司代码不能放公网，这是常态。**Gitea** 是个轻得多的选择：Go 写的，单二进制，几十 M 内存就能跑，能力上覆盖了 GitHub 的核心 90%。最关键的是 Gitea 1.19 之后内置了 **Gitea Actions**，**语法跟 GitHub Actions 完全一样**——你写过的 GHA workflow 直接搬过来基本就能跑。

为啥我在这里讲 Gitea 而不是 gitlab 或者说是 gogs？答案很简单，Gitea 是颜值最高的😉，另外的两个的 UI 压根不能看。同时 Gitea 的开源社区是非常活跃的，更新迭代都很快，而且原生使用 go 编写，是非常现代化的 git 托管平台。另外它也内置 CICD，而且和 drone 高度集成，是个人项目托管的不二之选。

---

## 2. 配置中心

### 2.1 为什么 Viper 不够用

第 10 节实战课用过 Viper:

```go
viper.SetConfigFile("config.yaml")
viper.ReadInConfig()
port := viper.GetInt("server.port")
```

单体小项目这样挺好。但项目稍微大一点、上了生产、跑了一段时间,就会撞到这几堵墙:

**多环境管理**。dev、test、staging、prod 各一份配置,改了 dev 忘了同步 prod,出事。

**动态变更**。比如限流阈值想从 100 调到 200,难道要改代码、提 PR、走 review、构建、灰度、全量?改个数字搞两天?

**灰度下发**。新配置只想让北京机房 10% 实例先用,观察一下没问题再全量。Viper 干不了。（其实能干，但是关键在于基础设施是否完备，有的地方有自己写的 k8s 平台，可以直接按照灰度将配置文件下发到容器内部，从而达到灰度的需求）

**审计**。"这个配置上周三是谁改的?为什么改?能回滚吗?"——同时我们配置文件一般不会明文进入到 git 仓库，一般是直接依赖专门的配置中心进行版本管理。

配置中心把这几件事就都干完了，当然，最后读文件解析配置还是需要用到 viper😀

### 2.2 它长什么样

抽象图很简单:

```
        ┌────────────────┐
        │   配置中心      │
        └───────┬────────┘
                │  长轮询/推送
        ┌───────┼────────┬────────┐
        ↓       ↓        ↓        ↓
     服务A    服务B    服务C    服务D
```

核心能力就五个:

- 集中存储,一个地方改全部生效
- 动态推送,改完几秒内运行中的服务收到新值
- 版本管理,能回滚
- 权限和审计,谁能改、改了什么都留痕
- 灰度,按集群/实例分组下发

国内主流就两个：**Apollo**(携程开源)和 **Nacos**(阿里开源)。下面分别讲。

---

### 2.3 Apollo

#### 它的架构

```
Portal (Web 界面,改配置)
   ↓
Admin Service (写,带审批)
   ↓
Config Service (读,客户端长轮询这个)
   ↓
Client SDK (嵌在你的服务里)
```

四个组件,部署确实有点重,但**权限和审批做得最好**,所以大公司喜欢用。需要严格走流程的场景(比如金融、合规),Apollo 基本是默认选项。

#### 核心概念

只有三个:

- **AppId**:你的应用名,比如 `user-service`
- **Cluster**:同一个应用的不同集群,比如 `default`、`beijing`、`shanghai`,可以让北京机房的服务读不一样的配置
- **Namespace**:配置分组。默认有个 `application`,你也可以建 `database`、`feature-flag` 这种,按业务拆开管理

#### Go 接入

用 `github.com/apolloconfig/agollo/v4`，核心就几行：连上 Config Service、注册变更监听、然后 `GetIntValue` 之类的读值。有一个关键配置 `IsBackupConfig: true`——它会把拉到的配置缓存到本地文件，Apollo 挂了或者网络抖动的时候，服务还能用上一份配置启动，**这是生产环境必开的**，不开的话 Apollo 一挂全公司服务连重启都重启不了。

---

### 2.4 Nacos

#### 一站式

Nacos 跟 Apollo 最大的区别:**它不光是配置中心,还是注册中心**。回想第 12 节服务发现讲的那些东西,Nacos 全包了。

这是它最大的卖点:**少装一个组件**。本来你要 Apollo 做配置 + Consul/Etcd 做注册,现在一个 Nacos 全搞定。运维少一种东西要管,问题少一半。

#### 数据模型

比 Apollo 还简单:

- **Namespace**:环境隔离,dev / prod 分开
- **Group**:业务分组
- **DataId**:配置文件名,比如 `user-service.yaml`

定位一份配置 = `Namespace + Group + DataId`,就这。

#### Go 接入

用 `github.com/nacos-group/nacos-sdk-go/v2`，也是连上之后拉配置 + 监听变更的模式。本地跑的话：

启动 Nacos:

```bash
docker run -d --name nacos -p 8848:8848 -p 9848:9848 \
  -e MODE=standalone nacos/nacos-server:v2.3.0
```

个人体会上是 apollo 更好用一点，nacos 其实也还不错，但是也没必要为了学习专门搭一套配置中心，了解即可。

---

## 3. 八股怎么学

### 3.1 八股是什么

大家在找实习的时候必然会历经八股文的折磨，你们可能对八股这个名词的第一印象很差，但是实际上八股更多是我们之前学习过程中忽略的，没有深入讲的一些底层知识，虽然，但是确实有一部分八股看着很难受。我们是学 Go 的，除了常规后端八股之外，当然也需要去学习 Go 的运行时机制，也就是 go 八股。

### 3.2 我们要看哪些八股，咋学？

首先是常见的计算机基础——操作系统、计算机网络，如果你之前学过我推荐的 mit6.s081，那么操作系统这块八股文会学的很轻松，计算机网络在初学可能会觉得只能靠死记硬背，但是你熟悉之后还是能够梳理出大致的知识体系的，这两个都可以在小林 coding 上面看，小林的八股掌握了，其实面试也没有什么大问题。

另外还有后端常用的中间件八股文了，Mysql、Redis、消息队列都是我们的老顾客了，可能有的面试不会问你计算机基础，但是这几个组件是包问的，这个也可以在小林coding里面解决，消息队列的八股，小林里面可能不全，多看几篇博客也可以解决了。

最后是我们的 go 八股文，和 jvav 相反，Go 相关的八股文相对比较少，但是与之对应的，是学习资源很少，我建议的学习路径是先学习一下操作系统相关的知识，不然看着会很吃力，然后看[《Go 语言设计与实现》][https://draven.co/golang/]先熟悉一下 Go 的八股文，在有比较完整的认知之后可以去自己梳理一下 Go 的源代码，了解一些新特性。然后这块面试就没啥问题了。

## 4. 结语

这集主要是补充作为合格的后端开发，我们缺少的一些知识面，当然其实你们还可以简单了解一下 k8s 相关的知识，但是对于单纯的 go 后端来说， k8s 并不需要深入学习，除非你目标是云原神高手。

这是最后一节课了，最后就祝大家学习顺利吧。
