docker run -d \
  --name ops-redis \
  --restart always \
  -p 6379:6379 \
  -v $(pwd)/docker-data/redis-data:/data \
  -v $(pwd)/conf/redis.conf:/usr/local/etc/redis/redis.conf \
  registry.1stcs.cn/zc/redis:7-alpine \
  redis-server /usr/local/etc/redis/redis.conf