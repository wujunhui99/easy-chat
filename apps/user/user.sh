goctl rpc protoc ./apps/user/rpc/user.proto --go_out=./apps/user/rpc --go-grpc_out=./apps/user/rpc --zrpc_out=./apps/user/rpc -style gozero

docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)

docker exec -it mysql sh
goctl api go -api apps/user/api/user.api -dir apps/user/api -style gozero
goctl api go -api apps/social/api/social.api -dir apps/social/api -style gozero
goctl rpc protoc ./apps/chat/rpc/chat.proto --go_out=./apps/chat/rpc --go-grpc_out=./apps/chat/rpc --zrpc_out=./apps/chat/rpc -style gozero

mysql -h 127.0.0.1 -P 13306 -u root -p
goctl rpc protoc ./apps/user/rpc/user.proto --go_out=./apps/user/rpc --go-grpc_out=./apps/user/rpc --zrpc_out=./apps/user/rpc -style gozero
goctl rpc protoc ./apps/social/rpc/social.proto --go_out=./apps/social/rpc --go-grpc_out=./apps/social/rpc --zrpc_out=./apps/social/rpc -style gozero

sudo docker login --username=344686925@qq.com crpi-6zxn5tvxgfe9vkci.cn-shenzhen.personal.cr.aliyuncs.com