goctl rpc protoc ./apps/user/rpc/user.proto --go_out=./apps/user/rpc --go-grpc_out=./apps/user/rpc --zrpc_out=./apps/user/rpc

docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)

docker exec -it mysql sh