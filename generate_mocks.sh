#!/bin/bash

cd api || exit
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