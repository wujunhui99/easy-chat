#!/bin/bash

docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
docker-compose down
docker stop $(docker ps -a -q)
export HOST_IP=$(ip route get 1 | awk '{print $7;exit}')
docker-compose up -d
# cd apps/user/rpc && go run . &
# cd /root/code/easy-chat
# cd apps/user/api && go run . &
# cd /root/code/easy-chat
# cd apps/social/rpc && go run . &
# cd /root/code/easy-chat
# cd apps/social/api && go run . &