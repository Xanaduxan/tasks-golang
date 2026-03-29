docker compose down

docker compose up -d


docker exec -i tasks-postgres psql -U tasks -d tasks < migrations/003_create_products_delivery_stocks.sql

docker run -d --name my-redis -p 6379:6379 redis

docker run -d   --name redis-ui   --network mynet   -p 8081:8081   -e REDIS_HOSTS=local:my-redis:6379   rediscommander/redis-commander

protoc \
--proto_path=. \
--go_out=task-service \
--go_opt=module=github.com/Xanaduxan/tasks-golang/task-service \
--go-grpc_out=task-service \
--go-grpc_opt=module=github.com/Xanaduxan/tasks-golang/task-service \
proto/task/v1/task.proto

protoc \
--proto_path=. \
--go_out=notification-service \
--go_opt=module=github.com/Xanaduxan/tasks-golang/notification-service \
--go-grpc_out=notification-service \
--go-grpc_opt=module=github.com/Xanaduxan/tasks-golang/notification-service \
proto/notification/v1/notification.proto