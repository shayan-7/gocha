#! /bin/bash

docker run \
    --detach \
    --name="gocha-db"\
    --env="MYSQL_ROOT_PASSWORD=my_password" \
    --env="MYSQL_DATABASE=gocha" \
    --env="MYSQL_USER=mysql" \
    --env="MYSQL_PASSWORD=mysql" \
    --env="MYSQL_ALLOW_EMPTY_PASSWORD=yes" \
    --publish 3306:3306 \
    mysql/mysql-server

docker run \
    --detach \
    --name="gocha-redis"\
    --publish 6379:6379 \
    redis

go run main.go $1
