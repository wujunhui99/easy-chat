#!/bin/bash

docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
docker-compose down
docker stop $(docker ps -a -q)
export HOST_IP=$(ip route get 1 | awk '{print $7;exit}')
docker-compose up -d