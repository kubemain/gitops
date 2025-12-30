docker run -d \
  --name ops-backend \
  --restart always \
  -p 8000:8000 \
  -e DB_HOST=ops-mysql \
  -e DB_PORT=3306 \
  -e DB_USER=ops_user \
  -e DB_PASSWORD=123456 \
  -e DB_NAME=ops_platform \
  -e REDIS_HOST=ops-redis \
  -e REDIS_PORT=6379 \
  -e JWT_SECRET=your-super-secret-key-change-it \
  -v $(pwd)/logs:/app/logs \
  --link ops-mysql:ops-mysql \
  --link ops-redis:ops-redis \
  registry.1stcs.cn/zc/ops-api:1.0.0