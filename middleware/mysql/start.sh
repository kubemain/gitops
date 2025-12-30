docker run -d --name ops-mysql --restart always -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -e MYSQL_USER=ops_user \
  -e MYSQL_PASSWORD=123456 \
  -v $(pwd)/my.cnf:/etc/mysql/conf.d/my.cnf \
  -v $(pwd)/mysql-data:/var/lib/mysql \
  harbor.1stcs.cn/x-ai-amd64/cloudstation-mysql:v8.0.34