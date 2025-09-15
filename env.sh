#!/bin/bash
# export HOST_IP=$(ip route get 1 | awk '{print $7;exit}')
export HOST_IP=$(ifconfig | grep -E 'inet ' | grep -v '127.0.0.1' | awk '{print $2}' | head -n 1)
# sudo chown -R 1001:1001 ./components/etcd/data ./components/etcd/logs
# sudo chmod -R 700 ./components/etcd/data ./components/etcd/logs


