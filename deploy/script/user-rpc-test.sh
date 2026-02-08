#!/bin/bash
reso_addr='registry.cn-hangzhou.aliyuncs.com/0rz/im-user-rpc-dev'
tag='latest'

container_name="user-rpc-test"
pod_ip="192.168.2.26"

# 停止容器
docker stop ${container_name}

# 删除容器
docker rm ${container_name}

# 删除镜像
docker rmi ${reso_addr}:${tag}

# 拉取镜像
docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# docker run -p 10001:8080 --network imooc_easy-chat -v /easy-chat/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run -p 10000:10000 -e POD_IP=${pod_ip} --network zeroim_im --name=${container_name} -d ${reso_addr}:${tag}
