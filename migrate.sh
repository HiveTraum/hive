#!/bin/bash

cd migrations || exit
goose postgres "user=hive dbname=hive_test sslmode=disable" up
goose postgres "user=hive dbname=hive sslmode=disable" up
