docker compose down

docker compose up -d


docker exec -i tasks-postgres psql -U tasks -d tasks < migrations/003_create_products_delivery_stocks.sql

docker run -d --name my-redis -p 6379:6379 redis

docker run -d   --name redis-ui   --network mynet   -p 8081:8081   -e REDIS_HOSTS=local:my-redis:6379   rediscommander/redis-commander

