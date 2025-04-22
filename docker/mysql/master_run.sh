docker stop mysql-master > /dev/null
docker rm -v -f mysql-master > /dev/null

docker run -p 3306:3306 \
--name mysql-master \
-v ./3306/log:/var/log/mysql \
-v ./3306/data:/var/lib/mysql \
-v ./3306/conf:/etc/mysql/conf.d \
-e MYSQL_ROOT_PASSWORD=123 \
-d mysql:5.7

