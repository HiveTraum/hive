#!/bin/bash

cd migrations || exit
goose postgres "user=auth dbname=auth_test sslmode=disable" up
goose postgres "user=auth dbname=auth sslmode=disable" up
