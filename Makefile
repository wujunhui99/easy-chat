# Colors for output
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m

.PHONY: help dev-start start-safe start-rpc start-ws start-api start-mq stop restart status logs logs-tail clean-logs clean-all-logs setup-logrotate test-deps test-logs install-server install-server-user-rpc install-server-user-api install-server-social-rpc install-server-social-api

LOG_DIRS_ALL := \
	logs/user/api \
	logs/user/rpc \
	logs/social/api \
	logs/social/rpc \
	logs/im/api \
	logs/im/rpc \
	logs/msg/msggate \
	logs/msg/mq

LOG_DIRS_RPC := logs/user/rpc logs/social/rpc logs/im/rpc
LOG_DIRS_API := logs/user/api logs/social/api logs/im/api

## help: 显示帮助信息
help:
	@echo "$(GREEN)Easy Chat 微服务管理命令$(NC)"
	@echo ""
	@echo "$(YELLOW)启动管理 (按依赖关系):$(NC)"
	@echo "  dev-start      - 一键启动所有微服务 (自动处理依赖)"
	@echo "  start-safe     - 分阶段安全启动 (推荐)"
	@echo "  stop           - 停止所有微服务"
	@echo "  restart        - 重启所有服务"
	@echo ""
	@echo "$(YELLOW)分阶段启动:$(NC)"
	@echo "  start-rpc      - 第1阶段: 启动所有RPC服务"
	@echo "  start-ws       - 第2阶段: 启动WebSocket服务"
	@echo "  start-api      - 第3阶段: 启动所有API服务"
	@echo "  start-mq       - 第4阶段: 启动消息队列"
	@echo ""
	@echo "$(YELLOW)监控管理:$(NC)"
	@echo "  status         - 查看服务状态"
	@echo "  logs           - 查看最新日志"
	@echo "  logs-tail      - 实时查看日志"
	@echo "  clean-logs     - 清理旧日志 (7天前)"
	@echo "  clean-all-logs - 删除所有日志"
	@echo ""
	@echo "$(YELLOW)测试验证:$(NC)"
	@echo "  test-deps      - 测试依赖关系和环境配置"
	@echo "  test-logs      - 测试日志系统是否正常工作"
	@echo ""
	@echo "$(YELLOW)原有构建命令:$(NC)"
	@echo "  release-test   - 构建所有服务"
	@echo "  install-server - 部署服务"
	@echo ""

## dev-start: 一键启动所有微服务 (goreman自动处理启动顺序)
dev-start:
	@echo "$(GREEN)启动所有微服务 (开发模式)...$(NC)"
	@command -v goreman >/dev/null 2>&1 || { echo "$(RED)未检测到 goreman，请运行 'go install github.com/mattn/goreman@latest'$(NC)"; exit 1; }
	@mkdir -p $(LOG_DIRS_ALL)
	@echo "$(YELLOW)日志将保存到 logs/ 目录$(NC)"
	@echo "$(YELLOW)注意: goreman会按Procfile中的顺序启动服务$(NC)"
	@goreman start

## start-safe: 分阶段安全启动 (推荐用于生产环境调试)
start-safe:
	@echo "$(GREEN)开始分阶段启动微服务...$(NC)"
	@make start-rpc
	@echo "$(YELLOW)等待RPC服务启动完成 (3秒)...$(NC)"
	@sleep 3
	@make start-ws
	@echo "$(YELLOW)等待WebSocket服务启动完成 (3秒)...$(NC)"
	@sleep 3
	@make start-api
	@echo "$(YELLOW)等待API服务启动完成 (3秒)...$(NC)"
	@sleep 3
	@make start-mq
	@echo "$(GREEN)所有服务已按依赖关系启动完成!$(NC)"

## start-rpc: 启动所有RPC服务
start-rpc:
	@echo "$(GREEN)第1阶段: 启动RPC服务...$(NC)"
	@command -v goreman >/dev/null 2>&1 || { echo "$(RED)未检测到 goreman，请运行 'go install github.com/mattn/goreman@latest'$(NC)"; exit 1; }
	@mkdir -p $(LOG_DIRS_RPC)
	@goreman start user-rpc social-rpc im-rpc &

## start-ws: 启动WebSocket服务
start-ws:
	@echo "$(GREEN)第2阶段: 启动WebSocket服务...$(NC)"
	@command -v goreman >/dev/null 2>&1 || { echo "$(RED)未检测到 goreman，请运行 'go install github.com/mattn/goreman@latest'$(NC)"; exit 1; }
	@mkdir -p logs/msg/msggate
	@goreman start im-ws &

## start-api: 启动所有API服务
start-api:
	@echo "$(GREEN)第3阶段: 启动API服务...$(NC)"
	@command -v goreman >/dev/null 2>&1 || { echo "$(RED)未检测到 goreman，请运行 'go install github.com/mattn/goreman@latest'$(NC)"; exit 1; }
	@mkdir -p $(LOG_DIRS_API)
	@goreman start user-api social-api im-api &

## start-mq: 启动消息队列服务
start-mq:
	@echo "$(GREEN)第4阶段: 启动消息队列...$(NC)"
	@command -v goreman >/dev/null 2>&1 || { echo "$(RED)未检测到 goreman，请运行 'go install github.com/mattn/goreman@latest'$(NC)"; exit 1; }
	@mkdir -p logs/msg/mq
	@goreman start msg-mq &

## stop: 停止所有微服务
stop:
	@echo "$(RED)停止所有微服务...$(NC)"
	@pkill -f "goreman" || true
	@pkill -f "go run" || true
	@echo "$(GREEN)所有服务已停止$(NC)"

## restart: 重启所有服务
restart:
	@echo "$(YELLOW)重启所有微服务...$(NC)"
	@make stop
	@sleep 3
	@make dev-start

## status: 查看服务状态
status:
	@echo "$(GREEN)查看服务状态...$(NC)"
	@ps aux | grep -E "(goreman|go run)" | grep -v grep || echo "$(YELLOW)没有运行的服务$(NC)"

## logs: 查看最新日志
logs:
	@echo "$(GREEN)最新的日志文件:$(NC)"
	@find logs -name "*.log" -type f -exec ls -lt {} + | head -10 2>/dev/null || echo "$(YELLOW)暂无日志文件$(NC)"

## logs-tail: 实时查看日志
logs-tail:
	@echo "$(GREEN)实时查看所有服务日志 (Ctrl+C 退出)...$(NC)"
	@if find logs -name "*.log" -type f -print -quit 2>/dev/null | grep -q .; then \
		find logs -name "*.log" -type f -print0 2>/dev/null | xargs -0 tail -f 2>/dev/null; \
	else \
		echo "$(YELLOW)暂无日志文件$(NC)"; \
	fi

## clean-logs: 清理7天前的日志
clean-logs:
	@echo "$(YELLOW)执行日志轮转和清理...$(NC)"
	@./scripts/log-rotation.sh

## clean-all-logs: 删除所有日志文件
clean-all-logs:
	@echo "$(RED)警告: 将删除所有日志文件...$(NC)"
	@read -p "确认删除所有日志吗? [y/N]: " confirm; \
		if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
			echo "$(YELLOW)删除所有日志文件...$(NC)"; \
			for dir in $(LOG_DIRS_ALL); do \
				if [ -d "$$dir" ]; then \
					rm -rf "$$dir"/*.log 2>/dev/null || true; \
					echo "已清理: $$dir"; \
				fi; \
			done; \
			echo "$(GREEN)所有日志已删除$(NC)"; \
		else \
			echo "$(GREEN)操作已取消$(NC)"; \
		fi

## setup-logrotate: 设置定时日志轮转
setup-logrotate:
	@echo "$(GREEN)设置定时日志轮转...$(NC)"
	@mkdir -p logs
	@{ crontab -l 2>/dev/null | grep -v "scripts/log-rotation.sh"; echo "0 0 * * * $(PWD)/scripts/log-rotation.sh >> $(PWD)/logs/rotation.log 2>&1"; } | crontab -
	@echo "$(GREEN)已设置每日凌晨自动日志轮转 (如需移除请手动编辑 crontab)$(NC)"

## test-deps: 测试依赖关系和环境配置
test-deps:
	@echo "$(GREEN)测试系统依赖和配置...$(NC)"
	@./scripts/quick-test.sh

## test-logs: 测试日志系统是否正常工作
test-logs:
	@echo "$(GREEN)测试日志系统...$(NC)"
	@echo 'test-log: echo "日志测试 $$(date)" 2>&1 | tee logs/test-$$(date +%Y%m%d-%H%M%S).log' > Procfile.test
	@mkdir -p logs
	@goreman -f Procfile.test start &
	@sleep 3
	@pkill -f "goreman.*Procfile.test" || true
	@rm -f Procfile.test
	@if ls logs/test-*.log 1>/dev/null 2>&1; then \
		echo "$(GREEN)✅ 日志系统工作正常$(NC)"; \
		echo "$(YELLOW)测试日志文件:$(NC)"; \
		ls -la logs/test-*.log | head -1; \
		rm -f logs/test-*.log; \
	else \
		echo "$(RED)❌ 日志系统有问题$(NC)"; \
	fi

user-rpc-dev:
	@make -f deploy/mk/user-rpc.mk release-test
user-api-dev:
	@make -f deploy/mk/user-api.mk release-test
im-rpc-dev:
	@make -f deploy/mk/im-rpc.mk release-test
im-api-dev:
	@make -f deploy/mk/im-api.mk release-test
im-ws-dev:
	@make -f deploy/mk/im-ws.mk release-test
task-mq-dev:
	@make -f deploy/mk/task-mq.mk release-test
social-rpc-dev:
	@make -f deploy/mk/social-rpc.mk release-test
social-api-dev:
	@make -f deploy/mk/social-api.mk release-test
release-test: user-rpc-dev user-api-dev social-rpc-dev social-api-dev im-rpc-dev im-api-dev im-ws-dev task-mq-dev

install-server:
	cd ./deploy/script && chmod +x release-test.sh && ./release-test.sh

install-server-user-rpc:
	cd ./deploy/script && chmod +x user-rpc-test.sh && ./user-rpc-test.sh
install-server-user-api:
	cd ./deploy/script && chmod +x user-api-test.sh && ./user-api-test.sh
install-server-social-rpc:
	cd ./deploy/script && chmod +x social-rpc-test.sh && ./social-rpc-test.sh
install-server-social-api:
	cd ./deploy/script && chmod +x social-api-test.sh && ./social-api-test.sh
