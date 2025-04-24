#!/bin/bash
reso_addr='crpi-6zxn5tvxgfe9vkci.cn-shenzhen.personal.cr.aliyuncs.com/easy-chat-junhui/user-api-dev'
tag='latest'

pod_ip=$(hostname -I | awk '{print $1}')

container_name="easy-chat-task-mq-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
docker run   -e POD_IP=${pod_ip} -e HOST_IP=${pod_ip} --network easy-chat  --name=${container_name} -d ${reso_addr}:${tag}