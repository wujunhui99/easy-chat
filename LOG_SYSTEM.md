# Easy Chat 日志系统说明

## 概述

新的日志系统解决了原有的问题：
- ✅ 一键启动所有微服务（基于 goreman 多进程管理）
- ✅ 日志自动输出到文件
- ✅ 按服务分类存储
- ✅ 自动日志轮转和清理
- ✅ 实时日志查看

## 目录结构

```
logs/
├── user/
│   ├── api/           # 用户API服务日志
│   └── rpc/           # 用户RPC服务日志
├── social/
│   ├── api/           # 社交API服务日志
│   └── rpc/           # 社交RPC服务日志
├── im/
│   ├── api/           # 即时消息API服务日志
│   ├── rpc/           # 即时消息RPC服务日志
│   └── ws/            # WebSocket服务日志
└── msg/
    └── mq/            # 消息队列服务日志
```

## 微服务依赖关系

微服务必须按以下顺序启动以确保依赖关系正确：

```
Phase 1: RPC Services (独立，无依赖)
├── user-rpc      (用户RPC服务)
├── social-rpc    (社交RPC服务)
└── im-rpc        (即时消息RPC服务)

Phase 2: WebSocket Service (无依赖)
└── im-ws         (WebSocket服务)

Phase 3: API Services (依赖RPC服务)
├── user-api      (依赖 user-rpc)
├── social-api    (依赖 social-rpc, user-rpc)
└── im-api        (依赖 im-rpc, social-rpc, user-rpc)

Phase 4: Message Queue (依赖其他所有服务)
└── msg-mq       (消息队列处理)
```

## 核心命令

### 1. 推荐启动方式
```bash
# 一键启动 (依赖 goreman)
make dev-start

# 分阶段安全启动 (推荐生产调试)
make start-safe
```

### 2. 分阶段手动启动
```bash
# 第1阶段：启动RPC服务
make start-rpc

# 第2阶段：启动WebSocket服务
make start-ws

# 第3阶段：启动API服务
make start-api

# 第4阶段：启动消息队列
make start-mq
```

### 3. 服务管理
```bash
# 停止所有服务
make stop

# 重启所有服务
make restart

# 查看服务状态
make status
```

### 4. 日志管理
```bash
# 查看最新日志文件
make logs

# 实时查看日志 (类似 tail -f)
make logs-tail

# 清理和轮转日志
make clean-logs
```

## 服务端口信息

| 服务 | 类型 | 端口 | 说明 |
|-----|------|------|------|
| user-api | HTTP API | 8888 | 用户相关API |
| user-rpc | gRPC | 10000 | 用户RPC服务 |
| social-api | HTTP API | 8881 | 社交功能API |
| social-rpc | gRPC | - | 社交RPC服务 |
| im-api | HTTP API | 8882 | 即时消息API |
| im-rpc | gRPC | - | 即时消息RPC |
| im-ws | WebSocket | 10090 | WebSocket连接 |
| msg-mq | MQ Consumer | - | 消息队列处理 |

## 日志轮转机制

### 自动轮转规则
- **保留天数**: 7天
- **单文件大小限制**: 100MB
- **自动压缩**: 大文件自动gzip压缩
- **定时清理**: 每日凌晨自动执行

### 设置定时轮转
```bash
make setup-logrotate
```
会在保留现有计划任务的前提下追加日志轮转任务，如需调整可手动编辑 `crontab -e`。

### 手动执行轮转
```bash
./scripts/log-rotation.sh
```

## 文件命名规则

日志文件按以下格式命名：
```
{service}-{component}-{YYYYMMDD-HHMMSS}.log
```

例如：
- `user-api-20240924-143022.log`
- `im-ws-20240924-143025.log`

## 故障排除

### 1. goreman 未安装
```bash
go install github.com/mattn/goreman@latest
```
请确保 `$GOPATH/bin` 已加入 `PATH`，`make` 目标会自动检查依赖。

运行 `make dev-start` 之前可手动设置 `HOST_IP`，若未设置脚本会自动检测；检测失败时默认使用 127.0.0.1。

### 2. 端口被占用
```bash
# 查看端口占用
lsof -i :8888
lsof -i :8881
lsof -i :8882
lsof -i :10090

# 杀死占用进程
kill -9 <PID>
```

### 3. 日志目录权限问题
```bash
chmod -R 755 logs/
```

### 4. 服务启动失败
1. 检查配置文件中的 `HOST_IP` 环境变量
2. 确保依赖服务(MySQL, Redis, etcd等)已启动
3. 查看具体服务的日志文件排查问题
4. 如果使用 `make dev-start`，请确认 `env.sh` 能在本机解析正确的 IP

## 环境变量

服务启动脚本会在运行前自动 source 根目录下的 `env.sh`，并在缺失时回退到 `127.0.0.1`。如需覆盖默认值，请提前导出：
```bash
export HOST_IP=<你的宿主机IP>
```

也可以运行 `source env.sh` 以自动探测宿主机 IP（脚本会优先使用无线网卡 en0/en1 的地址，若检测结果不准确，可手动导出 `HOST_IP` 覆盖）。

## 开发建议

1. **启动顺序**: 先启动基础设施(docker-compose up)，再启动微服务
2. **调试单个服务**: 可以cd到具体服务目录下单独运行
3. **日志分析**: 使用logs目录下的文件进行问题排查
4. **性能监控**: 可集成到现有的Jaeger链路追踪系统

## 与原系统的兼容性

- 保留原有的`make release-test`等构建命令
- 新增的开发命令不影响生产部署
- 可以与现有的Docker Compose环境并存使用