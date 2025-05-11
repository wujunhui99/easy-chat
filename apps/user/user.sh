goctl rpc protoc ./apps/user/rpc/user.proto --go_out=./apps/user/rpc --go-grpc_out=./apps/user/rpc --zrpc_out=./apps/user/rpc

docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)

docker exec -it mysql sh
goctl api go -api apps/user/api/user.api -dir apps/user/api -style gozero
goctl api go -api apps/social/api/social.api -dir apps/social/api -style gozero
goctl rpc protoc ./apps/im/rpc/im.proto --go_out=./apps/im/rpc --go-grpc_out=./apps/im/rpc --zrpc_out=./apps/im/rpc

mysql -h 127.0.0.1 -P 13306 -u root -p
goctl rpc protoc ./apps/user/rpc/user.proto --go_out=./apps/user/rpc --go-grpc_out=./apps/user/rpc --zrpc_out=./apps/user/rpc

sudo docker login --username=344686925@qq.com crpi-6zxn5tvxgfe9vkci.cn-shenzhen.personal.cr.aliyuncs.com