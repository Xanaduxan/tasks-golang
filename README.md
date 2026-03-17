docker compose down

docker compose up -d


docker exec -i tasks-postgres psql -U tasks -d tasks < migrations/003_create_products_delivery_stocks.sql

docker run -d --name my-redis -p 6379:6379 redis