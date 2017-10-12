#!/bin/bash
docker build -t interactive .
docker build -t follower -f Dockerfile-follower .
screen -d -m bash -c "docker run -it -p 7999 interactive"

for (( i = 0; i < 99; i++ )); do
  sleep 1
  docker run -d -p 7999 follower
done
screen -r
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
