#!/bin/bash

need_start_server_shell=(
  # rpc
  "user-rpc-test.sh"
  # api
)

# 确保是在正确的目录下
cd "$(dirname "$0")"

for i in "${need_start_server_shell[@]}"; do
    if [ -f "$i" ]; then
        echo "Executing $i"
        chmod +x "$i"
        ./"$i"
    else
        echo "Warning: $i not found"
    fi
done

echo "Checking running containers:"
docker ps

echo "Checking etcd keys:"
docker exec -it etcd etcdctl get --prefix ""
