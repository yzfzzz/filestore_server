docker stop mysql-slave > /dev/null
docker rm -v -f mysql-slave > /dev/null

docker run -p 3310:3306 \
--name mysql-slave \
-v ./3310/log:/var/log/mysql \
-v ./3310/data:/var/lib/mysql \
-v ./3310/conf:/etc/mysql/conf.d \
-e MYSQL_ROOT_PASSWORD=123 \
-d mysql:5.7

