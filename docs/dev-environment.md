# 开发环境指南

本文档介绍如何使用 Docker Compose 快速搭建本地开发环境。

## 服务清单

| 服务 | 版本 | 端口 | 用途 |
|------|------|------|------|
| MySQL | 8.0 | 3306 | 业务数据库 |
| Redis | 7 (Alpine) | 6379 | 缓存/消息队列 |
| Nacos | 2.3.0 | 8848, 9848, 9849 | 服务注册与配置中心 |
| Apollo | Quick Start | 8070, 8080, 8090 | 配置中心 |
| Apollo DB | MySQL 8.0 | 13306 | Apollo 专用数据库 |

## 前置要求

- Docker 20.10+
- Docker Compose v2+
- 至少 4GB 可用内存

## 快速开始

### 启动服务

```bash
make dev-up
```

或直接使用脚本：

```bash
./scripts/dev-env.sh start
```

首次启动会自动初始化所有数据库，Apollo 启动可能需要 1-2 分钟。

### 停止服务

```bash
make dev-down
```

## 服务访问

### MySQL（业务数据库）

- **地址**: `localhost:3306`
- **用户名**: `root`
- **密码**: `root`
- **字符集**: `utf8mb4`

连接示例：

```bash
# 命令行连接
mysql -h 127.0.0.1 -P 3306 -uroot -proot

# 通过 Docker 连接
docker compose exec mysql mysql -uroot -proot
```

### Redis

- **地址**: `localhost:6379`
- **密码**: 无

连接示例：

```bash
# 命令行连接
redis-cli -h 127.0.0.1 -p 6379

# 通过 Docker 连接
docker compose exec redis redis-cli
```

### Nacos

- **控制台**: http://localhost:8848/nacos
- **用户名**: `nacos`
- **密码**: `nacos`

gRPC 端口：
- 9848: 客户端 gRPC 请求
- 9849: 服务端 gRPC 请求

### Apollo

- **Portal 控制台**: http://localhost:8070
- **用户名**: `apollo`
- **密码**: `admin`

服务端口：
- 8080: ConfigService（客户端读取配置）
- 8090: AdminService（Portal 管理接口）
- 8070: Portal（Web 控制台）

Apollo 使用独立的 MySQL 数据库（端口 13306），数据库会在首次启动时自动初始化。

> **ARM64 用户 (Mac M1/M2)**：如果拉取镜像失败，请修改 `docker-compose.yml` 中的 Apollo 镜像为：
> ```yaml
> image: nobodyiam/apollo-quick-start:arm64
> ```

## 数据存储

数据通过 Docker named volumes 持久化存储：

| Volume | 用途 |
|--------|------|
| `mysql_data` | 业务 MySQL 数据 |
| `redis_data` | Redis AOF 持久化 |
| `nacos_data` | Nacos 数据 |
| `apollo_db_data` | Apollo MySQL 数据 |

查看 volumes：

```bash
docker volume ls | grep kratos-layout
```

## 常用命令

### 使用 Makefile

```bash
# 启动开发环境
make dev-up

# 停止开发环境
make dev-down
```

### 使用脚本

```bash
# 启动服务
./scripts/dev-env.sh start

# 停止服务
./scripts/dev-env.sh stop

# 重启服务
./scripts/dev-env.sh restart

# 查看状态
./scripts/dev-env.sh status

# 查看日志
./scripts/dev-env.sh logs

# 清理所有数据（谨慎使用）
./scripts/dev-env.sh clean
```

### 使用 Docker Compose

```bash
# 启动
docker compose up -d

# 停止
docker compose down

# 查看日志
docker compose logs -f

# 查看特定服务日志
docker compose logs -f apollo-quick-start

# 重启特定服务
docker compose restart nacos

# 进入容器
docker compose exec mysql bash
```

## 配置说明

### MySQL 配置

```yaml
environment:
  MYSQL_ROOT_PASSWORD: root    # root 密码
  MYSQL_ROOT_HOST: "%"         # 允许远程连接
  TZ: Asia/Shanghai            # 时区
command:
  - --character-set-server=utf8mb4
  - --collation-server=utf8mb4_unicode_ci
```

### Nacos 配置

```yaml
environment:
  MODE: standalone             # 单机模式
  JVM_XMS: 256m               # JVM 初始堆内存
  JVM_XMX: 256m               # JVM 最大堆内存
  NACOS_AUTH_ENABLE: "false"  # 禁用鉴权（开发环境）
```

### Apollo 配置

Apollo 使用官方 Quick Start 镜像 (`nobodyiam/apollo-quick-start`)，该镜像包含：
- ConfigService
- AdminService
- Portal

数据库初始化脚本位于 `scripts/sql/apollo/` 目录。

参考文档：
- [Apollo Quick Start Docker 部署](https://www.apolloconfig.com/#/zh/deployment/quick-start-docker)
- [GitHub: apollo-quick-start](https://github.com/apolloconfig/apollo-quick-start)

## 故障排查

### Apollo 无法启动

1. 检查 Apollo DB 是否健康：
   ```bash
   docker compose ps apollo-db
   ```

2. 查看 Apollo 日志：
   ```bash
   docker compose logs apollo-quick-start
   ```

3. 如果是 ARM64 架构（Mac M1/M2），使用 ARM64 镜像

### 端口冲突

如果端口被占用，修改 `docker-compose.yml` 中的端口映射：

```yaml
ports:
  - "13306:3306"  # 改为其他端口
```

### 内存不足

可以调整各服务的 JVM 参数：

```yaml
environment:
  JVM_XMS: 128m
  JVM_XMX: 128m
```

## 生产环境注意

本开发环境配置**仅适用于本地开发**，不适合生产环境：

1. MySQL/Redis 未配置密码或使用弱密码
2. Nacos 禁用了鉴权
3. Apollo 使用默认账号密码
4. 未配置 TLS/SSL
5. 未配置数据备份

生产环境请参考各组件官方文档进行安全配置。
