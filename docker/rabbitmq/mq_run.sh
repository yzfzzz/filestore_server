docker stop rabbit > /dev/null
docker rm -v -f rabbit > /dev/null

docker run -d --hostname rabbit-server \
--name rabbit \
-p 5672:5672 \
-p 15672:15672 \
-p 25672:25672 \
-v ./data:/var/lib/rabbitmq \
rabbitmq:management

docker exec -it rabbit rabbitmqctl add_user yzf 123456
docker exec -it rabbit  rabbitmqctl set_user_tags yzf administrator  # 赋予管理员角色
docker exec -it rabbit  rabbitmqctl set_permissions -p / yzf ".*" ".*" ".*"  # 授予所有权限