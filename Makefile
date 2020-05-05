#!make
include .env

run:
	cd src && go run hive

build:
	cd src && go build hive

protobuf:
	cd src/inout && protoc --go_out=. ./*.proto

coverage:
	cd src && go tool cover -html=coverage.out

mocks:
	cd src/api && mockgen -source=main.go -destination=../api/mocks.go -package=api
	cd src/controllers && mockgen -source=main.go -destination=../controllers/mocks.go -package=controllers
	cd src/eventDispatchers && mockgen -source=main.go -destination=../eventDispatchers/mocks.go -package=eventDispatchers
	cd src/passwordProcessors && mockgen -source=main.go -destination=../passwordProcessors/mocks.go -package=passwordProcessors
	cd src/stores && mockgen -source=main.go -destination=../stores/mocks.go -package=stores
	cd src/auth && mockgen -source=main.go -destination=../auth/mocks.go -package=auth
	cd src/repositories/inMemoryRepository && mockgen -source=main.go -destination=./mocks.go -package=inMemoryRepository
	cd src/repositories/postgresRepository && mockgen -source=main.go -destination=./mocks.go -package=postgresRepository
	cd src/repositories/redisRepository && mockgen -source=main.go -destination=./mocks.go -package=redisRepository