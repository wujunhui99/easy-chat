#!/usr/bin/env bash

set -euo pipefail

# Easy Chat 日志轮转和清理脚本
# 使用 cron 定时执行: 0 0 * * * /path/to/log-rotation.sh

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
LOG_DIR="${SCRIPT_DIR}/../logs"
KEEP_DAYS=${KEEP_DAYS:-7}
MAX_SIZE_MB=${MAX_SIZE_MB:-100}

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "========== Easy Chat 日志轮转开始 $(date) =========="

mkdir -p "$LOG_DIR"

count_files() {
    find "$LOG_DIR" -name "*.log" -type f 2>/dev/null | wc -l | tr -d ' '
}

get_total_size() {
    if command -v du >/dev/null 2>&1; then
        du -sm "$LOG_DIR" 2>/dev/null | cut -f1 | tr -d ' '
    else
        echo "0"
    fi
}

file_size_mb() {
    local file="$1"
    if stat --version >/dev/null 2>&1; then
        local size
        size=$(stat --format=%s "$file" 2>/dev/null || echo 0)
        echo $(( size / 1024 / 1024 ))
    else
        local size
        size=$(stat -f%z "$file" 2>/dev/null || echo 0)
        echo $(( size / 1024 / 1024 ))
    fi
}

echo -e "${GREEN}日志目录: $LOG_DIR${NC}"
echo -e "${GREEN}保留天数: $KEEP_DAYS 天${NC}"
echo -e "${GREEN}单文件最大: $MAX_SIZE_MB MB${NC}"

INITIAL_FILES=$(count_files)
INITIAL_SIZE=$(get_total_size)

echo -e "${YELLOW}轮转前状态:${NC}"
echo "  文件数量: ${INITIAL_FILES:-0}"
echo "  总大小: ${INITIAL_SIZE:-0}MB"
echo ""

echo -e "${YELLOW}1. 清理过期文件 (${KEEP_DAYS}天前)...${NC}"
DELETED_OLD=$(find "$LOG_DIR" -name "*.log" -type f -mtime +"$KEEP_DAYS" -print | wc -l | tr -d ' ')
find "$LOG_DIR" -name "*.log" -type f -mtime +"$KEEP_DAYS" -delete
echo "  删除过期文件: ${DELETED_OLD:-0} 个"

echo -e "${YELLOW}2. 压缩大文件 (>${MAX_SIZE_MB}MB)...${NC}"
COMPRESSED=0
while IFS= read -r -d '' file; do
    if [[ -f "$file" ]]; then
        SIZE_MB=$(file_size_mb "$file")
        if [[ "$SIZE_MB" -gt "$MAX_SIZE_MB" ]]; then
            if gzip "$file" 2>/dev/null; then
                echo "  压缩: $(basename "$file") (${SIZE_MB}MB)"
                ((COMPRESSED++))
            fi
        fi
    fi
done < <(find "$LOG_DIR" -name "*.log" -type f -print0)
echo "  压缩文件: $COMPRESSED 个"

echo -e "${YELLOW}3. 清理空目录...${NC}"
DELETED_DIRS=$(find "$LOG_DIR" -type d -empty -print | wc -l | tr -d ' ')
find "$LOG_DIR" -type d -empty -delete 2>/dev/null
echo "  删除空目录: ${DELETED_DIRS:-0} 个"

echo -e "${YELLOW}4. 日志文件统计...${NC}"
if compgen -G "$LOG_DIR/*" > /dev/null 2>&1; then
    for service_dir in "$LOG_DIR"/*; do
        if [[ -d "$service_dir" ]]; then
            service_name=$(basename "$service_dir")
            for component_dir in "$service_dir"/*; do
                if [[ -d "$component_dir" ]]; then
                    component_name=$(basename "$component_dir")
                    log_count=$(find "$component_dir" -name "*.log*" -type f | wc -l | tr -d ' ')
                    component_size=$(du -sm "$component_dir" 2>/dev/null | cut -f1 | tr -d ' ')
                    echo "    $service_name/$component_name: ${log_count:-0} 文件, ${component_size:-0}MB"
                fi
            done
        fi
    done
else
    echo "  尚无日志文件"
fi

FINAL_FILES=$(count_files)
FINAL_SIZE=$(get_total_size)

echo ""
echo -e "${YELLOW}轮转后状态:${NC}"
initial_files=${INITIAL_FILES:-0}
initial_size=${INITIAL_SIZE:-0}
final_files=${FINAL_FILES:-0}
final_size=${FINAL_SIZE:-0}
reduced_files=$(( initial_files - final_files ))
reduced_size=$(( initial_size - final_size ))
echo "  文件数量: ${final_files} (减少 ${reduced_files})"
echo "  总大小: ${final_size}MB (减少 ${reduced_size}MB)"

if [[ $final_size -gt 1000 ]]; then
    echo -e "${RED}警告: 日志总大小超过 1GB，执行紧急清理...${NC}"
    find "$LOG_DIR" -name "*.log" -type f -mtime +3 -delete
    echo "  紧急清理完成"
fi

echo -e "${GREEN}========== 日志轮转完成 $(date) ==========${NC}"
