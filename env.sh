#!/bin/bash
export HOST_IP=$(ip route get 1 | awk '{print $7;exit}')
sudo chown -R 1001:1001 ./components/etcd/data ./components/etcd/logs
sudo chmod -R 700 ./components/etcd/data ./components/etcd/logs


