redis:
	docker run -d -p 6379:6379 --name url-shortener-redis redis

postgres:
	docker run --name gorm -p 9920:5432 -e POSTGRES_USER=gorm -e POSTGRES_PASSWORD=gorm -d postgres