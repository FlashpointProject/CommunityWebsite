include .env
export $(shell sed 's/=.*//' .env)

db:
	docker-compose -p fpcomm -f dc-db.yml down
	docker-compose -p fpcomm -f dc-db.yml up -d

migrate:
	docker run --rm -v $(shell pwd)/postgres_migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}?sslmode=disable" up

migrate-to:
	docker run --rm -v $(shell pwd)/postgres_migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}?sslmode=disable" goto $(MIGRATION)

rebuild-postgres:
	docker-compose -p fpcomm down
	docker volume rm fpcomm_fpcomm_postgres_data
	docker-compose -p fpcomm -f dc-db.yml up -d
	sleep 6
	make migrate

run:
	/usr/local/go/bin/go run ./main/*.go