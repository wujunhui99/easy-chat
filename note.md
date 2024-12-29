### 编译go源码
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/user-rpc ./rpc/user.go
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/user-api ./api/user.go

### 编译镜像
docker build -t user-rpc -f ./Dockerfile_rpc .
docker build -t user-api -f ./Dockerfile_api .

### 启动镜像
docker run --name user-rpc -p 8080:8080 --net user --ip  168.10.0.70 --link etcd -d user-rpc
docker run --name user-api -p 8888:8888 --net user --ip  168.10.0.50 --link etcd -d user-api

### 接入mysql
mysql -h 127.0.0.1 -P 13306 -u root -p
