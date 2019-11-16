#!/bin/bash

cd infrastructure || exit
mockgen -source=infrastructure.go -destination=../mocks/infrastructure.go -package=mocks
