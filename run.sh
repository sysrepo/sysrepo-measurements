#!/bin/bash

docker pull sysrepo/sysrepo-netopeer2:devel

pwd_dir=$(pwd)
docker run -d -v $pwd_dir/docker:/opt/docker -p 830:830 --name sysrepo --rm sysrepo/sysrepo-netopeer2:devel
docker exec -d sysrepo bash -c /opt/docker/docker_entry_point.sh

# wait for sysrepo and netopeer-server to start inside docker
printf "\n############# wait for initialization #############\n"
sleep 2

printf "\n############# start the tests #############\n"
go run ./tests/simple_tests.go
go run ./tests/multiple_tests.go

printf "\n############# exit #############\n"
docker stop sysrepo
