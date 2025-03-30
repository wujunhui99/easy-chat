#!/bin/bash
reso_addr='crpi-6zxn5tvxgfe9vkci.cn-shenzhen.personal.cr.aliyuncs.com/easy-chat-junhui/user-rpc-dev'
tag='latest'

# pod_ip="192.168.117.24"
pod_ip=$(hostname -I | awk '{print $1}')

container_name="easy-chat-user-rpc-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# docker run -p 10001:8080 --network imooc_easy-im -v /easy-im/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run -p 10000:10000 -e POD_IP=${pod_ip} --network easy-chat  --name=${container_name} -d ${reso_addr}:${tag}