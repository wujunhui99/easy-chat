#!/bin/bash  
need_start_server_shell=(
  # user
  user-rpc-test.sh
  user-api-test.sh

  # social
  social-rpc-test.sh
  social-api-test.sh

)

for i in ${need_start_server_shell[*]} ; do
    chmod +x $i
    ./$i
done


docker ps

docker exec -it etcd etcdctl get --prefix ""