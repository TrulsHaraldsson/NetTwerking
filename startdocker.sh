#!/bin/bash
docker build -t interactive .
docker build -t follower -f Dockerfile-follower .
screen -d -m bash -c "docker run -it interactive"
sleep 2
for (( i = 0; i < 49; i++ )); do
  docker run -d follower
done
screen -r
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
