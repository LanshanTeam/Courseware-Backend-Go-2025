# Nginx 配置说明

## 配置作用

这个 demo 里，Nginx 作为统一入口网关，对外暴露宿主机 `8081`，容器内部监听 `8080`。它负责三件事：

- 按域名把请求分流到不同后端
- 按路径把请求分流到 `echo` 服务
- 在转发时补齐常见代理头，方便后端观察真实请求信息

对应的后端服务都不直接暴露到宿主机，而是在 Docker 内网中通过服务名互相访问：

- `api-1:8080`
- `api-2:8080`
- `admin:8080`
- `echo:8080`

## 整体结构

这份配置可以分成四块：

1. `worker_processes`
2. `events {}`
3. `http {}` 顶层容器
4. 三个 `upstream`
5. 三个 `server`


## `worker_processes`

```nginx
worker_processes auto;
```

这表示 Nginx 的 worker 进程数由运行环境自动决定。这个 demo 里保留 `auto`，目的是让配置更完整，也让教学示例覆盖到 Nginx 最常见的顶层参数。

在这个 demo 中，它不是路由逻辑的一部分，而是告诉 Nginx：启动多少个工作进程来处理请求。

## `events {}`

```nginx
events {
    worker_connections 1024;
}
```

这里定义的是每个 worker 在事件模型层面最多可同时处理的连接数。

这个 demo 没有继续展开更多事件模型调优项，只保留一个最常见的 `worker_connections 1024;`，原因是：

- 配置结构更完整
- 课堂上可以顺带解释“进程数”和“单进程连接数”这两个基础概念
- 不会引入超出当前示例目标的复杂调优内容

## `http {}`

```nginx
http {
    ...
}
```

所有 HTTP 相关配置都放在这里，包括：

- 上游后端定义
- 虚拟主机定义
- 反向代理规则

## 三个 `upstream`

### `api_backend`

```nginx
upstream api_backend {
    server api-1:8080;
    server api-2:8080;
}
```

这表示 `api` 服务有两个实例：

- `api-1`
- `api-2`

Nginx 会把发往 `api_backend` 的请求分配给这两个实例。这个 demo 里它的用途是演示负载均衡，因此连续访问 `api.nginx.localtest.me:8081/info` 时，返回里的 `instance` 会在 `api-1` 和 `api-2` 之间切换。

这里的 `api-1` 和 `api-2` 不是外部域名，而是 Docker Compose 默认网络里的服务名。

### `admin_backend`

```nginx
upstream admin_backend {
    server admin:8080;
}
```

这里只有一个后端实例 `admin:8080`，所以它不承担负载均衡演示，只用来展示基于域名把请求分到另一个服务。

### `echo_backend`

```nginx
upstream echo_backend {
    server echo:8080;
}
```

这里只有一个后端实例 `echo:8080`。它的作用不是业务处理，而是把收到的请求信息回显出来，方便你验证：

- 路径是否命中预期服务
- `Host` 是否被正确透传
- `X-Forwarded-*` 头是否被正确设置

## 三个 `server`

### `api.nginx.localtest.me`

```nginx
server {
    listen 8080;
    server_name api.nginx.localtest.me;

    location / {
        proxy_pass http://api_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Host $host;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

这段的含义是：

- 只要请求进入 Nginx 容器的 `8080`
- 并且 `Host` 是 `api.nginx.localtest.me`
- 就把请求转发给 `api_backend`

因此下面这些请求都会进 `api` 服务：

- `http://api.nginx.localtest.me:8081/info`
- `http://api.nginx.localtest.me:8081/ping`

其中宿主机访问时看到的是 `8081`，因为 compose 把宿主机 `8081` 映射到了 Nginx 容器内部的 `8080`。

### `admin.nginx.localtest.me`

```nginx
server {
    listen 8080;
    server_name admin.nginx.localtest.me;

    location / {
        proxy_pass http://admin_backend;
        ...
    }
}
```

这段和 `api` 的写法一样，只是目标换成了 `admin_backend`。它用于演示最直接的 Host 分流：

- `api.nginx.localtest.me` -> `api`
- `admin.nginx.localtest.me` -> `admin`

### `nginx.localtest.me` + `/echo/`

```nginx
server {
    listen 8080;
    server_name nginx.localtest.me;

    location /echo/ {
        proxy_pass http://echo_backend;
        ...
    }
}
```

这段不是按不同域名切服务，而是按路径切服务。

它的含义是：

- 当 `Host` 是 `nginx.localtest.me`
- 并且路径命中 `/echo/`
- 就把请求转给 `echo_backend`

所以这类请求会进 `echo`：

- `http://nginx.localtest.me:8081/echo/inspect`

这里的重点是演示“路径分流”，而不是“另一个独立域名服务”。

## 代理头为什么要这样设置

每个 `location` 里都有这些头：

```nginx
proxy_set_header Host $host;
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Host $host;
proxy_set_header X-Forwarded-Proto $scheme;
```

它们的作用分别是：

- `Host`
  让后端知道用户访问的是哪个域名
- `X-Real-IP`
  让后端看到直接连到 Nginx 的客户端地址
- `X-Forwarded-For`
  记录转发链路中的客户端 IP
- `X-Forwarded-Host`
  让后端知道原始请求的 Host
- `X-Forwarded-Proto`
  让后端知道原始请求是 `http` 还是 `https`

这个 demo 的 `api` 和 `echo` 会把这些值直接返回出来，因此这些头不是可有可无，而是演示的一部分。

## 一次请求是怎么走的

以 `http://api.nginx.localtest.me:8081/info` 为例：

1. 浏览器或 `curl` 访问宿主机 `8081`
2. Docker 把请求转发到 Nginx 容器内的 `8080`
3. Nginx 根据 `server_name api.nginx.localtest.me` 命中 `api` 这段配置
4. Nginx 把请求转发给 `api_backend`
5. `api_backend` 在 `api-1:8080` 和 `api-2:8080` 之间选择一个实例
6. 后端返回 JSON，里面能看到 `service=api` 和具体 `instance`

## 为什么这里不需要给每个后端分不同宿主机端口

因为 `api-1`、`api-2`、`admin`、`echo` 都运行在各自独立容器里，它们监听的是各自容器内部的 `8080`，不会互相冲突。

真正暴露到宿主机的只有 Nginx 入口：

- 宿主机 `8081` -> Nginx 容器 `8080`

后端之间的访问发生在 Docker 内网，不需要宿主机端口参与。

## 看这份配置时最容易误解的点

- `listen 8080` 指的是 Nginx 容器内部端口，不是宿主机端口
- `api-1:8080`、`admin:8080` 这些地址是 Docker 网络内地址，不是外部可访问地址
- `nginx.localtest.me` 这一段只处理 `/echo/` 前缀，不代表所有路径都会转到 `echo`
- `api_backend` 是显式写死的后端列表，因此新增实例时需要改这份配置
