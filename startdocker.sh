#!/bin/bash
docker build -t interactive .
docker build -t follower -f Dockerfile-follower .
screen -d -m bash -c "docker run -it interactive"
sleep 2
for (( i = 0; i < 10; i++ )); do
  echo "Starting one follower"
  docker run -d follower
done
screen -r
