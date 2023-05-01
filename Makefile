postgres:
	docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=rahasia -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root secondhand

dropdb:
	docker exec -it postgres15 dropdb secondhand

newmigrate:
	migrate create -ext sql -dir db/migration -seq

migrateup:
	migrate -path db/migration -database "postgresql://postgres:rahasia@localhost:5432/secondhand?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:rahasia@localhost:5432/secondhand?sslmode=disable" -verbose down 

test:
	go test -v -cover -short ./...

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: postgres createdb dropdb newmigrate migrateup migratedown test