# Nginx 与 Traefik Go 演示

演示同一组 `Hertz` Go 服务如何分别接入 `nginx` 和 `traefik`，重点覆盖：

- Host 分流
- 路径分流
- 多实例负载均衡
- 代理透传请求头

补充说明见 `deployments/nginx/nginx-config-explained.md` 和 `deployments/traefik/traefik-config-explained.md`。

## 架构

- `api-1`、`api-2`：两个 `api` 实例，用于负载均衡
- `admin`：管理端示例服务
- `echo`：请求回显服务
- `nginx`：静态配置代理，监听宿主机 `8081`
- `traefik`：Docker 动态发现代理，监听宿主机 `8092`

## 目录

- `cmd/api`：`api` 服务入口
- `cmd/admin`：`admin` 服务入口
- `cmd/echo`：`echo` 服务入口
- `internal/common`：共享响应和路由逻辑
- `docker-compose.yml`：整体编排
- `deployments/nginx/nginx.conf`：Nginx 路由
- `deployments/traefik/traefik.yml`：Traefik 静态配置

## 启动

Dockerfile 使用多阶段构建，在容器内自动编译 Go 程序，无需手动预编译。

```bash
docker compose up --build -d
```

或：

```bash
make up
```

## 停止

```bash
docker compose down
```

或：

```bash
make down
```

## 访问入口

为了能同时启动这两个，并没有让他们占用常用的80端口，而是各自分配了一个端口，因为使用域名时默认访问80端口，所以这里需要自己指定一下访问的端口

### Nginx

- `http://api.nginx.localtest.me:8081/info`
- `http://admin.nginx.localtest.me:8081/info`
- `http://nginx.localtest.me:8081/echo/inspect`

### Traefik

- `http://api.traefik.localtest.me:8092/info`
- `http://admin.traefik.localtest.me:8092/info`
- `http://traefik.localtest.me:8092/echo/inspect`
- `http://dashboard.traefik.localtest.me:8092/dashboard/`

## curl 验证

### Nginx Host 分流

```bash
curl -s http://api.nginx.localtest.me:8081/info
curl -s http://admin.nginx.localtest.me:8081/info
```

预期：
- 第一个返回 `"service":"api"`
- 第二个返回 `"service":"admin"`

### Nginx 路径分流

```bash
curl -s http://nginx.localtest.me:8081/echo/inspect
```

预期：
- 返回 `"service":"echo"`
- 返回 `"path":"/echo/inspect"`

### Traefik Host 分流

```bash
curl -s http://api.traefik.localtest.me:8092/info
curl -s http://admin.traefik.localtest.me:8092/info
```

预期：
- 第一个返回 `"service":"api"`
- 第二个返回 `"service":"admin"`

### Traefik 路径分流

```bash
curl -s http://traefik.localtest.me:8092/echo/inspect
```

预期：
- 返回 `"service":"echo"`
- 返回 `"path":"/echo/inspect"`

### 负载均衡

```bash
for i in $(seq 1 6); do curl -s http://api.nginx.localtest.me:8081/info; echo; done
for i in $(seq 1 6); do curl --noproxy '*' -s http://api.traefik.localtest.me:8092/info; echo; done
```

预期：
- 返回中的 `"instance"` 会在 `api-1` 和 `api-2` 之间切换

### 代理透传

返回 JSON 中应包含：

- `host`
- `client_ip`
- `forwarded.X-Forwarded-For`
- `forwarded.X-Forwarded-Proto`

