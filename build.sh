#! /bin/bash
docker run --rm -v "$PWD":/go/src/test -w /go/src/test golang:1.10 go build -o myapp
docker build -t test:latest .
rm -rf myapp