#!/bin/sh
sh ./build_go.sh
docker build -t ip2location-server -f server.Dockerfile .