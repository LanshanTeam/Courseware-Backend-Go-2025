# Traefik 配置说明

- 静态配置：[traefik.yml](/home/wang/code/Courseware-Backend-Go-2025/lesson14/getway_demo/deployments/traefik/traefik.yml)
- 动态路由规则来源：[docker-compose.yml](/home/wang/code/Courseware-Backend-Go-2025/lesson14/getway_demo/docker-compose.yml) 中各服务的 `labels`

只看 `traefik.yml` 不够，因为它只负责“Traefik 自己怎么启动”；真正的路由规则主要写在容器标签里。

## 这套配置负责什么

Traefik 在这个 demo 里承担的角色和 Nginx 一样，都是网关入口，但方式不同：

- Nginx：手写 upstream 和路由
- Traefik：先声明从 Docker 发现服务，再由容器标签提供路由规则

因此 Traefik 这里要分两层理解：

1. Traefik 自己监听哪个入口、从哪里发现服务
2. 哪些容器要暴露、各自匹配什么规则

## `traefik.yml` 解释

### `entryPoints`

```yaml
entryPoints:
  web:
    address: ":8080"
```

这表示 Traefik 在容器内部监听 `8080`，入口名叫 `web`。

后面 `docker-compose.yml` 里的各个 router 都会写：

- `traefik.http.routers.<name>.entrypoints: web`

也就是说，这些路由都挂在这个 `web` 入口上。

宿主机实际访问的是 `8092`，因为 compose 做了映射：

- 宿主机 `8092` -> Traefik 容器 `8080`

### `providers.docker`

```yaml
providers:
  docker:
    endpoint: "unix:///var/run/docker.sock"
    exposedByDefault: false
```

这是这份配置最核心的一段。

它的含义是：

- Traefik 从 Docker 读取容器信息
- 通过 Docker socket 拿到容器、标签、端口等元数据
- 默认不暴露所有容器

`exposedByDefault: false` 很重要，它表示：

- 只有显式写了 `traefik.enable=true` 的容器才会参与路由
- 没写的容器即使在同一个网络里，也不会自动暴露

这样可以避免把不该公开的容器意外挂出去。

### `api.dashboard`

```yaml
api:
  dashboard: true
```

这表示开启 Traefik 自带 dashboard 功能。

但“开启 dashboard”不等于“外部能访问 dashboard”。要对外访问，还需要额外给 Traefik 自己配置一条 router，这部分写在 `docker-compose.yml` 的 `traefik` 服务 `labels` 里。

### `log`

```yaml
log:
  level: INFO
```

日志级别设置为 `INFO`，方便观察 Traefik 的启动和路由行为。

### `accessLog`

```yaml
accessLog: {}
```

开启访问日志，便于验证请求是否真的经过 Traefik。

## `docker-compose.yml` 中的 Traefik 相关配置

Traefik 的动态规则并不在单独文件中，而是附着在容器上。

这份 demo 中有 5 类 Traefik 相关对象：

- `api-1`
- `api-2`
- `admin`
- `echo`
- `traefik` 自己的 dashboard router

## `api-1` 与 `api-2`

配置片段如下：

```yaml
labels:
  traefik.enable: "true"
  traefik.http.routers.api.rule: Host(`api.traefik.localtest.me`)
  traefik.http.routers.api.entrypoints: web
  traefik.http.routers.api.service: api
  traefik.http.services.api.loadbalancer.server.port: "8080"
```

`api-1` 和 `api-2` 两个服务都写了这一组标签。它们的组合含义是：

- `traefik.enable: "true"`
  允许这个容器被 Traefik 发现并接管
- `traefik.http.routers.api.rule`
  当请求的 Host 是 `api.traefik.localtest.me` 时命中这条路由
- `traefik.http.routers.api.entrypoints: web`
  这条路由挂在 `web` 入口上，也就是容器内 `8080`
- `traefik.http.routers.api.service: api`
  把请求交给名为 `api` 的后端服务池
- `traefik.http.services.api.loadbalancer.server.port: "8080"`
  告诉 Traefik：真正转发时访问容器内部 `8080`

最关键的一点是：

- `api-1` 和 `api-2` 共用了同一个 router 名 `api`
- 也共用了同一个 service 名 `api`

这样 Traefik 会把两个容器自动聚合成一个负载均衡后端池。因此你连续请求：

- `http://api.traefik.localtest.me:8092/info`

返回里的 `instance` 会在 `api-1` 和 `api-2` 之间切换。

## `admin`

配置片段如下：

```yaml
labels:
  traefik.enable: "true"
  traefik.http.routers.admin.rule: Host(`admin.traefik.localtest.me`)
  traefik.http.routers.admin.entrypoints: web
  traefik.http.services.admin.loadbalancer.server.port: "8080"
```

这组标签表示：

- Host 为 `admin.traefik.localtest.me` 时命中 `admin` 路由
- 转发到 `admin` 容器内部的 `8080`

这里只有一个实例，因此不承担负载均衡演示，而是用于展示最直接的 Host 分流。

## `echo`

配置片段如下：

```yaml
labels:
  traefik.enable: "true"
  traefik.http.routers.echo.rule: Host(`traefik.localtest.me`) && PathPrefix(`/echo`)
  traefik.http.routers.echo.entrypoints: web
  traefik.http.services.echo.loadbalancer.server.port: "8080"
```

这组标签的重点在于：

- 不只是匹配 Host
- 还要求路径命中 `/echo`

因此：

- `http://traefik.localtest.me:8092/echo/inspect`

会转发到 `echo`

它的作用是演示 Traefik 如何同时基于域名和路径前缀做路由。

## Traefik 自己的 dashboard

在 `traefik` 服务本身上还有这组标签：

```yaml
labels:
  traefik.enable: "true"
  traefik.http.routers.dashboard.rule: Host(`dashboard.traefik.localtest.me`)
  traefik.http.routers.dashboard.entrypoints: web
  traefik.http.routers.dashboard.service: api@internal
```

这里的关键不是转发到某个业务容器，而是：

- `api@internal`

它表示 Traefik 的内部服务，也就是 dashboard 自己。

因此这个入口：

- `http://dashboard.traefik.localtest.me:8092/dashboard/`

打开的是 Traefik 自带的管理界面，而不是业务服务。

## 一次请求是怎么走的

以 `http://api.traefik.localtest.me:8092/info` 为例：

1. 浏览器或 `curl` 访问宿主机 `8092`
2. Docker 把请求转发到 Traefik 容器内的 `8080`
3. Traefik 在 `web` 入口收到请求
4. Traefik 根据 Host 命中 `api` router
5. `api` router 关联到 `api` 这个 service
6. `api` service 下有两个实例：`api-1` 和 `api-2`
7. Traefik 选择其中一个实例，把请求转发到对应容器的 `8080`
8. 后端返回 JSON，里面能看到具体 `instance`

## 为什么 Traefik 这里看起来“配置少”

因为它把大量路由信息分散到了容器标签里。

从维护方式上看：

- `traefik.yml` 更像 Traefik 自己的启动参数
- `labels` 更像业务服务的挂载声明

这就是它和 Nginx 的最大区别：

- Nginx 把 upstream 和路由都集中写在一个文件里
- Traefik 把“发现规则”和“业务路由”拆开了

## 看这套配置时最容易误解的点

- `traefik.yml` 不是完整路由配置，真正的业务路由大多在 compose `labels`
- `entryPoints.web.address=:8080` 指的是 Traefik 容器内部端口，不是宿主机端口
- `loadbalancer.server.port=8080` 指的是目标业务容器内部端口
- `api-1` 和 `api-2` 能自动组成一个后端池，是因为它们共享了同一个 Traefik service 名
- 开启了 `dashboard: true` 之后，还必须有额外 router 才能从浏览器访问 dashboard
