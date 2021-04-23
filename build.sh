#! /bin/bash

mkdir bin
go build -o bin/main .
docker build -t gocha .
