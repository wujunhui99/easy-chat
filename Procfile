# Easy Chat Microservices Procfile
# 按依赖关系定义启动顺序，使用脚本包装器确保日志正确写入
# 正确启动顺序: user-rpc -> social-rpc -> im-rpc -> im-ws -> user-api -> social-api -> im-api -> task-mq

# Phase 1: RPC Services (独立，无依赖)
# 这些服务必须先启动，因为API服务依赖它们
user-rpc: ./scripts/run-service.sh apps/user/rpc logs/user/rpc user-rpc
social-rpc: ./scripts/run-service.sh apps/social/rpc logs/social/rpc social-rpc
im-rpc: ./scripts/run-service.sh apps/im/rpc logs/im/rpc im-rpc

# Phase 2: WebSocket Service (无依赖)
im-ws: ./scripts/run-service.sh apps/msg/msggate logs/msg/msggate im-ws

# Phase 3: API Services (依赖RPC服务)
# 需要等待RPC服务完全启动后再启动
user-api: ./scripts/run-service.sh apps/user/api logs/user/api user-api
social-api: ./scripts/run-service.sh apps/social/api logs/social/api social-api
im-api: ./scripts/run-service.sh apps/im/api logs/im/api im-api

# Phase 4: Message Queue (依赖ws服务)
# 最后启动，处理消息队列
msg-mq: ./scripts/run-service.sh apps/msg/mq logs/msg/mq msg-mq
