#!/usr/bin/env bash

set -euo pipefail

# 微服务启动包装器脚本，用于处理日志输出
# 用法: ./run-service.sh <service_path> <log_path> <service_name>

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
REPO_ROOT=$(cd "${SCRIPT_DIR}/.." && pwd)

SERVICE_PATH="${1:-}"
LOG_PATH="${2:-}"
SERVICE_NAME="${3:-}"

if [[ -z "$SERVICE_PATH" || -z "$LOG_PATH" || -z "$SERVICE_NAME" ]]; then
    echo "用法: $0 <service_path> <log_path> <service_name>"
    exit 1
fi

if [[ "$SERVICE_PATH" != /* ]]; then
    SERVICE_PATH="${REPO_ROOT}/${SERVICE_PATH}"
fi

if [[ "$LOG_PATH" != /* ]]; then
    LOG_PATH="${REPO_ROOT}/${LOG_PATH}"
fi

if [[ ! -d "$SERVICE_PATH" ]]; then
    echo "未找到服务目录: $SERVICE_PATH" >&2
    exit 1
fi

mkdir -p "$LOG_PATH"

TIMESTAMP=$(date +%Y%m%d-%H%M%S)
LOG_FILE="${LOG_PATH}/${SERVICE_NAME}-${TIMESTAMP}.log"

if [[ -z "${HOST_IP:-}" && -f "${REPO_ROOT}/env.sh" ]]; then
    # shellcheck disable=SC1090
    source "${REPO_ROOT}/env.sh"
fi

HOST_IP=${HOST_IP:-127.0.0.1}

echo "启动服务: $SERVICE_NAME"
echo "工作目录: $SERVICE_PATH"
echo "日志文件: $LOG_FILE"
echo "HOST_IP: $HOST_IP"

cd "$SERVICE_PATH"

cleanup() {
    pkill -P $$ 2>/dev/null || true
}

trap cleanup SIGINT SIGTERM

(
    env HOST_IP="$HOST_IP" go run .
) 2>&1 | tee "$LOG_FILE"

exit "${PIPESTATUS[0]}"
