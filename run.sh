#!/bin/bash

docker pull sysrepo/sysrepo-netopeer2:devel
docker pull patavee/scipy-matplotlib

pwd_dir=$(pwd)
docker run -d -v $pwd_dir/docker:/opt/docker -p 10830:830 --name sysrepo sysrepo/sysrepo-netopeer2:devel
docker exec -d sysrepo bash -c /opt/docker/docker_entry_point.sh

# wait for sysrepo and netopeer-server to start inside docker
printf "\n############# wait for initialization #############\n"
sleep 2

printf "\n############# start the tests #############\n"
go run ./tests/simple_tests.go
go run ./tests/multiple_tests.go

printf "\ngenerate graph\n"

docker run -i -t -v $pwd_dir:/opt/graph --rm patavee/scipy-matplotlib:latest bash -c "pip install configparser; cd /opt/graph; python config.py;"
printf "\n############# exit #############\n"
docker stop sysrepo
docker rm sysrepo
