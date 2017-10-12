#!/bin/bash
docker build -t interactive .
docker build -t follower -f Dockerfile-follower .
screen -d -m bash -c "docker run -it interactive"

for (( i = 0; i < 99; i++ )); do
  sleep 1
  docker run -d follower
done
screen -r
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
