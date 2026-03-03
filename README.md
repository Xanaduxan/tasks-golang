docker compose down

docker compose up -d


docker exec -i wallet-postgres psql -U wallet -d wallet < migrations/002_create_operations.sql