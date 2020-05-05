#!make
include .env

run:
	cd src && go run hive

build:
	cd src && go build hive

migrate: migrate-db-up migrate-test-up

migrate-db-up:
	cd src/migrations
	goose postgres "${DATABASE_URI}" up

migrate-db-down:
	cd src/migrations
	goose postgres "${DATABASE_URI}" down

migrate-test-up:
	cd src/migrations
	goose postgres "${DATABASE_URI}_test" up

migrate-test-down:
	cd src/migrations
	goose postgres "${DATABASE_URI}_test" down

protobuf:
	cd src/inout && protoc --go_out=. ./*.proto

coverage:
	go tool cover -html=coverage.out

mocks:
	cd src/api || exit
	mockgen -source=main.go -destination=../api/mocks.go -package=api
	cd ../controllers || exit
	mockgen -source=main.go -destination=../controllers/mocks.go -package=controllers
	cd ../eventDispatchers || exit
	mockgen -source=main.go -destination=../eventDispatchers/mocks.go -package=eventDispatchers
	cd ../passwordProcessors || exit
	mockgen -source=main.go -destination=../passwordProcessors/mocks.go -package=passwordProcessors
	cd ../stores || exit
	mockgen -source=main.go -destination=../stores/mocks.go -package=stores
	cd ../auth || exit
	mockgen -source=main.go -destination=../auth/mocks.go -package=auth
	cd ../repositories/inMemoryRepository || exit
	mockgen -source=main.go -destination=./mocks.go -package=inMemoryRepository
	cd ../repositories/postgresRepository || exit
	mockgen -source=main.go -destination=./mocks.go -package=postgresRepository
	cd ../repositories/redisRepository || exit
	mockgen -source=main.go -destination=./mocks.go -package=redisRepository